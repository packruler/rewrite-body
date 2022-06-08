package httputil_test

// import (
// 	"net/http"
// 	"testing"

// 	"github.com/packruler/plugin-utils/httputil"
// )

// func TestMonitoringConfigParsing(t *testing.T) {
// 	tests := []struct {
// 		desc            string
// 		inputTypes      string
// 		inputMethods    string
// 		expectedTypes   []string
// 		expectedMethods []string
// 	}{
// 		{
// 			desc:            "defaults will be supplied for empty strings",
// 			inputTypes:      "",
// 			inputMethods:    "",
// 			expectedTypes:   []string{"text/html"},
// 			expectedMethods: []string{http.MethodGet},
// 		},
// 		{
// 			desc:            "defaults will be supplied for empty types with populated methods",
// 			inputTypes:      "",
// 			inputMethods:    "POST",
// 			expectedTypes:   []string{"text/html"},
// 			expectedMethods: []string{http.MethodPost},
// 		},
// 		{
// 			desc:            "defaults will be supplied for empty methods with populated types",
// 			inputTypes:      "application/json",
// 			inputMethods:    "",
// 			expectedTypes:   []string{"application/json"},
// 			expectedMethods: []string{http.MethodGet},
// 		},
// 		{
// 			desc:            "proper cases will be used",
// 			inputTypes:      "TEXT/HTml",
// 			inputMethods:    "gEt",
// 			expectedTypes:   []string{"text/html"},
// 			expectedMethods: []string{http.MethodGet},
// 		},
// 		{
// 			desc:            "works with multiple values",
// 			inputTypes:      "text/html, application/json",
// 			inputMethods:    "get, post",
// 			expectedTypes:   []string{"text/html", "application/json"},
// 			expectedMethods: []string{http.MethodGet, http.MethodPost},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			config := httputil.ParseMonitoringConfig(test.inputTypes, test.inputMethods)

// 			if len(config.MonitoredTypes) != len(test.expectedTypes) {
// 				t.Errorf("Expected Types: '%v' | Got Types: '%v'", test.expectedTypes, config.MonitoredTypes)
// 			}

// 			for i, v := range test.expectedTypes {
// 				if v != config.MonitoredTypes[i] {
// 					t.Errorf("Expected Types: '%v' | Got Types: '%v'", test.expectedTypes, config.MonitoredTypes)
// 				}
// 			}

// 			if len(config.MonitoredMethods) != len(test.expectedMethods) {
// 				t.Errorf("Expected Methods: '%v' | Got Methods: '%v'", test.expectedMethods, config.MonitoredMethods)
// 			}

// 			for i, v := range test.expectedMethods {
// 				if v != config.MonitoredMethods[i] {
// 					t.Errorf("Expected Methods: '%v' | Got Methods: '%v'", test.expectedMethods, config.MonitoredMethods)
// 				}
// 			}
// 		})
// 	}
// }
