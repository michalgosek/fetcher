package config

import (
	"errors"
	"fetcher/internal/server"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server server.Config
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
	if err != nil {
		return nil, ErrConfigFileNotFound
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file %s failed: %v", c.path, err)
	}
	return &c, nil
}

var ErrConfigFileNotFound = errors.New("config file not exist or path to the file is invalid")
