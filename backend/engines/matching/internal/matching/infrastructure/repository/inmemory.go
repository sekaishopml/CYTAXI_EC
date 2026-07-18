package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/candidate"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/matching"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/ranking"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/valueobject"
)

var ErrMatchingNotFound = errors.New("matching not found")

// InMemoryMatchingRepository almacena matchings in-memory con sync.RWMutex
type InMemoryMatchingRepository struct {
	mu      sync.RWMutex
	byID    map[string]*matching.Matching
	byTrip  map[string][]matching.Matching
}

func NewInMemoryMatchingRepository() *InMemoryMatchingRepository {
	return &InMemoryMatchingRepository{
		byID:   make(map[string]*matching.Matching),
		byTrip: make(map[string][]matching.Matching),
	}
}

func (r *InMemoryMatchingRepository) FindByID(ctx context.Context, id valueobject.MatchingID) (*matching.Matching, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.byID[string(id)]
	if !ok {
		return nil, ErrMatchingNotFound
	}
	return m, nil
}

func (r *InMemoryMatchingRepository) FindByTripID(ctx context.Context, tripID valueobject.TripID) ([]matching.Matching, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]matching.Matching, len(r.byTrip[string(tripID)]))
	copy(out, r.byTrip[string(tripID)])
	return out, nil
}

func (r *InMemoryMatchingRepository) Save(ctx context.Context, m *matching.Matching) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[string(m.ID)] = m
	r.byTrip[string(m.TripID)] = append(r.byTrip[string(m.TripID)], *m)
	return nil
}

func (r *InMemoryMatchingRepository) Update(ctx context.Context, m *matching.Matching) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[string(m.ID)] = m
	return nil
}

// InMemoryCandidateRepository almacena candidates y rankings
type InMemoryCandidateRepository struct {
	mu       sync.RWMutex
	sets     map[string]*candidate.CandidateSet
	rankings map[string]*ranking.CandidateRanking
	attempts map[string]*candidate.AssignmentAttempt
}

func NewInMemoryCandidateRepository() *InMemoryCandidateRepository {
	return &InMemoryCandidateRepository{
		sets:     make(map[string]*candidate.CandidateSet),
		rankings: make(map[string]*ranking.CandidateRanking),
		attempts: make(map[string]*candidate.AssignmentAttempt),
	}
}

func (r *InMemoryCandidateRepository) FindByMatchingID(ctx context.Context, id valueobject.MatchingID) (*candidate.CandidateSet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sets[string(id)]
	if !ok {
		return nil, ErrMatchingNotFound
	}
	return s, nil
}

func (r *InMemoryCandidateRepository) FindRanking(ctx context.Context, id valueobject.MatchingID) (*ranking.CandidateRanking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rnk, ok := r.rankings[string(id)]
	if !ok {
		return nil, ErrMatchingNotFound
	}
	return rnk, nil
}

func (r *InMemoryCandidateRepository) Save(ctx context.Context, cs *candidate.CandidateSet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sets[string(cs.MatchingID)] = cs
	return nil
}

func (r *InMemoryCandidateRepository) SaveAttempt(ctx context.Context, a *candidate.AssignmentAttempt) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.attempts[string(a.MatchingID)] = a
	return nil
}
