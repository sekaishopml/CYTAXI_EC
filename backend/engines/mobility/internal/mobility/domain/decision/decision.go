package decision

import "time"

type DecisionID string

type Decision struct {
	ID              DecisionID      `json:"id"`
	TripID          string          `json:"trip_id"`
	Status          DecisionStatus  `json:"status"`
	SelectedDriver  string          `json:"selected_driver,omitempty"`
	StrategyUsed    string          `json:"strategy_used"`
	Score           float64         `json:"score"`
	PipelineSummary []StepSummary   `json:"pipeline_summary,omitempty"`
	DecidedAt       time.Time       `json:"decided_at"`
}

type DecisionStatus string

const (
	DecisionPending    DecisionStatus = "pending"
	DecisionApproved   DecisionStatus = "approved"
	DecisionRejected   DecisionStatus = "rejected"
	DecisionFallback   DecisionStatus = "fallback"
	DecisionFailed     DecisionStatus = "failed"
)

type StepSummary struct {
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	Duration  int64   `json:"duration_ms"`
}

type DecisionContext struct {
	TripID         string            `json:"trip_id"`
	UserID         string            `json:"user_id"`
	Origin         map[string]float64 `json:"origin"`
	Destination    map[string]float64 `json:"destination"`
	VehicleType    string            `json:"vehicle_type"`
	RequiresBaby   bool              `json:"requires_baby"`
	RequiresWheel  bool              `json:"requires_wheel"`
	MaxPassengers  int               `json:"max_passengers"`
	PricingData    map[string]any    `json:"pricing_data,omitempty"`
	PolicyData     map[string]any    `json:"policy_data,omitempty"`
	Custom         map[string]any    `json:"custom,omitempty"`
}

func NewDecision(tripID string) *Decision {
	return &Decision{
		ID:        DecisionID(tripID + "_" + id()),
		TripID:    tripID,
		Status:    DecisionPending,
		DecidedAt: time.Now(),
	}
}

func id() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
