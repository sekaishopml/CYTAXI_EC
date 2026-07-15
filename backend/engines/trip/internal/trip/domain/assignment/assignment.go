package assignment

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/valueobject"
)

type Assignment struct {
	ID            valueobject.TripID    `json:"id"`
	TripID        valueobject.TripID    `json:"trip_id"`
	DriverID      valueobject.DriverID  `json:"driver_id"`
	Status        AssignmentStatus      `json:"status"`
	Strategy      string                `json:"strategy"`
	Score         float64               `json:"score"`
	AssignedAt    time.Time             `json:"assigned_at"`
	UnassignedAt  *time.Time            `json:"unassigned_at,omitempty"`
}

type AssignmentStatus string

const (
	AssignmentPending   AssignmentStatus = "pending"
	AssignmentActive    AssignmentStatus = "active"
	AssignmentReplaced  AssignmentStatus = "replaced"
	AssignmentExpired   AssignmentStatus = "expired"
)
