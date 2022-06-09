package httputil

import (
	"net/http"
	"strings"
)

// MonitoringConfig structure of data for handling configuration for
// controlling what content is monitored.
type MonitoringConfig struct {
	MonitoredTypes   []string `json:"types,omitempty" yaml:"types,omitempty" toml:"types,omitempty" export:"true"`
	MonitoredMethods []string `json:"methods,omitempty" yaml:"methods,omitempty" toml:"methods,omitempty" export:"true"`
}

// ParseMonitoringConfig parse types and methods strings into MonitoringConfig.
func ParseMonitoringConfig(types string, methods string) MonitoringConfig {
	var typesArray []string

	if len(types) == 0 {
		typesArray = []string{"text/html"}
	} else {
		typesArray = strings.Split(
			strings.ToLower(types),
			", ",
		)
	}

	var methodsArray []string

	if len(methods) == 0 {
		methodsArray = []string{http.MethodGet}
	} else {
		methodsArray = strings.Split(
			strings.ToUpper(methods),
			", ",
		)
	}

	return MonitoringConfig{
		MonitoredTypes:   typesArray,
		MonitoredMethods: methodsArray,
	}
}
