package rewrite_body

import (
	"regexp"

	"github.com/packruler/plugin-utils/httputil"
)

// +k8s:deepcopy-gen=true

// Rewrite holds one rewrite body configuration.
type Rewrite struct {
	Regex       string `json:"regex" yaml:"regex" toml:"regex" export:"true"`
	Replacement string `json:"replacement" yaml:"replacement" toml:"replacement" export:"true"`
}

// +k8s:deepcopy-gen=true

// Config holds the plugin configuration.
type Config struct {
	LastModified      bool                      `json:"lastModified" toml:"lastModified" yaml:"lastModified" export:"true"`
	Rewrites          []Rewrite                 `json:"rewrites" toml:"rewrites" yaml:"rewrites" export:"true"`
	LogLevel          int8                      `json:"logLevel" toml:"logLevel" yaml:"logLevel" export:"true"`
	MonintoringConfig httputil.MonitoringConfig `json:"monitor" toml:"monitor" yaml:"monitor" export:"true"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// +k8s:deepcopy-gen=true

type rewrite struct {
	regex       *regexp.Regexp
	replacement []byte
}
