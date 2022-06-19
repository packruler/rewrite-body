package rewrite_body

// This file content is commented to ignore dependencies only needed for testing

// import (
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"reflect"
// 	"testing"

// 	"github.com/BurntSushi/toml"
// 	"github.com/packruler/rewrite-body/httputil"
// 	"github.com/stretchr/testify/assert"
// 	"gopkg.in/yaml.v3"
// )

// type testStruct struct {
// 	Config Config `json:"config,omitempty" toml:"config,omitempty" yaml:"config,omitempty"`
// }

// func TestConfigParsing(t *testing.T) {
// 	cfg := &testStruct{}

// 	_, err := toml.DecodeFile("./fixtures/sample.toml", &cfg)
// 	if err != nil {
// 		t.Errorf("Unable to decode file: %v", err)
// 	}

// 	reflection := reflect.ValueOf(cfg)
// 	log.Println("Decoded: ", reflection)

// 	cfgCopy := cfg
// 	assert.Equal(t, reflect.ValueOf(cfgCopy), reflect.ValueOf(cfg))
// 	assert.Equal(t, reflect.ValueOf(cfgCopy), reflect.ValueOf(cfg))
// 	assert.Equal(t, cfgCopy, cfg)

// 	expected := &testStruct{
// 		Config{
// 			LastModified: true,
// 			LogLevel:     0,
// 			MonintoringConfig: httputil.MonitoringConfig{
// 				MonitoredTypes:   []string{"text/html"},
// 				MonitoredMethods: []string{http.MethodGet},
// 			},
// 			Rewrites: []Rewrite{
// 				{
// 					Regex:       "test",
// 					Replacement: "test",
// 				},
// 			},
// 		},
// 	}
// 	assert.Equal(t, expected, cfg)
// }

// func TestConfigYamlParsing(t *testing.T) {
// 	cfg := &testStruct{}

// 	yamlFile, err := ioutil.ReadFile("./fixtures/sample.yaml")
// 	if err != nil {
// 		t.Error("Unable to read file.", err)
// 	}

// 	err = yaml.Unmarshal(yamlFile, &cfg)
// 	if err != nil {
// 		t.Errorf("Unable to decode file: %v", err)
// 	}

// 	reflection := reflect.ValueOf(cfg)
// 	log.Println("Decoded: ", reflection)

// 	cfgCopy := cfg
// 	assert.Equal(t, reflect.ValueOf(cfgCopy), reflect.ValueOf(cfg))
// 	assert.Equal(t, reflect.ValueOf(cfgCopy), reflect.ValueOf(cfg))
// 	assert.Equal(t, cfgCopy, cfg)

// 	expected := &testStruct{
// 		Config{
// 			LastModified: true,
// 			LogLevel:     0,
// 			MonintoringConfig: httputil.MonitoringConfig{
// 				MonitoredTypes:   []string{"text/html"},
// 				MonitoredMethods: []string{http.MethodGet},
// 			},
// 			Rewrites: []Rewrite{
// 				{
// 					Regex:       "test",
// 					Replacement: "test",
// 				},
// 			},
// 		},
// 	}
// 	assert.Equal(t, expected, cfg)
// }
