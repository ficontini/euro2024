package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Fetcher interface {
	FetchData() (*APIResponse, error)
}
type APIFetcher struct {
	APIHost string
	Path    string
	APIKey  string
}

func NewAPIFetcher(apiHost, apiKey, path string) Fetcher {
	return &APIFetcher{
		APIHost: apiHost,
		Path:    path,
		APIKey:  apiKey,
	}
}

func (f *APIFetcher) FetchData() (*APIResponse, error) {
	url := fmt.Sprintf("https://%s/%s", f.APIHost, f.Path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-rapidapi-key", f.APIKey)
	req.Header.Add("x-rapidapi-host", f.APIHost)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("api request failed")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var apiResponse *APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, nil
}
