package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfig_init(t *testing.T) {
	type fields struct {
		v *viper.Viper
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			tt.setup(t)

			// config
			c := &Config{
				v: tt.fields.v,
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

func Test_GetValue(t *testing.T) {
	// reset singleton for test
	instance = nil

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

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name  string
		want  map[string]any
		reset func()
		setup func()
	}{
		{
			name: "get config",
			want: map[string]any{
				"service.name": "TEST_NAME",
				"service.port": "TEST_PORT",
				"service.env":  "TEST_ENV",
			},
			reset: func() {
				instance = nil
			},
			setup: func() {
				tempDir := t.TempDir()
				subDir := filepath.Join(tempDir, "sub")
				if err := os.MkdirAll(subDir, 0755); err != nil {
					t.Fatalf("failed to create subdir: %v", err)
				}

				// write config.toml
				configContent := []byte(`
					[service]
					name = "TEST_NAME"
					port = "TEST_PORT"
					env = "TEST_ENV"
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

			// get config
			cfg := GetConfig()

			// check result
			for key, val := range tt.want {
				if cfg.Get(key) != val {
					t.Errorf("%s: expected %v, got %v", key, val, cfg.Get(key))
				}
			}

			// NOTE: add here if there are new configs
		})
	}
}

func TestSetTestConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]any
		setVals map[string]any
		setup   func()
	}{
		{
			name: "set test config",
			want: map[string]any{
				"test.test":   "test 1",
				"test.test_2": "test 2",
			},
			setVals: map[string]any{
				"test.test":   "test 1",
				"test.test_2": "test 2",
			},
		},
		{
			name: "reset config when being called",
			want: map[string]any{
				"test.test":   "test 1",
				"test.test_2": "test 2",
			},
			setVals: map[string]any{
				"test.test":   "test 1",
				"test.test_2": "test 2",
			},
			setup: func() {
				SetTestConfig(map[string]any{
					"test.abc": "test abc",
				})
			},
		},
	}

	for _, tt := range tests {
		if tt.setup != nil {
			tt.setup()
		}

		// set test config
		SetTestConfig(tt.setVals)

		for key, value := range tt.want {
			if GetConfig().Get(key) != value {
				t.Errorf("%s: expected %v, got %v", key, value, GetConfig().Get(key))
			}
		}
	}
}
