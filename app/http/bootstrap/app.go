// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package bootstrap

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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

func (app *App) Start() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	terminalHandler := make(chan os.Signal, 1)
	signal.Notify(
		terminalHandler,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		in := <-terminalHandler
		log.Printf("SYSTEM CALL: %+v", in)
		cancel()
	}()

	app.router.routerSetup(app.serverHTTP)
	app.serverHTTP.Start()

	<-ctx.Done()
}
