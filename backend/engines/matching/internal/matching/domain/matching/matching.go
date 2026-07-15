package matching

import (
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/valueobject"
)

type MatchingStatus string

const (
	StatusPending    MatchingStatus = "pending"
	StatusSearching  MatchingStatus = "searching"
	StatusEvaluating MatchingStatus = "evaluating"
	StatusRanking    MatchingStatus = "ranking"
	StatusSelecting  MatchingStatus = "selecting"
	StatusCompleted  MatchingStatus = "completed"
	StatusFailed     MatchingStatus = "failed"
	StatusCancelled  MatchingStatus = "cancelled"
)

type Matching struct {
	ID           valueobject.MatchingID `json:"id"`
	TripID       valueobject.TripID     `json:"trip_id"`
	Status       MatchingStatus         `json:"status"`
	Strategy     string                 `json:"strategy"`
	MaxCandidates int                   `json:"max_candidates"`
	RadiusMeters  float64               `json:"radius_meters"`
	Retries      int                    `json:"retries"`
	MaxRetries   int                    `json:"max_retries"`
	Session      MatchingSession        `json:"session"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty"`
}

type MatchingSession struct {
	SessionID     valueobject.SessionID     `json:"session_id"`
	PickupLat     float64                   `json:"pickup_lat"`
	PickupLng     float64                   `json:"pickup_lng"`
	VehicleType   string                    `json:"vehicle_type,omitempty"`
	Requirements  []string                  `json:"requirements,omitempty"`
}

func NewMatching(tripID valueobject.TripID, session MatchingSession) *Matching {
	now := time.Now()
	return &Matching{
		ID:           valueobject.NewMatchingID(),
		TripID:       tripID,
		Status:       StatusPending,
		Strategy:     "balanced",
		MaxCandidates: 10,
		RadiusMeters: 5000,
		MaxRetries:   3,
		Session:      session,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (m *Matching) StartSearching() error {
	if m.Status != StatusPending {
		return fmt.Errorf("cannot search in status %s", m.Status)
	}
	m.Status = StatusSearching
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Matching) StartEvaluating() error {
	m.Status = StatusEvaluating
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Matching) StartRanking() error {
	m.Status = StatusRanking
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Matching) Complete() error {
	now := time.Now()
	m.Status = StatusCompleted
	m.CompletedAt = &now
	m.UpdatedAt = now
	return nil
}

func (m *Matching) Fail() error {
	m.Status = StatusFailed
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Matching) Cancel() error {
	m.Status = StatusCancelled
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Matching) Retry() error {
	if m.Retries >= m.MaxRetries {
		return fmt.Errorf("max retries (%d) reached", m.MaxRetries)
	}
	m.Retries++
	m.Status = StatusSearching
	m.UpdatedAt = time.Now()
	return nil
}

type MatchingPolicy struct {
	Strategy      string  `json:"strategy"`
	MaxCandidates int     `json:"max_candidates"`
	RadiusMeters  float64 `json:"radius_meters"`
	MaxRetries    int     `json:"max_retries"`
	MinScore      float64 `json:"min_score"`
	MaxDistance   float64 `json:"max_distance_meters"`
}
