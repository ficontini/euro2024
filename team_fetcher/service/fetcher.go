package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	version = "v3"
	params  = "?league=4&season=2024"
)

type APIFetcher struct {
	APIHost string
	APIKey  string
}

func NewAPIFetcher(apiHost, apiKey string) Fetcher {
	return &APIFetcher{
		APIHost: apiHost,
		APIKey:  apiKey,
	}
}

func (f *APIFetcher) fetchTeams() (*ApiTeamResp, error) {
	url := fmt.Sprintf("https://%s/%s/%s/%s", f.APIHost, version, "teams", params)
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
	var apiResponse *ApiTeamResp
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, nil
}
func (f *APIFetcher) fetchPlayersByTeam(teamID int) (*ApiPlayerResp, error) {
	url := fmt.Sprintf("https://%s/%s/%s/%s&team=%v", f.APIHost, version, "players", params, teamID)
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
	var apiResponse *ApiPlayerResp
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse, nil
}
