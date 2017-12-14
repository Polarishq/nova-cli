package src

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

type NovaSearch struct {
	Auth    string
	NovaURL string
	ErrChan chan error
}

type NovaResults struct {
	NovaEvents []struct {
		Time string `json:"time"`
		Raw  string `json:"event.raw"`
	} `json:"events"`
}

type NovaResultsStats struct {
	NovaEvents []map[string]string `json:"events"`
}

// NewNovaSearch creates a new search obj
func NewNovaSearch(novaURL, auth string) *NovaSearch {
	return &NovaSearch{
		Auth:    auth,
		NovaURL: novaURL,
		ErrChan: make(chan error, 5),
	}
}

// BlockedErrorLogger blocks on the pipeline to complete and logs all errors
func (n *NovaSearch) BlockedErrorLogger() (errorsEncountered bool) {
	for e := range n.ErrChan {
		errorsEncountered = true
		log.Error(e)
	}
	return
}

func (n *NovaSearch) Search(keywords, transforms, report string) {
	defer close(n.ErrChan)

	log.Debugf("Searching keywords='%+v'", keywords)
	log.Debugf("Searching transforms='%+v'", transforms)
	log.Debugf("Searching report='%+v'", report)

	params := map[string]string{
		"keywords":   keywords,
		"transforms": transforms,
		"report":     report,
		"count":      defaultSearchResultsCount,
	}

	results, err := Get(n.NovaURL+eventsURLPath, params, n.Auth)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("Raw Results: %+v\n\n", string(results))

	if report == "" {
		n1 := NovaResults{}
		json.Unmarshal(results, &n1)
		for n, ne := range n1.NovaEvents {
			fmt.Printf("%5d %s %s\n", n, ne.Time, ne.Raw)
		}
	} else {
		n1 := NovaResultsStats{}
		json.Unmarshal(results, &n1)
		maxwidthK, maxwidthV := 1, 1
		for _, ne := range n1.NovaEvents {
			for k, v := range ne {
				if len(k) > maxwidthK {
					maxwidthK = len(k)
				}
				if len(v) > maxwidthV {
					maxwidthV = len(v)
				}
			}
			maxwidthK = maxwidthK + 2
			maxwidthV = maxwidthV + 2
			strFormat := fmt.Sprintf("|%%%ds | %%-%ds|\n", maxwidthK, maxwidthV)
			for k, v := range ne {
				fmt.Printf(strFormat, k, v)
			}
		}
	}
}
