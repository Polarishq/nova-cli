package src

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/urfave/cli.v1"

	"github.com/Polarishq/cli-suite/src/config"
	log "github.com/Sirupsen/logrus"
	"crypto/rand"
)

const AppVersion = "0.3.0"

func Foo(searchKeywords string) *cli.App {
	app := cli.NewApp()
	app.Name = "nova"
	app.Usage = "Tee stdin to Splunk NovaIngest. example: `cat hello.txt | nova` or `echo Splunk NovaIngest | nova`"
	app.Version = AppVersion
	cli.VersionFlag = cli.BoolFlag{Name: "version"}
	app.Compiled = time.Now()
	app.Authors = []cli.Author{{Name: "splunknova.com"}}
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "verbose, v"},
		cli.BoolFlag{Name: "tee"},
		cli.StringFlag{Name: "stats, s"},
		cli.StringFlag{Name: "transforms, t"},
		cli.StringFlag{Name: "report, r"},
		cli.StringFlag{Name: "apiKeyID, ki", EnvVar: "NOVA_API_KEY_ID"},
		cli.StringFlag{Name: "apiKeySecret, ks", EnvVar: "NOVA_API_KEY_SECRET"},
	}

	app.Action = func(clic *cli.Context) error {
		if clic.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}

		clientID := clic.String("apiKeyID")
		clientSecret := clic.String("apiKeySecret")
		if len(clientID) == 0 || len(clientSecret) == 0 {
			config.GetKeys()
			fmt.Fprint(os.Stderr, "Error: NOVA_API_KEY_ID or NOVA_API_KEY_SECRET either not set or passed in using -ki and -ks\n")
			os.Exit(1)
		}

		stat, _ := os.Stdin.Stat()

		if (stat.Mode() & os.ModeCharDevice) == 0 { // ingest mode
			var tr io.Reader
			if clic.Bool("tee") {
				tr = io.TeeReader(os.Stdin, os.Stdout)
			} else {
				tr = os.Stdin
			}

			hostname, _ := os.Hostname()
			source := "nova-cli-" + pseudoRandomID()
			auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))

			novaIngest := NewNovaIngest(source, hostname, auth)
			novaIngest.Start(tr)
			errorsEncountered := novaIngest.BlockedErrorLogger()
			if errorsEncountered {
				os.Exit(1)
			}
		} else { // search mode
			log.Infof("Searching keywords='%+v'\n", searchKeywords)
			log.Infof("Searching transforms='%+v'\n", clic.String("transforms"))
			log.Infof("Searching report='%+v'\n", clic.String("report"))

			novaSearch := NovaSearch{Auth: "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))}
			novaSearch.Search(searchKeywords, clic.String("transforms"), clic.String("report"))
		}
		return nil
	}
	app.Setup()
	app.Commands = nil
	return app
}

func pseudoRandomID() (string) {
	b := make([]byte, 7)
	_, err := rand.Read(b)
	if err != nil {
		return "00000000000000"
	}
	return fmt.Sprintf("%X", b)
}
