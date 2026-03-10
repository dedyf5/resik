// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package grpc

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dedyf5/resik/app/grpc/bootstrap"
)

var app *bootstrap.App
var cleanup func()
var err error

func Run() {
	c, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, cleanup, err = bootstrap.InitializeHTTP(c)

	if err != nil {
		panic(err)
	}

	if app == nil {
		panic("Failed to initialize app")
	}

	app.Start(c)

	defer cleanup()
}
