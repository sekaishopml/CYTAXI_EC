package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/domain/candidate"
	"github.com/sekaishopml/cytaxi/backend/engines/mobility/internal/mobility/domain/decision"
)

type DispatcherInputPort interface {
	Dispatch(ctx context.Context, ctxData decision.DecisionContext) (*decision.Decision, error)
	CancelDispatch(ctx context.Context, tripID string) error
}

type DispatcherOutputPort interface {
	OnDispatchStarted(ctx context.Context, tripID string) error
	OnDispatchCompleted(ctx context.Context, d decision.Decision) error
	OnDispatchFailed(ctx context.Context, tripID string, err error) error
}

type CandidateInputPort interface {
	FindCandidates(ctx context.Context, ctxData decision.DecisionContext) (*candidate.CandidateSet, error)
}

type AssignmentInputPort interface {
	Assign(ctx context.Context, driverID string, tripID string) error
	ConfirmAssignment(ctx context.Context, driverID string, tripID string) error
	RejectAssignment(ctx context.Context, driverID string, tripID string, reason string) error
}
