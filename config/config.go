package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	FileName = "config"
	FileType = "toml"

	DirectoryConfigName = "CONFIG_DIR"
	DirectoryConfigPath = "CONFIG_PATH"
	DirectoryConfigRoot = "CONFIG_ROOT"

	DefaultConfigPath = "./config"
)

func GetConfig() {
	// get path
	path := os.Getenv(DirectoryConfigPath)
	root := os.Getenv(DirectoryConfigRoot)

	if root == "" && path == "" {
		path = DefaultConfigPath
	} else if root != "" {
		path = root
	} else {
		path = "./" + path
	}

	// get directory
	dir := os.Getenv(DirectoryConfigName)

	if dir != "" {
		path += "/" + dir
	}

	viper.AddConfigPath(path)
	viper.SetConfigName(FileName)
	viper.SetConfigType(FileType)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %v", err))
	}
}

// mirroring viper function

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func Get(key string) any {
	return viper.Get(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func GetIntSlice(key string) []int {
	return viper.GetIntSlice(key)
}

func GetInt32(key string) int32 {
	return viper.GetInt32(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

func GetStringMapStringSlice(key string) map[string][]string {
	return viper.GetStringMapStringSlice(key)
}

func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func GetStringMap(key string) map[string]any {
	return viper.GetStringMap(key)
}

func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func GetUint(key string) uint {
	return viper.GetUint(key)
}

func GetUint32(key string) uint32 {
	return viper.GetUint32(key)
}

func GetUint64(key string) uint64 {
	return viper.GetUint64(key)
}

func IsParentKeyExists(key string) bool {
	return viper.Sub(key) != nil
}
