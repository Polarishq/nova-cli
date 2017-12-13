package src

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

const maxBufferSize = 1000000 // server side max is 1,048,576
const maxBufferTime = 1 * time.Second
const urlDefaultHost = "https://api.splunknova.com"
const urlDefaultPath = "/v1/events"


// NovaIngest defines metadata sent to log-input
type NovaIngest struct {
	Source string
	Entity string
	Auth   string
	ErrChan chan error
}

type novaEvent struct {
	Source string            `json:"source"`
	Entity string            `json:"entity"`
	Event  map[string]string `json:"event"`
}

// NewNovaIngest defines metadata sent to log-input
func NewNovaIngest(source, entity, auth string) *NovaIngest {
	return &NovaIngest{
		Source: source,
		Entity: entity,
		Auth: auth,
		ErrChan: make(chan error, 5),
	}
}

// Start sends lines from stdin to nova
func (n *NovaIngest) Start(r io.Reader) () {
	n.sendToNova(n.batchEvents(n.readIn(r)))
}

// BlockedErrorLogger blocks on the pipeline to complete and logs all errors
func (n *NovaIngest) BlockedErrorLogger() (errorsEncountered bool) {
	for e := range n.ErrChan {
		errorsEncountered = true
		log.Error(e)
	}
	return
}

func (n *NovaIngest) readIn(r io.Reader) (outChan chan string) {
	outChan = make(chan string)

	go func() {
		defer close(outChan)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			outChan <- line
		}
		if err := scanner.Err(); err != nil {
			n.ErrChan <- err
			return
		}
	}()
	return
}

func (n *NovaIngest) batchEvents(inChan chan string) (outChan chan *bytes.Buffer) {
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
					n.ErrChan <- err
					return
				}
				writer.Write(bytesArray)
				writer.Flush() // for accurately calculating buffer.Len()
			}
		}
	}()
	return
}

func (n *NovaIngest) sendToNova(inChan chan *bytes.Buffer) () {
	httpClient := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   10 * time.Second,
	}

	go func() {
		defer close(n.ErrChan)
		for buffer := range inChan {
			req, err := http.NewRequest("POST", urlDefaultHost+urlDefaultPath, buffer)
			if err != nil {
				n.ErrChan <- err
				return
			}
			req.Header.Set("Authorization", "Basic "+n.Auth)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "nova-cli-0.3.0")
			resp, err := httpClient.Do(req)
			if err != nil {
				n.ErrChan <- err
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				n.ErrChan <- err
				return
			}
			if resp.StatusCode != 200 {
				n.ErrChan <- fmt.Errorf("error sending to splunknova. X-SPLUNK-REQ-ID=%+v",
					resp.Header.Get("X-SPLUNK-REQ-ID"))
				log.Warnf("%+v", string(body))
				return
			}
			log.Infof("Response: %+v", string(body))
		}
	}()
	return
}