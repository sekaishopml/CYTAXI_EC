package validator

import (
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/driver"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/vehicle"
)

type DriverValidator interface {
	ValidateDriver(d *driver.Driver) error
	ValidateVehicle(v *vehicle.Vehicle) error
	CanGoOnline(driverID valueobject.DriverID) (bool, string)
}
