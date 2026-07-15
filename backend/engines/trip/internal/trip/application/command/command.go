package command

import (
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/destination"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/passenger"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/stop"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/valueobject"
)

type CreateTrip struct {
	CustomerID  valueobject.CustomerID
	Passenger   passenger.Passenger
	Pickup      stop.Stop
	Destination destination.Destination
}

type AssignDriver struct {
	TripID   valueobject.TripID
	DriverID valueobject.DriverID
	Strategy string
	Score    float64
}

type UnassignDriver struct {
	TripID   valueobject.TripID
	DriverID valueobject.DriverID
	Reason   string
}

type AcceptTrip struct {
	TripID   valueobject.TripID
	DriverID valueobject.DriverID
}

type RejectTrip struct {
	TripID valueobject.TripID
	Reason string
}

type StartTrip struct {
	TripID valueobject.TripID
}

type PauseTrip struct {
	TripID valueobject.TripID
}

type ResumeTrip struct {
	TripID valueobject.TripID
}

type CompleteTrip struct {
	TripID valueobject.TripID
}

type CancelTrip struct {
	TripID valueobject.TripID
	Reason string
	By     string
}

type AddStop struct {
	TripID  valueobject.TripID
	Stop    stop.Stop
}

type RemoveStop struct {
	TripID valueobject.TripID
	StopID valueobject.StopID
}

type ChangeDestination struct {
	TripID      valueobject.TripID
	Destination destination.Destination
}
