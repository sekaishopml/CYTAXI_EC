package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/trip"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/timeline"
	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/valueobject"
)

// ErrTripNotFound error when trip does not exist
var ErrTripNotFound = errors.New("trip not found")

// InMemoryTripRepository is a thread-safe in-memory implementation
type InMemoryTripRepository struct {
	mu        sync.RWMutex
	byID      map[string]*trip.Trip
	byCust    map[string][]*trip.Trip
	byDriver  map[string][]*trip.Trip
}

func NewInMemoryTripRepository() *InMemoryTripRepository {
	return &InMemoryTripRepository{
		byID:     make(map[string]*trip.Trip),
		byCust:   make(map[string][]*trip.Trip),
		byDriver: make(map[string][]*trip.Trip),
	}
}

func (r *InMemoryTripRepository) FindByID(ctx context.Context, id valueobject.TripID) (*trip.Trip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.byID[string(id)]
	if !ok {
		return nil, ErrTripNotFound
	}
	return t, nil
}

func (r *InMemoryTripRepository) FindByCustomerID(ctx context.Context, customerID valueobject.CustomerID) ([]*trip.Trip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := r.byCust[string(customerID)]
	return list, nil
}

func (r *InMemoryTripRepository) FindByDriverID(ctx context.Context, driverID valueobject.DriverID) ([]*trip.Trip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := r.byDriver[string(driverID)]
	return list, nil
}

func (r *InMemoryTripRepository) FindActive(ctx context.Context) ([]*trip.Trip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*trip.Trip, 0)
	for _, t := range r.byID {
		if t.IsActive() {
			out = append(out, t)
		}
	}
	return out, nil
}

func (r *InMemoryTripRepository) Save(ctx context.Context, t *trip.Trip) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[string(t.ID)] = t
	r.byCust[string(t.CustomerID)] = append(r.byCust[string(t.CustomerID)], t)
	return nil
}

func (r *InMemoryTripRepository) Update(ctx context.Context, t *trip.Trip) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[string(t.ID)] = t
	return nil
}

// InMemoryTimelineRepository stores timeline events in memory
type InMemoryTimelineRepository struct {
	mu   sync.RWMutex
	data map[string][]timeline.TimelineEntry
}

func NewInMemoryTimelineRepository() *InMemoryTimelineRepository {
	return &InMemoryTimelineRepository{data: make(map[string][]timeline.TimelineEntry)}
}

func (r *InMemoryTimelineRepository) FindByTripID(ctx context.Context, tripID valueobject.TripID) ([]timeline.TimelineEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data[string(tripID)], nil
}

func (r *InMemoryTimelineRepository) Save(ctx context.Context, e timeline.TimelineEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := string(e.TripID)
	r.data[id] = append(r.data[id], e)
	return nil
}

// InMemoryAssignmentRepository stores assignments
type InMemoryAssignmentRepository struct {
	mu   sync.RWMutex
	data map[string]any
}

func NewInMemoryAssignmentRepository() *InMemoryAssignmentRepository {
	return &InMemoryAssignmentRepository{data: make(map[string]any)}
}

func (r *InMemoryAssignmentRepository) FindByTripID(ctx context.Context, tripID valueobject.TripID) (any, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.data[string(tripID)]
	if !ok {
		return nil, ErrTripNotFound
	}
	return v, nil
}

func (r *InMemoryAssignmentRepository) Save(ctx context.Context, assignment any) error {
	// Try to extract tripID from assignment struct
	r.mu.Lock()
	defer r.mu.Unlock()
	// Use a key from the struct via reflection-free approach: use type assertion
	type tripIDer interface{ TripIDString() string }
	if t, ok := assignment.(tripIDer); ok {
		r.data[t.TripIDString()] = assignment
		return nil
	}
	// Fallback: store with a generic key
	r.data["_last"] = assignment
	return nil
}
