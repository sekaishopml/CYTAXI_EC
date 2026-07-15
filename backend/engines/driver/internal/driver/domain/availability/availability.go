package availability

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
)

type Status string

const (
	Available   Status = "available"
	Busy        Status = "busy"
	OnTrip      Status = "on_trip"
	Break       Status = "break"
	Unavailable Status = "unavailable"
)

type DriverAvailability struct {
	DriverID    valueobject.DriverID    `json:"driver_id"`
	Status      Status                  `json:"status"`
	Location    valueobject.Coordinates `json:"location"`
	LastPing    time.Time               `json:"last_ping"`
	ChangedAt   time.Time               `json:"changed_at"`
}

func NewAvailability(driverID valueobject.DriverID) *DriverAvailability {
	now := time.Now()
	return &DriverAvailability{
		DriverID:  driverID,
		Status:    Unavailable,
		LastPing:  now,
		ChangedAt: now,
	}
}

func (a *DriverAvailability) SetAvailable(loc valueobject.Coordinates) {
	a.Status = Available
	a.Location = loc
	a.LastPing = time.Now()
	a.ChangedAt = time.Now()
}

func (a *DriverAvailability) SetBusy() {
	a.Status = Busy
	a.ChangedAt = time.Now()
}

func (a *DriverAvailability) SetOnTrip() {
	a.Status = OnTrip
	a.ChangedAt = time.Now()
}

func (a *DriverAvailability) SetBreak() {
	a.Status = Break
	a.ChangedAt = time.Now()
}

func (a *DriverAvailability) Ping(loc valueobject.Coordinates) {
	a.Location = loc
	a.LastPing = time.Now()
}
