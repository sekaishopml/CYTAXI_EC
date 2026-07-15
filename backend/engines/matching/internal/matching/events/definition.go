package events

const (
	EventMatchingStarted     = "matching.started"
	EventCandidatesFound     = "matching.candidates_found"
	EventCandidateRejected   = "matching.candidate_rejected"
	EventCandidateAccepted   = "matching.candidate_accepted"
	EventDriverAssigned      = "matching.driver_assigned"
	EventDriverDeclined      = "matching.driver_declined"
	EventMatchingRetried     = "matching.retried"
	EventMatchingCancelled   = "matching.cancelled"
	EventMatchingCompleted   = "matching.completed"
)

type MatchingStartedPayload struct {
	MatchingID string  `json:"matching_id"`
	TripID     string  `json:"trip_id"`
	PickupLat  float64 `json:"pickup_lat"`
	PickupLng  float64 `json:"pickup_lng"`
	Strategy   string  `json:"strategy"`
}

type CandidatesFoundPayload struct {
	MatchingID   string `json:"matching_id"`
	Total        int    `json:"total"`
	RadiusMeters float64 `json:"radius_meters"`
}

type CandidateRejectedPayload struct {
	MatchingID string  `json:"matching_id"`
	DriverID   string  `json:"driver_id"`
	Reason     string  `json:"reason"`
	Score      float64 `json:"score"`
}

type DriverAssignedPayload struct {
	MatchingID string  `json:"matching_id"`
	TripID     string  `json:"trip_id"`
	DriverID   string  `json:"driver_id"`
	Score      float64 `json:"score"`
	Distance   float64 `json:"distance_meters"`
	ETA        int     `json:"eta_seconds"`
}

type DriverDeclinedPayload struct {
	MatchingID string `json:"matching_id"`
	DriverID   string `json:"driver_id"`
	Reason     string `json:"reason"`
}

type MatchingCancelledPayload struct {
	MatchingID string `json:"matching_id"`
	Reason     string `json:"reason"`
}

type MatchingCompletedPayload struct {
	MatchingID string  `json:"matching_id"`
	DriverID   string  `json:"driver_id"`
	Score      float64 `json:"score"`
	Attempts   int     `json:"attempts"`
}
