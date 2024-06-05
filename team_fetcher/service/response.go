package service

type ApiTeamResp struct {
	Response []TeamResp `json:"response"`
}

type TeamResp struct {
	Team  Team  `json:"team"`
	Venue Venue `json:"venue"`
}

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Venue struct {
	Name string `json:"name"`
	City string `json:"city"`
}

type ApiPlayerResp struct {
	Response []PlayerResp `json:"response"`
}
type PlayerResp struct {
	Player     Player       `json:"player"`
	Statistics []Statistics `json:"statistics"`
}
type Player struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
}
type Statistics struct {
	Shots Shots `json:"shots"`
	Goals Goals `json:"goals"`
	Cards Cards `json:"cards"`
}
type Shots struct {
	Total int `json:"total"`
}
type Goals struct {
	Total   int `json:"total"`
	Assists int `json:"assists"`
}
type Cards struct {
	Yellow int `json:"yellow"`
	Red    int `json:"red"`
}
