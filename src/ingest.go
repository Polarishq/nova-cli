package src

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"crypto/rand"
	log "github.com/Sirupsen/logrus"
	"compress/gzip"
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
	Gzip    bool
	Type    int
	TotalRequests int
}

type novaEventFormat struct {
	Source string            `json:"source"`
	Entity string            `json:"entity"`
	Event  map[string]string `json:"event"`
}

func NewNovaIngestForEvents(novaURL, entity, auth string, gzip bool) *NovaIngest {
	return &NovaIngest{
		Source:  novaCLISourcePrefix + pseudoRandomID(),
		Entity:  entity,
		Auth:    auth,
		NovaURL: novaURL+eventsURLPath,
		ErrChan: make(chan error, 5),
		Gzip: gzip,
		Type: EventIngestor,
	}
}

func NewNovaIngestForMetrics(novaURL, entity, auth string, gzip bool) *NovaIngest {
	return &NovaIngest{
		Source:  novaCLISourcePrefix + pseudoRandomID(),
		Entity:  entity,
		Auth:    auth,
		NovaURL: novaURL+ metricsURLIngestPath,
		ErrChan: make(chan error, 5),
		Gzip: gzip,
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
	log.Debug("Total Ingest Requests Sent: ", n.TotalRequests)
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

type MyBuffer struct {
	Buffer *bytes.Buffer
	Writer io.Writer
	Gzip   bool
}

func NewMB(gz bool) MyBuffer {
	m := MyBuffer{Gzip: gz, Buffer: &bytes.Buffer{}}
	if gz {
		m.Writer = gzip.NewWriter(m.Buffer)
	} else {
		m.Writer = bufio.NewWriter(m.Buffer)
	}
	return m
}

func (m *MyBuffer) Finalize() (outBytes *bytes.Buffer) {
	outBytes = m.Buffer
	if m.Gzip {
		m.Writer.(*gzip.Writer).Flush()
		m.Writer.(*gzip.Writer).Close()

		m.Buffer = &bytes.Buffer{}
		m.Writer = gzip.NewWriter(m.Buffer)
	} else {
		m.Writer.(*bufio.Writer).Flush()

		m.Buffer = &bytes.Buffer{}
		m.Writer = bufio.NewWriter(m.Buffer)
	}
	return
}

func (m *MyBuffer) AddData(data []byte) {
	m.Writer.Write(data)
}

func (m *MyBuffer) Len() int {
	// todo: is this needed?
	//if m.Gzip {
	//	m.Writer.(*gzip.Writer).Flush()
	//} else {
	//	m.Writer.(*bufio.Writer).Flush()
	//}
	return m.Buffer.Len()
}


func (n *NovaIngest) batchEvents(inChan chan string) (outChan chan *bytes.Buffer) {
	ticker := time.Tick(ingestionBufferTimeout)

	mb := NewMB(n.Gzip)

	outChan = make(chan *bytes.Buffer, 10) // buffer at most 10 http requests

	go func() {
		defer close(outChan)
		for {
			if mb.Len() > ingestionBufferSizeBytes {
				outChan <- mb.Finalize()
				continue
			}
			select {
			case <-ticker:
				outChan <- mb.Finalize()
				log.Info("1")
				continue
			default:
				//log.Info("2")
				line, ok := <-inChan
				if !ok {
					outChan <- mb.Finalize()
					return
				}
				bytesArray, err := n.marshal(line)
				if err != nil {
					n.ErrChan <- err
					return
				}
				mb.AddData(bytesArray)
			}
		}
	}()
	return
}


func (n *NovaIngest) sendToNova(inChan chan *bytes.Buffer) {
	go func() {
		defer close(n.ErrChan)
		for buffer := range inChan {
			_, err := Post(n.NovaURL, buffer, n.Auth, n.Gzip)
			if err != nil {
				n.ErrChan <- err
				return
			}
			n.TotalRequests++
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
