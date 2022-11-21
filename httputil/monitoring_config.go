package httputil

import "net/http"
import "strings"

// MonitoringConfig structure of data for handling configuration for
// controlling what content is monitored.
type MonitoringConfig struct {
	Types   []string `json:"types,omitempty" yaml:"types,omitempty" toml:"types,omitempty" export:"true"`
	Methods []string `json:"methods,omitempty" yaml:"methods,omitempty" toml:"methods,omitempty" export:"true"`
}

// EnsureDefaults check Types and Methods for empty arrays and apply default values if found.
func (config *MonitoringConfig) EnsureDefaults() {
	if len(config.Methods) == 0 {
		config.Methods = []string{http.MethodGet}
	}

	if len(config.Types) == 0 {
		config.Types = []string{"text/html"}
	}

	// handle mangled/flattened monitoring config
	if len(config.Methods) == 1 {
		for _, configMethod := range config.Methods {
			if strings.Contains(configMethod,"║") {
				config.Methods = strings.Split(strings.ReplaceAll(configMethod,"║24║",""),"║")
			}
		}
	}

	if len(config.Types) == 1 {
		for _, configType := range config.Types {
			if strings.Contains(configType,"║") {
				config.Types = strings.Split(strings.ReplaceAll(configType,"║24║",""),"║")
			}
		}
	}
}
