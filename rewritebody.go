// Package rewrite_body a plugin to rewrite response body.
package rewrite_body

import (
	"context"
	"net/http"

	"github.com/packruler/rewrite-body/handler"
)

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *handler.Config {
	return &handler.Config{}
}

// New creates and returns a new rewrite body plugin instance.
func New(context context.Context, next http.Handler, config *handler.Config, name string) (http.Handler, error) {
	return handler.New(context, next, config, name)
}
