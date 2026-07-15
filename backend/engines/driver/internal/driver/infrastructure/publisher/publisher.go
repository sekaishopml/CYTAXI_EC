package publisher

import "context"

type DriverEventPublisher interface {
	PublishDriverCreated(ctx context.Context, driverID, phone, name string) error
	PublishDriverApproved(ctx context.Context, driverID string) error
	PublishDriverRejected(ctx context.Context, driverID, reason string) error
	PublishDriverOnline(ctx context.Context, driverID string, lat, lng float64) error
	PublishDriverOffline(ctx context.Context, driverID string) error
	PublishDriverSuspended(ctx context.Context, driverID, reason string) error
	PublishVehicleUpdated(ctx context.Context, driverID, vehicleID, plate, vehicleType string) error
	PublishLicenseUpdated(ctx context.Context, driverID, licenseNumber, category, expiresAt string) error
	PublishAvailabilityChanged(ctx context.Context, driverID, status string, lat, lng float64) error
	PublishCapabilitiesChanged(ctx context.Context, driverID string, capabilities []string) error
	PublishPreferencesUpdated(ctx context.Context, driverID string, changes any) error
	PublishDocumentUploaded(ctx context.Context, driverID, docType, docID string) error
}

type LogPublisher struct{}

func NewLogPublisher() *LogPublisher {
	return &LogPublisher{}
}

func (p *LogPublisher) PublishDriverCreated(ctx context.Context, driverID, phone, name string) error {
	return nil
}

func (p *LogPublisher) PublishDriverApproved(ctx context.Context, driverID string) error {
	return nil
}

func (p *LogPublisher) PublishDriverRejected(ctx context.Context, driverID, reason string) error {
	return nil
}

func (p *LogPublisher) PublishDriverOnline(ctx context.Context, driverID string, lat, lng float64) error {
	return nil
}

func (p *LogPublisher) PublishDriverOffline(ctx context.Context, driverID string) error {
	return nil
}

func (p *LogPublisher) PublishDriverSuspended(ctx context.Context, driverID, reason string) error {
	return nil
}

func (p *LogPublisher) PublishVehicleUpdated(ctx context.Context, driverID, vehicleID, plate, vehicleType string) error {
	return nil
}

func (p *LogPublisher) PublishLicenseUpdated(ctx context.Context, driverID, licenseNumber, category, expiresAt string) error {
	return nil
}

func (p *LogPublisher) PublishAvailabilityChanged(ctx context.Context, driverID, status string, lat, lng float64) error {
	return nil
}

func (p *LogPublisher) PublishCapabilitiesChanged(ctx context.Context, driverID string, capabilities []string) error {
	return nil
}

func (p *LogPublisher) PublishPreferencesUpdated(ctx context.Context, driverID string, changes any) error {
	return nil
}

func (p *LogPublisher) PublishDocumentUploaded(ctx context.Context, driverID, docType, docID string) error {
	return nil
}
