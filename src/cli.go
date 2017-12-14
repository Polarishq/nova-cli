package src

import (
	"io"
	"os"
	"time"

	"gopkg.in/urfave/cli.v1"

	log "github.com/Sirupsen/logrus"
	"net/url"
)

func NewCLI(searchKeywords string) *cli.App {
	app := cli.NewApp()
	app.Name = "nova-cli"
	app.Usage = "ingest and search logs using splunknova.com. " +
		"\n     ingest example: `tail -f /var/log/system.log | nova` or `echo hello world | nova`" +
		"\n     search example: `nova error` or `nova error -r 'stats count'`"
	app.Version = AppVersion
	cli.VersionFlag = cli.BoolFlag{Name: "version"}
	app.Compiled = time.Now()
	app.Authors = []cli.Author{{Name: "join us on slack: community.splunknova.com"}}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "login",
			Usage: "validate and save credentials to disk",
		},
		cli.BoolFlag{
			Name:  "count, c",
			Usage: "shorthand for -r 'stats count', takes precedence over -s and -r",
		},
		cli.StringFlag{
			Name:  "stats, s",
			Usage: "shorthand for -r 'stats .....', takes precedence over -r",
		},
		cli.StringFlag{
			Name:  "report, r",
			Usage: "apply aggregations to the search results. e.g. -r 'stats avg(mb) perc90(mb)'",
		},
		cli.StringFlag{
			Name:  "transforms, t",
			Usage: "apply transformations to each matching event. e.g. -t 'eval mb = gb * 1024'",
		},
		cli.BoolFlag{
			Name:  "tee",
			Usage: "tee to stdout after sending data to splunknova. only valid when piping stdin into nova-cli",
		},
		cli.StringFlag{
			Name:  "novaURL",
			Value: defaultNovaURL,
			Usage: "point to a different splunknova URL (used for testing)",
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "turn on debug information",
			EnvVar: "NOVA_CLI_VERBOSE",
		},
	}

	app.Action = func(cliContext *cli.Context) error {
		var err error

		if cliContext.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}

		novaUrl := cliContext.String("novaURL")
		_, err = url.ParseRequestURI(novaUrl)
		if err != nil {
			log.Errorf("novaURL=%s isn't a valid URL", novaUrl)
			os.Exit(1)
		}

		var clientID, clientSecret string

		if cliContext.Bool("login") {
			clientID, clientSecret, err = SaveCredentials(novaUrl)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			} else {
				log.Infof("Login succeeded, keys saved to %s", getConfigFilePath())
				os.Exit(0)
			}
		} else {
			clientID, clientSecret, err = GetCredentials(novaUrl)
			if err != nil {
				log.Error(err)
				log.Infof("Please run `nova --login`")
				os.Exit(1)
			}
		}

		authHeader := GetBasicAuthHeader(clientID, clientSecret)

		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 { // ingest mode
			var tr io.Reader
			if cliContext.Bool("tee") {
				tr = io.TeeReader(os.Stdin, os.Stdout)
			} else {
				tr = os.Stdin
			}

			hostname, _ := os.Hostname()

			novaIngest := NewNovaIngest(novaUrl, hostname, authHeader)
			novaIngest.Start(tr)
			errorsEncountered := novaIngest.WaitAndLogErrors()
			if errorsEncountered {
				os.Exit(1)
			}
		} else if searchKeywords == "" {
			cli.ShowAppHelp(cliContext)
		} else { // search mode
			var reportStr string
			if cliContext.Bool("count") {
				reportStr = "stats count"
			} else if cliContext.String("stats") != "" {
				reportStr = "stats " + cliContext.String("stats")
			} else {
				reportStr = cliContext.String("report")
			}
			novaSearch := NewNovaSearch(novaUrl, authHeader)
			novaSearch.Search(searchKeywords, cliContext.String("transforms"), reportStr)
			errorsEncountered := novaSearch.WaitAndLogErrors()
			if errorsEncountered {
				os.Exit(1)
			}
		}
		return nil
	}
	app.Setup()
	app.Commands = nil
	return app
}
