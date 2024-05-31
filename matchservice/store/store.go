package store

import (
	"context"
	"sync"

	"github.com/ficontini/euro2024/types"
)

type Store interface {
	Add(context.Context, *types.Match) error
	Get(context.Context) ([]*types.Match, error)
	Clean(context.Context) error
}

type InMemoryStore struct {
	mu            sync.RWMutex
	matches       []*types.Match
	matchesByTeam map[string][]*types.Match
}

func NewInMemoryStore() Store {
	return &InMemoryStore{
		matchesByTeam: make(map[string][]*types.Match),
		matches:       []*types.Match{},
	}
}
func (s *InMemoryStore) Add(_ context.Context, match *types.Match) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.matches = append(s.matches, match)

	s.matchesByTeam[match.Home] = append(s.matchesByTeam[match.Home], match)
	s.matchesByTeam[match.Away] = append(s.matchesByTeam[match.Away], match)

	return nil
}
func (s *InMemoryStore) Get(_ context.Context) ([]*types.Match, error) {
	return s.matches, nil
}
func (s *InMemoryStore) GetMatchesByTeam(_ context.Context, team string) ([]*types.Match, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.matchesByTeam[team], nil
}
func (s *InMemoryStore) Clean(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.matches = []*types.Match{}
	s.matchesByTeam = make(map[string][]*types.Match)

	return nil
}
