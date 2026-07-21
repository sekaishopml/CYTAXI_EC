package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/availability"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/driver"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/license"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/vehicle"
)

type DriverInputPort interface {
	Register(ctx context.Context, phone valueobject.Phone, name string) (*driver.Driver, error)
	GetDriver(ctx context.Context, id valueobject.DriverID) (*driver.Driver, error)
	Approve(ctx context.Context, id valueobject.DriverID) error
	Reject(ctx context.Context, id valueobject.DriverID, reason string) error
	Suspend(ctx context.Context, id valueobject.DriverID, reason string) error
}

type AvailabilityInputPort interface {
	GoOnline(ctx context.Context, id valueobject.DriverID, loc valueobject.Coordinates) error
	GoOffline(ctx context.Context, id valueobject.DriverID) error
	Ping(ctx context.Context, id valueobject.DriverID, loc valueobject.Coordinates) error
	GetAvailability(ctx context.Context, id valueobject.DriverID) (*availability.DriverAvailability, error)
	FindAvailableNearby(ctx context.Context, lat, lng, radiusMeters float64) ([]availability.DriverAvailability, error)
}

type VehicleInputPort interface {
	AddVehicle(ctx context.Context, driverID valueobject.DriverID, v *vehicle.Vehicle) error
	GetVehicles(ctx context.Context, driverID valueobject.DriverID) ([]vehicle.Vehicle, error)
	UpdateVehicle(ctx context.Context, v *vehicle.Vehicle) error
}

type LicenseInputPort interface {
	AddLicense(ctx context.Context, l *license.License) error
	GetLicenses(ctx context.Context, driverID valueobject.DriverID) ([]license.License, error)
	VerifyLicense(ctx context.Context, licenseID valueobject.LicenseNumber) error
}

type DriverService interface {
	DriverInputPort
	AvailabilityInputPort
	VehicleInputPort
	LicenseInputPort
}
