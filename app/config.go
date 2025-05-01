package app

import (
	"github.com/evenyosua18/ego/config"
	"strings"
)

type Config struct {
	ServiceConfig *Service
}

type Service struct {
	Name string
	Port string
	Env  string
}

const (
	ServiceName        = "service.name"
	ServicePort        = "service.port"
	ServiceEnv         = "service.env"
	DefaultServiceName = "local"
	DefaultServicePort = ":8080"
	DefaultServiceEnv  = "dev"
)

func BuildConfiguration() *Config {
	cfg := &Config{}

	// setup service configuration
	cfg.ServiceConfig = &Service{
		Name: config.GetString(ServiceName),
		Port: config.GetString(ServicePort),
		Env:  config.GetString(ServiceEnv),
	}

	// set default value
	if cfg.ServiceConfig.Name == "" {
		cfg.ServiceConfig.Name = DefaultServiceName
	}

	if cfg.ServiceConfig.Env == "" {
		cfg.ServiceConfig.Env = DefaultServiceEnv
	}

	if cfg.ServiceConfig.Port == "" {
		cfg.ServiceConfig.Port = DefaultServicePort
	} else if !strings.HasPrefix(cfg.ServiceConfig.Port, ":") {
		cfg.ServiceConfig.Port = ":" + cfg.ServiceConfig.Port
	}

	return cfg
}
