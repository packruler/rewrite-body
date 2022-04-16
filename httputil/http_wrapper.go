// Package httputil a package for handling http data tasks
package httputil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/packruler/rewrite-body/compressutil"
)

// HTTPWrapper a struct to be used for handling data manipulation by this plugin.
type HTTPWrapper struct {
	Request  *http.Request
	Response *http.Response

	buffer bytes.Buffer

	lastModified bool
	wroteHeader  bool

	http.ResponseWriter
}

// SetLastModified set value in struct from external package.
func (wrapper *HTTPWrapper) SetLastModified(value bool) {
	wrapper.lastModified = value
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
func (wrapper *HTTPWrapper) GetContentEncoding() string {
	return wrapper.Header().Get("Content-Encoding")
}

// SupportsProcessing determine if HttpWrapper is supported by this plugin based on encoding.
func (wrapper *HTTPWrapper) SupportsProcessing() bool {
	switch wrapper.GetContentEncoding() {
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

// WriteHeader write the status code and other content for the wrapper.
func (wrapper *HTTPWrapper) WriteHeader(statusCode int) {
	if !wrapper.lastModified {
		wrapper.Header().Del("Last-Modified")
	}

	wrapper.wroteHeader = true

	// wrapper.Header().WriteHeader(statusCode)

	// Delegates the Content-Length Header creation to the final body write.
	wrapper.Header().Del("Content-Length")
}

// DecompressError an error that occurred in decompression process.
type DecompressError struct {
	error
}

// GetContent load []byte for uncompressed data in response.
// Inspiration from https://github.com/andybalholm/redwood/blob/master/proxy.go.
func (wrapper *HTTPWrapper) GetContent(maxLength int) ([]byte, error) {
	log.Printf("Response: %+v", wrapper.Response)

	if wrapper.buffer.Len() > maxLength {
		return nil, fmt.Errorf("content too large: %d", wrapper.Request.ContentLength)
	}

	// limitedReader := &io.LimitedReader{
	// 	R: &wrapper.buffer,
	// 	N: int64(maxLength),
	// }

	content := wrapper.buffer.Bytes()
	log.Printf("Content: %s", content)

	// Servers that use broken chunked Transfer-Encoding can give us unexpected EOFs,
	// even if we got all the content.
	// if errors.Is(err, io.ErrUnexpectedEOF) && wrapper.Response.ContentLength == -1 {
	// 	err = nil
	// }

	// log.Println("1")

	// if err != nil {
	// 	return nil, err
	// }

	// log.Println("2")

	// if limitedReader.N == 0 {
	// 	// We read maxLen without reaching the end.
	// 	return nil, nil
	// }

	// if wrapper.GetContentEncoding() == "" {
	// 	wrapper.Response.ContentLength = int64(len(content))
	// }

	// wrapper.Response.Body = io.NopCloser(bytes.NewReader(content))

	log.Println("3")

	if contentEncoding := wrapper.GetContentEncoding(); contentEncoding != "" && len(content) > 0 {
		br := bytes.NewReader(content)

		decompressed, err := compressutil.Decode(br, wrapper.GetContentEncoding())
		if err != nil {
			return content, &DecompressError{}
		}

		return decompressed, nil
	}

	log.Println("4")

	return content, nil
}

// SetContent set the content of a response based to be the data supplied encoded with the encoding supplied.
func (wrapper *HTTPWrapper) SetContent(data []byte, encoding string) {
	wrapper.Header().Set("Content-Encoding", encoding)

	if encoding != "" && encoding != "identity" {
		readCloser, err := compressutil.Encode(data, encoding)
		if err != nil {
			log.Printf("Unable to encode data: %v", err)

			return
		}

		if readCloser != nil {
			encodedData, err := io.ReadAll(readCloser)
			if err != nil {
				log.Printf("Error encoding data: %v", err)
			}

			_, _ = wrapper.Write(encodedData)
			wrapper.Header().Set("Content-Encoding", encoding)
		}
	}
}

// Hijack run http.Hijacker.Hijack().
func (wrapper *HTTPWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := wrapper.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("%T is not a http.Hijacker", wrapper.ResponseWriter)
	}

	return hijacker.Hijack()
}

// BodyBytes test.
func (wrapper *HTTPWrapper) BodyBytes() []byte {
	return wrapper.buffer.Bytes()
}
