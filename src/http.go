package src

import (
	"net/url"
	"bytes"
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"time"
	"fmt"
	"io/ioutil"
)

var HTTPClient *http.Client

func init() {
	HTTPClient = &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   10 * time.Second,
	}
}

// Post makes an HTTP POST
func Post(targetURL string, params map[string]string, authHeader string) ([]byte, error) {
	log.Debugf("POST targetURL=%s, params=%+v, auth=%s", targetURL, params, authHeader)
	var payload io.Reader
	if params != nil {
		payload = bytes.NewBufferString(convertToValues(params).Encode())
	}
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
	resp, err := HTTPClient.Do(request)
	if err != nil {
		log.WithFields(log.Fields{"code": resp.StatusCode, "request": request, "err": err}).Errorf("error dialing to splunk")
		return nil, fmt.Errorf("server error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"code": resp.StatusCode, "body": body, "err": err}).Errorf("error reading response body")
		return nil, fmt.Errorf("server error")
	}
	if resp.StatusCode > 399 {
		log.WithFields(log.Fields{"code": resp.StatusCode, "body": body}).Errorf("non 2xx or non 3xx response code")
		return nil, fmt.Errorf("server error")
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
