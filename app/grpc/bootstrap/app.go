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
}

func newApp(serverHTTP *ServerHTTP) (*App, func(), error) {
	app := &App{
		serverHTTP: serverHTTP,
	}
	return app, func() {
		serverHTTP.Close()
	}, nil
}

func (app *App) Start(c context.Context) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app.serverHTTP.Start()

	<-c.Done()

	log.Println("received shutdown signal, stopping application...")
}
