package store

import (
	"context"
	"sort"
	"sync"

	"github.com/ficontini/euro2024/types"
)

type Store interface {
	Add(context.Context, *types.Match) error
	Get(context.Context) ([]*types.Match, error)
	GetMatchesByTeam(context.Context, string) ([]*types.Match, error)
	GetMatchesByRound(context.Context, types.RoundStatus) ([]*types.Match, error)
	Clean(context.Context) error
}

type InMemoryStore struct {
	mu      sync.RWMutex
	matches []*types.Match
}

func NewInMemoryStore() Store {
	return &InMemoryStore{
		matches: []*types.Match{},
	}

}
func (s *InMemoryStore) Add(_ context.Context, match *types.Match) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.matches = append(s.matches, match)
	sort.Slice(s.matches, func(i, j int) bool {
		return s.matches[i].Date.Before(s.matches[j].Date)
	})

	return nil
}
func (s *InMemoryStore) Get(_ context.Context) ([]*types.Match, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	matches := make([]*types.Match, len(s.matches))

	copy(matches, s.matches)

	return matches, nil
}
func (s *InMemoryStore) GetMatchesByTeam(_ context.Context, team string) ([]*types.Match, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var matches []*types.Match
	for _, m := range s.matches {
		if m.Home.Name == team || m.Away.Name == team {
			matches = append(matches, m)
		}
	}
	return matches, nil
}
func (s *InMemoryStore) GetMatchesByRound(_ context.Context, round types.RoundStatus) ([]*types.Match, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var matches []*types.Match
	for _, m := range s.matches {
		if m.Round == round {
			matches = append(matches, m)
		}
	}
	return matches, nil
}
func (s *InMemoryStore) Clean(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.matches = []*types.Match{}

	return nil
}
