package rewrite_body

import (
	"log"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Config Config `json:"config,omitempty" toml:"config,omitempty" yaml:"config,omitempty"`
}

func TestConfigParsing(t *testing.T) {
	cfg := &testStruct{}

	_, err := toml.DecodeFile("./fixtures/sample.toml", &cfg)
	if err != nil {
		t.Errorf("Unable to decode file: %v", err)
	}

	log.Printf("Decoded: %v", cfg)

	cfgCopy := cfg
	assert.Equal(t, reflect.ValueOf(cfgCopy), reflect.ValueOf(cfg))
	assert.Equal(t, reflect.ValueOf(cfgCopy), reflect.ValueOf(cfg))
	assert.Equal(t, cfgCopy, cfg)
}
