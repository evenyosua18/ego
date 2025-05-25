package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
	"sync"
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

func Test_normalizePort(t *testing.T) {
	type args struct {
		port string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "without ':'",
			args: args{
				port: "8080",
			},
			want: ":8080",
		},
		{
			name: "with ':'",
			args: args{
				port: ":8080",
			},
			want: ":8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizePort(tt.args.port); got != tt.want {
				t.Errorf("normalizePort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_build(t *testing.T) {
	type fields struct {
		ServiceConfig *Service
		v             *viper.Viper
	}
	tests := []struct {
		name         string
		fields       fields
		setVals      map[string]string
		expectedVals map[string]string
		expectedConf Config
	}{
		{
			name: "build config, all values filled",
			fields: fields{
				ServiceConfig: &Service{},
				v:             viper.New(),
			},
			setVals: map[string]string{
				ServiceName: "test-name",
				ServiceEnv:  "test-env",
				ServicePort: ":8080",
			},
			expectedConf: Config{
				ServiceConfig: &Service{
					Name: "test-name",
					Port: ":8080",
					Env:  "test-env",
				},
			},
		},
		{
			name: "build config, use default value",
			fields: fields{
				ServiceConfig: &Service{},
				v:             viper.New(),
			},
			setVals: nil,
			expectedConf: Config{
				ServiceConfig: &Service{
					Name: DefaultServiceName,
					Port: DefaultServicePort,
					Env:  DefaultServiceEnv,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ServiceConfig: tt.fields.ServiceConfig,
				v:             tt.fields.v,
			}

			// set val
			if tt.setVals != nil {
				for key, value := range tt.setVals {
					c.v.Set(key, value)
				}
			}

			c.build()

			// expected service config
			if !reflect.DeepEqual(c.ServiceConfig, tt.expectedConf.ServiceConfig) {
				t.Errorf("ServiceConfig mismatch: got %+v, want %+v", c.ServiceConfig, tt.expectedConf.ServiceConfig)
			}

			// NOTE: add here if there are new configs
		})
	}
}

func TestConfig_getOrDefault(t *testing.T) {
	type fields struct {
		ServiceConfig *Service
		v             *viper.Viper
	}

	type args struct {
		key        string
		defaultVal string
		setVal     string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      string
		keyValues map[string]string
	}{
		{
			name: "value is empty",
			fields: fields{
				ServiceConfig: &Service{},
			},
			args: args{
				key:        "test",
				defaultVal: "test",
				setVal:     "",
			},
			want: "test",
		},
		{
			name: "value is not empty",
			fields: fields{
				ServiceConfig: &Service{},
			},
			args: args{
				key:        "test",
				defaultVal: "test",
				setVal:     "test-123",
			},
			want: "test-123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// new viper
			v := viper.New()
			if tt.args.setVal != "" {
				v.Set(tt.args.key, tt.args.setVal)
			}

			// config
			c := &Config{
				ServiceConfig: tt.fields.ServiceConfig,
				v:             v,
			}

			if got := c.getOrDefault(tt.args.key, tt.args.defaultVal); got != tt.want {
				t.Errorf("getOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	resetConfig := func() {
		// reset singleton for test
		instance = nil
		once = sync.Once{}
	}

	tests := []struct {
		name  string
		want  *Config
		reset func()
		setup func()
	}{
		{
			name: "get config",
			want: &Config{
				ServiceConfig: &Service{
					Name: DefaultServiceName,
					Port: DefaultServicePort,
					Env:  DefaultServiceEnv,
				},
				v: viper.New(),
			},
			reset: resetConfig,
			setup: func() {
				tempDir := t.TempDir()
				subDir := filepath.Join(tempDir, "sub")
				if err := os.MkdirAll(subDir, 0755); err != nil {
					t.Fatalf("failed to create subdir: %v", err)
				}

				// write config.toml
				configContent := []byte(`
					[service]
					name = ""
					port = ""
					env = ""
				`)

				if err := os.WriteFile(filepath.Join(subDir, "config.toml"), configContent, 0644); err != nil {
					t.Fatalf("failed to write config.toml: %v", err)
				}

				// set env vars
				t.Setenv(DirectoryConfigPath, tempDir)
				t.Setenv(DirectoryConfigName, "sub")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset
			tt.reset()

			// setup
			tt.setup()

			// check result
			cfg := GetConfig()

			if !reflect.DeepEqual(cfg.ServiceConfig, tt.want.ServiceConfig) {
				t.Errorf("service config = %v, want %v", cfg.ServiceConfig, tt.want.ServiceConfig)
			}

			// NOTE: add here if there are new configs
		})
	}
}

func Test_GetValue(t *testing.T) {
	// reset singleton for test
	instance = nil
	once = sync.Once{}

	// setup
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	// write config.toml
	configContent := []byte(`
		string_val = "TEST"
		int_val = 123
		bool_val = true
		float_val = 1.23
		string_slice_val = ["a", "b", "c"]
		int_slice_val = [1, 2, 3]

		[string_map]
		key1 = "value1"
		key2 = "value2"

		[string_map_slice]
		key1 = ["x", "y"]
		key2 = ["z"]
	`)

	if err := os.WriteFile(filepath.Join(subDir, "config.toml"), configContent, 0644); err != nil {
		t.Fatalf("failed to write config.toml: %v", err)
	}

	// set env vars
	t.Setenv(DirectoryConfigPath, tempDir)
	t.Setenv(DirectoryConfigName, "sub")

	// get config
	cfg := GetConfig()

	tests := []struct {
		name     string
		fn       func() any
		expected any
	}{
		{
			name: "Get",
			fn: func() any {
				return cfg.Get("string_val")
			},
			expected: "TEST",
		},
		{
			name: "GetString",
			fn: func() any {
				return cfg.GetString("string_val")
			},
			expected: "TEST",
		},
		{
			name: "GetInt",
			fn: func() any {
				return cfg.GetInt("int_val")
			},
			expected: 123,
		},
		{
			name: "GetBool",
			fn: func() any {
				return cfg.GetBool("bool_val")
			},
			expected: true,
		},
		{
			name: "GetStringSlice",
			fn: func() any {
				return cfg.GetStringSlice("string_slice_val")
			},
			expected: []string{"a", "b", "c"},
		},
		{
			name: "GetIntSlice",
			fn: func() any {
				return cfg.GetIntSlice("int_slice_val")
			},
			expected: []int{1, 2, 3},
		},
		{
			name: "GetInt32",
			fn: func() any {
				return cfg.GetInt32("int_val")
			},
			expected: int32(123),
		},
		{
			name: "GetInt64",
			fn: func() any {
				return cfg.GetInt64("int_val")
			},
			expected: int64(123),
		},
		{
			name: "GetStringMapStringSlice",
			fn: func() any {
				return cfg.GetStringMapStringSlice("string_map_slice")
			},
			expected: map[string][]string{
				"key1": {"x", "y"},
				"key2": {"z"},
			},
		},
		{
			name: "GetFloat64",
			fn: func() any {
				return cfg.GetFloat64("float_val")
			},
			expected: 1.23,
		},
		{
			name: "GetStringMap",
			fn: func() any {
				return cfg.GetStringMap("string_map")
			},
			expected: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "GetStringMapString",
			fn: func() any {
				return cfg.GetStringMapString("string_map")
			},
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "GetUint",
			fn: func() any {
				return cfg.GetUint("int_val")
			},
			expected: uint(123),
		},
		{
			name: "GetUint32",
			fn: func() any {
				return cfg.GetUint32("int_val")
			},
			expected: uint32(123),
		},
		{
			name: "GetUint64",
			fn: func() any {
				return cfg.GetUint64("int_val")
			},
			expected: uint64(123),
		},
		{
			name: "IsParentKeyExists - value exists",
			fn: func() any {
				return cfg.IsParentKeyExists("string_map")
			},
			expected: true,
		},
		{
			name: "IsParentKeyExists - value not exists",
			fn: func() any {
				return cfg.IsParentKeyExists("string_map_not_exists")
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		got := tt.fn()
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("%s: expected %v (%T), got %v (%T)", tt.name, tt.expected, tt.expected, got, got)
		}
	}
}
