package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/availability"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/driver"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/license"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/vehicle"
	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/infrastructure/repository"
)

type DriverService struct {
	driverRepo       repository.DriverRepository
	vehicleRepo      repository.VehicleRepository
	licenseRepo      repository.LicenseRepository
	availabilityRepo repository.AvailabilityRepository
	logger           *slog.Logger
}

func NewDriverService(
	driverRepo repository.DriverRepository,
	vehicleRepo repository.VehicleRepository,
	licenseRepo repository.LicenseRepository,
	availabilityRepo repository.AvailabilityRepository,
	logger *slog.Logger,
) *DriverService {
	return &DriverService{
		driverRepo:       driverRepo,
		vehicleRepo:      vehicleRepo,
		licenseRepo:      licenseRepo,
		availabilityRepo: availabilityRepo,
		logger:           logger,
	}
}

func (s *DriverService) Register(ctx context.Context, phone valueobject.Phone, name string) (*driver.Driver, error) {
	d := driver.NewDriver(phone, name)
	if err := s.driverRepo.Save(ctx, d); err != nil {
		return nil, fmt.Errorf("save driver: %w", err)
	}
	return d, nil
}

func (s *DriverService) GetDriver(ctx context.Context, id valueobject.DriverID) (*driver.Driver, error) {
	return s.driverRepo.FindByID(ctx, id)
}

func (s *DriverService) Approve(ctx context.Context, id valueobject.DriverID) error {
	d, err := s.driverRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find driver: %w", err)
	}
	d.Approve()
	return s.driverRepo.Update(ctx, d)
}

func (s *DriverService) Reject(ctx context.Context, id valueobject.DriverID, reason string) error {
	d, err := s.driverRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find driver: %w", err)
	}
	d.Reject()
	return s.driverRepo.Update(ctx, d)
}

func (s *DriverService) Suspend(ctx context.Context, id valueobject.DriverID, reason string) error {
	d, err := s.driverRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find driver: %w", err)
	}
	d.Suspend()
	return s.driverRepo.Update(ctx, d)
}

func (s *DriverService) GoOnline(ctx context.Context, id valueobject.DriverID, loc valueobject.Coordinates) error {
	d, err := s.driverRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find driver: %w", err)
	}
	if d.Status != driver.DriverApproved && d.Status != driver.DriverOffline {
		return fmt.Errorf("driver %s cannot go online from status %s", id, d.Status)
	}
	d.GoOnline()

	avail := availability.NewAvailability(id)
	avail.SetAvailable(loc)
	if err := s.availabilityRepo.Save(ctx, avail); err != nil {
		return fmt.Errorf("save availability: %w", err)
	}

	return s.driverRepo.Update(ctx, d)
}

func (s *DriverService) GoOffline(ctx context.Context, id valueobject.DriverID) error {
	d, err := s.driverRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find driver: %w", err)
	}
	d.GoOffline()
	return s.driverRepo.Update(ctx, d)
}

func (s *DriverService) Ping(ctx context.Context, id valueobject.DriverID, loc valueobject.Coordinates) error {
	avail, err := s.availabilityRepo.FindByDriverID(ctx, id)
	if err != nil {
		return fmt.Errorf("find availability: %w", err)
	}
	avail.Ping(loc)
	return s.availabilityRepo.Save(ctx, avail)
}

func (s *DriverService) GetAvailability(ctx context.Context, id valueobject.DriverID) (*availability.DriverAvailability, error) {
	return s.availabilityRepo.FindByDriverID(ctx, id)
}

func (s *DriverService) FindAvailableNearby(ctx context.Context, lat, lng, radiusMeters float64) ([]availability.DriverAvailability, error) {
	return s.availabilityRepo.FindAvailableInRadius(ctx, lat, lng, radiusMeters)
}

func (s *DriverService) AddVehicle(ctx context.Context, driverID valueobject.DriverID, v *vehicle.Vehicle) error {
	return s.vehicleRepo.Save(ctx, v)
}

func (s *DriverService) GetVehicles(ctx context.Context, driverID valueobject.DriverID) ([]vehicle.Vehicle, error) {
	return s.vehicleRepo.FindByDriverID(ctx, driverID)
}

func (s *DriverService) UpdateVehicle(ctx context.Context, v *vehicle.Vehicle) error {
	return s.vehicleRepo.Update(ctx, v)
}

func (s *DriverService) AddLicense(ctx context.Context, l *license.License) error {
	return s.licenseRepo.Save(ctx, l)
}

func (s *DriverService) GetLicenses(ctx context.Context, driverID valueobject.DriverID) ([]license.License, error) {
	return s.licenseRepo.FindByDriverID(ctx, driverID)
}

func (s *DriverService) VerifyLicense(ctx context.Context, licenseID valueobject.LicenseNumber) error {
	return fmt.Errorf("not implemented")
}
