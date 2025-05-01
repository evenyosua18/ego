package app

import (
	"github.com/evenyosua18/ego/config"
	"github.com/evenyosua18/ego/router"
)

type app struct {
	router *router.Router
	cfg    *Config
}

func (a *app) RunRest() {
	// db connection

	// get router
	a.router = router.NewRouter()

	// listen
	a.router.Listen(a.cfg.ServiceConfig.Port)
}

func GetApp() IApp {
	//  init config
	config.GetConfig()

	return &app{cfg: BuildConfiguration()}
}
