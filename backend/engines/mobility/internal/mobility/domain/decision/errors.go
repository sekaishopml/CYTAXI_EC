package decision

import "errors"

var (
	ErrNoCandidates      = errors.New("dispatch: no candidates available")
	ErrStrategyNotFound  = errors.New("dispatch: strategy not found")
	ErrPipelineFailed    = errors.New("dispatch: pipeline execution failed")
	ErrAssignmentTimeout = errors.New("dispatch: assignment timed out")
	ErrInvalidContext    = errors.New("dispatch: invalid decision context")
	ErrDuplicateDecision = errors.New("dispatch: duplicate decision for trip")
)
