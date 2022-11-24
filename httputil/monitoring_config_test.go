package httputil_test

import (
	"net/http"
	"testing"

	"github.com/packruler/rewrite-body/httputil"
)

func TestMonitoringConfigParsing(t *testing.T) {
	tests := []struct {
		desc           string
		inputConfig    httputil.MonitoringConfig
		expectedConfig httputil.MonitoringConfig
	}{
		{
			desc: "defaults will be supplied for empty arrays",
			inputConfig: httputil.MonitoringConfig{
				Types:   []string{},
				Methods: []string{},
			},
			expectedConfig: httputil.MonitoringConfig{
				Types:   []string{"text/html"},
				Methods: []string{http.MethodGet},
			},
		},
		{
			desc: "defaults will be supplied for empty types with populated methods",
			inputConfig: httputil.MonitoringConfig{
				Types:   []string{},
				Methods: []string{"POST"},
			},
			expectedConfig: httputil.MonitoringConfig{
				Types:   []string{"text/html"},
				Methods: []string{http.MethodPost},
			},
		},
		{
			desc: "defaults will be supplied for populated types with empty methods",
			inputConfig: httputil.MonitoringConfig{
				Types:   []string{"application/json"},
				Methods: []string{},
			},
			expectedConfig: httputil.MonitoringConfig{
				Types:   []string{"application/json"},
				Methods: []string{http.MethodGet},
			},
		},
		{
			desc: "handle weird yaml parsing",
			inputConfig: httputil.MonitoringConfig{
				Types:   []string{"║24║application/javascript║application/json"},
				Methods: []string{},
			},
			expectedConfig: httputil.MonitoringConfig{
				Types:   []string{"application/javascript", "application/json"},
				Methods: []string{http.MethodGet},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			config := test.inputConfig

			config.EnsureDefaults()
			config.EnsureProperFormat()

			if len(config.Types) != len(test.expectedConfig.Types) {
				t.Errorf("Expected Types: '%v' | Got Types: '%v'", test.expectedConfig.Types, config.Types)
			}

			for i, v := range test.expectedConfig.Types {
				if v != config.Types[i] {
					t.Errorf("Expected Types: '%v' | Got Types: '%v'", test.expectedConfig.Types, config.Types)
				}
			}

			if len(config.Methods) != len(test.expectedConfig.Methods) {
				t.Errorf("Expected Methods: '%v' | Got Methods: '%v'", test.expectedConfig.Methods, config.Methods)
			}

			for i, v := range test.expectedConfig.Methods {
				if v != config.Methods[i] {
					t.Errorf("Expected Methods: '%v' | Got Methods: '%v'", test.expectedConfig.Methods, config.Methods)
				}
			}
		})
	}
}
