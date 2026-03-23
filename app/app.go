package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/evenyosua18/ego/auth"
	"github.com/evenyosua18/ego/cache"
	"github.com/evenyosua18/ego/cache/redis_adapter"
	"github.com/evenyosua18/ego/code"
	"github.com/evenyosua18/ego/config"
	"github.com/evenyosua18/ego/http"
	"github.com/evenyosua18/ego/logger"
	"github.com/evenyosua18/ego/request"
	"github.com/evenyosua18/ego/sqldb"
	"github.com/evenyosua18/ego/tracer"
)

type app struct {
	httpRouter http.IHttpRouter
}

var (
	osExit         = os.Exit
	panicFunc      = func(v any) { panic(v) }
	loadCodesFunc  = code.LoadCodes
	runSentryFunc  = tracer.RunSentry
	dbConnectFunc  = sqldb.Connect
	redisNewFunc   = redis_adapter.NewRedisAdapter
	authManageFunc = auth.ManageAccessToken
	newRouterFunc  = func(cfg http.RouteConfig) http.IHttpRouter { return http.NewRouter(cfg) }
	signalNotify   = signal.Notify
)

func (a *app) RunRest() {
	// init config (local values)
	appConfig := Config{}
	appConfig.build()

	// initial remote config fetch
	remoteUrl := config.GetConfig().GetString(RemoteConfigUrl)
	providerName := config.GetConfig().GetString(RemoteConfigProviderName)
	var remoteProvider config.RemoteConfigProvider

	applyRemoteConfig := func(values map[string]any) {
		serviceName := config.GetConfig().GetString(ServiceName)
		if serviceName == "" {
			serviceName = DefaultServiceName
		}

		// If payload is nested by service name, apply only that service's config.
		if svcConfig, ok := values[serviceName].(map[string]any); ok {
			config.GetConfig().Merge(svcConfig)
		} else {
			// Fallback: apply the whole payload if it isn't nested by service.
			config.GetConfig().Merge(values)
		}
	}

	if remoteUrl != "" && providerName != "" {
		switch providerName {
		case "firebase":
			var creds []byte
			if credStr := config.GetConfig().GetString("feature_flag.credentials"); credStr != "" {
				creds = []byte(credStr)
			}
			remoteProvider = config.NewFirebaseRemoteConfig(creds)
		default:
			fmt.Printf("warning: unsupported remote config provider: %s\n", providerName)
		}

		if remoteProvider != nil {
			values, err := remoteProvider.Fetch(context.Background())
			if err != nil {
				fmt.Printf("warning: remote config fetch failed: %v\n", err)
			} else {
				applyRemoteConfig(values)
				// Rebuild appConfig so subsequent setup logic gets the latest overridden values
				appConfig.build()
			}
		}
	}

	// init breaker config
	if appConfig.BreakerConfig != nil {
		request.InitBreakerConfig(request.BreakerConfig{
			MaxRequests:         uint32(appConfig.BreakerConfig.MaxRequest),
			Interval:            appConfig.BreakerConfig.Interval,
			Timeout:             appConfig.BreakerConfig.Timeout,
			ConsecutiveFailures: uint32(appConfig.BreakerConfig.ConsecutiveFailures),
		})
	}

	// setup logger
	logger.SetLogger(logger.NewDefaultLogger(logger.ParseLevel(appConfig.LoggerConfig.Level)))

	// load custom codes
	if appConfig.CodeConfig != nil && appConfig.CodeConfig.Filename != "" {
		if err := loadCodesFunc(config.GetConfig().GetConfigPath() + "/" + appConfig.CodeConfig.Filename); err != nil {
			panicFunc(err)
		}
	}

	// tracer
	if appConfig.TracerConfig != nil && appConfig.TracerConfig.DSN != "" {
		flushFunction, err := runSentryFunc(tracer.Config{
			Dsn:             appConfig.TracerConfig.DSN,
			Env:             appConfig.AppConfig.Env,
			TraceSampleRate: appConfig.TracerConfig.SampleRate,
			FlushTime:       appConfig.TracerConfig.FlushTime,
		})
		if err != nil {
			panicFunc(err)
		}

		defer flushFunction(appConfig.TracerConfig.FlushTime)
	}

	// db connection
	if appConfig.DatabaseConfig != nil && appConfig.DatabaseConfig.Name != "" {
		// trying to connect database
		db, err := dbConnectFunc(appConfig.DatabaseConfig.Driver, appConfig.getDBUri(), &sqldb.Config{
			MaxOpenConns:    appConfig.DatabaseConfig.MaxOpenConns,
			MaxIdleConns:    appConfig.DatabaseConfig.MaxIdleConns,
			ConnMaxLifetime: appConfig.DatabaseConfig.ConnMaxLifetime,
			ConnMaxIdleTime: appConfig.DatabaseConfig.ConnMaxIdleTime,
		})
		if err != nil {
			panicFunc(err)
		}

		// set db connection
		sqldb.SetDB(db)
	}

	// cache connection
	if appConfig.CacheConfig.Redis != nil && appConfig.CacheConfig.Redis.Addr != "" {
		// trying to connect redis
		redis, err := redisNewFunc(redis_adapter.RedisConfig{
			Addr:         appConfig.CacheConfig.Redis.Addr,
			Password:     appConfig.CacheConfig.Redis.Password,
			DB:           appConfig.CacheConfig.Redis.DB,
			MaxRetries:   appConfig.CacheConfig.Redis.MaxRetries,
			MinIdleConns: appConfig.CacheConfig.Redis.MinIdleConns,
			PoolSize:     appConfig.CacheConfig.Redis.PoolSize,
			IdleTimeout:  appConfig.CacheConfig.Redis.IdleTimeout,
			ReadTimeout:  appConfig.CacheConfig.Redis.ReadTimeout,
			WriteTimeout: appConfig.CacheConfig.Redis.WriteTimeout,
			DialTimeout:  appConfig.CacheConfig.Redis.DialTimeout,
			PoolTimeout:  appConfig.CacheConfig.Redis.PoolTimeout,
		})
		if err != nil {
			fmt.Println(appConfig.CacheConfig.Redis.Addr)
			panicFunc(err)
		}

		// set to cache manager
		cache.InitCacheManager(redis)
	}

	// get router
	a.httpRouter = newRouterFunc(http.RouteConfig{
		MainPrefix:          appConfig.RouterConfig.Prefix,
		ShowRegisteredRoute: appConfig.RouterConfig.ShowRegistered,
		HtmlPath:            appConfig.RouterConfig.HtmlPath,
		DisableAuthChecker:  appConfig.RouterConfig.DisableAuthChecker,
		ReadTimeout:         appConfig.RouterConfig.ReadTimeout,
		WriteTimeout:        appConfig.RouterConfig.WriteTimeout,
		IdleTimeout:         appConfig.RouterConfig.IdleTimeout,
		CORS: http.CORSConfig{
			AllowOrigins:     appConfig.RouterConfig.AllowOrigins,
			AllowMethods:     appConfig.RouterConfig.AllowMethods,
			AllowHeaders:     appConfig.RouterConfig.AllowHeaders,
			AllowCredentials: appConfig.RouterConfig.AllowCredentials,
		},
		Doc: http.DocumentationConfig{
			Path: appConfig.RouterConfig.DocPath,
		},
		RateLimit: http.RateLimitConfig{
			MaxLimit: appConfig.RouterConfig.MaxLimit,
		},
	})

	// init background context
	bgCtx, cancelBg := context.WithCancel(context.Background())
	defer cancelBg()

	// init remote config auto refresh
	if remoteUrl != "" && remoteProvider != nil {
		period := config.GetConfig().GetDuration(RemoteConfigRefreshPeriod)
		config.AutoRefresh(bgCtx, remoteProvider, period, applyRemoteConfig)
	}

	// init token manager
	if appConfig.AuthSvcConfig.ClientId != "" && appConfig.AuthSvcConfig.ClientSecret != "" {
		if err := authManageFunc(bgCtx, appConfig.AuthSvcConfig.BaseUrl, appConfig.AuthSvcConfig.ClientId, appConfig.AuthSvcConfig.ClientSecret); err != nil {
			panicFunc(fmt.Errorf("failed allocating service access token on startup: %w", err))
		}
	}

	// graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signalNotify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		// listen
		if err := a.httpRouter.Listen(appConfig.RouterConfig.Port); err != nil {
			panicFunc(err)
		}
	}()

	<-quit
	logger.Info("Shutting down server...")

	// cancel background services (e.g., token manager)
	cancelBg()
	logger.Info(fmt.Sprintf("Active connections: %d", a.httpRouter.ActiveConnections()))

	ctx, cancel := context.WithTimeout(context.Background(), appConfig.AppConfig.ShutdownTimeout) // 5s timeout
	defer cancel()

	if err := a.httpRouter.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		osExit(0)
	} else {
		log.Println("Cleanup finished. Exiting...")
		osExit(0)
	}
}

func GetApp() IApp {
	return &app{}
}
