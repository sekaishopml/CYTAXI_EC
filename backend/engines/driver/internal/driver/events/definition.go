package events

const (
	EventDriverCreated       = "driver.created"
	EventDriverApproved      = "driver.approved"
	EventDriverRejected      = "driver.rejected"
	EventDriverOnline        = "driver.online"
	EventDriverOffline       = "driver.offline"
	EventDriverSuspended     = "driver.suspended"
	EventVehicleUpdated      = "driver.vehicle_updated"
	EventLicenseUpdated      = "driver.license_updated"
	EventAvailabilityChanged = "driver.availability_changed"
	EventCapabilitiesChanged = "driver.capabilities_changed"
	EventPreferencesUpdated  = "driver.preferences_updated"
	EventDocumentUploaded    = "driver.document_uploaded"
)

type DriverCreatedPayload struct {
	DriverID string `json:"driver_id"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
}

type DriverApprovedPayload struct {
	DriverID string `json:"driver_id"`
}

type DriverRejectedPayload struct {
	DriverID string `json:"driver_id"`
	Reason   string `json:"reason"`
}

type DriverOnlinePayload struct {
	DriverID string  `json:"driver_id"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

type DriverOfflinePayload struct {
	DriverID string `json:"driver_id"`
}

type DriverSuspendedPayload struct {
	DriverID string `json:"driver_id"`
	Reason   string `json:"reason"`
}

type VehicleUpdatedPayload struct {
	DriverID  string `json:"driver_id"`
	VehicleID string `json:"vehicle_id"`
	Plate     string `json:"plate"`
	Type      string `json:"type"`
}

type LicenseUpdatedPayload struct {
	DriverID    string `json:"driver_id"`
	LicenseNumber string `json:"license_number"`
	Category    string `json:"category"`
	ExpiresAt   string `json:"expires_at"`
}

type AvailabilityChangedPayload struct {
	DriverID string  `json:"driver_id"`
	Status   string  `json:"status"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

type CapabilitiesChangedPayload struct {
	DriverID     string   `json:"driver_id"`
	Capabilities []string `json:"capabilities"`
}

type PreferencesUpdatedPayload struct {
	DriverID string `json:"driver_id"`
	Changes  any    `json:"changes"`
}

type DocumentUploadedPayload struct {
	DriverID string `json:"driver_id"`
	DocType  string `json:"doc_type"`
	DocID    string `json:"doc_id"`
}
