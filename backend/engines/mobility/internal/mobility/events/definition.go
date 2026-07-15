package events

const (
	EventDispatchStarted    = "mobility.dispatch_started"
	EventDispatchCompleted  = "mobility.dispatch_completed"
	EventDispatchFailed     = "mobility.dispatch_failed"
	EventDispatchCancelled  = "mobility.dispatch_cancelled"
	EventCandidateFound     = "mobility.candidate_found"
	EventCandidateFiltered  = "mobility.candidate_filtered"
	EventCandidateSelected  = "mobility.candidate_selected"
	EventPipelineStepPassed = "mobility.pipeline_step_passed"
	EventPipelineStepFailed = "mobility.pipeline_step_failed"
	EventAssignmentSent     = "mobility.assignment_sent"
	EventAssignmentAccepted = "mobility.assignment_accepted"
	EventAssignmentRejected = "mobility.assignment_rejected"
	EventNoCandidates       = "mobility.no_candidates"
	EventStrategyApplied    = "mobility.strategy_applied"
)

type DispatchStartedPayload struct {
	TripID    string `json:"trip_id"`
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
}

type DispatchCompletedPayload struct {
	TripID   string  `json:"trip_id"`
	DriverID string  `json:"driver_id"`
	Strategy string  `json:"strategy"`
	Score    float64 `json:"score"`
}

type CandidateSelectedPayload struct {
	TripID   string  `json:"trip_id"`
	DriverID string  `json:"driver_id"`
	Strategy string  `json:"strategy"`
	Score    float64 `json:"score"`
	Distance float64 `json:"distance_meters"`
}

type PipelineStepPayload struct {
	TripID  string `json:"trip_id"`
	Step    string `json:"step"`
	Status  string `json:"status"`
	Elapsed int64  `json:"duration_ms"`
}
