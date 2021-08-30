package config_test

import (
	"errors"
	"fetcher/internal/config"
	"testing"
	"time"
)

func Test_ConfigReadWrongPath_Unit(t *testing.T) {
	const path = "./random_file.yml"
	cfg, err := config.New(config.WithPath(path))
	if !errors.Is(err, config.ErrConfigFileNotFound) {
		t.Fatalf("expected to get err type %T; got %T", config.ErrConfigFileNotFound, err)
	}

	if cfg != nil {
		t.Fatalf("expected to get nil cfg, got %v", cfg)
	}
}

func Test_ConfigReadFromFile_Unit(t *testing.T) {
	const path = "./example_cfg.yml"
	cfg, err := config.New(config.WithPath(path))
	if err != nil {
		t.Fatalf("expected to get nil err; got %v", err)
	}

	if cfg.Server.Addr != "localhost:8090" {
		t.Fatalf("expected to get addr equals localhost:8090, got %s", cfg.Server.Addr)
	}

	if cfg.Server.ShutdownTime != 10*time.Second {
		t.Fatalf("expected to get shudtown time equals 10s; got %v", cfg.Server.ShutdownTime.Seconds())
	}
}
