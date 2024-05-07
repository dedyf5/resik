// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package main

import "github.com/dedyf5/resik/app/grpc/bootstrap"

var app *bootstrap.App
var cleanup func()
var err error

func main() {
	app, cleanup, err = bootstrap.InitializeHTTP()

	if err != nil {
		panic(err)
	}

	if app == nil {
		panic("Failed to initialize app")
	}

	app.Start()
	defer cleanup()
}
