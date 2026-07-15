package timeline

import (
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/valueobject"
)

type TimelineEntry struct {
	ID        string                 `json:"id"`
	TripID    valueobject.TripID     `json:"trip_id"`
	Timestamp time.Time              `json:"timestamp"`
	Event     string                 `json:"event"`
	Details   map[string]any         `json:"details,omitempty"`
}

func NewEntry(tripID valueobject.TripID, event string) TimelineEntry {
	return TimelineEntry{
		ID:        fmt.Sprintf("tme_%d", time.Now().UnixNano()),
		TripID:    tripID,
		Timestamp: time.Now(),
		Event:     event,
		Details:   make(map[string]any),
	}
}
