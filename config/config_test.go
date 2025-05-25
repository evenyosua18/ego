package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_init(t *testing.T) {
	type fields struct {
		ServiceConfig *Service
		v             *viper.Viper
	}

	tests := []struct {
		name                string
		fields              fields
		wantErr             bool
		setup               func(t *testing.T) string
		expectedConfigValue map[string]string
	}{
		{
			name: "init config",
			fields: fields{
				v: viper.New(),
			},
			wantErr: false,
			setup: func(t *testing.T) string {
				// get temporary directory
				tempDir := t.TempDir()

				// set subdirectory
				subDir := filepath.Join(tempDir, "sub")
				if err := os.MkdirAll(subDir, 0755); err != nil {
					t.Fatalf("failed to create subdir: %v", err)
				}

				// config bytes
				configBytes := []byte(`
					[service]
					name = "test-name"
					port = "test-port"
					env = "test-env"
				`)

				// set config path
				configPath := filepath.Join(subDir, "config.toml")
				if err := os.WriteFile(configPath, configBytes, 0644); err != nil {
					t.Fatalf("failed to write config.toml: %v", err)
				}

				// set env
				_ = os.Setenv(DirectoryConfigPath, tempDir)
				_ = os.Setenv(DirectoryConfigName, "sub")

				return configPath
			},
			expectedConfigValue: map[string]string{
				ServiceName: "test-name",
				ServiceEnv:  "test-env",
				ServicePort: "test-port",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			tt.setup(t)

			// config
			c := &Config{
				ServiceConfig: tt.fields.ServiceConfig,
				v:             tt.fields.v,
			}

			// initiate
			if err := c.init(); (err != nil) != tt.wantErr {
				t.Errorf("init() error = %v, wantErr %v", err, tt.wantErr)
			}

			// check the value
			if !tt.wantErr {
				for key, value := range tt.expectedConfigValue {
					if got := c.v.GetString(key); got != value {
						t.Errorf("expected %s = %s, got = '%s'", key, value, got)
					}
				}
			}
		})
	}
}
