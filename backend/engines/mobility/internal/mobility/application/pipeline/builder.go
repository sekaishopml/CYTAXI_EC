package pipeline

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/domain/candidate"
	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/domain/decision"
)

type CandidateBuilder struct {
	filters    []FilterStep
	scorers    []ScorerStep
}

type FilterStep interface {
	PipelineStep
	Filter(ctx context.Context, candidates []candidate.Candidate, ctxData decision.DecisionContext) ([]candidate.Candidate, error)
}

type ScorerStep interface {
	PipelineStep
	Score(ctx context.Context, candidate *candidate.Candidate, ctxData decision.DecisionContext) (float64, error)
}

func NewCandidateBuilder() *CandidateBuilder {
	return &CandidateBuilder{}
}

func (b *CandidateBuilder) AddFilter(f FilterStep) {
	b.filters = append(b.filters, f)
}

func (b *CandidateBuilder) AddScorer(s ScorerStep) {
	b.scorers = append(b.scorers, s)
}

func (b *CandidateBuilder) Build(ctx context.Context, raw []candidate.Candidate, ctxData decision.DecisionContext) (*candidate.CandidateSet, error) {
	current := raw

	for _, filter := range b.filters {
		filtered, err := filter.Filter(ctx, current, ctxData)
		if err != nil {
			return nil, err
		}
		current = filtered
	}

	set := &candidate.CandidateSet{}
	for _, c := range current {
		var totalScore float64
		for _, scorer := range b.scorers {
			score, err := scorer.Score(ctx, &c, ctxData)
			if err != nil {
				continue
			}
			totalScore += score
		}
		c.Score = totalScore
		c.Status = candidate.CandidateAvailable
		set.Add(c)
	}

	return set, nil
}

type ProximityFilter struct{}

func (f *ProximityFilter) Name() string { return "proximity_filter" }

func (f *ProximityFilter) Execute(ctx context.Context, candidates *candidate.CandidateSet, ctxData decision.DecisionContext) (*candidate.CandidateSet, error) {
	return candidates, nil
}

func (f *ProximityFilter) Filter(ctx context.Context, candidates []candidate.Candidate, ctxData decision.DecisionContext) ([]candidate.Candidate, error) {
	return candidates, nil
}

type AvailabilityFilter struct{}

func (f *AvailabilityFilter) Name() string { return "availability_filter" }

func (f *AvailabilityFilter) Execute(ctx context.Context, candidates *candidate.CandidateSet, ctxData decision.DecisionContext) (*candidate.CandidateSet, error) {
	return candidates, nil
}

func (f *AvailabilityFilter) Filter(ctx context.Context, candidates []candidate.Candidate, ctxData decision.DecisionContext) ([]candidate.Candidate, error) {
	return candidates, nil
}
