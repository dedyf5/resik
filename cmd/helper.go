// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package cmd

import (
	"fmt"

	"github.com/dedyf5/resik/app/grpc"
	"github.com/dedyf5/resik/app/rest"
	"github.com/dedyf5/resik/config"
	"github.com/spf13/cobra"
)

var runGRPC = &cobra.Command{
	Use:   "grpc",
	Short: "Run gRPC app",
	Run: func(cmd *cobra.Command, args []string) {
		grpc.Run()
	},
}

var runREST = &cobra.Command{
	Use:   "rest",
	Short: "Run REST app",
	Run: func(cmd *cobra.Command, args []string) {
		rest.Run()
	},
}

var runHelp = func(_ *cobra.Command, _ []string) {
	fmt.Print(config.AppLogoASCII + usage)
}

var usage = `
Usage:
  [command]

Available Commands:
  grpc 		To run gRPC app
  rest 		To run REST app

`
