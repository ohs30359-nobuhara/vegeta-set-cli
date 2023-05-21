package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Scenario struct {
	Url    string  `yaml:"url"`
	Method string  `yaml:"method"`
	Ratio  string  `yaml:"ratio"`
	Value  *string `yaml:"value"`
}

type Config struct {
	Scenario []Scenario `yaml:"scenario"`
	Rate     int        `yaml:"rate"`
}

func Load(path string) (Config, error) {
	content, e := os.ReadFile(path)
	if e != nil {
		return Config{}, e
	}

	var config Config
	if e := yaml.Unmarshal(content, &config); e != nil {
		return Config{}, e
	}
	return config, nil
}
