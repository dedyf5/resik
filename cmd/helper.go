// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package cmd

import (
	"github.com/dedyf5/resik/app/grpc"
	"github.com/dedyf5/resik/app/rest"
	"github.com/dedyf5/resik/cmd/migrator"
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

var runMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  `Commands to apply, rollback, or check database migrations.`,
}

var runMigrateUp = &cobra.Command{
	Use:   "up [steps]",
	Short: "Apply pending migrations",
	Long:  `Apply all pending migrations, or a specific number of steps if provided.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		steps := ""
		if len(args) > 0 {
			steps = args[0]
		}
		return migrator.RunUp(steps)
	},
}

var runMigrateDown = &cobra.Command{
	Use:   "down [steps]",
	Short: "Rollback applied migrations",
	Long:  `Rollback the last applied migration, or a specific number of steps if provided. Default is 1 step if no number is provided.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		steps := ""
		if len(args) > 0 {
			steps = args[0]
		}
		return migrator.RunDown(steps)
	},
}

var runMigrateVersion = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	Long:  `Display the current database migration version.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrator.RunVersion()
	},
}

var mainHelpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}Usage:
  {{.UseLine}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

Use "{{.CommandPath}} [command] --help" for more information about a command.
`
