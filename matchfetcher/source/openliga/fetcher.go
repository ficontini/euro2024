package openliga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/ficontini/euro2024/matchfetcher/service"
)

const lastPage = 7

type APIFetcher struct {
	APIHost string
	client  *http.Client
	errch   chan error
	matchch chan []*Match
}

func NewAPIFetcher(apiHost string, client *http.Client) service.Fetcher {
	if client == nil {
		client = http.DefaultClient
	}
	return &APIFetcher{
		APIHost: apiHost,
		client:  client,
		errch:   make(chan error, lastPage),
		matchch: make(chan []*Match, lastPage),
	}
}

func (f *APIFetcher) FetchData() (any, error) {
	wg := sync.WaitGroup{}
	wg.Add(lastPage)
	for i := 1; i <= lastPage; i++ {
		go func(page int) {
			resp, err := f.makeRequest(fmt.Sprintf("%s/%v", f.APIHost, page))
			if err != nil {
				f.errch <- err
				return
			}
			f.matchch <- resp
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(f.errch)
	close(f.matchch)
	for err := range f.errch {
		if err != nil {
			return nil, err
		}
	}

	var matches []*Match
	for m := range f.matchch {
		matches = append(matches, m...)
	}
	return matches, nil
}

func (f *APIFetcher) makeRequest(url string) ([]*Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := f.client.Do(req)
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
