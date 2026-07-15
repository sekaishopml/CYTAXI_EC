package assignment

import "time"

type AssignmentID string

type Assignment struct {
	ID            AssignmentID    `json:"id"`
	TripID        string          `json:"trip_id"`
	DriverID      string          `json:"driver_id"`
	Status        AssignmentStatus `json:"status"`
	Strategy      string          `json:"strategy"`
	Score         float64         `json:"score"`
	Reason        string          `json:"reason"`
	AssignedAt    time.Time       `json:"assigned_at"`
	CompletedAt   *time.Time      `json:"completed_at,omitempty"`
	PipelineSteps []StepResult    `json:"pipeline_steps,omitempty"`
}

type AssignmentStatus string

const (
	AssignmentPending    AssignmentStatus = "pending"
	AssignmentProposed   AssignmentStatus = "proposed"
	AssignmentAccepted   AssignmentStatus = "accepted"
	AssignmentRejected   AssignmentStatus = "rejected"
	AssignmentExpired    AssignmentStatus = "expired"
	AssignmentCompleted  AssignmentStatus = "completed"
	AssignmentCancelled  AssignmentStatus = "cancelled"
)

type StepResult struct {
	Name      string  `json:"name"`
	Passed    bool    `json:"passed"`
	Score     float64 `json:"score"`
	Reason    string  `json:"reason,omitempty"`
}

type AssignmentResponse struct {
	Assignment *Assignment `json:"assignment,omitempty"`
	Success    bool        `json:"success"`
	Error      string      `json:"error,omitempty"`
}
