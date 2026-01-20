package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/evenyosua18/ego/code"
	"github.com/evenyosua18/ego/config"
	"github.com/evenyosua18/ego/http"
	"github.com/evenyosua18/ego/logger"
	"github.com/evenyosua18/ego/sqldb"
	"github.com/evenyosua18/ego/tracer"
)

type app struct {
	httpRouter http.IHttpRouter
}

func (a *app) RunRest() {
	// init config
	appConfig := Config{}
	appConfig.build()

	// setup logger
	logger.SetLogger(logger.NewDefaultLogger(logger.ParseLevel(appConfig.LoggerConfig.Level)))

	// load custom codes
	if appConfig.CodeConfig.Filename != "" {
		if err := code.LoadCodes(config.GetConfig().GetConfigPath() + "/" + appConfig.CodeConfig.Filename); err != nil {
			panic(err)
		}
	}

	// tracer
	if appConfig.TracerConfig.DSN != "" {
		flushFunction, err := tracer.RunSentry(tracer.Config{
			Dsn:             appConfig.TracerConfig.DSN,
			Env:             appConfig.AppConfig.Env,
			TraceSampleRate: appConfig.TracerConfig.SampleRate,
			FlushTime:       appConfig.TracerConfig.FlushTime,
		})
		if err != nil {
			panic(err)
		}

		defer flushFunction(appConfig.TracerConfig.FlushTime)
	}

	// db connection
	if appConfig.DatabaseConfig.Name != "" {
		// trying to connect database
		db, err := sqldb.Connect(appConfig.DatabaseConfig.Driver, appConfig.getDBUri(), &sqldb.Config{
			MaxOpenConns:    appConfig.DatabaseConfig.MaxOpenConns,
			MaxIdleConns:    appConfig.DatabaseConfig.MaxIdleConns,
			ConnMaxLifetime: appConfig.DatabaseConfig.ConnMaxLifetime,
			ConnMaxIdleTime: appConfig.DatabaseConfig.ConnMaxIdleTime,
		})
		if err != nil {
			panic(err)
		}

		// set db connection
		sqldb.SetDB(db)
	}

	// get router
	a.httpRouter = http.NewRouter(http.RouteConfig{
		MainPrefix:          appConfig.RouterConfig.Prefix,
		ShowRegisteredRoute: appConfig.RouterConfig.ShowRegistered,
		HtmlPath:            appConfig.RouterConfig.HtmlPath,
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

	// graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		// listen
		if err := a.httpRouter.Listen(appConfig.RouterConfig.Port); err != nil {
			panic(err)
		}
	}()

	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), appConfig.AppConfig.ShutdownTimeout) // 5s timeout
	defer cancel()

	if err := a.httpRouter.ShutdownWithContext(ctx); err != nil {
		panic(fmt.Sprintf("Server forced to shutdown: %v", err))
	}
}

func GetApp() IApp {
	return &app{}
}
