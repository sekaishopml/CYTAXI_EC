package candidate

import "time"

type AssignmentAttempt struct {
	ID         string    `json:"id"`
	MatchingID string    `json:"matching_id"`
	DriverID   string    `json:"driver_id"`
	Attempt    int       `json:"attempt"`
	Status     AttemptStatus `json:"status"`
	Reason     string    `json:"reason,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type AttemptStatus string

const (
	AttemptDelivered AttemptStatus = "delivered"
	AttemptAccepted  AttemptStatus = "accepted"
	AttemptDeclined  AttemptStatus = "declined"
	AttemptTimeout   AttemptStatus = "timeout"
)

type AssignmentResult struct {
	MatchingID string       `json:"matching_id"`
	TripID     string       `json:"trip_id"`
	DriverID   string       `json:"driver_id"`
	Success    bool         `json:"success"`
	Strategy   string       `json:"strategy"`
	Attempts   int          `json:"attempts"`
	FinalScore float64      `json:"final_score"`
	CreatedAt  time.Time    `json:"created_at"`
}
