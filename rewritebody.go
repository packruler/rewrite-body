// Package rewrite_body a plugin to rewrite response body.
package rewrite_body

import (
	"context"
	"net/http"

	"github.com/packruler/rewrite-body/handler"
)

// New creates and returns a new rewrite body plugin instance.
func New(context context.Context, next http.Handler, config *handler.Config, name string) (http.Handler, error) {
	return handler.New(context, next, config, name)
}
