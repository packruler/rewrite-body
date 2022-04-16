// Package httputil a package for handling http data tasks
package httputil

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/packruler/rewrite-body/compressutil"
)

// ResponseWrapper a wrapper used to simplify ResponseWriter data access and manipulation.
type ResponseWrapper struct {
	buffer       bytes.Buffer
	lastModified bool
	wroteHeader  bool

	http.ResponseWriter
}

// WriteHeader into wrapped ResponseWriter.
func (wrapper *ResponseWrapper) WriteHeader(statusCode int) {
	if !wrapper.lastModified {
		wrapper.ResponseWriter.Header().Del("Last-Modified")
	}

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
	encoding := wrapper.GetContentEncoding()

	return compressutil.Decode(wrapper.GetBuffer(), encoding)
}

// SetContent write data to the internal ResponseWriter buffer
// and match initial encoding.
func (wrapper *ResponseWrapper) SetContent(data []byte) {
	encoding := wrapper.GetContentEncoding()

	bodyBytes, _ := compressutil.Encode(data, encoding)

	if !wrapper.wroteHeader {
		wrapper.WriteHeader(http.StatusOK)
	}

	if _, err := wrapper.ResponseWriter.Write(bodyBytes); err != nil {
		log.Printf("unable to write rewrited body: %v", err)
	}
}

// SupportsProcessing determine if http.Request is supported by this plugin.
func SupportsProcessing(request *http.Request) bool {
	// Ignore non GET requests
	if request.Method != "GET" {
		return false
	}

	if strings.Contains(request.Header.Get("Upgrade"), "websocket") {
		log.Printf("Ignoring websocket request for %s", request.RequestURI)

		return false
	}

	return true
}

// GetContentEncoding get the Content-Encoding header value.
func (wrapper *ResponseWrapper) GetContentEncoding() string {
	return wrapper.Header().Get("Content-Encoding")
}

// GetContentType get the Content-Encoding header value.
func (wrapper *ResponseWrapper) GetContentType() string {
	return wrapper.Header().Get("Content-Type")
}

// SupportsProcessing determine if HttpWrapper is supported by this plugin based on encoding.
func (wrapper *ResponseWrapper) SupportsProcessing() bool {
	contentType := wrapper.GetContentType()

	// If content type does not match return values with false
	if contentType != "" && !strings.Contains(contentType, "text") {
		return false
	}

	encoding := wrapper.GetContentEncoding()

	// If content type is supported validate encoding as well
	switch encoding {
	case "gzip":
		fallthrough
	case "deflate":
		fallthrough
	case "identity":
		fallthrough
	case "":
		return true
	default:
		return false
	}
}

// SetLastModified update the local lastModified variable from non-package-based users.
func (wrapper *ResponseWrapper) SetLastModified(value bool) {
	wrapper.lastModified = value
}
