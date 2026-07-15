package dispatcher

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/domain/candidate"
	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/domain/decision"
	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/application/pipeline"
)

type DispatcherCoordinator struct {
	pipeline       *pipeline.DecisionPipeline
	candidateFinder CandidateFinder
	logger         *slog.Logger
}

type CandidateFinder interface {
	FindCandidates(ctx context.Context, ctxData decision.DecisionContext) (*candidate.CandidateSet, error)
}

func NewDispatcherCoordinator(
	pipeline *pipeline.DecisionPipeline,
	candidateFinder CandidateFinder,
	logger *slog.Logger,
) *DispatcherCoordinator {
	return &DispatcherCoordinator{
		pipeline:       pipeline,
		candidateFinder: candidateFinder,
		logger:         logger,
	}
}

func (c *DispatcherCoordinator) Dispatch(ctx context.Context, ctxData decision.DecisionContext) (*decision.Decision, error) {
	start := time.Now()
	c.logger.Info("dispatch started", "trip_id", ctxData.TripID)

	if ctxData.TripID == "" {
		return nil, fmt.Errorf("%w: trip_id required", decision.ErrInvalidContext)
	}

	d := decision.NewDecision(ctxData.TripID)

	candidates, err := c.candidateFinder.FindCandidates(ctx, ctxData)
	if err != nil {
		d.Status = decision.DecisionFailed
		return d, fmt.Errorf("find candidates: %w", err)
	}

	if len(candidates.Candidates) == 0 {
		d.Status = decision.DecisionRejected
		return d, fmt.Errorf("%w for trip %s", decision.ErrNoCandidates, ctxData.TripID)
	}

	result, err := c.pipeline.Execute(ctx, ctxData, candidates)
	if err != nil {
		d.Status = decision.DecisionFailed
		c.logger.Error("pipeline failed", "trip_id", ctxData.TripID, "error", err)
		return d, fmt.Errorf("pipeline: %w", err)
	}

	d.Status = decision.DecisionApproved
	d.SelectedDriver = result.DriverID
	d.StrategyUsed = result.Strategy
	d.Score = result.Score
	d.PipelineSummary = result.Steps

	c.logger.Info("dispatch completed",
		"trip_id", ctxData.TripID,
		"driver", result.DriverID,
		"strategy", result.Strategy,
		"score", result.Score,
		"duration", time.Since(start).String(),
	)

	return d, nil
}

func (c *DispatcherCoordinator) CancelDispatch(ctx context.Context, tripID string) error {
	c.logger.Info("dispatch cancelled", "trip_id", tripID)
	return nil
}
