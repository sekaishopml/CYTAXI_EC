package context

import "time"

type CustomerContext struct {
	CustomerID   string            `json:"customer_id"`
	Phone        string            `json:"phone"`
	Name         string            `json:"name"`
	Preferences  map[string]any    `json:"preferences"`
	RecentTrips  int               `json:"recent_trips"`
	LastTripAt   *time.Time        `json:"last_trip_at,omitempty"`
	Custom       map[string]any    `json:"custom,omitempty"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func NewCustomerContext(customerID, phone string) *CustomerContext {
	return &CustomerContext{
		CustomerID:  customerID,
		Phone:       phone,
		Preferences: make(map[string]any),
		Custom:      make(map[string]any),
		UpdatedAt:   time.Now(),
	}
}

func (c *CustomerContext) SetPreference(key string, value any) {
	c.Preferences[key] = value
	c.UpdatedAt = time.Now()
}

func (c *CustomerContext) GetPreference(key string) (any, bool) {
	v, ok := c.Preferences[key]
	return v, ok
}

func (c *CustomerContext) SetCustom(key string, value any) {
	c.Custom[key] = value
	c.UpdatedAt = time.Now()
}
