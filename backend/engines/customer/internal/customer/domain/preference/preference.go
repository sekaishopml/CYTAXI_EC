package preference

import "time"

type Preferences struct {
	CustomerID       string    `json:"customer_id"`
	VehicleType      string    `json:"vehicle_type,omitempty"`
	MaxPassengers    int       `json:"max_passengers,omitempty"`
	RequiresBabySeat bool      `json:"requires_baby_seat"`
	RequiresWheelchair bool    `json:"requires_wheelchair"`
	PaymentMethod    string    `json:"payment_method,omitempty"`
	Notifications    bool      `json:"notifications"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewPreferences(customerID string) *Preferences {
	return &Preferences{
		CustomerID:    customerID,
		Notifications: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

type PreferenceUpdate struct {
	VehicleType      *string `json:"vehicle_type,omitempty"`
	MaxPassengers    *int    `json:"max_passengers,omitempty"`
	RequiresBabySeat *bool   `json:"requires_baby_seat,omitempty"`
	RequiresWheelchair *bool `json:"requires_wheelchair,omitempty"`
	PaymentMethod    *string `json:"payment_method,omitempty"`
	Notifications    *bool   `json:"notifications,omitempty"`
}

func (p *Preferences) Apply(update PreferenceUpdate) {
	if update.VehicleType != nil {
		p.VehicleType = *update.VehicleType
	}
	if update.MaxPassengers != nil {
		p.MaxPassengers = *update.MaxPassengers
	}
	if update.RequiresBabySeat != nil {
		p.RequiresBabySeat = *update.RequiresBabySeat
	}
	if update.RequiresWheelchair != nil {
		p.RequiresWheelchair = *update.RequiresWheelchair
	}
	if update.PaymentMethod != nil {
		p.PaymentMethod = *update.PaymentMethod
	}
	if update.Notifications != nil {
		p.Notifications = *update.Notifications
	}
	p.UpdatedAt = time.Now()
}
