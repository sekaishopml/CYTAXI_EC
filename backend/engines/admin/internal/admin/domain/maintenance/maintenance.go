package maintenance

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/valueobject"
)

type MaintenanceWindow struct {
	ID          valueobject.MaintenanceID `json:"id"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	StartAt     time.Time                 `json:"start_at"`
	EndAt       time.Time                 `json:"end_at"`
	Affects     []string                  `json:"affects"` // engines affected
	Status      MaintenanceStatus         `json:"status"`
	CreatedAt   time.Time                 `json:"created_at"`
}

type MaintenanceStatus string

const (
	MaintScheduled  MaintenanceStatus = "scheduled"
	MaintInProgress MaintenanceStatus = "in_progress"
	MaintCompleted  MaintenanceStatus = "completed"
	MaintCancelled  MaintenanceStatus = "cancelled"
)

func NewMaintenanceWindow(title string, startAt, endAt time.Time) *MaintenanceWindow {
	return &MaintenanceWindow{
		ID:        valueobject.NewMaintenanceID(),
		Title:     title,
		StartAt:   startAt,
		EndAt:     endAt,
		Status:    MaintScheduled,
		CreatedAt: time.Now(),
	}
}
