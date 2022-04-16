// Package plugin_rewritebody a plugin to rewrite response body.
package plugin_rewritebody

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
)

// Rewrite holds one rewrite body configuration.
type Rewrite struct {
	Regex       string `json:"regex,omitempty"`
	Replacement string `json:"replacement,omitempty"`
}

// Config holds the plugin configuration.
type Config struct {
	LastModified bool      `json:"lastModified,omitempty"`
	Rewrites     []Rewrite `json:"rewrites,omitempty"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type rewrite struct {
	regex       *regexp.Regexp
	replacement []byte
}

type rewriteBody struct {
	name         string
	next         http.Handler
	rewrites     []rewrite
	lastModified bool
}

// New creates and returns a new rewrite body plugin instance.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	rewrites := make([]rewrite, len(config.Rewrites))

	for index, rewriteConfig := range config.Rewrites {
		regex, err := regexp.Compile(rewriteConfig.Regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex %q: %w", rewriteConfig.Regex, err)
		}

		rewrites[index] = rewrite{
			regex:       regex,
			replacement: []byte(rewriteConfig.Replacement),
		}
	}

	return &rewriteBody{
		name:         name,
		next:         next,
		rewrites:     rewrites,
		lastModified: config.LastModified,
	}, nil
}

func (bodyRewrite *rewriteBody) ServeHTTP(response http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Upgrade") == "websocket" {
		log.Printf("Ignoring websocket upgrade| Host: \"%s\" | Path: \"%s\"", req.Host, req.URL.Path)

		return
	}

	wrappedWriter := &responseWriter{
		lastModified:   bodyRewrite.lastModified,
		ResponseWriter: response,
	}

	bodyRewrite.next.ServeHTTP(wrappedWriter, req)

	encoding, _, isSupported := wrappedWriter.getHeaderContent()

	if !isSupported {
		if _, err := response.Write(wrappedWriter.buffer.Bytes()); err != nil {
			log.Printf("unable to write body: %v", err)
		}

		return
	}

	bodyBytes, ok := wrappedWriter.decompressBody(encoding)
	if ok {
		for _, rwt := range bodyRewrite.rewrites {
			bodyBytes = rwt.regex.ReplaceAll(bodyBytes, rwt.replacement)
		}

		bodyBytes = prepareBodyBytes(bodyBytes, encoding)
	} else {
		bodyBytes = wrappedWriter.buffer.Bytes()
	}

	if _, err := response.Write(bodyBytes); err != nil {
		log.Printf("unable to write rewrited body: %v", err)
	}
}

func (wrappedWriter *responseWriter) getHeaderContent() (encoding string, contentType string, isSupported bool) {
	encoding = wrappedWriter.Header().Get("Content-Encoding")
	contentType = wrappedWriter.Header().Get("Content-Type")

	// If content type does not match return values with false
	if contentType != "" && !strings.Contains(contentType, "text") {
		return encoding, contentType, false
	}

	// If content type is supported validate encoding as well
	switch encoding {
	case "gzip":
		fallthrough
	case "deflate":
		fallthrough
	case "identity":
		fallthrough
	case "":
		return encoding, contentType, true
	default:
		return encoding, contentType, false
	}
}

func (wrappedWriter *responseWriter) decompressBody(encoding string) ([]byte, bool) {
	switch encoding {
	case "gzip":
		return getBytesFromGzip(wrappedWriter.buffer)

	case "deflate":
		return getBytesFromZlib(wrappedWriter.buffer)

	default:
		return wrappedWriter.buffer.Bytes(), true
	}
}

func getBytesFromZlib(buffer bytes.Buffer) ([]byte, bool) {
	zlibReader, err := zlib.NewReader(&buffer)
	if err != nil {
		log.Printf("Failed to load body reader: %v", err)

		return buffer.Bytes(), false
	}

	bodyBytes, err := io.ReadAll(zlibReader)
	if err != nil {
		log.Printf("Failed to read body: %s", err)

		return buffer.Bytes(), false
	}

	err = zlibReader.Close()

	if err != nil {
		log.Printf("Failed to close reader: %v", err)

		return buffer.Bytes(), false
	}

	return bodyBytes, true
}

func getBytesFromGzip(buffer bytes.Buffer) ([]byte, bool) {
	gzipReader, err := gzip.NewReader(&buffer)
	if err != nil {
		log.Printf("Failed to load body reader: %v", err)

		return buffer.Bytes(), false
	}

	bodyBytes, err := io.ReadAll(gzipReader)
	if err != nil {
		log.Printf("Failed to read body: %s", err)

		return buffer.Bytes(), false
	}

	err = gzipReader.Close()

	if err != nil {
		log.Printf("Failed to close reader: %v", err)

		return buffer.Bytes(), false
	}

	return bodyBytes, true
}

func prepareBodyBytes(bodyBytes []byte, encoding string) []byte {
	switch encoding {
	case "gzip":
		return compressWithGzip(bodyBytes)

	case "deflate":
		return compressWithZlib(bodyBytes)

	default:
		return bodyBytes
	}
}

func compressWithGzip(bodyBytes []byte) []byte {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	if _, err := gzipWriter.Write(bodyBytes); err != nil {
		log.Printf("unable to recompress rewrited body: %v", err)

		return bodyBytes
	}

	if err := gzipWriter.Close(); err != nil {
		log.Printf("unable to close gzip writer: %v", err)

		return bodyBytes
	}

	return buf.Bytes()
}

func compressWithZlib(bodyBytes []byte) []byte {
	var buf bytes.Buffer
	zlibWriter := zlib.NewWriter(&buf)

	if _, err := zlibWriter.Write(bodyBytes); err != nil {
		log.Printf("unable to recompress rewrited body: %v", err)

		return bodyBytes
	}

	if err := zlibWriter.Close(); err != nil {
		log.Printf("unable to close zlib writer: %v", err)

		return bodyBytes
	}

	return buf.Bytes()
}

type responseWriter struct {
	buffer       bytes.Buffer
	lastModified bool
	wroteHeader  bool

	http.ResponseWriter
}

func (wrappedWriter *responseWriter) WriteHeader(statusCode int) {
	if !wrappedWriter.lastModified {
		wrappedWriter.ResponseWriter.Header().Del("Last-Modified")
	}

	wrappedWriter.wroteHeader = true

	// Delegates the Content-Length Header creation to the final body write.
	wrappedWriter.ResponseWriter.Header().Del("Content-Length")

	wrappedWriter.ResponseWriter.WriteHeader(statusCode)
}

func (wrappedWriter *responseWriter) Write(p []byte) (int, error) {
	if !wrappedWriter.wroteHeader {
		wrappedWriter.WriteHeader(http.StatusOK)
	}

	return wrappedWriter.buffer.Write(p)
}

func (wrappedWriter *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := wrappedWriter.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("%T is not a http.Hijacker", wrappedWriter.ResponseWriter)
	}

	return hijacker.Hijack()
}

func (wrappedWriter *responseWriter) Flush() {
	if flusher, ok := wrappedWriter.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
