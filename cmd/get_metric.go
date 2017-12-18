package cmd

import (
	"github.com/spf13/cobra"
	"github.com/splunknova/nova-cli/src"
	"strings"
	log "github.com/Sirupsen/logrus"
	"os"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get metrics",
	Args: cobra.RangeArgs(1, 3),
	Run: func(cmd *cobra.Command, args []string) {
		clientID, clientSecret, err := src.GetCredentials(NovaURL)
		if err != nil {
			log.Error(err)
			log.Infof("Please run `nova login`")
			os.Exit(1)
		}
		authHeader := src.GetBasicAuthHeader(clientID, clientSecret)

		m := src.NewNovaMetricsSearch(NovaURL, authHeader)

		aggregations, _ := cmd.Flags().GetString("aggregations")
		span, _ := cmd.Flags().GetString("span")
		group, _ := cmd.Flags().GetString("group")

		data, err := m.GetAggregations(strings.Join(args, ","), aggregations, group, span)
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
	getCmd.Flags().StringP("aggregations", "a", "", "stats aggregations to run on metrics (e.g. avg, min, max, etc.)")
	getCmd.MarkFlagRequired("aggregations")
	getCmd.Flags().StringP("group", "g", "", "group aggregations by dimensions")
	getCmd.Flags().StringP("span", "s", "", "group aggregations by time span (1m, 1s, 1d)")
}
