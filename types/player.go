package types

type Player struct {
	FirstName   string `json:"firstName" dynamodbav:"firstName"`
	LastName    string `json:"lastName" dynamodbav:"lastName"`
	Team        string `json:"team" dynamodbav:"team"`
	Age         int    `json:"age" dynamodbav:"age"`
	Shots       int    `json:"shots" dynamodbav:"shots"`
	Goals       int    `json:"goals" dynamodbav:"goals"`
	Assists     int    `json:"assists" dynamodbav:"assists"`
	YellowCards int    `json:"yellowCards" dynamodbav:"yellowCards"`
	RedCards    int    `json:"redCards" dynamodbav:"redCards"`
}

func NewPlayer(firstName, lastName, team string, age int, statistics *Statistics) *Player {
	return &Player{
		FirstName:   firstName,
		LastName:    lastName,
		Team:        team,
		Age:         age,
		Shots:       statistics.Performance.Shots,
		Goals:       statistics.Performance.Goals,
		Assists:     statistics.Performance.Assists,
		YellowCards: statistics.Cards.Yellow,
		RedCards:    statistics.Cards.Red,
	}
}

type Statistics struct {
	Performance *Performance
	Cards       *Cards
}

func NewStatistics(performance *Performance, cards *Cards) *Statistics {
	return &Statistics{
		Performance: performance,
		Cards:       cards,
	}
}

type Performance struct {
	Shots   int
	Goals   int
	Assists int
}

func NewPerformance(shots, goals, assists int) *Performance {
	return &Performance{
		Shots:   shots,
		Goals:   goals,
		Assists: assists,
	}
}

type Cards struct {
	Yellow int
	Red    int
}

func NewCards(yellow, red int) *Cards {
	return &Cards{
		Yellow: yellow,
		Red:    red,
	}
}
