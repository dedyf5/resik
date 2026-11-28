// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package grpc provides the main entry point for initializing and running the
// gRPC service along with OS termination signal handling.
package grpc

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dedyf5/resik/app/grpc/bootstrap"
)

// app holds the initialized gRPC application instance.
var app *bootstrap.App

// cleanup holds the dependency cleanup function returned during bootstrap.
var cleanup func()

// err holds any error encountered while initializing the application.
var err error

// Run initializes the gRPC application, listens for OS interrupt and
// termination signals (SIGINT and SIGTERM), and starts the server until a
// shutdown signal is received. It panics if initialization fails or
// the application instance is nil.
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
