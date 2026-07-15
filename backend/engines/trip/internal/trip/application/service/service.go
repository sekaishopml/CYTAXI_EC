package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/trip"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/infrastructure/repository"
)

type TripService struct {
	tripRepo     repository.TripRepository
	timelineRepo repository.TimelineRepository
	assignmentRepo repository.AssignmentRepository
	logger       *slog.Logger
}

func NewTripService(
	tripRepo repository.TripRepository,
	timelineRepo repository.TimelineRepository,
	assignmentRepo repository.AssignmentRepository,
	logger *slog.Logger,
) *TripService {
	return &TripService{
		tripRepo:       tripRepo,
		timelineRepo:   timelineRepo,
		assignmentRepo: assignmentRepo,
		logger:         logger,
	}
}

func (s *TripService) Create(ctx context.Context, cmd command.CreateTrip) error {
	t := trip.NewTrip(cmd.CustomerID, cmd.Passenger, cmd.Pickup, cmd.Destination)
	return s.tripRepo.Save(ctx, t)
}

func (s *TripService) AssignDriver(ctx context.Context, cmd command.AssignDriver) error {
	t, err := s.tripRepo.FindByID(ctx, cmd.TripID)
	if err != nil {
		return fmt.Errorf("find trip: %w", err)
	}
	if err := t.StartSearching(); err != nil {
		return err
	}
	if err := t.AssignDriver(cmd.DriverID); err != nil {
		return err
	}
	return s.tripRepo.Update(ctx, t)
}

func (s *TripService) UnassignDriver(ctx context.Context, cmd command.UnassignDriver) error {
	t, err := s.tripRepo.FindByID(ctx, cmd.TripID)
	if err != nil {
		return fmt.Errorf("find trip: %w", err)
	}
	t.UnassignDriver()
	return s.tripRepo.Update(ctx, t)
}

func (s *TripService) Accept(ctx context.Context, cmd command.AcceptTrip) error {
	t, err := s.tripRepo.FindByID(ctx, cmd.TripID)
	if err != nil {
		return fmt.Errorf("find trip: %w", err)
	}
	if err := t.AcceptTrip(); err != nil {
		return err
	}
	return s.tripRepo.Update(ctx, t)
}

func (s *TripService) Reject(ctx context.Context, cmd command.RejectTrip) error {
	t, err := s.tripRepo.FindByID(ctx, cmd.TripID)
	if err != nil {
		return fmt.Errorf("find trip: %w", err)
	}
	if err := t.RejectTrip(cmd.Reason); err != nil {
		return err
	}
	return s.tripRepo.Update(ctx, t)
}

func (s *TripService) Start(ctx context.Context, cmd command.StartTrip) error {
	t, err := s.tripRepo.FindByID(ctx, cmd.TripID)
	if err != nil {
		return fmt.Errorf("find trip: %w", err)
	}
	if err := t.DriverArrived(); err != nil {
		return err
	}
	if err := t.StartTrip(); err != nil {
		return err
	}
	return s.tripRepo.Update(ctx, t)
}

func (s *TripService) Pause(ctx context.Context, cmd command.PauseTrip) error {
	t, err := s.tripRepo.FindByID(ctx, cmd.TripID)
	if err != nil {
		return fmt.Errorf("find trip: %w", err)
	}
	if err := t.PauseTrip(); err != nil {
		return err
	}
	return s.tripRepo.Update(ctx, t)
}

func (s *TripService) Resume(ctx context.Context, cmd command.ResumeTrip) error {
	t, err := s.tripRepo.FindByID(ctx, cmd.TripID)
	if err != nil {
		return fmt.Errorf("find trip: %w", err)
	}
	if err := t.ResumeTrip(); err != nil {
		return err
	}
	return s.tripRepo.Update(ctx, t)
}

func (s *TripService) Complete(ctx context.Context, cmd command.CompleteTrip) error {
	t, err := s.tripRepo.FindByID(ctx, cmd.TripID)
	if err != nil {
		return fmt.Errorf("find trip: %w", err)
	}
	if err := t.CompleteTrip(); err != nil {
		return err
	}
	return s.tripRepo.Update(ctx, t)
}

func (s *TripService) Cancel(ctx context.Context, cmd command.CancelTrip) error {
	t, err := s.tripRepo.FindByID(ctx, cmd.TripID)
	if err != nil {
		return fmt.Errorf("find trip: %w", err)
	}
	if err := t.CancelTrip(cmd.Reason); err != nil {
		return err
	}
	return s.tripRepo.Update(ctx, t)
}

func (s *TripService) AddStop(ctx context.Context, cmd command.AddStop) error {
	t, err := s.tripRepo.FindByID(ctx, cmd.TripID)
	if err != nil {
		return fmt.Errorf("find trip: %w", err)
	}
	if err := t.AddStop(cmd.Stop); err != nil {
		return err
	}
	return s.tripRepo.Update(ctx, t)
}

func (s *TripService) RemoveStop(ctx context.Context, cmd command.RemoveStop) error {
	t, err := s.tripRepo.FindByID(ctx, cmd.TripID)
	if err != nil {
		return fmt.Errorf("find trip: %w", err)
	}
	if err := t.RemoveStop(cmd.StopID); err != nil {
		return err
	}
	return s.tripRepo.Update(ctx, t)
}

func (s *TripService) ChangeDestination(ctx context.Context, cmd command.ChangeDestination) error {
	t, err := s.tripRepo.FindByID(ctx, cmd.TripID)
	if err != nil {
		return fmt.Errorf("find trip: %w", err)
	}
	if err := t.ChangeDestination(cmd.Destination); err != nil {
		return err
	}
	return s.tripRepo.Update(ctx, t)
}

func (s *TripService) GetTrip(ctx context.Context, q query.GetTrip) (*query.TripResult, error) {
	t, err := s.tripRepo.FindByID(ctx, q.TripID)
	if err != nil {
		return nil, fmt.Errorf("find trip: %w", err)
	}
	return &query.TripResult{Trip: t}, nil
}

func (s *TripService) GetTripHistory(ctx context.Context, q query.GetTripHistory) (*query.TripListResult, error) {
	trips, err := s.tripRepo.FindByCustomerID(ctx, q.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("find trips: %w", err)
	}
	return &query.TripListResult{Trips: trips, TotalCount: len(trips)}, nil
}

func (s *TripService) GetActiveTrips(ctx context.Context, q query.GetActiveTrips) (*query.TripListResult, error) {
	trips, err := s.tripRepo.FindActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("find active trips: %w", err)
	}
	return &query.TripListResult{Trips: trips, TotalCount: len(trips)}, nil
}
