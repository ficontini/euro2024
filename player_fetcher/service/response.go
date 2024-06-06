package service

type ApiResponse struct {
	Response []PlayerResp `json:"response"`
	Paging   Paging       `json:"paging"`
}
type Paging struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}
type PlayerResp struct {
	Player     Player       `json:"player"`
	Statistics []Statistics `json:"statistics"`
}
type Player struct {
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
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
