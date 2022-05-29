// Package rewrite_body a plugin to rewrite response body.
package rewrite_body

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/packruler/rewrite-body/httputil"
	"github.com/packruler/rewrite-body/logger"
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
	LogLevel     int8      `json:"logLevel,omitempty"`
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
	logger       logger.LogWriter
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

	switch config.LogLevel {
	// Convert default 0 to Info level
	case 0:
		config.LogLevel = int8(logger.Info)
	// Allow -1 to be call for Trace level
	case -1:
		config.LogLevel = int8(logger.Trace)
	default:
	}

	logWriter := *logger.CreateLogger(logger.LogLevel(config.LogLevel))

	return &rewriteBody{
		name:         name,
		next:         next,
		rewrites:     rewrites,
		lastModified: config.LastModified,
		logger:       logWriter,
	}, nil
}

func (bodyRewrite *rewriteBody) ServeHTTP(response http.ResponseWriter, req *http.Request) {
	defer handlePanic()

	// allow default http.ResponseWriter to handle calls targeting WebSocket upgrades and non GET methods
	if !httputil.SupportsProcessing(req) {
		bodyRewrite.logger.LogDebugf("Ignoring unsupported request: %v", req)
		bodyRewrite.next.ServeHTTP(response, req)

		return
	}

	bodyRewrite.logger.LogDebugf("Starting supported request: %v", req)

	wrappedWriter := &httputil.ResponseWrapper{
		ResponseWriter: response,
	}

	wrappedWriter.SetLastModified(bodyRewrite.lastModified)

	// look into using https://pkg.go.dev/net/http#RoundTripper
	bodyRewrite.next.ServeHTTP(wrappedWriter, req)

	if !wrappedWriter.SupportsProcessing() {
		// We are ignoring these any errors because the content should be unchanged here.
		// This could "error" if writing is not supported but content will return properly.
		_, _ = response.Write(wrappedWriter.GetBuffer().Bytes())
		bodyRewrite.logger.LogDebugf("Ignoring unsupported response: %v", wrappedWriter)

		return
	}

	bodyBytes, err := wrappedWriter.GetContent()
	if err != nil {
		bodyRewrite.logger.LogErrorf("Error loading content: %v", err)

		if _, err := response.Write(wrappedWriter.GetBuffer().Bytes()); err != nil {
			bodyRewrite.logger.LogErrorf("unable to write error content: %v", err)
		}

		return
	}

	bodyRewrite.logger.LogDebugf("Response body: %s", bodyBytes)

	if len(bodyBytes) == 0 {
		// If the body is empty there is no purpose in continuing this process.
		return
	}

	for _, rwt := range bodyRewrite.rewrites {
		bodyBytes = rwt.regex.ReplaceAll(bodyBytes, rwt.replacement)
	}

	bodyRewrite.logger.LogDebugf("Transformed body: %s", bodyBytes)
	wrappedWriter.SetContent(bodyBytes)
}

func handlePanic() {
	if recovery := recover(); recovery != nil {
		if err, ok := recovery.(error); ok {
			logError(err)
		} else {
			log.Printf("Unhandled error: %v", recovery)
		}
	}
}

func logError(err error) {
	// Ignore http.ErrAbortHandler because they are expected errors that do not require handling
	if errors.Is(err, http.ErrAbortHandler) {
		return
	}

	log.Printf("Recovered from: %v", err)
}
