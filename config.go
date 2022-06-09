package rewrite_body

import (
	"regexp"

	"github.com/packruler/plugin-utils/httputil"
)

// Rewrite holds one rewrite body configuration.
type Rewrite struct {
	Regex       string `json:"regex" yaml:"regex" toml:"regex"`
	Replacement string `json:"replacement" yaml:"replacement" toml:"replacement"`
}

// Config holds the plugin configuration.
type Config struct {
	LastModified      bool                      `json:"lastModified" toml:"lastModified" yaml:"lastModified"`
	Rewrites          []Rewrite                 `json:"rewrites" toml:"rewrites" yaml:"rewrites"`
	LogLevel          int8                      `json:"logLevel" toml:"logLevel" yaml:"logLevel"`
	MonintoringConfig httputil.MonitoringConfig `json:"monitor" toml:"monitor" yaml:"monitor"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type rewrite struct {
	regex       *regexp.Regexp
	replacement []byte
}
