package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/packruler/rewrite-body/compress"
)

type HttpWrapper struct {
	Request  *http.Request
	Response *http.Response

	http.ResponseWriter

	lastModified bool
	wroteHeader  bool
}

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

func (wrapper *HttpWrapper) GetContentEncoding() string {
	return wrapper.ResponseWriter.Header().Get("Content-Encoding")
}

func (wrapper *HttpWrapper) SupportsProcessing() bool {
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

func (wrapper *HttpWrapper) WriteHeader(statusCode int) {
	if !wrapper.lastModified {
		wrapper.ResponseWriter.Header().Del("Last-Modified")
	}

	wrapper.wroteHeader = true

	// Delegates the Content-Length Header creation to the final body write.
	wrapper.ResponseWriter.Header().Del("Content-Length")

	wrapper.ResponseWriter.WriteHeader(statusCode)
}

type DecompressError struct {
	error
}

// Inspiration from https://github.com/andybalholm/redwood/blob/master/proxy.go.
// GetContent load []byte for uncompressed data in response.
func (wrapper *HttpWrapper) GetContent(maxLength int) ([]byte, error) {
	if wrapper.Response.ContentLength > int64(maxLength) {
		return nil, fmt.Errorf("content too large: %d", wrapper.Request.ContentLength)
	}

	lr := &io.LimitedReader{
		R: wrapper.Response.Body,
		N: int64(maxLength),
	}
	content, err := ioutil.ReadAll(lr)

	// Servers that use broken chunked Transfer-Encoding can give us unexpected EOFs,
	// even if we got all the content.
	if err == io.ErrUnexpectedEOF && wrapper.Response.ContentLength == -1 {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	if lr.N == 0 {
		// We read maxLen without reaching the end.
		wrapper.Response.Body = io.NopCloser(io.MultiReader(bytes.NewReader(content), wrapper.Response.Body))
		return nil, nil
	}

	if wrapper.GetContentEncoding() == "" {
		wrapper.Response.ContentLength = int64(len(content))
	}
	wrapper.Response.Body = io.NopCloser(bytes.NewReader(content))

	if ce := wrapper.GetContentEncoding(); ce != "" && len(content) > 0 {
		br := bytes.NewReader(content)
		decompressed, err := compress.Decode(br, wrapper.GetContentEncoding())
		if err != nil {
			return content, &DecompressError{}
		}

		return decompressed, nil
	}

	return content, nil
}
