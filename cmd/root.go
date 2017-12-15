package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/splunknova/nova-cli/src"
	"io"
	log "github.com/Sirupsen/logrus"
	"net/url"
)

var NovaURL string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nova",
	Short: "A convenient command line tool to pipe logs to splunknova.com and search them",
	Run: func(cmd *cobra.Command, args []string) {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 { // ingest mode
			log.Debugf("Will attempt to pipe data from stdin to splunknova")
			var tr io.Reader

			if tee, _ := cmd.Flags().GetBool("tee"); tee {
				tr = io.TeeReader(os.Stdin, os.Stdout)
			} else {
				tr = os.Stdin
			}

			clientID, clientSecret, err := src.GetCredentials(NovaURL)
			if err != nil {
				log.Error(err)
				log.Infof("Please run `nova login`")
				os.Exit(1)
			}
			authHeader := src.GetBasicAuthHeader(clientID, clientSecret)

			hostname, _ := os.Hostname()

			novaIngest := src.NewNovaIngest(NovaURL, hostname, authHeader)
			novaIngest.Start(tr)
			errorsEncountered := novaIngest.WaitAndLogErrors()
			if errorsEncountered {
				os.Exit(1)
			}
		} else {
			fmt.Println(cmd.Short)
			fmt.Println()
			fmt.Printf(cmd.UsageString())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&NovaURL, "novaurl", "https://api.splunknova.com", "point to a different splunknova URL (used for testing)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "print debug information")
	rootCmd.Flags().BoolP("tee", "t", false, "tee to stdout after sending data to splunknova. only valid when piping stdin into nova-cli")

}

func initConfig() {
	if v,_ := rootCmd.Flags().GetBool("verbose"); v {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	_, err := url.ParseRequestURI(NovaURL)
	if err != nil {
		log.Errorf("%s isn't a valid NovaURL", NovaURL)
		os.Exit(1)
	}
}