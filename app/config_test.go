package app

import (
	"github.com/evenyosua18/ego/config"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

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
		AppConfig *App
		v         *viper.Viper
	}
	tests := []struct {
		name         string
		fields       fields
		setVals      map[string]any
		expectedVals map[string]string
		expectedConf Config
	}{
		{
			name: "build config, all values filled",
			fields: fields{
				AppConfig: &App{},
			},
			setVals: map[string]any{
				ServiceName: "test-name",
				ServiceEnv:  "local",
				ServicePort: ":8080",
			},
			expectedConf: Config{
				AppConfig: &App{
					Name: "test-name",
					Port: ":8080",
					Env:  "local",
				},
			},
		},
		{
			name: "build config, use default value",
			fields: fields{
				AppConfig: &App{},
				v:         viper.New(),
			},
			setVals: map[string]any{
				ServiceEnv: "local",
			},
			expectedConf: Config{
				AppConfig: &App{
					Name: DefaultServiceName,
					Port: DefaultServicePort,
					Env:  "local",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.SetTestConfig(tt.setVals)

			c := &Config{
				AppConfig: tt.fields.AppConfig,
			}

			c.build()

			// expected service config
			if !reflect.DeepEqual(c.AppConfig, tt.expectedConf.AppConfig) {
				t.Errorf("ServiceConfig mismatch: got %+v, want %+v", c.AppConfig, tt.expectedConf.AppConfig)
			}

			// NOTE: add here if there are new configs
		})
	}
}

func TestConfig_getOrDefault(t *testing.T) {
	type args struct {
		key        string
		defaultVal string
		setVal     string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "value is empty",
			args: args{
				key:        "test",
				defaultVal: "test",
				setVal:     "",
			},
			want: "test",
		},
		{
			name: "value is not empty",
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
			// set config
			config.SetTestConfig(map[string]any{
				tt.args.key: tt.args.setVal,
			})

			// build app config
			appConfig := Config{}

			if got := appConfig.getOrDefault(tt.args.key, tt.args.defaultVal); got != tt.want {
				t.Errorf("getOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
