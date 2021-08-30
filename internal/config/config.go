package config

import (
	"errors"
	"fetcher/internal/server"
	"fetcher/internal/worker"
	"fmt"
	"io/fs"

	"github.com/spf13/viper"
)

type Config struct {
	Server server.Config
	Worker worker.Config
	path   string
}

type Option func(c *Config)

func WithPath(s string) Option {
	return func(c *Config) {
		c.path = s
	}
}

func New(opts ...Option) (*Config, error) {
	c := Config{
		path: "./config.yml",
	}

	for _, o := range opts {
		o(&c)
	}

	viper.SetConfigFile(c.path)
	err := viper.ReadInConfig()
	var target *fs.PathError
	if err != nil && errors.As(err, &target) {
		return nil, ErrConfigFileNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("reading config failed: %v", err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file %s failed: %v", c.path, err)
	}
	return &c, nil
}

var ErrConfigFileNotFound = errors.New("config file not exist or path to the file is invalid")
