package swimmy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (ap *api) do(req *http.Request) (resp *http.Response, err error) {
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

	req, err := http.NewRequest(
		"POST",
		ap.urlFor(fmt.Sprintf("/api/v0/services/%s/tsdb", service)).String(),
		bytes.NewReader(requestJSON),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := ap.do(req)
	if err != nil {
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

func (ap *api) urlFor(path string) *url.URL {
	newURL, err := url.Parse(ap.baseURL.String())
	if err != nil {
		panic("invalid url passed")
	}
	newURL.Path = path
	return newURL
}
