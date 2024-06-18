package openliga

type Match struct {
	DateUTC    string   `json:"matchDateTimeUTC"`
	Date       string   `json:"matchDateTime"`
	IsFinished bool     `json:"matchIsFinished"`
	Team1      Team     `json:"team1"`
	Team2      Team     `json:"team2"`
	Location   Location `json:"location"`
	Goals      []Goal   `json:"goals"`
}
type Team struct {
	Name string `json:"teamName"`
}
type Location struct {
	City    string `json:"locationCity"`
	Stadium string `json:"locationStadium"`
}
type Goal struct {
	Team1 int `json:"scoreTeam1"`
	Team2 int `json:"scoreTeam2"`
}
