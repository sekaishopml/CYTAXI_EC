package vehicle

import (
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
)

type Vehicle struct {
	ID         valueobject.VehicleID   `json:"id"`
	DriverID   valueobject.DriverID    `json:"driver_id"`
	Plate      valueobject.PlateNumber `json:"plate"`
	Brand      string                  `json:"brand"`
	Model      string                  `json:"model"`
	Year       int                     `json:"year"`
	Color      string                  `json:"color"`
	Type       VehicleType             `json:"type"`
	MaxPass    int                     `json:"max_passengers"`
	BabySeat   bool                    `json:"baby_seat"`
	Wheelchair bool                    `json:"wheelchair"`
	Active     bool                    `json:"active"`
	CreatedAt  time.Time               `json:"created_at"`
	UpdatedAt  time.Time               `json:"updated_at"`
}

type VehicleType string

const (
	TypeStandard VehicleType = "standard"
	TypeXL       VehicleType = "xl"
	TypePremium  VehicleType = "premium"
	TypeElectric VehicleType = "electric"
)

func NewVehicle(driverID valueobject.DriverID, plate valueobject.PlateNumber) *Vehicle {
	now := time.Now()
	return &Vehicle{
		ID:        valueobject.VehicleID(fmt.Sprintf("veh_%d", now.UnixNano())),
		DriverID:  driverID,
		Plate:     plate,
		Active:    true,
		Type:      TypeStandard,
		MaxPass:   4,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
