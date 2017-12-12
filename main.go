package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/urfave/cli.v1"

	"github.com/Polarishq/cli-suite/src/config"
	"github.com/Polarishq/middleware/framework/log"
)

const maxBufferSize = 1000000 // server side max is 1,048,576
const maxBufferTime = 1 * time.Second
const urlDefaultHost = "https://api.splunknova.com"
const urlDefaultPath = "/v1/events"

func main() {
	app := cli.NewApp()
	app.Name = "nova"
	app.Usage = "Tee stdin to Splunk Nova. example: `cat hello.txt | nova` or `echo Splunk Nova | nova`"
	app.Version = "0.3.0"
	cli.VersionFlag = cli.BoolFlag{Name: "version"}
	app.Compiled = time.Now()
	app.Authors = []cli.Author{{Name: "Splunk Nova"}}
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "verbose, v"},
		cli.BoolFlag{Name: "debug, vv"},
		cli.BoolFlag{Name: "quiet, q"},
		cli.StringFlag{Name: "apiKeyID, ki", EnvVar: "NOVA_API_KEY_ID"},
		cli.StringFlag{Name: "apiKeySecret, ks", EnvVar: "NOVA_API_KEY_SECRET"},
	}
	app.Action = func(clic *cli.Context) error {
		if clic.Bool("debug") {
			log.SetDebug(true)
		} else if clic.Bool("verbose") {
			log.SetDebug(false)
		} else {
			log.SetError()
		}

		clientID := clic.String("apiKeyID")
		clientSecret := clic.String("apiKeySecret")
		if len(clientID) == 0 || len(clientSecret) == 0 {
			config.GetKeys()
			fmt.Fprint(os.Stderr, "Error: NOVA_API_KEY_ID or NOVA_API_KEY_SECRET either not set or passed in using -ki and -ks\n")
			os.Exit(1)
		}

		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			var tr io.Reader
			if clic.Bool("quiet") {
				tr = os.Stdin
			} else {
				tr = io.TeeReader(os.Stdin, os.Stdout)
			}

			hostname, _ := os.Hostname()
			auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
			i := Nova{"nova-cli", hostname, auth}
			doneChan := i.Start(tr)
			<-doneChan
		}
		return nil
	}
	app.Setup()
	app.Commands = nil
	app.Run(os.Args)
}

// Input defines metadata sent to log-input
type Nova struct {
	Source string
	Entity string
	Auth   string
}

type novaEvent struct {
	Source string            `json:"source"`
	Entity string            `json:"entity"`
	Event  map[string]string `json:"event"`
}

// Start sends lines from stdin to nova
func (n *Nova) Start(r io.Reader) (doneChan chan struct{}) {
	return n.sendToNova(n.batchEvents(n.readFromStdin(r)))
}

func (n *Nova) readFromStdin(r io.Reader) (outChan chan string) {
	outChan = make(chan string)

	go func() {
		defer close(outChan)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			outChan <- line
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}()
	return
}

func (n *Nova) batchEvents(inChan chan string) (outChan chan *bytes.Buffer) {
	ticker := time.Tick(maxBufferTime)

	buffer := &bytes.Buffer{}
	writer := bufio.NewWriter(buffer)

	outChan = make(chan *bytes.Buffer, 10) // buffer at most 10 http requests

	go func() {
		defer close(outChan)
		for {
			if buffer.Len() > maxBufferSize {
				outChan <- buffer
				buffer = &bytes.Buffer{}
				writer.Reset(buffer)
			}

			select {
			case <-ticker:
				outChan <- buffer
				buffer = &bytes.Buffer{}
				writer.Reset(buffer)
			default:
				line, ok := <-inChan
				if !ok {
					outChan <- buffer
					return
				}
				nEvent := novaEvent{Source: n.Source, Entity: n.Entity, Event: map[string]string{"raw": line}}
				bytesArray, err := json.Marshal(nEvent)
				if err != nil {
					panic(err)
				}
				writer.Write(bytesArray)
				writer.Flush() // for accurately calculating buffer.Len()
			}
		}
	}()
	return
}

func (n *Nova) sendToNova(inChan chan *bytes.Buffer) (doneChan chan struct{}) {
	httpClient := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   10 * time.Second,
	}

	doneChan = make(chan struct{})
	go func() {
		for buffer := range inChan {
			req, err := http.NewRequest("POST", urlDefaultHost+urlDefaultPath, buffer)
			if err != nil {
				panic(err)
			}
			req.Header.Set("Authorization", "Basic "+n.Auth)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "nova-cli-0.3.0")
			resp, err := httpClient.Do(req)
			if err != nil {
				panic(err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				panic(err)
			}
			fmt.Printf("Response: %+v", string(body))
		}
		doneChan <- struct{}{}
	}()
	return
}
