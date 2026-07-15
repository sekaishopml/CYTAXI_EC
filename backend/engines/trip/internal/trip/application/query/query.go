package query

import (
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/trip"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/valueobject"
)

type GetTrip struct {
	TripID valueobject.TripID
}

type GetTripHistory struct {
	CustomerID valueobject.CustomerID
	Limit      int
	Offset     int
}

type GetTripTimeline struct {
	TripID valueobject.TripID
}

type GetDriverTrips struct {
	DriverID valueobject.DriverID
	From     int64
	To       int64
}

type GetCustomerTrips struct {
	CustomerID valueobject.CustomerID
	From       int64
	To         int64
}

type GetActiveTrips struct{}

type TripResult struct {
	Trip     *trip.Trip
	Timeline any
}

type TripListResult struct {
	Trips      []*trip.Trip
	TotalCount int
}
