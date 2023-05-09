package rewrite_body

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/joinrepublic/traefik-rewrite-body-csp/compressutil"
	"github.com/joinrepublic/traefik-rewrite-body-csp/handler"
)

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		desc            string
		contentEncoding string
		contentType     string `default:"text/html"`
                placeholder     string
                nonceGenerator  func(string) []byte
		lastModified    bool
		resBody         string
		expResBody      string
		expLastModified bool
                headers         map[string]string
                expHeaders      map[string]string
	}{
		{
			desc: "should replace foo by bar",
                        placeholder: "foo",
                        nonceGenerator: func(_ string) []byte {
                          return []byte("bar")
                        },
			contentType: "text/html",
			resBody:     "foo is the new bar",
			expResBody:  "bar is the new bar",
                        headers: map[string]string{
                          "content-security-policy": "script-src 'nonce-foo'",
                          "content-security-policy-report-only": "script-src 'strict-dynamic' 'nonce-foo'",
                          "baz": "script-src 'nonce-foo'",
                        },
                        expHeaders: map[string]string{
                          "content-security-policy": "script-src 'nonce-bar'",
                          "content-security-policy-report-only": "script-src 'strict-dynamic' 'nonce-bar'",
                          "baz": "script-src 'nonce-foo'", // control group
                        },
		},
		{
			desc: "should not replace anything if content encoding is not identity or empty",
                        placeholder: "foo",
                        nonceGenerator: func(_ string) []byte {
                          return []byte("bar")
                        },
			contentEncoding: "other",
			contentType:     "text/html",
			resBody:         "foo is the new bar",
			expResBody:      "foo is the new bar",
                        headers: map[string]string{
                          "content-security-policy": "script-src 'nonce-foo'",
                        },
		},
		{
			desc: "should not replace anything if content type does not contain text or is not empty",
                        placeholder: "foo",
                        nonceGenerator: func(_ string) []byte {
                          return []byte("bar")
                        },
			contentType: "image",
			resBody:     "foo is the new bar",
			expResBody:  "foo is the new bar",
                        headers: map[string]string{
                          "content-security-policy": "script-src 'nonce-foo'",
                        },
		},
		{
			desc: "should replace foo by bar if content encoding is identity",
                        placeholder: "foo",
                        nonceGenerator: func(_ string) []byte {
                          return []byte("bar")
                        },
			contentEncoding: "identity",
			contentType:     "text/html",
			resBody:         "foo is the new bar",
			expResBody:      "bar is the new bar",
                        headers: map[string]string{
                          "content-security-policy": "script-src 'nonce-foo'",
                        },
		},
		{
			desc: "should not remove the last modified header",
                        placeholder: "foo",
                        nonceGenerator: func(_ string) []byte {
                          return []byte("bar")
                        },
			contentEncoding: "identity",
			contentType:     "text/html",
			lastModified:    true,
			resBody:         "foo is the new bar",
			expResBody:      "bar is the new bar",
			expLastModified: true,
                        headers: map[string]string{
                          "content-security-policy": "script-src 'nonce-foo'",
                        },
		},
		{
			desc: "should support gzip encoding",
                        placeholder: "foo",
                        nonceGenerator: func(_ string) []byte {
                          return []byte("bar")
                        },
			contentEncoding: "gzip",
			contentType:     "text/html",
			lastModified:    true,
			resBody:         compressString("foo is the new bar", "gzip"),
			expResBody:      compressString("bar is the new bar", "gzip"),
			expLastModified: true,
                        headers: map[string]string{
                          "content-security-policy": "script-src 'nonce-foo'",
                        },
		},
		{
			desc: "should support deflate encoding",
                        placeholder: "foo",
                        nonceGenerator: func(_ string) []byte {
                          return []byte("bar")
                        },
			contentEncoding: "deflate",
			contentType:     "text/html",
			lastModified:    true,
			resBody:         compressString("foo is the new bar", "deflate"),
			expResBody:      compressString("bar is the new bar", "deflate"),
			expLastModified: true,
                        headers: map[string]string{
                          "content-security-policy": "script-src 'nonce-foo'",
                        },
		},
		{
			desc: "should ignore unsupported encoding",
                        placeholder: "foo",
                        nonceGenerator: func(_ string) []byte {
                          return []byte("bar")
                        },
			contentEncoding: "br",
			contentType:     "text/html",
			lastModified:    true,
			resBody:         "foo is the new bar",
			expResBody:      "foo is the new bar",
			expLastModified: true,
                        headers: map[string]string{
                          "content-security-policy": "script-src 'nonce-foo'",
                        },
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			config := &handler.Config{
				LastModified: test.lastModified,
                                Placeholder:  test.placeholder,
                                NonceGenerator: test.nonceGenerator,
				LogLevel:     -1,
			}

			next := func(responseWriter http.ResponseWriter, req *http.Request) {
				responseWriter.Header().Set("Content-Encoding", test.contentEncoding)
				responseWriter.Header().Set("Content-Type", test.contentType)
				responseWriter.Header().Set("Last-Modified", "Thu, 02 Jun 2016 06:01:08 GMT")
				responseWriter.Header().Set("Content-Length", strconv.Itoa(len(test.resBody)))
                                for headerName, value := range test.headers {
                                  responseWriter.Header().Set(headerName, value)
                                }

				responseWriter.WriteHeader(http.StatusOK)

				_, _ = fmt.Fprintf(responseWriter, test.resBody)
			}

			rewriteBody, err := New(context.Background(), http.HandlerFunc(next), config, "rewriteBody")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept", "text/html")

			rewriteBody.ServeHTTP(recorder, req)

			if _, exists := recorder.Result().Header["Last-Modified"]; exists != test.expLastModified {
				t.Errorf("got last-modified header %v, want %v", exists, test.expLastModified)
			}

			if _, exists := recorder.Result().Header["Content-Length"]; exists {
				t.Error("The Content-Length Header must be deleted")
			}

                        for headerName, expectedValue := range test.expHeaders {
                          actualValue := recorder.Result().Header.Get(headerName)
                          if actualValue != expectedValue {
                                  t.Errorf("expected %v header to be %v, got %v", headerName, expectedValue, actualValue)
                          }
                        }

			if !bytes.Equal([]byte(test.expResBody), recorder.Body.Bytes()) {
				t.Errorf("got body: %v\n wanted: %v", recorder.Body.Bytes(), []byte(test.expResBody))
			}
		})
	}
}

func compressString(value string, encoding string) string {
	compressed, _ := compressutil.Encode([]byte(value), encoding)

	return string(compressed)
}

