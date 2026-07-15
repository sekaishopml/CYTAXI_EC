package capability

import (
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
)

type Capability string

const (
	CapBabySeat       Capability = "baby_seat"
	CapWheelchair     Capability = "wheelchair"
	CapXLPassengers   Capability = "xl_passengers"
	CapPremiumVehicle Capability = "premium_vehicle"
	CapElectric       Capability = "electric"
	CapCashPayment    Capability = "cash_payment"
	CapCardPayment    Capability = "card_payment"
	CapPetsAllowed    Capability = "pets_allowed"
)

type DriverCapabilities struct {
	DriverID     valueobject.DriverID `json:"driver_id"`
	Capabilities []Capability         `json:"capabilities"`
	MaxPass      int                  `json:"max_passengers"`
	VehicleType  string               `json:"vehicle_type"`
}

func NewDriverCapabilities(driverID valueobject.DriverID) *DriverCapabilities {
	return &DriverCapabilities{
		DriverID: driverID,
		MaxPass:  4,
	}
}

func (c *DriverCapabilities) Has(cap Capability) bool {
	for _, v := range c.Capabilities {
		if v == cap {
			return true
		}
	}
	return false
}

func (c *DriverCapabilities) Add(cap Capability) {
	if !c.Has(cap) {
		c.Capabilities = append(c.Capabilities, cap)
	}
}
