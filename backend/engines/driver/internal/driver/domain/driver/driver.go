package driver

import (
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
)

type Driver struct {
	ID         valueobject.DriverID `json:"id"`
	Phone      valueobject.Phone    `json:"phone"`
	Email      valueobject.Email    `json:"email,omitempty"`
	Name       string               `json:"name"`
	Status     DriverStatus         `json:"status"`
	Rating     float64              `json:"rating"`
	TripsCount int                  `json:"trips_count"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

type DriverStatus string

const (
	DriverPending   DriverStatus = "pending"
	DriverApproved  DriverStatus = "approved"
	DriverRejected  DriverStatus = "rejected"
	DriverOnline    DriverStatus = "online"
	DriverOffline   DriverStatus = "offline"
	DriverSuspended DriverStatus = "suspended"
)

func NewDriver(phone valueobject.Phone, name string) *Driver {
	now := time.Now()
	return &Driver{
		ID:        valueobject.DriverID(fmt.Sprintf("drv_%d", now.UnixNano())),
		Phone:     phone,
		Name:      name,
		Status:    DriverPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (d *Driver) Approve()  { d.Status = DriverApproved; d.UpdatedAt = time.Now() }
func (d *Driver) Reject()   { d.Status = DriverRejected; d.UpdatedAt = time.Now() }
func (d *Driver) GoOnline() { d.Status = DriverOnline; d.UpdatedAt = time.Now() }
func (d *Driver) GoOffline() { d.Status = DriverOffline; d.UpdatedAt = time.Now() }
func (d *Driver) Suspend()  { d.Status = DriverSuspended; d.UpdatedAt = time.Now() }
