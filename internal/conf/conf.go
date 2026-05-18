package conf

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Bootstrap struct {
	Server    Server            `yaml:"server"`
	Discovery map[string]string `yaml:"discovery"`
	DB        string            `yaml:"db"`
	NATS      string            `yaml:"nats"`
}

type Server struct {
	HTTP HTTP `yaml:"http"`
}

type HTTP struct {
	Addr string `yaml:"addr"`
}

func Load(path string) (*Bootstrap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var bc Bootstrap
	if err := yaml.Unmarshal(data, &bc); err != nil {
		return nil, err
	}

	return &bc, nil
}
