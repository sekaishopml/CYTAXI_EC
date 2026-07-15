package repository

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/trip"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/timeline"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/valueobject"
)

type TripRepository interface {
	FindByID(ctx context.Context, id valueobject.TripID) (*trip.Trip, error)
	FindByCustomerID(ctx context.Context, customerID valueobject.CustomerID) ([]*trip.Trip, error)
	FindByDriverID(ctx context.Context, driverID valueobject.DriverID) ([]*trip.Trip, error)
	FindActive(ctx context.Context) ([]*trip.Trip, error)
	Save(ctx context.Context, t *trip.Trip) error
	Update(ctx context.Context, t *trip.Trip) error
}

type TimelineRepository interface {
	FindByTripID(ctx context.Context, tripID valueobject.TripID) ([]timeline.TimelineEntry, error)
	Save(ctx context.Context, e timeline.TimelineEntry) error
}

type AssignmentRepository interface {
	FindByTripID(ctx context.Context, tripID valueobject.TripID) (any, error)
	Save(ctx context.Context, assignment any) error
}
