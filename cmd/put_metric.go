package cmd

import (
	"github.com/spf13/cobra"
	"github.com/splunknova/nova-cli/source"

	"os"
	"strings"
	"encoding/json"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Create a metric, e.g. (nova metric put cpu.usage 10 -d 'region:us-east-1')",
	Args: cobra.ExactArgs(2),
	PreRun: Authorize,
	Run: func(cmd *cobra.Command, args []string) {
		dims, _ := cmd.Flags().GetString("dimensions")

		// this is so hacky :-(
		metricBody := map[string]string{"metric_name": args[0], "_value": args[1], "entity": Hostname, "source": "nova-cli"}
		for k, v := range splitDims(dims) {
			metricBody[k] = v
		}
		b, _ := json.Marshal(map[string]interface{}{"fields": metricBody})
		tr := strings.NewReader(string(b))

		novaIngest := source.NewNovaIngestForMetrics(NovaURL, Hostname, AuthHeader)
		novaIngest.Start(tr)
		errorsEncountered := novaIngest.WaitAndLogErrors()
		if errorsEncountered {
			os.Exit(1)
		}
	},
}

func splitDims(dims string) map[string]string {
	splitDims := map[string]string{}
	dimArray := strings.Split(dims, ",")
	for _, dim := range dimArray {
		kv := strings.Split(dim, ":")
		if len(kv) != 2 {
			continue
		}
		splitDims[kv[0]] = kv[1]
	}
	return splitDims
}

func init() {
	metricCmd.AddCommand(putCmd)
	putCmd.Flags().StringP("dimensions", "d", "", "Comma separated dimensions to tag the metric with (e.g. 'foo:bar,baz:qux')")
}
