package app

import (
	"github.com/evenyosua18/ego/config"
	"github.com/evenyosua18/ego/router"
)

type app struct {
	router router.IRouter
}

func (a *app) RunRest() {
	// db connection

	// get router
	a.router = router.NewRouter()

	// listen
	if err := a.router.Listen(config.GetConfig().ServiceConfig.Port); err != nil {
		panic(err)
	}
}

func GetApp() IApp {
	return &app{}
}
