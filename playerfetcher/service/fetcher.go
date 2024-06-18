package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const path = "v3/players?league=4&season=2024"

type Fetcher interface {
	FetchData() ([]*ApiResponse, error)
}
type APIFetcher struct {
	APIHost string
	APIKey  string
	url     string
}

func NewAPIFetcher(apiHost, apiKey string) Fetcher {
	return &APIFetcher{
		APIHost: apiHost,
		APIKey:  apiKey,
		url:     fmt.Sprintf("https://%s/%s&page", apiHost, path),
	}
}
func (f *APIFetcher) FetchData() ([]*ApiResponse, error) {
	var (
		total    = 2
		response []*ApiResponse
	)
	for i := 1; i <= total; i++ {
		res, err := f.fetch(i)
		if err != nil {
			return nil, err
		}
		total = res.Paging.Total
		response = append(response, res)
		time.Sleep(6 * time.Second)
	}
	fmt.Println("len:", len(response))
	return response, nil
}

func (f *APIFetcher) fetch(page int) (*ApiResponse, error) {
	fmt.Println("page:", page)
	req, err := f.newRequest(page)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("api request failed")
	}

	return f.unmarshallResponse(resp)
}

func (f *APIFetcher) newRequest(page int) (*http.Request, error) {
	url := fmt.Sprintf("%s=%v", f.url, page)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-rapidapi-key", f.APIKey)
	req.Header.Add("x-rapidapi-host", f.APIHost)
	return req, nil
}
func (f *APIFetcher) unmarshallResponse(resp *http.Response) (*ApiResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse *ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, nil
}
