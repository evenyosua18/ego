package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	FileName = "config"
	FileType = "toml"

	DirectoryConfigName = "CONFIG_DIR"
	DirectoryConfigPath = "CONFIG_PATH"
	DirectoryConfigRoot = "CONFIG_ROOT"

	ServiceName        = "service.name"
	ServicePort        = "service.port"
	ServiceEnv         = "service.env"
	DefaultServiceName = "local"
	DefaultServicePort = ":8080"
	DefaultServiceEnv  = "dev"

	DefaultConfigPath = "./config"
)

var (
	once     sync.Once
	instance *Config
)

type Config struct {
	ServiceConfig *Service

	v *viper.Viper
}

type Service struct {
	Name string
	Port string
	Env  string
}

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{v: viper.New()}
		if err := instance.init(); err != nil {
			panic(fmt.Errorf("error initializing config: %v", err))
		}
		instance.build()
	})

	return instance
}

// mirroring viper function

func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *Config) GetInt(key string) int {
	return c.v.GetInt(key)
}

func (c *Config) GetBool(key string) bool {
	return c.v.GetBool(key)
}

func (c *Config) Get(key string) any {
	return c.v.Get(key)
}

func (c *Config) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(key)
}

func (c *Config) GetIntSlice(key string) []int {
	return c.v.GetIntSlice(key)
}

func (c *Config) GetInt32(key string) int32 {
	return c.v.GetInt32(key)
}

func (c *Config) GetInt64(key string) int64 {
	return c.v.GetInt64(key)
}

func (c *Config) GetStringMapStringSlice(key string) map[string][]string {
	return c.v.GetStringMapStringSlice(key)
}

func (c *Config) GetFloat64(key string) float64 {
	return c.v.GetFloat64(key)
}

func (c *Config) GetStringMap(key string) map[string]any {
	return c.v.GetStringMap(key)
}

func (c *Config) GetStringMapString(key string) map[string]string {
	return c.v.GetStringMapString(key)
}

func (c *Config) GetUint(key string) uint {
	return c.v.GetUint(key)
}

func (c *Config) GetUint32(key string) uint32 {
	return c.v.GetUint32(key)
}

func (c *Config) GetUint64(key string) uint64 {
	return c.v.GetUint64(key)
}

func (c *Config) IsParentKeyExists(key string) bool {
	return c.v.Sub(key) != nil
}

func (c *Config) init() error {
	// get path
	path := os.Getenv(DirectoryConfigPath)
	root := os.Getenv(DirectoryConfigRoot)

	if root == "" && path == "" {
		path = DefaultConfigPath
	} else if root != "" {
		path = root
	} else if !filepath.IsAbs(path) {
		path = "./" + path
	}

	// get directory
	dir := os.Getenv(DirectoryConfigName)

	if dir != "" {
		path += "/" + dir
	}

	c.v.AddConfigPath(path)
	c.v.SetConfigName(FileName)
	c.v.SetConfigType(FileType)
	c.v.AutomaticEnv()
	c.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return c.v.ReadInConfig()
}

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
	val := c.v.GetString(key)
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
