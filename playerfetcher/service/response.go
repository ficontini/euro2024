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
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
}
type Statistics struct {
	Team   Team   `json:"team"`
	Games  Games  `json:"games"`
	Goals  Goals  `json:"goals"`
	Cards  Cards  `json:"cards"`
	Passes Passes `json:"passes"`
}
type Games struct {
	Position string `json:"position"`
}
type Team struct {
	Name string `json:"name"`
}

type Goals struct {
	Total   int `json:"total"`
	Assists int `json:"assists"`
}
type Cards struct {
	Yellow int `json:"yellow"`
	Red    int `json:"red"`
}

type Passes struct {
	Total    int `json:"total"`
	Accuracy int `json:"accuracy"`
}
