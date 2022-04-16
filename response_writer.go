package rewrite_body

import (
	"bytes"
	"log"
	"net/http"
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
func (wrappedWriter *ResponseWrapper) GetContent(encoding string) ([]byte, bool) {
	return wrappedWriter.decompressBody(encoding)
}

// SetContent stuff.
func (wrappedWriter *ResponseWrapper) SetContent(data []byte, encoding string) {
	bodyBytes, _ := prepareBodyBytes(data, encoding)

	if !wrappedWriter.wroteHeader {
		wrappedWriter.WriteHeader(http.StatusOK)
	}

	if _, err := wrappedWriter.ResponseWriter.Write(bodyBytes); err != nil {
		log.Printf("unable to write rewrited body: %v", err)
	}
}
