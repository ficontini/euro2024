package types

type Player struct {
	FirstName  string
	LastName   string
	Team       string
	Age        int
	Statistics *Statistics
}

func NewPlayer(firstName, lastName, team string, age int, statistics *Statistics) *Player {
	return &Player{
		FirstName:  firstName,
		LastName:   lastName,
		Team:       team,
		Age:        age,
		Statistics: statistics,
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
