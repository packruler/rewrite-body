// Package rewrite_body a plugin to rewrite response body.
package rewrite_body

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/packruler/rewrite-body/httputil"
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
	if !httputil.SupportsProcessing(req) {
		return
	}

	wrappedWriter := &httputil.ResponseWrapper{
		ResponseWriter: response,
	}

	wrappedWriter.SetLastModified(bodyRewrite.lastModified)

	bodyRewrite.next.ServeHTTP(wrappedWriter, req)

	isSupported := wrappedWriter.SupportsProcessing()
	if !isSupported {
		if _, err := response.Write(wrappedWriter.GetBuffer().Bytes()); err != nil {
			log.Printf("unable to write body: %v", err)
		}

		return
	}

	bodyBytes, err := wrappedWriter.GetContent()
	if err == nil {
		for _, rwt := range bodyRewrite.rewrites {
			bodyBytes = rwt.regex.ReplaceAll(bodyBytes, rwt.replacement)
		}
	} else {
		bodyBytes = wrappedWriter.GetBuffer().Bytes()
	}

	wrappedWriter.SetContent(bodyBytes)
}
