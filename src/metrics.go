package src

import (
	//"encoding/json"
	//"fmt"
	log "github.com/Sirupsen/logrus"
	//"strings"
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


func (n *NovaMetricsSearch) GetAggergations(metric_names, aggregations, groupBy, span string) {
	defer close(n.ErrChan)

	log.Debugf("Searching metric_names='%+v'", metric_names)
	log.Debugf("Searching aggregations='%+v'", aggregations)
	log.Debugf("Searching groupBy='%+v'", groupBy)
	log.Debugf("Searching span='%+v'", span)


	params := map[string]string{
		"group_by":       groupBy,
		"span":           span,
	}

	urlFinal := n.NovaURL + metricsURLPath + "/" + metric_names + "/" + aggregations

	results, err := Get(urlFinal, params, n.Auth)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("Raw Results: %+v\n\n", string(results))

	log.Info(string(results))
}

//const field2MaxWidth = 100
//
//func breakString(bigString string) []string {
//	if len(bigString) < field2MaxWidth {
//		return []string{bigString}
//	} else {
//		return append([]string{bigString[0:field2MaxWidth-1]}, breakString(bigString[field2MaxWidth:])...)
//	}
//}
//
//func printTable2(data StrMatrix) {
//	maxwidthF1 := 30
//	maxwidthF2 := field2MaxWidth
//	strFormat := fmt.Sprintf("│%%%ds │ %%-%ds│\n", maxwidthF1, maxwidthF2)
//	fmt.Println("┌" + strings.Repeat("─", maxwidthF1+1) + "┬" + strings.Repeat("─", maxwidthF2+1) + "┐")
//	for _, datum := range data {
//		f1 := datum[0]
//		f2 := breakString(datum[1])
//		for _, ff2 := range f2 {
//			fmt.Printf(strFormat, f1, ff2)
//			f1 = ""
//		}
//	}
//	fmt.Println("└" + strings.Repeat("─", maxwidthF1+1) + "┴" + strings.Repeat("─", maxwidthF2+1) + "┘")
//}
//
//
//func printTable(data map[string]string) {
//	maxwidthK, maxwidthV := 1, 1
//
//	for k, v := range data {
//		if len(k) > maxwidthK {
//			maxwidthK = len(k)
//		}
//		if len(v) > maxwidthV {
//			maxwidthV = len(v)
//		}
//	}
//	maxwidthK = maxwidthK + 2
//	maxwidthV = maxwidthV + 2
//	strFormat := fmt.Sprintf("│%%%ds │ %%-%ds│\n", maxwidthK, maxwidthV)
//	fmt.Println("┌" + strings.Repeat("─", maxwidthK+1) + "┬" + strings.Repeat("─", maxwidthV+1) + "┐")
//	for k, v := range data {
//		fmt.Printf(strFormat, k, v)
//	}
//	fmt.Println("└" + strings.Repeat("─", maxwidthK+1) + "┴" + strings.Repeat("─", maxwidthV+1) + "┘")
//}
