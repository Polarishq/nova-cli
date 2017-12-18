package src

import (
	//"encoding/json"
	//"fmt"
	log "github.com/Sirupsen/logrus"
	//"strings"
	"encoding/json"
	"fmt"
)

type NovaMetricsSearch struct {
	Auth    string
	NovaURL string
	ErrChan chan error
}

type MetricsLSResponse struct {
	Metrics []string `json:"metrics"`
}

type MetricsGetResponse struct {
	Aggregations []string `json:"aggregations"`
	Dimensions []string `json:"dimensions"`
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

	urlFinal := n.NovaURL + metricsURLSearchPath

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


func (n *NovaMetricsSearch) GetAggregations(metric_names, aggregations, groupBy, span string) (StrMatrix, error) {
	defer close(n.ErrChan)

	log.Debugf("Searching metric_names='%+v'", metric_names)
	log.Debugf("Searching aggregations='%+v'", aggregations)
	log.Debugf("Searching groupBy='%+v'", groupBy)
	log.Debugf("Searching span='%+v'", span)


	params := map[string]string{
		"group_by":       groupBy,
		"span":           span,
	}

	urlFinal := n.NovaURL + metricsURLSearchPath + "/" + metric_names + "/" + aggregations

	results, err := Get(urlFinal, params, n.Auth)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("Raw Results: %+v\n\n", string(results))

	m := []map[string]string{}
	json.Unmarshal(results, &m)
	data := StrMatrix{}
	if len(m) < 1 {
		return data, nil
	}
	for k, v := range m[0] {
		data = append(data, []string{k, v})
	}
	return data, nil
}
