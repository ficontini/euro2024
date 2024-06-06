package types

import "time"

type MatchStatus string

const (
	NS             MatchStatus = "No Started"
	LIVE           MatchStatus = "Live"
	FINISH         MatchStatus = "Finish"
	UNKNOWN        MatchStatus = "Unknown"
	default_action             = "sendMessage"
)

type Match struct {
	Date     time.Time   `json:"date"`
	Location *Location   `json:"location"`
	Home     *MatchTeam  `json:"home"`
	Away     *MatchTeam  `json:"away"`
	Status   MatchStatus `json:"status"`
}

func NewMatch(date time.Time, location *Location, home, away *MatchTeam, status MatchStatus) *Match {
	return &Match{
		Date:     date,
		Location: location,
		Home:     home,
		Away:     away,
		Status:   status,
	}
}
func (m *Match) IsLive() bool {
	return m.Status == LIVE
}
func (m *Match) IsUpcoming() bool {
	return m.Date.After(time.Now()) && m.Status != FINISH && m.Status != LIVE
}

type Location struct {
	City    string `json:"city"`
	Stadium string `json:"stadium"`
}

func NewLocation(city, stadium string) *Location {
	return &Location{
		City:    city,
		Stadium: stadium,
	}
}

type MatchTeam struct {
	Name  string `json:"name"`
	Goals int    `json:"goals"`
}

func NewMatchTeam(name string, goals int) *MatchTeam {
	return &MatchTeam{
		Name:  name,
		Goals: goals,
	}
}

type Message struct {
	Action  string   `json:"action"`
	Matches []*Match `json:"matches"`
}

func NewMessage(matches []*Match) Message {
	return Message{
		Action:  default_action,
		Matches: matches,
	}
}
