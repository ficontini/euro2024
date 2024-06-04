package openliga

type Match struct {
	Date     string   `json:"matchDateTime"`
	Team1    Team     `json:"team1"`
	Team2    Team     `json:"team2"`
	Location Location `json:"location"`
	Goals    []Goal   `json:"matchResults"`
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
