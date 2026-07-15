package events

const (
	EventTripRequested       = "trip.requested"
	EventTripCreated         = "trip.created"
	EventDriverAssigned      = "trip.driver_assigned"
	EventDriverUnassigned    = "trip.driver_unassigned"
	EventTripAccepted        = "trip.accepted"
	EventTripRejected        = "trip.rejected"
	EventDriverArrived       = "trip.driver_arrived"
	EventTripStarted         = "trip.started"
	EventTripPaused          = "trip.paused"
	EventTripResumed         = "trip.resumed"
	EventTripCompleted       = "trip.completed"
	EventTripCancelled       = "trip.cancelled"
	EventDestinationChanged  = "trip.destination_changed"
	EventStopAdded           = "trip.stop_added"
	EventStopRemoved         = "trip.stop_removed"
)

type TripRequestedPayload struct {
	TripID        string  `json:"trip_id"`
	CustomerID    string  `json:"customer_id"`
	OriginLat     float64 `json:"origin_lat"`
	OriginLng     float64 `json:"origin_lng"`
	OriginAddress string  `json:"origin_address"`
	DestLat       float64 `json:"dest_lat"`
	DestLng       float64 `json:"dest_lng"`
	DestAddress   string  `json:"dest_address"`
}

type TripCreatedPayload struct {
	TripID        string `json:"trip_id"`
	CustomerID    string `json:"customer_id"`
	Status        string `json:"status"`
	OriginAddress string `json:"origin_address"`
	DestAddress   string `json:"dest_address"`
}

type DriverAssignedPayload struct {
	TripID    string  `json:"trip_id"`
	DriverID  string  `json:"driver_id"`
	Strategy  string  `json:"strategy"`
	Score     float64 `json:"score"`
	Distance  float64 `json:"distance_meters"`
	ETA       int     `json:"eta_seconds"`
}

type DriverUnassignedPayload struct {
	TripID   string `json:"trip_id"`
	DriverID string `json:"driver_id"`
	Reason   string `json:"reason"`
}

type TripAcceptedPayload struct {
	TripID    string `json:"trip_id"`
	DriverID  string `json:"driver_id"`
}

type TripRejectedPayload struct {
	TripID   string `json:"trip_id"`
	Reason   string `json:"reason"`
}

type TripStartedPayload struct {
	TripID    string `json:"trip_id"`
	DriverID  string `json:"driver_id"`
	StartedAt string `json:"started_at"`
}

type TripCompletedPayload struct {
	TripID      string  `json:"trip_id"`
	DistanceKM  float64 `json:"distance_km"`
	DurationSec int     `json:"duration_sec"`
	Fare        float64 `json:"fare"`
}

type TripCancelledPayload struct {
	TripID   string `json:"trip_id"`
	Reason   string `json:"reason"`
	By       string `json:"cancelled_by"`
}

type DestinationChangedPayload struct {
	TripID     string  `json:"trip_id"`
	NewLat     float64 `json:"new_lat"`
	NewLng     float64 `json:"new_lng"`
	NewAddress string  `json:"new_address"`
}

type StopAddedPayload struct {
	TripID   string  `json:"trip_id"`
	StopID   string  `json:"stop_id"`
	Address  string  `json:"address"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}
