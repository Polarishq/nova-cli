package src

import (
	"net/url"
	"io"
	"net/http"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

var HTTPClient *http.Client

func init() {
	HTTPClient = &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   httpTimeout,
	}
}

// Post makes an HTTP POST
func Post(targetURL string, payload io.Reader, authHeader string) ([]byte, error) {
	log.Debugf("POST targetURL=%s, payload=%+v, auth=%s", targetURL, payload, authHeader)
	req, _ := http.NewRequest("POST", targetURL, payload)
	if authHeader != "" {
		req.Header.Add("Authorization", authHeader)
	}
	return doRequest(req)
}

// Get makes and HTTP GET
func Get(targetURL string, params map[string]string, authHeader string) ([]byte, error) {
	log.Debugf("GET targetURL=%s, params=%+v, auth=%s", targetURL, params, authHeader)

	req, _ := http.NewRequest("GET", targetURL, nil)
	req.URL.RawQuery = convertToValues(params).Encode()
	if authHeader != "" {
		req.Header.Add("Authorization", authHeader)
	}
	return doRequest(req)
}

func doRequest(request *http.Request) ([]byte, error) {
	request.Header.Set("User-Agent", "nova-cli-" + AppVersion)
	request.Header.Set("Content-Type", "application/json")
	resp, err := HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error dialing to splunknova: %+v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Debugf("responseBody=%+v, responseCode=%+v, err=%+v", string(body), resp.StatusCode, err)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: code:%+v, body:%+v, err:%+v", resp.StatusCode, string(body), err)
	}
	if resp.StatusCode > 399 {
		return nil, fmt.Errorf("error communicating with splunknova. X-SPLUNK-REQ-ID=%+v code:%+v, body:%+v",
			resp.Header.Get("X-SPLUNK-REQ-ID"), resp.StatusCode, string(body))
	}
	return body, nil
}

func convertToValues(data map[string]string) (url.Values) {
	values := make(url.Values)
	for k, v := range data {
		values.Add(k, v)
	}
	return values
}
