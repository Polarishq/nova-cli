package source

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strings"
	"bytes"
)

type NovaSearch struct {
	Auth    string
	NovaURL string
	ErrChan chan error
}

// NewNovaSearch creates a new search obj
func NewNovaSearch(novaURL, auth string) *NovaSearch {
	return &NovaSearch{
		Auth:    auth,
		NovaURL: novaURL,
		ErrChan: make(chan error, 5),
	}
}

// WaitAndLogErrors blocks on the pipeline to complete and logs all errors
func (n *NovaSearch) WaitAndLogErrors() (errorsEncountered bool) {
	for e := range n.ErrChan {
		errorsEncountered = true
		log.Error(e)
	}
	return
}

func (n *NovaSearch) Search(searchTerms, transforms, report string) (data StrMatrix) {
	defer close(n.ErrChan)

	log.Debugf("Searching searchTerms='%+v'", searchTerms)
	log.Debugf("Searching transforms='%+v'", transforms)
	log.Debugf("Searching report='%+v'", report)

	searchTerms = fmt.Sprintf("source=%s* %s", novaCLISourcePrefix, searchTerms)

	searchQuery := NovaSearchEventsQuery{
		Blocking: true,
		SearchTerms: searchTerms,
		Transforms: strings.Split(transforms, ","),
		Reports: strings.Split(report, ","),
		Mode: "raw_1000",
	}

	searchQueryJSON, _ := json.Marshal(searchQuery)
	bytes := &bytes.Buffer{}
	bytes.Write(searchQueryJSON)

	results, err := Post(n.NovaURL+eventsSearchPath, bytes, n.Auth)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Raw Results: %+v", string(results))

	itemsJSON, err := ParseSearchResults(results)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("ItemsJSON: %+v", string(itemsJSON))

	if report != "" {
		events := NovaIncomingEventReporting{}
		err = json.Unmarshal(itemsJSON, &events)
		if err != nil {
			log.Error(err)
		}
		log.Debugf("reporting: %+v", events)
		if len(events) > 0 {
			if kvFoo, ok := events[0].Payload.(map[string]interface{}); ok {
				for k, v := range kvFoo {
					if vStr, ok := v.(string); ok {
						data = append(data, []string{k, vStr})
					}
				}
			}
		}
	} else {
		events := NovaIncomingEventNonReporting{}
		err = json.Unmarshal(itemsJSON, &events)
		if err != nil {
			log.Error(err)
		}
		log.Debugf("non reporting: %+v", events)
		for _, e1 := range events {
			data = append(data, []string{e1.Time ,e1.Payload.Event.Raw})
		}
	}

	log.Debugf("Processed Results: %+v", data)

	return data
}


