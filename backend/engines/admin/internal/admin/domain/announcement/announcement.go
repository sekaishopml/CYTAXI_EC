package announcement

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/valueobject"
)

type Announcement struct {
	ID          valueobject.AnnouncementID `json:"id"`
	Title       string                     `json:"title"`
	Body        string                     `json:"body"`
	Priority    valueobject.Priority       `json:"priority"`
	PublishAt   time.Time                  `json:"publish_at"`
	ExpiresAt   *time.Time                 `json:"expires_at,omitempty"`
	CreatedAt   time.Time                  `json:"created_at"`
}

func NewAnnouncement(title, body string, priority valueobject.Priority, publishAt time.Time) *Announcement {
	return &Announcement{
		ID:        valueobject.NewAnnouncementID(),
		Title:     title,
		Body:      body,
		Priority:  priority,
		PublishAt: publishAt,
		CreatedAt: time.Now(),
	}
}
