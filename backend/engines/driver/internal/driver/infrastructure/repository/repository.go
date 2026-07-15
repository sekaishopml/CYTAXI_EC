package repository

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/availability"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/document"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/driver"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/license"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/preference"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/vehicle"
)

type DriverRepository interface {
	FindByID(ctx context.Context, id valueobject.DriverID) (*driver.Driver, error)
	FindByPhone(ctx context.Context, phone valueobject.Phone) (*driver.Driver, error)
	Save(ctx context.Context, d *driver.Driver) error
	Update(ctx context.Context, d *driver.Driver) error
	FindByStatus(ctx context.Context, status driver.DriverStatus) ([]driver.Driver, error)
}

type VehicleRepository interface {
	FindByID(ctx context.Context, id valueobject.VehicleID) (*vehicle.Vehicle, error)
	FindByDriverID(ctx context.Context, driverID valueobject.DriverID) ([]vehicle.Vehicle, error)
	Save(ctx context.Context, v *vehicle.Vehicle) error
	Update(ctx context.Context, v *vehicle.Vehicle) error
}

type LicenseRepository interface {
	FindByDriverID(ctx context.Context, driverID valueobject.DriverID) ([]license.License, error)
	Save(ctx context.Context, l *license.License) error
}

type AvailabilityRepository interface {
	FindByDriverID(ctx context.Context, driverID valueobject.DriverID) (*availability.DriverAvailability, error)
	Save(ctx context.Context, a *availability.DriverAvailability) error
	FindAvailableInRadius(ctx context.Context, lat, lng, radiusMeters float64) ([]availability.DriverAvailability, error)
}

type PreferenceRepository interface {
	FindByDriverID(ctx context.Context, driverID valueobject.DriverID) (*preference.Preferences, error)
	Save(ctx context.Context, p *preference.Preferences) error
}

type DocumentRepository interface {
	FindByDriverID(ctx context.Context, driverID valueobject.DriverID) ([]document.Document, error)
	Save(ctx context.Context, d *document.Document) error
}
