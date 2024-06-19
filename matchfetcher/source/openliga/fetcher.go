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
	errch   chan error
	matchch chan []*Match
}

func NewAPIFetcher(apiHost string) service.Fetcher {
	return &APIFetcher{
		APIHost: apiHost,
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
			f.errch <- err
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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.New("request time out")
		}
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
