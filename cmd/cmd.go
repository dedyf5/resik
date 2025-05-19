// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package cmd

import (
	"os"

	"github.com/dedyf5/resik/config"
	"github.com/spf13/cobra"
)

var mainCMD = &cobra.Command{
	Short: config.AppLogoASCII,
}

func Execute() {
	if err := mainCMD.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	mainCMD.SetHelpTemplate(mainHelpTemplate)
	mainCMD.AddCommand(runGRPC)
	mainCMD.AddCommand(runREST)
	mainCMD.AddCommand(runMigrate)

	runMigrate.AddCommand(runMigrateUp)
	runMigrate.AddCommand(runMigrateDown)
	runMigrate.AddCommand(runMigrateVersion)
}
