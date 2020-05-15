
package cmd

import (
	"fmt"
	"github.com/deimosfr/jeedom-status/pkg"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s", pkg.GetCurrentVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}