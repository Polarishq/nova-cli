package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// metricCmd represents the metric command
var metricCmd = &cobra.Command{
	Use:   "metric",
	Short: "List, get, and put metrics",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.UsageString())
	},
}

func init() {
	rootCmd.AddCommand(metricCmd)
}
