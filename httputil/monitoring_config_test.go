package httputil_test

import (
	"net/http"
	"testing"

	"github.com/packruler/rewrite-body/httputil"
)

func TestMonitoringConfigParsing(t *testing.T) {
	tests := []struct {
		desc            string
		inputTypes      []string
		inputMethods    []string
		expectedTypes   []string
		expectedMethods []string
	}{
		{
			desc:            "defaults will be supplied for empty arrays",
			inputTypes:      []string{},
			inputMethods:    []string{},
			expectedTypes:   []string{"text/html"},
			expectedMethods: []string{http.MethodGet},
		},
		// {
		// 	desc:            "defaults will be supplied for empty types with populated methods",
		// 	inputTypes:      "",
		// 	inputMethods:    "POST",
		// 	expectedTypes:   []string{"text/html"},
		// 	expectedMethods: []string{http.MethodPost},
		// },
		// {
		// 	desc:            "defaults will be supplied for empty methods with populated types",
		// 	inputTypes:      "application/json",
		// 	inputMethods:    "",
		// 	expectedTypes:   []string{"application/json"},
		// 	expectedMethods: []string{http.MethodGet},
		// },
		// {
		// 	desc:            "proper cases will be used",
		// 	inputTypes:      "TEXT/HTml",
		// 	inputMethods:    "gEt",
		// 	expectedTypes:   []string{"text/html"},
		// 	expectedMethods: []string{http.MethodGet},
		// },
		// {
		// 	desc:            "works with multiple values",
		// 	inputTypes:      "text/html, application/json",
		// 	inputMethods:    "get, post",
		// 	expectedTypes:   []string{"text/html", "application/json"},
		// 	expectedMethods: []string{http.MethodGet, http.MethodPost},
		// },
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			config := httputil.MonitoringConfig{
				Types:   test.inputTypes,
				Methods: test.inputMethods,
			}

			config.EnsureDefaults()

			if len(config.Types) != len(test.expectedTypes) {
				t.Errorf("Expected Types: '%v' | Got Types: '%v'", test.expectedTypes, config.Types)
			}

			for i, v := range test.expectedTypes {
				if v != config.Types[i] {
					t.Errorf("Expected Types: '%v' | Got Types: '%v'", test.expectedTypes, config.Types)
				}
			}

			if len(config.Methods) != len(test.expectedMethods) {
				t.Errorf("Expected Methods: '%v' | Got Methods: '%v'", test.expectedMethods, config.Methods)
			}

			for i, v := range test.expectedMethods {
				if v != config.Methods[i] {
					t.Errorf("Expected Methods: '%v' | Got Methods: '%v'", test.expectedMethods, config.Methods)
				}
			}
		})
	}
}
