package source

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"crypto/rand"
	log "github.com/Sirupsen/logrus"
)

const (
	EventIngestor = iota
	MetricIngestor = iota
)

// NovaIngest creates a new ingest obj
type NovaIngest struct {
	Source  string
	Entity  string
	Auth    string
	NovaURL string
	ErrChan chan error
	Type    int
}

type novaEventFormat struct {
	Source string            `json:"source"`
	Entity string            `json:"entity"`
	Event  map[string]string `json:"event"`
}

func NewNovaIngestForEvents(novaURL, entity, auth string) *NovaIngest {
	return &NovaIngest{
		Source:  novaCLISourcePrefix + pseudoRandomID(),
		Entity:  entity,
		Auth:    auth,
		NovaURL: novaURL+eventsURLPath,
		ErrChan: make(chan error, 5),
		Type: EventIngestor,
	}
}

func NewNovaIngestForMetrics(novaURL, entity, auth string) *NovaIngest {
	return &NovaIngest{
		Source:  novaCLISourcePrefix + pseudoRandomID(),
		Entity:  entity,
		Auth:    auth,
		NovaURL: novaURL+ metricsURLIngestPath,
		ErrChan: make(chan error, 5),
		Type: MetricIngestor,
	}
}

// Start sends lines from stdin to nova
func (n *NovaIngest) Start(r io.Reader) {
	n.sendToNova(n.batchEvents(n.readIn(r)))
}

// WaitAndLogErrors blocks on the pipeline to complete and logs all errors
func (n *NovaIngest) WaitAndLogErrors() (errorsEncountered bool) {
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
	ticker := time.Tick(ingestionBufferTimeout)

	buffer := &bytes.Buffer{}
	writer := bufio.NewWriter(buffer)

	outChan = make(chan *bytes.Buffer, 10) // buffer at most 10 http requests

	go func() {
		defer close(outChan)
		for {
			if buffer.Len() > ingestionBufferSizeBytes {
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
				bytesArray, err := n.marshal(line)
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

func (n *NovaIngest) sendToNova(inChan chan *bytes.Buffer) {
	go func() {
		defer close(n.ErrChan)
		for buffer := range inChan {
			_, err := Post(n.NovaURL, buffer, n.Auth)
			if err != nil {
				n.ErrChan <- err
				return
			}
		}
	}()
	return
}

func (n *NovaIngest) marshal(line string) ([]byte, error) {
	if n.Type == EventIngestor {
		nEvent := novaEventFormat{Source: n.Source, Entity: n.Entity, Event: map[string]string{"raw": line}}
		return json.Marshal(nEvent)
	} else if n.Type == MetricIngestor {
		return []byte(line), nil
	} else {
		return nil, nil
	}
}

func pseudoRandomID() string {
	b := make([]byte, 7)
	_, err := rand.Read(b)
	if err != nil {
		return "00000000000000"
	}
	return fmt.Sprintf("%X", b)
}
