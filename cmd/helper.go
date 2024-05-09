// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	configEntity "github.com/dedyf5/resik/entities/config"
	"github.com/spf13/cobra"
)

var runGRPC = &cobra.Command{
	Use:   "grpc",
	Short: "Run gRPC app",
	Run: func(cmd *cobra.Command, args []string) {
		appHandler(configEntity.ModuleGRPC)
	},
}

var runREST = &cobra.Command{
	Use:   "rest",
	Short: "Run REST app",
	Run: func(cmd *cobra.Command, args []string) {
		appHandler(configEntity.ModuleREST)
	},
}

func appHandler(module configEntity.Module) {
	dirName := module.DirectoryName()
	runArgs := []string{"run", "-tags", dirName, fmt.Sprintf("./app/%s", dirName)}
	command := exec.Command("go", runArgs...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		log.Fatalln("Error running app:", err)
		os.Exit(1)
	}
}

var runHelp = func(_ *cobra.Command, _ []string) {
	fmt.Print(AppLogoASCII + usage)
}

var usage = `
Usage:
  [command]

Available Commands:
  grpc 		To run gRPC app
  rest 		To run REST app

`
