package app

import (
	"github.com/evenyosua18/ego/http"
)

type app struct {
	httpRouter http.IHttpRouter
}

func (a *app) RunRest() {
	// init config
	appConfig := Config{}
	appConfig.build()

	// db connection

	// get router
	a.httpRouter = http.NewRouter()

	// listen
	if err := a.httpRouter.Listen(appConfig.ServiceConfig.Port); err != nil {
		panic(err)
	}
}

func GetApp() IApp {
	return &app{}
}
