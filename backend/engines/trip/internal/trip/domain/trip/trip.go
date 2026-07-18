package trip

import (
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/destination"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/passenger"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/stop"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/valueobject"
)

type Trip struct {
	ID           valueobject.TripID      `json:"id"`
	CustomerID   valueobject.CustomerID  `json:"customer_id"`
	DriverID     valueobject.DriverID    `json:"driver_id,omitempty"`
	Status       valueobject.TripStatus  `json:"status"`
	Passenger    passenger.Passenger     `json:"passenger"`
	Pickup       stop.Stop               `json:"pickup"`
	Destination  destination.Destination `json:"destination"`
	Stops        []stop.Stop             `json:"stops,omitempty"`
	Estimate     TripEstimate            `json:"estimate,omitempty"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
	StartedAt    *time.Time              `json:"started_at,omitempty"`
	CompletedAt  *time.Time              `json:"completed_at,omitempty"`
	CancelledAt  *time.Time              `json:"cancelled_at,omitempty"`
	CancelReason string                  `json:"cancel_reason,omitempty"`
}

type TripEstimate struct {
	Distance valueobject.Distance `json:"distance"`
	ETA      valueobject.ETA      `json:"eta"`
	Fare     valueobject.Money    `json:"fare"`
}

func NewTrip(customerID valueobject.CustomerID, p passenger.Passenger, pickup stop.Stop, dest destination.Destination) *Trip {
	now := time.Now()
	return &Trip{
		ID:          valueobject.NewTripID(),
		CustomerID:  customerID,
		Status:      valueobject.StatusCreated,
		Passenger:   p,
		Pickup:      pickup,
		Destination: dest,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *Trip) AssignDriver(driverID valueobject.DriverID) error {
	if t.Status != valueobject.StatusSearching {
		return fmt.Errorf("cannot assign driver in status %s", t.Status)
	}
	t.DriverID = driverID
	t.Status = valueobject.StatusDriverAssigned
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Trip) UnassignDriver() error {
	t.DriverID = ""
	t.Status = valueobject.StatusSearching
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Trip) StartSearching() error {
	if t.Status != valueobject.StatusCreated {
		return fmt.Errorf("cannot search in status %s", t.Status)
	}
	t.Status = valueobject.StatusSearching
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Trip) AcceptTrip() error {
	if t.Status != valueobject.StatusDriverAssigned {
		return fmt.Errorf("cannot accept in status %s", t.Status)
	}
	t.Status = valueobject.StatusAccepted
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Trip) RejectTrip(reason string) error {
	if t.Status != valueobject.StatusDriverAssigned {
		return fmt.Errorf("cannot reject in status %s", t.Status)
	}
	t.Status = valueobject.StatusSearching
	t.DriverID = ""
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Trip) DriverArrived() error {
	if t.Status != valueobject.StatusAccepted {
		return fmt.Errorf("driver cannot arrive in status %s", t.Status)
	}
	t.Status = valueobject.StatusArrived
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Trip) StartTrip() error {
	if t.Status != valueobject.StatusArrived {
		return fmt.Errorf("cannot start in status %s", t.Status)
	}
	now := time.Now()
	t.Status = valueobject.StatusStarted
	t.StartedAt = &now
	t.UpdatedAt = now
	return nil
}

func (t *Trip) PauseTrip() error {
	if t.Status != valueobject.StatusStarted {
		return fmt.Errorf("cannot pause in status %s", t.Status)
	}
	t.Status = valueobject.StatusPaused
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Trip) ResumeTrip() error {
	if t.Status != valueobject.StatusPaused {
		return fmt.Errorf("cannot resume in status %s", t.Status)
	}
	t.Status = valueobject.StatusResumed
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Trip) CompleteTrip() error {
	if t.Status != valueobject.StatusStarted && t.Status != valueobject.StatusResumed {
		return fmt.Errorf("cannot complete in status %s", t.Status)
	}
	now := time.Now()
	t.Status = valueobject.StatusCompleted
	t.CompletedAt = &now
	t.UpdatedAt = now
	return nil
}

func (t *Trip) CancelTrip(reason string) error {
	if t.Status == valueobject.StatusCompleted || t.Status == valueobject.StatusCancelled {
		return fmt.Errorf("cannot cancel in status %s", t.Status)
	}
	now := time.Now()
	t.Status = valueobject.StatusCancelled
	t.CancelledAt = &now
	t.CancelReason = reason
	t.UpdatedAt = now
	return nil
}

func (t *Trip) AddStop(s stop.Stop) error {
	if t.Status != valueobject.StatusStarted {
		return fmt.Errorf("cannot add stop in status %s", t.Status)
	}
	t.Stops = append(t.Stops, s)
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Trip) RemoveStop(stopID valueobject.StopID) error {
	if t.Status != valueobject.StatusStarted {
		return fmt.Errorf("cannot remove stop in status %s", t.Status)
	}
	for i, s := range t.Stops {
		if s.ID == stopID {
			t.Stops = append(t.Stops[:i], t.Stops[i+1:]...)
			t.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("stop %s not found", stopID)
}

func (t *Trip) ChangeDestination(dest destination.Destination) error {
	if t.Status == valueobject.StatusCompleted || t.Status == valueobject.StatusCancelled {
		return fmt.Errorf("cannot change destination in status %s", t.Status)
	}
	t.Destination = dest
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Trip) SetEstimate(est TripEstimate) {
	t.Estimate = est
	t.UpdatedAt = time.Now()
}

// IsActive returns true if trip is in an active state
func (t *Trip) IsActive() bool {
	return t.Status != valueobject.StatusCompleted && t.Status != valueobject.StatusCancelled
}
