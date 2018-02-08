package source

import (
	//"encoding/json"
	//"fmt"
	log "github.com/Sirupsen/logrus"
	//"strings"
	"encoding/json"
	"fmt"
	"bytes"
)

type NovaMetricsSearch struct {
	Auth    string
	NovaURL string
	ErrChan chan error
}

// NewNovaSearch creates a new search obj
func NewNovaMetricsSearch(novaURL, auth string) *NovaMetricsSearch {
	return &NovaMetricsSearch{
		Auth:    auth,
		NovaURL: novaURL,
		ErrChan: make(chan error, 5),
	}
}

// WaitAndLogErrors blocks on the pipeline to complete and logs all errors
func (n *NovaMetricsSearch) WaitAndLogErrors() (errorsEncountered bool) {
	for e := range n.ErrChan {
		errorsEncountered = true
		log.Error(e)
	}
	return
}

func (n *NovaMetricsSearch) GetLs() (StrMatrix, error) {
	defer close(n.ErrChan)

	urlFinal := n.NovaURL + metricsListPath

	results, err := Get(urlFinal, nil, n.Auth)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("Raw Results: %+v\n\n", string(results))

	m := MetricsLSResponse{}
	json.Unmarshal(results, &m)
	data := StrMatrix{}
	for k, v := range m.Metrics {
		data = append(data, []string{fmt.Sprint(k+1), v})
	}
	return data, nil
}


func (n *NovaMetricsSearch) GetStats(metric_names, stats, groupBy []string, span string) (data StrMatrix, err error) {
	defer close(n.ErrChan)

	log.Debugf("Searching metric_names='%+v'", metric_names)
	log.Debugf("Searching stats='%+v'", stats)
	log.Debugf("Searching groupBy='%+v'", groupBy)
	log.Debugf("Searching span='%+v'", span)

	urlFinal := n.NovaURL + "/v1/search/stats"

	searchQuery := NovaSearchStatsQuery{
		Fields:     metric_names,
		FieldsType: "metrics",
		GroupBy:    groupBy,
		Stats:      stats,
		Span:       span,
		Blocking:   true,
	}
	searchQueryJSON, _ := json.Marshal(searchQuery)
	bytes := &bytes.Buffer{}
	bytes.Write(searchQueryJSON)

	results, err := Post(urlFinal, bytes, n.Auth)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Raw Results: %+v", string(results))
	itemsJSON, err := ParseSearchResults(results)
	if err != nil {
		return
	}

	items := []map[string]interface{}{}
	err = json.Unmarshal(itemsJSON, &items)
	if err != nil{
		log.Error(err)
		return
	}

	if len(items) > 0 {
		for k, v := range items[0] {
			if vStr, ok := v.(string); ok {
				data = append(data, []string{k, vStr})
			} else {
				log.Warningf("Skipping non-string value %+v", v)
			}
		}
		log.Debugf("Processed Results: %+v", data)
	}
	return
}
