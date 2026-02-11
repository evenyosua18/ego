package app

import (
	"reflect"
	"testing"
	"time"

	"github.com/evenyosua18/ego/config"
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
		AppConfig      *App
		CodeConfig     *Code
		DatabaseConfig *Database
		TracerConfig   *Tracer
		RouterConfig   *Router
		LoggerConfig   *Logger
		CacheConfig    *Cache
	}
	tests := []struct {
		name         string
		fields       fields
		setVals      map[string]any
		expectedConf Config
	}{
		{
			name:   "build config, all values filled",
			fields: fields{},
			setVals: map[string]any{
				ServiceName:             "test-service",
				ServiceEnv:              "production",
				ServiceShutdownTimeout:  "60s",
				CustomCodeFilePath:      "test.yaml",
				DatabaseDriver:          "postgres",
				DatabaseUser:            "user",
				DatabasePassword:        "password",
				DatabaseHost:            "localhost",
				DatabasePort:            "5432",
				DatabaseName:            "dbname",
				DatabaseProtocol:        "tcp",
				DatabaseParams:          "sslmode=disable",
				DatabaseMaxOpenConns:    50,
				DatabaseMaxIdleConns:    10,
				DatabaseConnMaxLifetime: "1h",
				DatabaseConnMaxIdleTime: "10m",
				TracerDSN:               "http://tracer",
				TracerSampleRate:        0.5,
				TracerFlushTime:         "2s",
				RouterPrefix:            "/api",
				RouterPort:              ":9090",
				RouterShowRegistered:    true,
				RouterHtmlPath:          "/html",
				RouterAllowOrigins:      []string{"*"},
				RouterAllowMethods:      []string{"GET"},
				RouterAllowHeaders:      []string{"Content-Type"},
				RouterAllowCredentials:  true,
				RouterDocPath:           "/docs",
				RouterMaxConnection:     1000,
				RouterMaxLimit:          50,
				LoggerLevel:             "debug",
				CacheRedisAddr:          "localhost:6379",
				CacheRedisDB:            1,
				CacheRedisMaxRetries:    5,
				CacheRedisMinIdleConns:  5,
				CacheRedisPoolSize:      20,
				CacheRedisIdleTimeout:   "10m",
				CacheRedisReadTimeout:   "1s",
				CacheRedisWriteTimeout:  "1s",
				CacheRedisDialTimeout:   "1s",
				CacheRedisPoolTimeout:   "1s",
			},
			expectedConf: Config{
				AppConfig: &App{
					Name:            "test-service",
					Env:             "production",
					ShutdownTimeout: 60 * time.Second,
				},
				CodeConfig: &Code{
					Filename: "test.yaml",
				},
				DatabaseConfig: &Database{
					Driver:          "postgres",
					User:            "user",
					Password:        "password",
					Host:            "localhost",
					Port:            "5432",
					Name:            "dbname",
					Protocol:        "tcp",
					Params:          "sslmode=disable",
					MaxOpenConns:    50,
					MaxIdleConns:    10,
					ConnMaxLifetime: 1 * time.Hour,
					ConnMaxIdleTime: 10 * time.Minute,
				},
				TracerConfig: &Tracer{
					DSN:        "http://tracer",
					SampleRate: 0.5,
					FlushTime:  "2s",
				},
				RouterConfig: &Router{
					Prefix:           "/api",
					Port:             ":9090",
					ShowRegistered:   true,
					HtmlPath:         "/html",
					AllowOrigins:     []string{"*"},
					AllowMethods:     []string{"GET"},
					AllowHeaders:     []string{"Content-Type"},
					AllowCredentials: true,
					DocPath:          "/docs",
					MaxConnection:    1000,
					MaxLimit:         50,
				},
				LoggerConfig: &Logger{
					Level: "debug",
				},
				CacheConfig: &Cache{
					Redis: &Redis{
						Addr:         "localhost:6379",
						DB:           1,
						MaxRetries:   5,
						MinIdleConns: 5,
						PoolSize:     20,
						IdleTimeout:  10 * time.Minute,
						ReadTimeout:  1 * time.Second,
						WriteTimeout: 1 * time.Second,
						DialTimeout:  1 * time.Second,
						PoolTimeout:  1 * time.Second,
					},
				},
			},
		},
		{
			name:   "build config, use default value",
			fields: fields{},
			setVals: map[string]any{
				ServiceEnv: "local",
			},
			expectedConf: Config{
				AppConfig: &App{
					Name:            DefaultServiceName,
					Env:             "local",
					ShutdownTimeout: DefaultServiceShutdownTimeout,
				},
				CodeConfig: &Code{
					Filename: "",
				},
				DatabaseConfig: &Database{
					Driver:          DefaultDatabaseDriver,
					User:            DefaultDatabaseUser,
					Password:        "",
					Host:            DefaultDatabaseHost,
					Port:            DefaultDatabasePort,
					Name:            "",
					Protocol:        DefaultDatabaseProtocol,
					Params:          "",
					MaxOpenConns:    DefaultDatabaseMaxOpenConns,
					MaxIdleConns:    DefaultDatabaseMaxIdleConns,
					ConnMaxLifetime: DefaultDatabaseConnMaxLifetime,
					ConnMaxIdleTime: DefaultDatabaseConnMaxIdleTime,
				},
				TracerConfig: &Tracer{
					DSN:        "",
					SampleRate: DefaultTracerSampleRate,
					FlushTime:  DefaultTracerFlushTime,
				},
				RouterConfig: &Router{
					Prefix:           "",
					Port:             DefaultRouterPort,
					ShowRegistered:   false,
					HtmlPath:         "",
					AllowOrigins:     nil,
					AllowMethods:     nil,
					AllowHeaders:     nil,
					AllowCredentials: false,
					DocPath:          "",
					MaxConnection:    DefaultRouterMaxConnection,
					MaxLimit:         DefaultRouterMaxLimit,
				},
				LoggerConfig: &Logger{
					Level: DefaultLoggerLevel,
				},
				CacheConfig: &Cache{
					Redis: &Redis{
						Addr:         "",
						DB:           DefaultCacheRedisDB,
						MaxRetries:   DefaultCacheRedisMaxRetries,
						MinIdleConns: DefaultCacheRedisMinIdleConns,
						PoolSize:     DefaultCacheRedisPoolSize,
						IdleTimeout:  DefaultCacheRedisIdleTimeout,
						ReadTimeout:  DefaultCacheRedisReadTimeout,
						WriteTimeout: DefaultCacheRedisWriteTimeout,
						DialTimeout:  DefaultCacheRedisDialTimeout,
						PoolTimeout:  DefaultCacheRedisPoolTimeout,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.SetTestConfig(tt.setVals)

			c := &Config{}
			c.build()

			if !reflect.DeepEqual(c.AppConfig, tt.expectedConf.AppConfig) {
				t.Errorf("AppConfig mismatch:\n got %+v\n want %+v", c.AppConfig, tt.expectedConf.AppConfig)
			}
			if !reflect.DeepEqual(c.CodeConfig, tt.expectedConf.CodeConfig) {
				t.Errorf("CodeConfig mismatch:\n got %+v\n want %+v", c.CodeConfig, tt.expectedConf.CodeConfig)
			}
			if !reflect.DeepEqual(c.DatabaseConfig, tt.expectedConf.DatabaseConfig) {
				t.Errorf("DatabaseConfig mismatch:\n got %+v\n want %+v", c.DatabaseConfig, tt.expectedConf.DatabaseConfig)
			}
			if !reflect.DeepEqual(c.TracerConfig, tt.expectedConf.TracerConfig) {
				t.Errorf("TracerConfig mismatch:\n got %+v\n want %+v", c.TracerConfig, tt.expectedConf.TracerConfig)
			}
			if !reflect.DeepEqual(c.RouterConfig, tt.expectedConf.RouterConfig) {
				t.Errorf("RouterConfig mismatch:\n got %+v\n want %+v", c.RouterConfig, tt.expectedConf.RouterConfig)
			}
			if !reflect.DeepEqual(c.LoggerConfig, tt.expectedConf.LoggerConfig) {
				t.Errorf("LoggerConfig mismatch:\n got %+v\n want %+v", c.LoggerConfig, tt.expectedConf.LoggerConfig)
			}
			if !reflect.DeepEqual(c.CacheConfig, tt.expectedConf.CacheConfig) {
				t.Errorf("CacheConfig mismatch:\n got %+v\n want %+v", c.CacheConfig, tt.expectedConf.CacheConfig)
			}
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

func TestConfig_getOrDefaultDuration(t *testing.T) {
	type args struct {
		key        string
		defaultVal time.Duration
		setVal     string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "value is empty",
			args: args{
				key:        "test",
				defaultVal: 5 * time.Second,
				setVal:     "",
			},
			want: 5 * time.Second,
		},
		{
			name: "value is invalid",
			args: args{
				key:        "test",
				defaultVal: 5 * time.Second,
				setVal:     "invalid",
			},
			want: 5 * time.Second,
		},
		{
			name: "value is valid",
			args: args{
				key:        "test",
				defaultVal: 5 * time.Second,
				setVal:     "10s",
			},
			want: 10 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.SetTestConfig(map[string]any{
				tt.args.key: tt.args.setVal,
			})

			c := &Config{}
			if got := c.getOrDefaultDuration(tt.args.key, tt.args.defaultVal); got != tt.want {
				t.Errorf("getOrDefaultDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_getOrDefaultInt(t *testing.T) {
	type args struct {
		key        string
		defaultVal int
		setVal     int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "value is 0",
			args: args{
				key:        "test",
				defaultVal: 10,
				setVal:     0,
			},
			want: 10,
		},
		{
			name: "value is valid",
			args: args{
				key:        "test",
				defaultVal: 10,
				setVal:     20,
			},
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.SetTestConfig(map[string]any{
				tt.args.key: tt.args.setVal,
			})

			c := &Config{}
			if got := c.getOrDefaultInt(tt.args.key, tt.args.defaultVal); got != tt.want {
				t.Errorf("getOrDefaultInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_getOrDefaultFloat(t *testing.T) {
	type args struct {
		key        string
		defaultVal float64
		setVal     float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "value is 0",
			args: args{
				key:        "test",
				defaultVal: 1.5,
				setVal:     0,
			},
			want: 1.5,
		},
		{
			name: "value is valid",
			args: args{
				key:        "test",
				defaultVal: 1.5,
				setVal:     2.5,
			},
			want: 2.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.SetTestConfig(map[string]any{
				tt.args.key: tt.args.setVal,
			})

			c := &Config{}
			if got := c.getOrDefaultFloat(tt.args.key, tt.args.defaultVal); got != tt.want {
				t.Errorf("getOrDefaultFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_getDBUri(t *testing.T) {
	type fields struct {
		DatabaseConfig *Database
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "without params",
			fields: fields{
				DatabaseConfig: &Database{
					User:     "user",
					Password: "password",
					Protocol: "tcp",
					Host:     "localhost",
					Port:     "3306",
					Name:     "db",
					Params:   "",
				},
			},
			want: "user:password@tcp(localhost:3306)/db?parseTime=true",
		},
		{
			name: "with params",
			fields: fields{
				DatabaseConfig: &Database{
					User:     "user",
					Password: "password",
					Protocol: "tcp",
					Host:     "localhost",
					Port:     "3306",
					Name:     "db",
					Params:   "charset=utf8",
				},
			},
			want: "user:password@tcp(localhost:3306)/db?parseTime=true&charset=utf8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				DatabaseConfig: tt.fields.DatabaseConfig,
			}
			if got := c.getDBUri(); got != tt.want {
				t.Errorf("Config.getDBUri() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeRoutePrefix(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty prefix",
			args: args{
				prefix: "",
			},
			want: "",
		},
		{
			name: "prefix with -svc",
			args: args{
				prefix: "user-svc",
			},
			want: "user",
		},
		{
			name: "prefix without -svc",
			args: args{
				prefix: "user",
			},
			want: "user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeRoutePrefix(tt.args.prefix); got != tt.want {
				t.Errorf("normalizeRoutePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
