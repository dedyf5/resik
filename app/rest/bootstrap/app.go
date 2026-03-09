// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package bootstrap

import (
	"context"
	"log"
	"runtime"
)

type App struct {
	serverHTTP *ServerHTTP
	router     *Router
}

func newApp(
	serverHTTP *ServerHTTP,
	router *Router,
) (*App, func(), error) {
	app := App{
		serverHTTP: serverHTTP,
		router:     router,
	}
	return &app, func() {
		if err := serverHTTP.Close(); err != nil {
			log.Printf("ERROR NewApp: %s", err.Error())
		}
	}, nil
}

func (app *App) Start(c context.Context) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app.router.routerSetup(app.serverHTTP)
	app.serverHTTP.Start(c)

	<-c.Done()

	log.Println("received shutdown signal, stopping application...")
}
