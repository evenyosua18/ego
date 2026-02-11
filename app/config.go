package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/evenyosua18/ego/config"
)

const (
	LocalEnv = "local" // used for validation: db password can be empty if local env

	ServiceName            = "service.name"
	ServiceEnv             = "service.env"
	ServiceShutdownTimeout = "service.shutdown_timeout"

	DefaultServiceName            = "temporary-service"
	DefaultServiceEnv             = "local"
	DefaultServiceShutdownTimeout = 30 * time.Second

	CustomCodeFilePath = "code.filename"

	DatabaseName            = "database.name"
	DatabasePort            = "database.port"
	DatabaseDriver          = "database.driver"
	DatabasePassword        = "database.password"
	DatabaseHost            = "database.host"
	DatabaseUser            = "database.user"
	DatabaseProtocol        = "database.protocol"
	DatabaseParams          = "database.params"
	DatabaseMaxOpenConns    = "database.max_open_conns"
	DatabaseMaxIdleConns    = "database.max_idle_conns"
	DatabaseConnMaxLifetime = "database.conn_max_lifetime"
	DatabaseConnMaxIdleTime = "database.conn_max_idle_time"

	DefaultDatabasePort            = "3306"
	DefaultDatabaseDriver          = "mysql"
	DefaultDatabaseHost            = "localhost"
	DefaultDatabaseUser            = "root"
	DefaultDatabaseProtocol        = "tcp"
	DefaultDatabaseMaxOpenConns    = 100
	DefaultDatabaseMaxIdleConns    = 20
	DefaultDatabaseConnMaxLifetime = 30 * time.Minute
	DefaultDatabaseConnMaxIdleTime = 5 * time.Minute

	TracerDSN        = "tracer.dsn"
	TracerSampleRate = "tracer.sample_rate"
	TracerFlushTime  = "tracer.flush_time"

	DefaultTracerSampleRate = 1.0
	DefaultTracerFlushTime  = "1"

	RouterPrefix           = "router.prefix"
	RouterPort             = "router.port"
	RouterShowRegistered   = "router.show_registered"
	RouterHtmlPath         = "router.html_path"
	RouterAllowOrigins     = "router.allow_origins"
	RouterAllowMethods     = "router.allow_methods"
	RouterAllowHeaders     = "router.allow_headers"
	RouterAllowCredentials = "router.allow_credentials"
	RouterDocPath          = "router.doc_path"
	RouterMaxConnection    = "router.max_connection"
	RouterMaxLimit         = "router.rate_limit.max_limit"

	DefaultRouterPort          = ":8080"
	DefaultRouterMaxLimit      = 100
	DefaultRouterMaxConnection = 5000

	LoggerLevel        = "logger.level"
	DefaultLoggerLevel = "info"

	CacheRedisAddr         = "cache.redis.addr"
	CacheRedisPassword     = "cache.redis.password"
	CacheRedisDB           = "cache.redis.db"
	CacheRedisMaxRetries   = "cache.redis.max_retries"
	CacheRedisMinIdleConns = "cache.redis.min_idle_conns"
	CacheRedisPoolSize     = "cache.redis.pool_size"
	CacheRedisIdleTimeout  = "cache.redis.idle_timeout"
	CacheRedisReadTimeout  = "cache.redis.read_timeout"
	CacheRedisWriteTimeout = "cache.redis.write_timeout"
	CacheRedisDialTimeout  = "cache.redis.dial_timeout"
	CacheRedisPoolTimeout  = "cache.redis.pool_timeout"

	DefaultCacheRedisDB           = 0
	DefaultCacheRedisMaxRetries   = 3
	DefaultCacheRedisMinIdleConns = 10
	DefaultCacheRedisPoolSize     = 10
	DefaultCacheRedisIdleTimeout  = 5 * time.Minute
	DefaultCacheRedisReadTimeout  = 5 * time.Second
	DefaultCacheRedisWriteTimeout = 5 * time.Second
	DefaultCacheRedisDialTimeout  = 5 * time.Second
	DefaultCacheRedisPoolTimeout  = 5 * time.Second
)

var ErrEmptyDBPassword = fmt.Errorf("db password can't be empty for non local environment")

type (
	Config struct {
		AppConfig      *App
		CodeConfig     *Code
		DatabaseConfig *Database
		TracerConfig   *Tracer
		RouterConfig   *Router
		LoggerConfig   *Logger
		CacheConfig    *Cache
	}

	Code struct {
		Filename string
	}

	App struct {
		Name            string
		Port            string
		Env             string
		ShutdownTimeout time.Duration
	}

	Logger struct {
		Level string
	}

	Database struct {
		Driver   string
		User     string
		Password string
		Host     string
		Port     string
		Name     string
		Protocol string
		Params   string

		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime time.Duration
		ConnMaxIdleTime time.Duration
	}

	Tracer struct {
		DSN        string
		SampleRate float64
		FlushTime  string
	}

	Router struct {
		Prefix         string
		Port           string
		ShowRegistered bool
		HtmlPath       string
		MaxConnection  int

		// CORS
		AllowOrigins     []string
		AllowMethods     []string
		AllowHeaders     []string
		AllowCredentials bool

		// docs
		DocPath string

		// rate limit
		MaxLimit int
	}

	Cache struct {
		Redis *Redis
	}

	Redis struct {
		Addr         string
		Password     string
		DB           int
		MaxRetries   int
		MinIdleConns int
		PoolSize     int
		IdleTimeout  time.Duration
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		DialTimeout  time.Duration
		PoolTimeout  time.Duration
	}
)

func (c *Config) build() {
	// setup service configuration
	c.AppConfig = &App{
		Name:            c.getOrDefault(ServiceName, DefaultServiceName),
		Env:             c.getOrDefault(ServiceEnv, DefaultServiceEnv),
		ShutdownTimeout: c.getOrDefaultDuration(ServiceShutdownTimeout, DefaultServiceShutdownTimeout),
	}

	// setup logger configuration
	c.LoggerConfig = &Logger{
		Level: c.getOrDefault(LoggerLevel, DefaultLoggerLevel),
	}

	// setup code configuration
	c.CodeConfig = &Code{
		Filename: c.getOrDefault(CustomCodeFilePath, ""),
	}

	// setup db configuration
	c.DatabaseConfig = &Database{
		Driver:          c.getOrDefault(DatabaseDriver, DefaultDatabaseDriver),
		User:            c.getOrDefault(DatabaseUser, DefaultDatabaseUser),
		Password:        c.getOrDefault(DatabasePassword, ""),
		Host:            c.getOrDefault(DatabaseHost, DefaultDatabaseHost),
		Port:            c.getOrDefault(DatabasePort, DefaultDatabasePort),
		Name:            c.getOrDefault(DatabaseName, ""),
		Protocol:        c.getOrDefault(DatabaseProtocol, DefaultDatabaseProtocol),
		Params:          c.getOrDefault(DatabaseParams, ""),
		MaxOpenConns:    c.getOrDefaultInt(DatabaseMaxOpenConns, DefaultDatabaseMaxOpenConns),
		MaxIdleConns:    c.getOrDefaultInt(DatabaseMaxIdleConns, DefaultDatabaseMaxIdleConns),
		ConnMaxLifetime: c.getOrDefaultDuration(DatabaseConnMaxLifetime, DefaultDatabaseConnMaxLifetime),
		ConnMaxIdleTime: c.getOrDefaultDuration(DatabaseConnMaxIdleTime, DefaultDatabaseConnMaxIdleTime),
	}

	if c.AppConfig.Env != LocalEnv && c.DatabaseConfig.Password == "" {
		// validate db password can't be empty if not local env
		panic(ErrEmptyDBPassword)
	}

	// tracer
	c.TracerConfig = &Tracer{
		DSN:        c.getOrDefault(TracerDSN, ""),
		SampleRate: c.getOrDefaultFloat(TracerSampleRate, DefaultTracerSampleRate),
		FlushTime:  c.getOrDefault(TracerFlushTime, DefaultTracerFlushTime),
	}

	// route
	c.RouterConfig = &Router{
		MaxLimit:         c.getOrDefaultInt(RouterMaxLimit, DefaultRouterMaxLimit),
		Prefix:           normalizeRoutePrefix(c.getOrDefault(RouterPrefix, "")),
		Port:             normalizePort(c.getOrDefault(RouterPort, DefaultRouterPort)),
		ShowRegistered:   config.GetConfig().GetBool(RouterShowRegistered), // if not true or 1, will return false, no need to set default value anymore
		HtmlPath:         c.getOrDefault(RouterHtmlPath, ""),
		AllowOrigins:     config.GetConfig().GetStringSlice(RouterAllowOrigins),
		AllowMethods:     config.GetConfig().GetStringSlice(RouterAllowMethods),
		AllowHeaders:     config.GetConfig().GetStringSlice(RouterAllowHeaders),
		AllowCredentials: config.GetConfig().GetBool(RouterAllowCredentials),
		DocPath:          c.getOrDefault(RouterDocPath, ""),
		MaxConnection:    c.getOrDefaultInt(RouterMaxConnection, DefaultRouterMaxConnection),
	}

	// cache redis
	c.CacheConfig = &Cache{
		Redis: &Redis{
			Addr:         c.getOrDefault(CacheRedisAddr, ""),
			Password:     c.getOrDefault(CacheRedisPassword, ""),
			DB:           c.getOrDefaultInt(CacheRedisDB, DefaultCacheRedisDB),
			MaxRetries:   c.getOrDefaultInt(CacheRedisMaxRetries, DefaultCacheRedisMaxRetries),
			MinIdleConns: c.getOrDefaultInt(CacheRedisMinIdleConns, DefaultCacheRedisMinIdleConns),
			PoolSize:     c.getOrDefaultInt(CacheRedisPoolSize, DefaultCacheRedisPoolSize),
			IdleTimeout:  c.getOrDefaultDuration(CacheRedisIdleTimeout, DefaultCacheRedisIdleTimeout),
			ReadTimeout:  c.getOrDefaultDuration(CacheRedisReadTimeout, DefaultCacheRedisReadTimeout),
			WriteTimeout: c.getOrDefaultDuration(CacheRedisWriteTimeout, DefaultCacheRedisWriteTimeout),
			DialTimeout:  c.getOrDefaultDuration(CacheRedisDialTimeout, DefaultCacheRedisDialTimeout),
			PoolTimeout:  c.getOrDefaultDuration(CacheRedisPoolTimeout, DefaultCacheRedisPoolTimeout),
		},
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

func (c *Config) getOrDefaultDuration(key string, defaultVal time.Duration) time.Duration {
	// set value
	val := config.GetConfig().GetString(key)
	if val == "" {
		return defaultVal
	}

	// convert to duration
	valDuration, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}

	return valDuration
}

func (c *Config) getOrDefaultInt(key string, defaultVal int) int {
	// set value
	val := config.GetConfig().GetInt(key)
	if val == 0 {
		return defaultVal
	}

	return val
}

func (c *Config) getOrDefaultFloat(key string, defaultVal float64) float64 {
	val := config.GetConfig().GetFloat64(key)

	if val == 0 {
		return defaultVal
	}

	return val
}

func (c *Config) getDBUri() string {
	dbUri := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true", c.DatabaseConfig.User, c.DatabaseConfig.Password, c.DatabaseConfig.Protocol, c.DatabaseConfig.Host, c.DatabaseConfig.Port, c.DatabaseConfig.Name)

	if c.DatabaseConfig.Params != "" {
		return dbUri + "&" + c.DatabaseConfig.Params
	}

	return dbUri
}

func normalizeRoutePrefix(prefix string) string {
	if prefix == "" {
		return ""
	}

	return fmt.Sprintf("%s", strings.ReplaceAll(prefix, "-svc", ""))
}
