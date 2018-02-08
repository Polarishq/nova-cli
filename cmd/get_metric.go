package cmd

import (
	"github.com/spf13/cobra"
	"github.com/splunknova/nova-cli/source"
	"strings"
	"os"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get metrics",
	Args: cobra.RangeArgs(1, 3),
	PreRun: Authorize,
	Run: func(cmd *cobra.Command, args []string) {
		m := source.NewNovaMetricsSearch(NovaURL, AuthHeader)

		stats, _ := cmd.Flags().GetString("stats")
		span, _ := cmd.Flags().GetString("span")
		group, _ := cmd.Flags().GetString("group")

		data, err := m.GetStats(args, strings.Split(stats, ","), strings.Split(group, ","), span)

		if err != nil {
			os.Exit(1)
		}
		if table, _ := rootCmd.Flags().GetBool("table"); table {
			data.PrintTable()
		} else {
			data.PrintList()
		}
	},
}

func init() {
	metricCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("stats", "s", "", "stats to run on metrics (e.g. avg, min, max, etc.)")
	getCmd.MarkFlagRequired("stats")
	getCmd.Flags().StringP("group", "g", "", "group stats by dimensions")
	getCmd.Flags().String("span",  "", "group stats by time span (1m, 1s, 1d)")
}
