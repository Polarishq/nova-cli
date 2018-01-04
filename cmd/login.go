package cmd

import (
	"github.com/spf13/cobra"
	"github.com/splunknova/nova-cli/source"
	log "github.com/Sirupsen/logrus"
	"os"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Validate and save credentials to disk",
	Run: func(cmd *cobra.Command, args []string) {
		_, _, err := source.SaveCredentials(NovaURL)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
