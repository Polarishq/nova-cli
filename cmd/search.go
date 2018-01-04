package cmd

import (
	"github.com/spf13/cobra"
	"github.com/splunknova/nova-cli/source"
	"os"
	log "github.com/Sirupsen/logrus"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search splunknova for events",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clientID, clientSecret, err := source.GetCredentials(NovaURL)
		if err != nil {
			log.Error(err)
			log.Infof("Please run `nova login`")
			os.Exit(1)
		}
		authHeader := source.GetBasicAuthHeader(clientID, clientSecret)

		novaSearch := source.NewNovaSearch(NovaURL, authHeader)

		reportStr := ""
		count, _ := cmd.LocalFlags().GetBool("count")
		statsVal, _ := cmd.LocalFlags().GetString("stats")
		reportVal, _ := cmd.LocalFlags().GetString("report")
		transformsVal, _ := cmd.LocalFlags().GetString("transforms")

		if count {
			reportStr = "stats count"
		} else if statsVal != "" {
			reportStr = "stats " + statsVal
		} else {
			reportStr = reportVal
		}

		data := novaSearch.Search(args[0], transformsVal, reportStr)
		errorsEncountered := novaSearch.WaitAndLogErrors()
		if errorsEncountered {
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
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolP("count", "c", false, "shorthand for -r 'stats count', takes precedence over -s and -r")
	searchCmd.Flags().StringP("stats", "s", "", "shorthand for -r 'stats ...', takes precedence over -r")
	searchCmd.Flags().StringP("report", "r", "", "apply aggregations to the search results. e.g. -r 'stats avg(mb) perc90(mb)'")
	searchCmd.Flags().StringP("transforms", "t", "", "apply transformations to each matching event. e.g. -t 'eval mb = gb * 1024'")
}
