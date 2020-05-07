
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(GetCurrentVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func GetCurrentVersion() string {
	return "v0.1" // ci-version-check
}