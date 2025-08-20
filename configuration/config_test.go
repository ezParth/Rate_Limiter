package configuration

import (
	"testing"
)

// just so you noobs know, these test cases are not for custom yaml configurations. Write your own tests :/

func TestLoadYaml(t *testing.T) {
	cfg, err := LoadYAML("../config.yml")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Server.Listen != 8080 {
		t.Errorf("expected listen=8080, got %d", cfg.Server.Listen)
	}

	if len(cfg.Server.Upstream) != 2 {
		t.Errorf("expected 2 upstreams, got %d", len(cfg.Server.Upstream))
	}
}

func TestInvalidYAML(t *testing.T) {
	_, err := LoadYAML("invalid.yml")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
