package domain

import "time"

type Decision struct {
	PolicyID    PolicyID        `json:"policy_id"`
	RuleID      RuleID          `json:"rule_id"`
	Actions     []Action        `json:"actions"`
	Matched     bool            `json:"matched"`
	Confidence  float64         `json:"confidence"`
	EvaluatedAt time.Time       `json:"evaluated_at"`
	Context     map[string]any  `json:"context,omitempty"`
}

type DecisionContext struct {
	UserID        string         `json:"user_id,omitempty"`
	Role          string         `json:"role,omitempty"`
	Phone         string         `json:"phone,omitempty"`
	Location      map[string]any `json:"location,omitempty"`
	TripData      map[string]any `json:"trip_data,omitempty"`
	DriverData    map[string]any `json:"driver_data,omitempty"`
	PricingData   map[string]any `json:"pricing_data,omitempty"`
	TimeOfDay     string         `json:"time_of_day,omitempty"`
	DayOfWeek     int            `json:"day_of_week,omitempty"`
	Custom        map[string]any `json:"custom,omitempty"`
}

func (ctx *DecisionContext) Get(key string) any {
	switch key {
	case "user_id":
		return ctx.UserID
	case "role":
		return ctx.Role
	case "phone":
		return ctx.Phone
	case "time_of_day":
		return ctx.TimeOfDay
	case "day_of_week":
		return ctx.DayOfWeek
	default:
		if ctx.Custom != nil {
			return ctx.Custom[key]
		}
		return nil
	}
}

type EvaluationResult struct {
	Decisions []Decision   `json:"decisions"`
	Matched   bool         `json:"matched"`
	PolicyIDs []PolicyID   `json:"policy_ids"`
	Errors    []error      `json:"errors,omitempty"`
}

func NewEvaluationResult() *EvaluationResult {
	return &EvaluationResult{
		Decisions: make([]Decision, 0),
		PolicyIDs: make([]PolicyID, 0),
	}
}

func (r *EvaluationResult) AddDecision(d Decision) {
	r.Decisions = append(r.Decisions, d)
	r.PolicyIDs = append(r.PolicyIDs, d.PolicyID)
	if d.Matched {
		r.Matched = true
	}
}

func (r *EvaluationResult) AddError(err error) {
	r.Errors = append(r.Errors, err)
}
