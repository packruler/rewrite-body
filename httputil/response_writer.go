// Package httputil a package for handling http data tasks
package httputil

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strings"
        "regexp"

	"github.com/joinrepublic/traefik-rewrite-body-csp/compressutil"
	"github.com/joinrepublic/traefik-rewrite-body-csp/logger"
)

// ResponseWrapper a wrapper used to simplify ResponseWriter data access and manipulation.
type ResponseWrapper struct {
	buffer       bytes.Buffer
	lastModified bool `default:"true"`
	wroteHeader  bool

	code int `default:"200"`

	logWriter  logger.LogWriter
	monitoring MonitoringConfig

        cspPlaceholder *regexp.Regexp
        generateNonce func(string) []byte

	http.ResponseWriter
}

// WrapWriter create a ResponseWrapper for provided configuration.
func WrapWriter(
	responseWriter http.ResponseWriter,
	monitoringConfig MonitoringConfig,
	logWriter logger.LogWriter,
	lastModified bool,
        cspPlaceholder *regexp.Regexp,
        generateNonce func(string) []byte,
) *ResponseWrapper {
	return &ResponseWrapper{
		buffer:         bytes.Buffer{},
		lastModified:   lastModified,
		wroteHeader:    false,
		code:           http.StatusOK,
		logWriter:      logWriter,
		monitoring:     monitoringConfig,
		ResponseWriter: responseWriter,
                cspPlaceholder: cspPlaceholder,
                generateNonce:  generateNonce,
	}
}

func (wrapper *ResponseWrapper) ContainsCSP() bool {
        return wrapper.GetHeader("content-security-policy") != "" || wrapper.GetHeader("content-security-policy-report-only") != ""
}
func (wrapper *ResponseWrapper) overrideCSPHeaders() {
        nonce := generateNonceString()

        csp := wrapper.GetHeader("content-security-policy")
        cspReportOnly := wrapper.GetHeader("content-security-policy-report-only")

        replacement := wrapper.generateNonce(nonce)

        if csp != "" {
                wrapper.Header().Del("content-security-policy")
                wrapper.SetHeader(
                        "content-security-policy", 
                        string(wrapper.cspPlaceholder.ReplaceAll([]byte(csp), replacement)),
                )
        }

        if cspReportOnly != "" {
                wrapper.Header().Del("content-security-policy-report-only")
                wrapper.SetHeader(
                        "content-security-policy-report-only", 
                        string(wrapper.cspPlaceholder.ReplaceAll([]byte(cspReportOnly), replacement)),
                )
        }
}

// WriteHeader into wrapped ResponseWriter.
func (wrapper *ResponseWrapper) WriteHeader(statusCode int) {
	if wrapper.wroteHeader {
		return
	}

        if wrapper.ContainsCSP() {
          wrapper.overrideCSPHeaders()
        }

	if !wrapper.lastModified {
		wrapper.ResponseWriter.Header().Del("Last-Modified")
	}

	wrapper.code = statusCode
	wrapper.wroteHeader = true

	// Delegates the Content-Length Header creation to the final body write.
	wrapper.ResponseWriter.Header().Del("Content-Length")

	wrapper.ResponseWriter.WriteHeader(statusCode)
}

// Write data to internal buffer and mark the status code as http.StatusOK.
func (wrapper *ResponseWrapper) Write(data []byte) (int, error) {
	if !wrapper.wroteHeader {
		wrapper.WriteHeader(http.StatusOK)
	}

	return wrapper.buffer.Write(data)
}

// GetBuffer get a pointer to the ResponseWriter buffer.
func (wrapper *ResponseWrapper) GetBuffer() *bytes.Buffer {
	return &wrapper.buffer
}

// GetContent load the content currently in the internal buffer
// accounting for possible encoding.
func (wrapper *ResponseWrapper) GetContent() ([]byte, error) {
	encoding := wrapper.getContentEncoding()

	return compressutil.Decode(wrapper.GetBuffer(), encoding)
}

// SetContent write data to the internal ResponseWriter buffer
// and match initial encoding.
func (wrapper *ResponseWrapper) SetContent(data []byte, encoding string) {
	bodyBytes, _ := compressutil.Encode(data, encoding)

	if !wrapper.wroteHeader {
		wrapper.WriteHeader(http.StatusOK)
	}

	if _, err := wrapper.ResponseWriter.Write(bodyBytes); err != nil {
		wrapper.logWriter.LogErrorf("unable to write rewriten body: %v", err)
		wrapper.LogHeaders()
	}
}

func (wrapper *ResponseWrapper) GetHeader(headerName string) string {
	return wrapper.ResponseWriter.Header().Get(headerName)
}

func (wrapper *ResponseWrapper) SetHeader(headerName string, newValue string) {
	wrapper.ResponseWriter.Header().Set(headerName, newValue)
}

// LogHeaders writes current response headers.
func (wrapper *ResponseWrapper) LogHeaders() {
	wrapper.logWriter.LogDebugf("Error Headers: %v", wrapper.ResponseWriter.Header())
}

// getContentEncoding get the Content-Encoding header value.
func (wrapper *ResponseWrapper) getContentEncoding() string {
	return wrapper.GetHeader("Content-Encoding")
}

// getContentType get the Content-Encoding header value.
func (wrapper *ResponseWrapper) getContentType() string {
	return wrapper.GetHeader("Content-Type")
}

// SupportsProcessing determine if HttpWrapper is supported by this plugin based on encoding.
func (wrapper *ResponseWrapper) SupportsProcessing() bool {
	foundContentType := false

	// If content type does not match return values with false
	contentType := wrapper.getContentType()
	for _, monitoredType := range wrapper.monitoring.Types {
		if strings.Contains(contentType, monitoredType) {
			foundContentType = true

			break
		}
	}

	if !foundContentType {
		return false
	}

	encoding := wrapper.getContentEncoding()

	// If content type is supported validate encoding as well
	switch encoding {
	case compressutil.Gzip, compressutil.Deflate, compressutil.Identity, "":
		return true
	default:
		return false
	}
}

// SetLastModified update the local lastModified variable from non-package-based users.
func (wrapper *ResponseWrapper) SetLastModified(value bool) {
	wrapper.lastModified = value
}

// CloseNotify returns a channel that receives at most a
// single value (true) when the client connection has gone away.
func (wrapper *ResponseWrapper) CloseNotify() <-chan bool {
	if w, ok := wrapper.ResponseWriter.(http.CloseNotifier); ok {
		return w.CloseNotify()
	}

	return make(<-chan bool)
}

// Hijack hijacks the connection.
func (wrapper *ResponseWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := wrapper.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}

	return nil, nil, fmt.Errorf("%T is not a http.Hijacker", wrapper.ResponseWriter)
}

// Flush sends any buffered data to the client.
func (wrapper *ResponseWrapper) Flush() {
	// If WriteHeader was already called from the caller, this is a NOOP.
	// Otherwise, codeCatcher.code is actually a 200 here.
	wrapper.WriteHeader(wrapper.code)

	if flusher, ok := wrapper.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
