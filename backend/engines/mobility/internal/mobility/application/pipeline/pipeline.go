package pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/domain/candidate"
	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/domain/decision"
)

type PipelineStep interface {
	Name() string
	Execute(ctx context.Context, candidates *candidate.CandidateSet, ctxData decision.DecisionContext) (*candidate.CandidateSet, error)
}

type Strategy interface {
	Name() string
	Select(ctx context.Context, candidates *candidate.CandidateSet, ctxData decision.DecisionContext) (*candidate.Candidate, error)
}

type PipelineResult struct {
	DriverID string        `json:"driver_id"`
	Strategy string        `json:"strategy"`
	Score    float64       `json:"score"`
	Steps    []decision.StepSummary `json:"steps"`
}

type DecisionPipeline struct {
	steps    []PipelineStep
	strategy Strategy
}

func NewDecisionPipeline(strategy Strategy) *DecisionPipeline {
	return &DecisionPipeline{
		steps:    make([]PipelineStep, 0),
		strategy: strategy,
	}
}

func (p *DecisionPipeline) AddStep(step PipelineStep) {
	p.steps = append(p.steps, step)
}

func (p *DecisionPipeline) Execute(ctx context.Context, ctxData decision.DecisionContext, candidates *candidate.CandidateSet) (*PipelineResult, error) {
	current := candidates
	var steps []decision.StepSummary

	for _, step := range p.steps {
		start := time.Now()
		result, err := step.Execute(ctx, current, ctxData)
		elapsed := time.Since(start).Milliseconds()

		status := "passed"
		if err != nil {
			status = "failed"
		}

		steps = append(steps, decision.StepSummary{
			Name:     step.Name(),
			Status:   status,
			Duration: elapsed,
		})

		if err != nil {
			return &PipelineResult{
				Steps: steps,
			}, fmt.Errorf("pipeline step %s: %w", step.Name(), err)
		}

		current = result
	}

	selected, err := p.strategy.Select(ctx, current, ctxData)
	if err != nil {
		return &PipelineResult{Steps: steps}, fmt.Errorf("strategy %s: %w", p.strategy.Name(), err)
	}

	return &PipelineResult{
		DriverID: selected.DriverID,
		Strategy: p.strategy.Name(),
		Score:    selected.Score,
		Steps:    steps,
	}, nil
}
