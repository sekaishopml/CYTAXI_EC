package preference

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
)

type Preferences struct {
	DriverID    valueobject.DriverID `json:"driver_id"`
	MaxDistance int                  `json:"max_distance"` // meters
	MinFare     float64              `json:"min_fare"`
	AutoAccept  bool                 `json:"auto_accept"`
	RadiusKm    float64              `json:"radius_km"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

func NewPreferences(driverID valueobject.DriverID) *Preferences {
	return &Preferences{
		DriverID:    driverID,
		MaxDistance: 5000,
		RadiusKm:    5.0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
