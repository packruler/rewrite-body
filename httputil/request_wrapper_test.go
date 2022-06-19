package httputil

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/packruler/rewrite-body/logger"
)

func TestGetEncodingTarget(t *testing.T) {
	tests := []struct {
		desc           string
		acceptEncoding string
		expectedTarget string
	}{
		{
			desc:           "Supports gzip",
			acceptEncoding: "gzip",
			expectedTarget: "gzip",
		},
		{
			desc:           "Supports deflate",
			acceptEncoding: "deflate",
			expectedTarget: "deflate",
		},
		{
			desc:           "Supports identity",
			acceptEncoding: "identity",
			expectedTarget: "identity",
		},
		{
			desc:           "Ignores brotli",
			acceptEncoding: "br, gzip",
			expectedTarget: "gzip",
		},
		{
			desc:           "Wildcard to gzip",
			acceptEncoding: "*",
			expectedTarget: "gzip",
		},
		{
			desc:           "Respects quality in order",
			acceptEncoding: "gzip;q=0.8, deflate;q=0.6",
			expectedTarget: "gzip",
		},
		{
			desc:           "Respects quality out of order",
			acceptEncoding: "gzip;q=0.8, deflate;q=0.9",
			expectedTarget: "deflate",
		},
	}

	defaultMonitoring := MonitoringConfig{
		Types:   []string{"text/html"},
		Methods: []string{"GET"},
	}

	defaultLogWriter := logger.CreateLogger(logger.Error)

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			request, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				"http://google.com",
				&bytes.Reader{})
			if err != nil {
				t.Errorf("Error creating request: %v", err)
			}
			request.Header.Set("Accept-Encoding", test.acceptEncoding)

			wrappedRequest := WrapRequest(request, defaultMonitoring, *defaultLogWriter)
			target := wrappedRequest.GetEncodingTarget()
			if target != test.expectedTarget {
				t.Errorf("Expected: '%s' | Got: '%s'", test.expectedTarget, target)
			}
		})
	}
}

func TestRemoveUnuspportedEncoding(t *testing.T) {
	tests := []struct {
		desc           string
		acceptEncoding string
		expectedTarget string
	}{
		{
			desc:           "Supports gzip",
			acceptEncoding: "gzip",
			expectedTarget: "gzip",
		},
		{
			desc:           "Supports deflate",
			acceptEncoding: "deflate",
			expectedTarget: "deflate",
		},
		{
			desc:           "Supports identity",
			acceptEncoding: "identity",
			expectedTarget: "identity",
		},
		{
			desc:           "Ignores brotli",
			acceptEncoding: "br, gzip",
			expectedTarget: " gzip",
		},
		{
			desc:           "Wildcard is dropped",
			acceptEncoding: "*",
			expectedTarget: "",
		},
		{
			desc:           "Respects quality in order",
			acceptEncoding: "gzip;q=0.8, deflate;q=0.6",
			expectedTarget: "gzip;q=0.8, deflate;q=0.6",
		},
	}

	defaultMonitoring := MonitoringConfig{
		Types:   []string{"text/html"},
		Methods: []string{"GET"},
	}

	defaultLogWriter := logger.CreateLogger(logger.Error)

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			request, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				"http://google.com",
				&bytes.Reader{})
			if err != nil {
				t.Errorf("Error creating request: %v", err)
			}
			request.Header.Set("Accept-Encoding", test.acceptEncoding)

			wrappedRequest := WrapRequest(request, defaultMonitoring, *defaultLogWriter)
			target := wrappedRequest.CloneWithSupportedEncoding().Header.Get("Accept-Encoding")

			if target != test.expectedTarget {
				t.Errorf("Expected: '%s' | Got: '%s'", test.expectedTarget, target)
			}
		})
	}
}

func TestSupportsProcessing(t *testing.T) {
	tests := []struct {
		desc             string
		inputType        string
		inputMethod      string
		monitoringConfig MonitoringConfig
		expectedSupport  bool
	}{
		{
			desc:            "Supports default config",
			inputType:       "text/html",
			inputMethod:     http.MethodGet,
			expectedSupport: true,
			monitoringConfig: MonitoringConfig{
				Types:   []string{"text/html"},
				Methods: []string{"GET"},
			},
		},
		{
			desc:            "Supports default browser load",
			inputType:       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			inputMethod:     http.MethodGet,
			expectedSupport: true,
			monitoringConfig: MonitoringConfig{
				Types:   []string{"text/html"},
				Methods: []string{"GET"},
			},
		},
		{
			desc:            "Supports when types includes unsupported type first",
			inputType:       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			inputMethod:     http.MethodGet,
			expectedSupport: true,
			monitoringConfig: MonitoringConfig{
				Types:   []string{"application/javascript", "text/html"},
				Methods: []string{"GET"},
			},
		},
		{
			desc:            "Supports multiple methods",
			inputType:       "text/html",
			inputMethod:     http.MethodPost,
			expectedSupport: true,
			monitoringConfig: MonitoringConfig{
				Types:   []string{"text/html"},
				Methods: []string{"GET", "POST"},
			},
		},
		{
			desc:            "Does not support type not included",
			inputType:       "application/javascript",
			inputMethod:     http.MethodGet,
			expectedSupport: false,
			monitoringConfig: MonitoringConfig{
				Types:   []string{"text/html"},
				Methods: []string{"GET"},
			},
		},
		{
			desc:            "Does not support method not included",
			inputType:       "text/html",
			inputMethod:     http.MethodPost,
			expectedSupport: false,
			monitoringConfig: MonitoringConfig{
				Types:   []string{"text/html"},
				Methods: []string{"GET"},
			},
		},
	}

	defaultLogWriter := logger.CreateLogger(logger.Error)

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			request, err := http.NewRequestWithContext(
				context.Background(),
				test.inputMethod,
				"http://google.com",
				&bytes.Reader{})
			if err != nil {
				t.Errorf("Error creating request: %v", err)
			}
			request.Header.Set("Accept", test.inputType)

			wrappedRequest := WrapRequest(request, test.monitoringConfig, *defaultLogWriter)

			if test.expectedSupport != wrappedRequest.SupportsProcessing() {
				t.Errorf("Test input: '%v'", test)
			}
		})
	}
}
