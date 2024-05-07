// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const AppLogoASCII string = `
 ____           _ _    
|  _ \ ___  ___(_) | __
| |_) / _ \/ __| | |/ /
|  _ <  __/\__ \ |   < 
|_| \_\___||___/_|_|\_\

`

var mainCMD = &cobra.Command{
	Short: AppLogoASCII,
}

func Execute() {
	if err := mainCMD.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	mainCMD.SetHelpFunc(runHelp)
	mainCMD.AddCommand(runGRPC)
	mainCMD.AddCommand(runREST)
}
