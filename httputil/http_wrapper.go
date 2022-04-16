// Package httputil a package for handling http data tasks
package httputil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/packruler/rewrite-body/compressutil"
)

// HTTPWrapper a struct to be used for handling data manipulation by this plugin.
type HTTPWrapper struct {
	Request  *http.Request
	Response *http.Response

	http.ResponseWriter

	lastModified bool
	wroteHeader  bool
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
	return wrapper.ResponseWriter.Header().Get("Content-Encoding")
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
		wrapper.ResponseWriter.Header().Del("Last-Modified")
	}

	wrapper.wroteHeader = true

	// Delegates the Content-Length Header creation to the final body write.
	wrapper.ResponseWriter.Header().Del("Content-Length")

	wrapper.ResponseWriter.WriteHeader(statusCode)
}

// DecompressError an error that occurred in decompression process.
type DecompressError struct {
	error
}

// GetContent load []byte for uncompressed data in response.
// Inspiration from https://github.com/andybalholm/redwood/blob/master/proxy.go.
func (wrapper *HTTPWrapper) GetContent(maxLength int) ([]byte, error) {
	if wrapper.Response.ContentLength > int64(maxLength) {
		return nil, fmt.Errorf("content too large: %d", wrapper.Request.ContentLength)
	}

	limitedReader := &io.LimitedReader{
		R: wrapper.Response.Body,
		N: int64(maxLength),
	}
	content, err := ioutil.ReadAll(limitedReader)

	// Servers that use broken chunked Transfer-Encoding can give us unexpected EOFs,
	// even if we got all the content.
	if errors.Is(err, io.ErrUnexpectedEOF) && wrapper.Response.ContentLength == -1 {
		err = nil
	}

	if err != nil {
		return nil, err
	}

	if limitedReader.N == 0 {
		// We read maxLen without reaching the end.
		wrapper.Response.Body = io.NopCloser(io.MultiReader(bytes.NewReader(content), wrapper.Response.Body))

		return nil, nil
	}

	if wrapper.GetContentEncoding() == "" {
		wrapper.Response.ContentLength = int64(len(content))
	}

	wrapper.Response.Body = io.NopCloser(bytes.NewReader(content))

	if contentEncoding := wrapper.GetContentEncoding(); contentEncoding != "" && len(content) > 0 {
		br := bytes.NewReader(content)

		decompressed, err := compressutil.Decode(br, wrapper.GetContentEncoding())
		if err != nil {
			return content, &DecompressError{}
		}

		return decompressed, nil
	}

	return content, nil
}

// SetContent set the content of a response based to be the data supplied encoded with the encoding supplied.
func (wrapper *HTTPWrapper) SetContent(data []byte, encoding string) {
	wrapper.Response.Header.Set("Content-Encoding", encoding)

	wrapper.Response.ContentLength = int64(len(data))
	wrapper.Response.Body = io.NopCloser(bytes.NewReader(data))

	if encoding != "" && encoding != "identity" {
		readCloser, err := compressutil.Encode(data, encoding)
		if err != nil {
			log.Printf("Unable to encode data: %v", err)

			return
		}

		if readCloser != nil {
			wrapper.Response.Body = readCloser
			wrapper.Response.Header.Set("Content-Encoding", encoding)
			wrapper.Response.ContentLength = -1

			return
		}
	}

	wrapper.Response.ContentLength = int64(len(data))
	wrapper.Response.Body = io.NopCloser(bytes.NewReader(data))
}