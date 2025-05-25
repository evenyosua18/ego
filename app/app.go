package app

import (
	"github.com/evenyosua18/ego/config"
	"github.com/evenyosua18/ego/http"
)

type app struct {
	httpRouter http.IHttpRouter
}

func (a *app) RunRest() {
	// db connection

	// get router
	a.httpRouter = http.NewRouter()

	// listen
	if err := a.httpRouter.Listen(config.GetConfig().ServiceConfig.Port); err != nil {
		panic(err)
	}
}

func GetApp() IApp {
	return &app{}
}
