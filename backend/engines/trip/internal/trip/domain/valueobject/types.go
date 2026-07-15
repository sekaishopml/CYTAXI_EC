package valueobject

import (
	"fmt"
	"time"
)

type TripID string
type CustomerID string
type DriverID string
type StopID string

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Distance struct {
	Meters float64 `json:"meters"`
	KM     float64 `json:"km"`
}

type ETA struct {
	Duration time.Duration `json:"duration"`
	Seconds  int           `json:"seconds"`
}

type Money struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type TripStatus string

const (
	StatusRequested      TripStatus = "requested"
	StatusCreated        TripStatus = "created"
	StatusSearching      TripStatus = "searching"
	StatusDriverAssigned TripStatus = "driver_assigned"
	StatusAccepted       TripStatus = "accepted"
	StatusArrived        TripStatus = "arrived"
	StatusStarted        TripStatus = "started"
	StatusPaused         TripStatus = "paused"
	StatusResumed        TripStatus = "resumed"
	StatusCompleted      TripStatus = "completed"
	StatusCancelled      TripStatus = "cancelled"
)

func NewTripID() TripID { return TripID(fmt.Sprintf("trip_%d", time.Now().UnixNano())) }

func (c Coordinates) IsZero() bool { return c.Lat == 0 && c.Lng == 0 }

func NewDistance(meters float64) Distance {
	return Distance{Meters: meters, KM: meters / 1000}
}

func NewETA(seconds int) ETA {
	return ETA{Duration: time.Duration(seconds) * time.Second, Seconds: seconds}
}

func NewMoney(amount float64, currency string) Money {
	if currency == "" {
		currency = "USD"
	}
	return Money{Amount: amount, Currency: currency}
}

func (m Money) IsZero() bool { return m.Amount == 0 }
func (m Money) Multiply(factor float64) Money { return Money{Amount: m.Amount * factor, Currency: m.Currency} }

func IsValidTransition(from, to TripStatus) bool {
	return validateTransition(from, to) == nil
}

func validateTransition(from, to TripStatus) error {
	allowed := validTransitions[from]
	for _, s := range allowed {
		if s == to {
			return nil
		}
	}
	return fmt.Errorf("invalid transition %s → %s", from, to)
}

var validTransitions = map[TripStatus][]TripStatus{
	StatusRequested:      {StatusCreated, StatusCancelled},
	StatusCreated:        {StatusSearching, StatusCancelled},
	StatusSearching:      {StatusDriverAssigned, StatusCancelled},
	StatusDriverAssigned: {StatusAccepted, StatusCancelled},
	StatusAccepted:       {StatusArrived, StatusCancelled},
	StatusArrived:        {StatusStarted, StatusCancelled},
	StatusStarted:        {StatusPaused, StatusCompleted, StatusCancelled},
	StatusPaused:         {StatusResumed},
	StatusResumed:        {StatusPaused, StatusCompleted, StatusCancelled},
	StatusCompleted:      {},
	StatusCancelled:      {},
}
