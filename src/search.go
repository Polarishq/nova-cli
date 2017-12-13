package src

import (
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"fmt"
)

type NovaSearch struct {
	Auth string
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

func (n *NovaSearch) Search(keywords, transforms, report string) {
	params := map[string]string{
		"keywords": keywords,
		"transforms": transforms,
		"report": report,
		"count": "20",
	}

	results, err := Get(urlDefaultHost + urlDefaultPath, params, n.Auth)
	if err != nil {
		panic(err)
	}

	if report == "" {
		n1 := NovaResults{}
		json.Unmarshal(results, &n1)
		log.Debugf("All Results: %+v\n\n", string(results))
		for n, ne := range n1.NovaEvents {
			fmt.Printf("%5d %s %s\n", n, ne.Time, ne.Raw)
		}
	} else {
		n1 := NovaResultsStats{}
		json.Unmarshal(results, &n1)
		log.Debugf("All Results: %+v\n\n", string(results))
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
			//fmt.Printf("Results %d\n", i+1)
			for k, v := range ne {
				fmt.Printf(strFormat, k, v)
			}
		}
	}
}
