package openliga

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ficontini/euro2024/match_fetcher/fetcher"
)

type APIFetcher struct {
	APIHost string
}

func NewAPIFetcher(apiHost string) fetcher.Fetcher {
	return &APIFetcher{
		APIHost: apiHost,
	}
}

func (f *APIFetcher) FetchData() (any, error) {
	req, err := http.NewRequest(http.MethodGet, f.APIHost, nil)
	if err != nil {
		return nil, err
	}
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
	var matches []*Match
	err = json.Unmarshal(body, &matches)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
