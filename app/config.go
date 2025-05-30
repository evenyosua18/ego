package app

import (
	"github.com/evenyosua18/ego/config"
	"strings"
)

const (
	ServiceName = "service.name"
	ServicePort = "service.port"
	ServiceEnv  = "service.env"

	DefaultServiceName = "local"
	DefaultServicePort = ":8080"
	DefaultServiceEnv  = "dev"
)

type (
	Config struct {
		ServiceConfig *Service
		CodeConfig    *Code
	}

	Code struct {
		FilePath string
	}

	Service struct {
		Name string
		Port string
		Env  string
	}
)

func (c *Config) build() {
	// setup service configuration
	c.ServiceConfig = &Service{
		Name: c.getOrDefault(ServiceName, DefaultServiceName),
		Port: normalizePort(c.getOrDefault(ServicePort, DefaultServicePort)),
		Env:  c.getOrDefault(ServiceEnv, DefaultServiceEnv),
	}

	return
}

func (c *Config) getOrDefault(key string, defaultVal string) string {
	val := config.GetConfig().GetString(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func normalizePort(port string) string {
	if !strings.HasPrefix(port, ":") {
		return ":" + port
	}
	return port
}
