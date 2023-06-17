package util

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

type ElasticsearchInfo struct {
	// `yaml:"address"` is a "field tag"
	//    - This tag is instructing the YAML decoder (from a library, in this case, gopkg.in/yaml.v3) that when it sees a field named "address" in your YAML data, it should populate that data into the "Address" field in your struct.
	//
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type GCSInfo struct {
	Bucket string `yaml:"bucket"`
}

type TokenInfo struct {
	Secret string `yaml:"secret"`
}

type ApplicationConfig struct {
	ElasticsearchConfig *ElasticsearchInfo `yaml:"elasticsearch"`
	GCSConfig           *GCSInfo           `yaml:"gcs"`
	TokenConfig         *TokenInfo         `yaml:"token"`
}

func LoadApplicationConfig(configDir, configFile string) (*ApplicationConfig, error) {
	content, err := ioutil.ReadFile(filepath.Join(configDir, configFile))
	if err != nil {
		return nil, err
	}

	var config ApplicationConfig
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
