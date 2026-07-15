package strategy

import (
	"context"
	"math"

	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/domain/candidate"
	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/domain/decision"
)

type NearestDriver struct{}

func (s *NearestDriver) Name() string { return "nearest_driver" }

func (s *NearestDriver) Select(ctx context.Context, candidates *candidate.CandidateSet, ctxData decision.DecisionContext) (*candidate.Candidate, error) {
	available := candidates.FilterByStatus(candidate.CandidateAvailable)
	if len(available) == 0 {
		return nil, decision.ErrNoCandidates
	}

	best := &available[0]
	for i := range available {
		if available[i].DistanceMeters < best.DistanceMeters {
			best = &available[i]
		}
	}

	best.Status = candidate.CandidateSelected
	return best, nil
}

type HighestRated struct{}

func (s *HighestRated) Name() string { return "highest_rated" }

func (s *HighestRated) Select(ctx context.Context, candidates *candidate.CandidateSet, ctxData decision.DecisionContext) (*candidate.Candidate, error) {
	available := candidates.FilterByStatus(candidate.CandidateAvailable)
	if len(available) == 0 {
		return nil, decision.ErrNoCandidates
	}

	best := &available[0]
	for i := range available {
		if available[i].Score > best.Score {
			best = &available[i]
		}
	}

	best.Status = candidate.CandidateSelected
	return best, nil
}

type BalancedScore struct {
	distanceWeight float64
	scoreWeight    float64
}

func NewBalancedScore(distanceWeight, scoreWeight float64) *BalancedScore {
	return &BalancedScore{
		distanceWeight: distanceWeight,
		scoreWeight:    scoreWeight,
	}
}

func (s *BalancedScore) Name() string { return "balanced_score" }

func (s *BalancedScore) Select(ctx context.Context, candidates *candidate.CandidateSet, ctxData decision.DecisionContext) (*candidate.Candidate, error) {
	available := candidates.FilterByStatus(candidate.CandidateAvailable)
	if len(available) == 0 {
		return nil, decision.ErrNoCandidates
	}

	var maxDist, maxScore float64
	for _, c := range available {
		if c.DistanceMeters > maxDist {
			maxDist = c.DistanceMeters
		}
		if c.Score > maxScore {
			maxScore = c.Score
		}
	}

	best := &available[0]
	bestComposite := math.MaxFloat64

	for i := range available {
		normDist := 1.0
		if maxDist > 0 {
			normDist = available[i].DistanceMeters / maxDist
		}
		normScore := 0.0
		if maxScore > 0 {
			normScore = available[i].Score / maxScore
		}

		composite := s.distanceWeight*normDist - s.scoreWeight*normScore
		if composite < bestComposite {
			bestComposite = composite
			best = &available[i]
		}
	}

	best.Status = candidate.CandidateSelected
	return best, nil
}

type StrategyRegistry struct {
	strategies map[string]Strategy
}

type Strategy interface {
	Name() string
	Select(ctx context.Context, candidates *candidate.CandidateSet, ctxData decision.DecisionContext) (*candidate.Candidate, error)
}

func NewStrategyRegistry() *StrategyRegistry {
	return &StrategyRegistry{
		strategies: make(map[string]Strategy),
	}
}

func (r *StrategyRegistry) Register(s Strategy) {
	r.strategies[s.Name()] = s
}

func (r *StrategyRegistry) Get(name string) (Strategy, error) {
	s, ok := r.strategies[name]
	if !ok {
		return nil, decision.ErrStrategyNotFound
	}
	return s, nil
}
