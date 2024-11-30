package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	RootDir string `yaml:"rootDir"`
	HTTP    HTTP   `yaml:"http"`
	IPXE    []IPXE `yaml:"ipxe"`
}

type HTTP struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type IPXE struct {
	Name       string            `yaml:"name"`
	IPs        []string          `yaml:"ips"`
	KernelArgs map[string]string `yaml:"kernelArgs"`
}

func FromFile(path string) (*Config, error) {
	cfg := &Config{}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) GetRootDir() string {
	if c.RootDir == "" {
		return "./"
	}
	return c.RootDir
}
