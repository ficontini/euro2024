package types

type Team struct {
	ID      int
	Name    string
	Stadium *Stadium
	Players []*Player
}

func NewTeam(id int, name string, stadium *Stadium) *Team {
	return &Team{
		ID:      id,
		Name:    name,
		Stadium: stadium,
	}
}

type Stadium struct {
	Name string
	City string
}

func NewStadium(name, city string) *Stadium {
	return &Stadium{
		Name: name,
		City: city,
	}
}

type Player struct {
	FirstName  string
	LastName   string
	Age        int
	Statistics *Statistics
}

func NewPlayer(firstName, lastName string, age int, statistics *Statistics) *Player {
	return &Player{
		FirstName:  firstName,
		LastName:   lastName,
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
