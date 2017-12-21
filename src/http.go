package src

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"runtime"
)

var HTTPClient *http.Client

func init() {
	HTTPClient = &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   httpTimeout,
	}
}

// Post makes an HTTP POST
func Post(targetURL string, payload io.Reader, authHeader string, gzip bool) ([]byte, error) {
	req, _ := http.NewRequest("POST", targetURL, payload)
	if authHeader != "" {
		req.Header.Add("Authorization", authHeader)
	}
	if gzip {
		req.Header.Add("Content-Encoding", "gzip")
	}
	return doRequest(req)
}

// Get makes and HTTP GET
func Get(targetURL string, params map[string]string, authHeader string) ([]byte, error) {
	req, _ := http.NewRequest("GET", targetURL, nil)
	req.URL.RawQuery = convertToValues(params).Encode()
	if authHeader != "" {
		req.Header.Add("Authorization", authHeader)
	}
	return doRequest(req)
}

func doRequest(request *http.Request) ([]byte, error) {
	request.Header.Set("User-Agent", GetUserAgent())
	request.Header.Set("Content-Type", "application/json")
	//log.Debug("Request Full: ", request) // too verbose
	log.Debug("Request Headers: ", request.Header)
	log.Debug("Request Length: ", request.ContentLength)
	resp, err := HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error dialing to splunknova: %+v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Debug("Response: ", resp)
	log.Debug("Response Body: ", string(body))
	if err != nil {
		return nil, fmt.Errorf("error reading response body: code:%+v, body:%+v, err:%+v", resp.StatusCode, string(body), err)
	}
	if resp.StatusCode > 399 {
		return nil, fmt.Errorf("error communicating with splunknova. X-SPLUNK-REQ-ID=%+v code:%+v, body:%+v",
			resp.Header.Get("X-SPLUNK-REQ-ID"), resp.StatusCode, string(body))
	}
	return body, nil
}

func convertToValues(data map[string]string) url.Values {
	values := make(url.Values)
	for k, v := range data {
		values.Add(k, v)
	}
	return values
}

func GetUserAgent() string {
	return "nova-cli-" + AppVersion + "/" + runtime.GOOS + "-" + runtime.GOARCH
}