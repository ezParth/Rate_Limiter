package configuration

type YAMLCONFIG struct {
	Server SERVER `yaml:"server"`
}

type SERVER struct {
	Listen   int         `yaml:"listen"`
	Workers  int         `yaml:"workers"`
	Upstream []UPSTREAMS `yaml:"upstream"`
	Rules    []RULES     `yaml:"rules"`
}

type UPSTREAMS struct {
	ID     string `yaml:"id"`
	Server string `yaml:"server"`
}

type RULES struct {
	Path     string   `yaml:"path"`
	Upstream []string `yaml:"upstream"`
}
