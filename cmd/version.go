package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/splunknova/nova-cli/source"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of nova-cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(source.GetUserAgent())
	},
}