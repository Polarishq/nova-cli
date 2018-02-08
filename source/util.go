package source

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
)

func ParseSearchResults(results []byte) (jsonBytes []byte, err error) {
	novaResponse := NovaSearchResponse{}
	err = json.Unmarshal(results, &novaResponse)
	if err != nil {
		log.Error(err)
		return
	}
	for _, e := range novaResponse.Errors {
		log.Warningf("Error from splunknova.com: %+v", e)
	}
	if novaResponse.Job.Status != "complete" {
		log.Errorf("Exiting, job failed or not complete: %+v", novaResponse.Job)
		return
	}
	if len(novaResponse.Items) < 1 {
		return
	}
	jsonBytes, err = json.Marshal(novaResponse.Items)
	return
}