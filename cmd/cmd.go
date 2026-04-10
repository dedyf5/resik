// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package cmd

import (
	"fmt"
	"os"

	"github.com/dedyf5/resik/buildinfo"
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/pkg/color"
	"github.com/spf13/cobra"
)

var mainCMD = &cobra.Command{
	Short: buildinfo.FrameworkLogoASCIIVersion,
}

func Execute() {
	if err := mainCMD.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	frameworkName := color.Format(color.GREEN, buildinfo.FrameworkName)

	app := config.GetApp()
	appName := color.Format(color.GREEN, app.Name())
	appVersion := color.Format(color.YELLOW, app.Version())

	versionTemplate := fmt.Sprintf(
		"%s Version: {{.Version}}\n%s Version: %s\nGit Commit: %s\nBuild Time: %s\n",
		frameworkName,
		appName,
		appVersion,
		buildinfo.AppGitCommit,
		buildinfo.AppBuildTime,
	)

	mainCMD.Version = color.Format(color.YELLOW, buildinfo.FrameworkVersion)
	mainCMD.SetVersionTemplate(versionTemplate)
	mainCMD.SetHelpTemplate(mainHelpTemplate)

	mainCMD.AddCommand(runGRPC)
	mainCMD.AddCommand(runREST)
	mainCMD.AddCommand(runMigrate)

	runMigrate.AddCommand(runMigrateUp)
	runMigrate.AddCommand(runMigrateDown)
	runMigrate.AddCommand(runMigrateVersion)
}
