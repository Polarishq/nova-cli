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

	log "github.com/Sirupsen/logrus"
	"os"
	"strings"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Create a metric",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		clientID, clientSecret, err := src.GetCredentials(NovaURL)
		if err != nil {
			log.Error(err)
			log.Infof("Please run `nova login`")
			os.Exit(1)
		}
		authHeader := src.GetBasicAuthHeader(clientID, clientSecret)

		hostname, _ := os.Hostname()

		tr := strings.NewReader("")

		novaIngest := src.NewNovaIngest(NovaURL, hostname, authHeader)
		novaIngest.Start(tr)
		errorsEncountered := novaIngest.WaitAndLogErrors()
		if errorsEncountered {
			os.Exit(1)
		}
	},
}

func init() {
	metricCmd.AddCommand(putCmd)
	putCmd.Flags().StringP("dimensions", "d", "", "Comma separated dimensions to tag the metric with (e.g. 'foo:bar,baz:qux')")
}
