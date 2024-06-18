package openliga

import (
	"embed"
	"encoding/json"
	"log"
	"time"

	"github.com/ficontini/euro2024/matchfetcher/service"
	"github.com/ficontini/euro2024/types"
)

const (
	layout           = "2006-01-02T15:04:05"
	layoutUTC        = "2006-01-02T15:04:05Z"
	translation_file = "translations.json"
)

//go:embed translations.json
var translationsFile embed.FS

type APIProcessor struct {
	translations map[string]string
}

func NewApiProcessor() service.Processor {
	translations, err := loadTranslations(translation_file)
	if err != nil {
		log.Fatal(err)
	}
	return &APIProcessor{
		translations: translations,
	}
}
func (p *APIProcessor) ProcessData(data any) ([]*types.Match, error) {
	resp := data.([]*Match)
	var matches []*types.Match
	for _, r := range resp {
		dateUTC, err := time.Parse(layoutUTC, r.DateUTC)
		if err != nil {
			return nil, err
		}
		date, err := time.Parse(layout, r.Date)
		if err != nil {
			return nil, err
		}
		result := calculateResult(r.Goals)
		match := types.NewMatch(
			date,
			types.NewLocation(r.Location.City, r.Location.Stadium),
			types.NewMatchTeam(p.translateTeamName(r.Team1.Name), result.Team1),
			types.NewMatchTeam(p.translateTeamName(r.Team2.Name), result.Team2),
			calculateState(dateUTC.UTC(), r.IsFinished),
		)
		matches = append(matches, match)
	}
	return matches, nil
}
func (p *APIProcessor) translateTeamName(name string) string {
	if translated, exists := p.translations[name]; exists {
		return translated
	}
	return name
}

func calculateState(date time.Time, isFinished bool) types.MatchStatus {
	now := time.Now().UTC()
	if date.After(now) {
		return types.NS
	}

	if isFinished {
		return types.FINISH
	}

	return types.LIVE
}
func calculateResult(goals []Goal) *Result {
	if len(goals) == 0 {
		return newResult(0, 0)
	}
	lastGoal := goals[len(goals)-1]
	return newResult(lastGoal.Team1, lastGoal.Team2)
}

func loadTranslations(filepath string) (map[string]string, error) {
	bytes, err := translationsFile.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var translations map[string]string
	if err := json.Unmarshal(bytes, &translations); err != nil {
		return nil, err
	}
	return translations, nil
}

type Result struct {
	Team1 int
	Team2 int
}

func newResult(team1, team2 int) *Result {
	return &Result{
		Team1: team1,
		Team2: team2,
	}
}
