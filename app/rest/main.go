// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package rest

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dedyf5/resik/app/rest/bootstrap"
)

var app *bootstrap.App
var cleanup func()
var err error

// @securityDefinitions.apikey	BearerAuth
// @in 							header
// @name 						Authorization
// @description					Type "Bearer " followed by a space and then your API token.
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
