package apifotball

import "time"

type APIResponse struct {
	Response []Match `json:"response"`
}
type Match struct {
	Fixture Fixture `json:"fixture"`
	Teams   Teams   `json:"teams"`
	Goals
}
type Fixture struct {
	Date   time.Time `json:"date"`
	Venue  Venue     `json:"venue"`
	Status Status    `json:"status"`
}
type Venue struct {
	Name string `json:"name"`
	City string `json:"city"`
}
type Status struct {
	Short string `json:"short"`
}
type Teams struct {
	Home Team `json:"home"`
	Away Team `json:"away"`
}

type Team struct {
	Name   string `json:"name"`
	Winner *bool  `json:"winner"`
}
type Goals struct {
	Home int `json:"home"`
	Away int `json:"away"`
}
