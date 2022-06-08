package rewrite_body

import (
	"regexp"

	"github.com/packruler/rewrite-body/httputil"
)

// Rewrite holds one rewrite body configuration.
type Rewrite struct {
	Regex       string `json:"regex,omitempty" yaml:"regex,omitempty" toml:"regex,omitempty" export:"true"`
	Replacement string `json:"replacement,omitempty" yaml:"replacement,omitempty" toml:"replacement,omitempty" export:"true"`
}

// Config holds the plugin configuration.
type Config struct {
	LastModified      bool                       `json:"lastModified,omitempty" toml:"lastModified,omitempty" yaml:"lastModified,omitempty" export:"true"`
	Rewrites          []Rewrite                  `json:"rewrites,omitempty" toml:"rewrites,omitempty" yaml:"rewrites,omitempty" export:"true"`
	LogLevel          int8                       `json:"logLevel,omitempty" toml:"logLevel,omitempty" yaml:"logLevel,omitempty" export:"true"`
	MonintoringConfig *httputil.MonitoringConfig `json:"monitor,omitempty" toml:"monitor,omitempty" yaml:"monitor,omitempty" export:"true"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type rewrite struct {
	regex       *regexp.Regexp
	replacement []byte
}
