package swimmy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type api struct {
	baseURL *url.URL
	apiKey  string
	debug   bool
}

func newAPI(rawURL, apiKey string, debug bool) (*api, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return &api{u, apiKey, debug}, nil
}

var apiRequestTimeout = 30 * time.Second

func (ap *api) do(method, url string, body []byte) (resp *http.Response, err error) {
	if ap.debug {
		log.Printf("[DEBUG] url:%s, body:%s", url, string(body))
		return nil, nil
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", ap.apiKey)
	req.Header.Set("User-Agent", "mackerel-swimmy")

	client := &http.Client{} // same as http.DefaultClient
	client.Timeout = apiRequestTimeout
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ap *api) postServiceMetrics(service string, values []metricValue) error {
	requestJSON, err := json.Marshal(values)
	if err != nil {
		return err
	}

	resp, err := ap.do("POST",
		ap.urlFor(fmt.Sprintf("/api/v0/services/%s/tsdb", service)),
		requestJSON,
	)

	if err != nil || ap.debug {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("API result failed: %s", resp.Status)
	}
	return nil
}

func (ap *api) urlFor(path string) string {
	newURL, _ := url.Parse(ap.baseURL.String())
	newURL.Path = path
	return newURL.String()
}
