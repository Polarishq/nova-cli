// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

		m.GetAggergations(strings.Join(args, ","), "avg", "", "")

	},
}

func init() {
	metricCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("aggregations", "a", "", "stats aggregations to run on metrics (e.g. avg, min, max, etc.)")
	getCmd.MarkFlagRequired("aggregations")
	getCmd.Flags().StringP("group", "g", "", "group aggregations by dimensions")
	getCmd.Flags().StringP("span", "s", "", "group aggregations by time span (1m, 1s, 1d)")
}
