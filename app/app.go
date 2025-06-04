package app

import (
	"github.com/evenyosua18/ego/http"
	"github.com/evenyosua18/ego/sqldb"
)

type app struct {
	httpRouter http.IHttpRouter
}

func (a *app) RunRest() {
	// init config
	appConfig := Config{}
	appConfig.build()

	// load custom codes

	// db connection
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

	// get router
	a.httpRouter = http.NewRouter()

	// listen
	if err := a.httpRouter.Listen(appConfig.AppConfig.Port); err != nil {
		panic(err)
	}
}

func GetApp() IApp {
	return &app{}
}
