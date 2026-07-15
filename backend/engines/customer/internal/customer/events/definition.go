package events

const (
	EventCustomerCreated       = "customer.created"
	EventCustomerUpdated       = "customer.updated"
	EventCustomerBlocked       = "customer.blocked"
	EventProfileUpdated        = "customer.profile_updated"
	EventPreferencesUpdated    = "customer.preferences_updated"
	EventFavoritePlaceAdded    = "customer.favorite_place_added"
	EventFavoritePlaceRemoved  = "customer.favorite_place_removed"
	EventCustomerContextChanged = "customer.context_changed"
)

type CustomerCreatedPayload struct {
	CustomerID string `json:"customer_id"`
	Phone      string `json:"phone"`
	Name       string `json:"name,omitempty"`
	Timestamp  string `json:"timestamp"`
}

type CustomerUpdatedPayload struct {
	CustomerID string `json:"customer_id"`
	Changes    any    `json:"changes"`
}

type PreferencesUpdatedPayload struct {
	CustomerID string `json:"customer_id"`
	Changes    any    `json:"changes"`
}

type FavoritePlaceAddedPayload struct {
	CustomerID string  `json:"customer_id"`
	PlaceID    string  `json:"place_id"`
	Name       string  `json:"name"`
	Category   string  `json:"category,omitempty"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
}

type CustomerContextChangedPayload struct {
	CustomerID string `json:"customer_id"`
	Changes    any    `json:"changes"`
}
