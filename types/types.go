package types

import "time"

type MatchStatus string

const (
	NS      MatchStatus = "No Started"
	LIVE    MatchStatus = "Live"
	FINISH  MatchStatus = "Finish"
	UNKNOWN MatchStatus = "Unknown"
)

type Match struct {
	Date     time.Time
	Location *Location
	Home     string
	Away     string
	Status   MatchStatus
	Result   *Result
}

func NewMatch(date time.Time, location *Location, home, away string, status MatchStatus, result *Result) *Match {
	return &Match{
		Date:     date,
		Location: location,
		Home:     home,
		Away:     away,
		Status:   status,
		Result:   result,
	}
}
func (m *Match) IsLive() bool {
	return m.Status == LIVE
}
func (m *Match) IsUpcoming() bool {
	return m.Date.After(time.Now())
}

type Location struct {
	City    string
	Stadium string
}

func NewLocation(city, stadium string) *Location {
	return &Location{
		City:    city,
		Stadium: stadium,
	}
}

type Result struct {
	Home int
	Away int
}

func NewResult(home, away int) *Result {
	return &Result{
		Home: home,
		Away: away,
	}
}
