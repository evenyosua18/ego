package app

import (
	"github.com/evenyosua18/ego/code"
	"github.com/evenyosua18/ego/config"
	"github.com/evenyosua18/ego/http"
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
		MaxLimit:            appConfig.RouterConfig.MaxLimit,
		MainPrefix:          appConfig.RouterConfig.Prefix,
		ShowRegisteredRoute: appConfig.RouterConfig.ShowRegistered,
	})

	// listen
	if err := a.httpRouter.Listen(appConfig.RouterConfig.Port); err != nil {
		panic(err)
	}
}

func GetApp() IApp {
	return &app{}
}
