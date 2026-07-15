package preference

import "github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/valueobject"

type NotificationPreference struct {
	RecipientID valueobject.RecipientID `json:"recipient_id"`
	Preferences []ChannelPreference     `json:"preferences"`
}

type ChannelPreference struct {
	Channel valueobject.ChannelType `json:"channel"`
	Enabled bool                    `json:"enabled"`
	QuietHours *QuietHours          `json:"quiet_hours,omitempty"`
}

type QuietHours struct {
	Start string `json:"start"` // HH:MM
	End   string `json:"end"`   // HH:MM
	Timezone string `json:"timezone"`
}

func NewNotificationPreference(recipientID valueobject.RecipientID) *NotificationPreference {
	return &NotificationPreference{
		RecipientID: recipientID,
		Preferences: make([]ChannelPreference, 0),
	}
}

func (p *NotificationPreference) IsChannelEnabled(channel valueobject.ChannelType) bool {
	for _, pref := range p.Preferences {
		if pref.Channel == channel {
			return pref.Enabled
		}
	}
	return true
}
