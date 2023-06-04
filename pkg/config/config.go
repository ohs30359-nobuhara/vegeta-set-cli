package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Scenario struct {
	Url    string  `yaml:"url"`
	Method string  `yaml:"method"`
	Ratio  int     `yaml:"ratio"`
	Value  *string `yaml:"value"`
}

type Tester struct {
	Limit int `yaml:"limit"`
}

type Config struct {
	Scenario []Scenario    `yaml:"scenario"`
	Rate     int           `yaml:"rate"`
	Duration time.Duration `yaml:"duration"`
	Tester   Tester        `yaml:"tester"`
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
