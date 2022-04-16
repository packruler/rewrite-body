package rewrite_body

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/packruler/rewrite-body/compressutil"
)

// ResponseWrapper stuff.
type ResponseWrapper struct {
	buffer       bytes.Buffer
	lastModified bool
	wroteHeader  bool

	http.ResponseWriter
}

// WriteHeader stuff.
func (wrappedWriter *ResponseWrapper) WriteHeader(statusCode int) {
	if !wrappedWriter.lastModified {
		wrappedWriter.ResponseWriter.Header().Del("Last-Modified")
	}

	wrappedWriter.wroteHeader = true

	// Delegates the Content-Length Header creation to the final body write.
	wrappedWriter.ResponseWriter.Header().Del("Content-Length")

	wrappedWriter.ResponseWriter.WriteHeader(statusCode)
}

// Write stuff.
func (wrappedWriter *ResponseWrapper) Write(p []byte) (int, error) {
	if !wrappedWriter.wroteHeader {
		wrappedWriter.WriteHeader(http.StatusOK)
	}

	return wrappedWriter.buffer.Write(p)
}

// Hijack stuff.
// func (wrappedWriter *ResponseWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
// 	hijacker, ok := wrappedWriter.ResponseWriter.(http.Hijacker)
// 	if !ok {
// 		return nil, nil, fmt.Errorf("%T is not a http.Hijacker", wrappedWriter.ResponseWriter)
// 	}

// 	return hijacker.Hijack()
// }

// Flush stuff.
// func (wrappedWriter *ResponseWrapper) Flush() {
// 	if flusher, ok := wrappedWriter.ResponseWriter.(http.Flusher); ok {
// 		flusher.Flush()
// 	}
// }

// GetBuffer stuff.
func (wrappedWriter *ResponseWrapper) GetBuffer() *bytes.Buffer {
	return &wrappedWriter.buffer
}

// GetContent stuff.
func (wrappedWriter *ResponseWrapper) GetContent(encoding string) ([]byte, error) {
	return compressutil.Decode(wrappedWriter.GetBuffer(), encoding)
}

// SetContent stuff.
func (wrappedWriter *ResponseWrapper) SetContent(data []byte, encoding string) {
	bodyBytes, _ := compressutil.Encode(data, encoding)

	if !wrappedWriter.wroteHeader {
		wrappedWriter.WriteHeader(http.StatusOK)
	}

	if _, err := wrappedWriter.ResponseWriter.Write(bodyBytes); err != nil {
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
func (wrappedWriter *ResponseWrapper) GetContentEncoding() string {
	return wrappedWriter.Header().Get("Content-Encoding")
}

// GetContentType get the Content-Encoding header value.
func (wrappedWriter *ResponseWrapper) GetContentType() string {
	return wrappedWriter.Header().Get("Content-Type")
}

// SupportsProcessing determine if HttpWrapper is supported by this plugin based on encoding.
func (wrappedWriter *ResponseWrapper) SupportsProcessing() bool {
	contentType := wrappedWriter.GetContentType()

	// If content type does not match return values with false
	if contentType != "" && !strings.Contains(contentType, "text") {
		return false
	}

	encoding := wrappedWriter.GetContentEncoding()

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

// SetLastModified stuff.
func (wrappedWriter *ResponseWrapper) SetLastModified(value bool) {
	wrappedWriter.lastModified = value
}
