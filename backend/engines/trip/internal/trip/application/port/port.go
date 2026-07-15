package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/application/query"
)

type TripService interface {
	Create(ctx context.Context, cmd command.CreateTrip) error
	AssignDriver(ctx context.Context, cmd command.AssignDriver) error
	UnassignDriver(ctx context.Context, cmd command.UnassignDriver) error
	Accept(ctx context.Context, cmd command.AcceptTrip) error
	Reject(ctx context.Context, cmd command.RejectTrip) error
	Start(ctx context.Context, cmd command.StartTrip) error
	Pause(ctx context.Context, cmd command.PauseTrip) error
	Resume(ctx context.Context, cmd command.ResumeTrip) error
	Complete(ctx context.Context, cmd command.CompleteTrip) error
	Cancel(ctx context.Context, cmd command.CancelTrip) error
	AddStop(ctx context.Context, cmd command.AddStop) error
	RemoveStop(ctx context.Context, cmd command.RemoveStop) error
	ChangeDestination(ctx context.Context, cmd command.ChangeDestination) error
	GetTrip(ctx context.Context, q query.GetTrip) (*query.TripResult, error)
	GetTripHistory(ctx context.Context, q query.GetTripHistory) (*query.TripListResult, error)
	GetActiveTrips(ctx context.Context, q query.GetActiveTrips) (*query.TripListResult, error)
}
