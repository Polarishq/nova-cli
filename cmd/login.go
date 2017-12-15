package cmd

import (
	"github.com/spf13/cobra"
	"github.com/splunknova/nova-cli/src"
	log "github.com/Sirupsen/logrus"
	"os"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "validate and save credentials to disk",
	Run: func(cmd *cobra.Command, args []string) {
		_, _, err := src.SaveCredentials(NovaURL)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
