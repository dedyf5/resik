// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package rest provides the main entry point for initializing and running the
// REST HTTP service along with OS termination signal handling.
package rest

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dedyf5/resik/app/rest/bootstrap"
)

// app holds the initialized REST application instance.
var app *bootstrap.App

// cleanup holds the dependency cleanup function returned during bootstrap.
var cleanup func()

// err holds any error encountered while initializing the application.
var err error

// Run initializes the REST application, listens for OS interrupt and
// termination signals (SIGINT and SIGTERM), and starts the HTTP server until
// a shutdown signal is received. It panics if initialization fails or the
// application instance is nil.
//
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
