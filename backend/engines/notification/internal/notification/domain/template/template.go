package template

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/valueobject"
)

type NotificationTemplate struct {
	ID          valueobject.TemplateID `json:"id"`
	Name        string                 `json:"name"`
	Channel     valueobject.ChannelType `json:"channel"`
	Subject     string                 `json:"subject,omitempty"`
	Body        string                 `json:"body"`
	Locale      valueobject.Locale     `json:"locale"`
	Variables   []string               `json:"variables,omitempty"`
	Enabled     bool                   `json:"enabled"`
	Version     int                    `json:"version"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

func NewTemplate(id valueobject.TemplateID, name string, channel valueobject.ChannelType, body string) *NotificationTemplate {
	now := time.Now()
	return &NotificationTemplate{
		ID:        id,
		Name:      name,
		Channel:   channel,
		Body:      body,
		Locale:    "es",
		Enabled:   true,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (t *NotificationTemplate) Apply(params map[string]string) string {
	result := t.Body
	for k, v := range params {
		result = replaceParam(result, k, v)
	}
	return result
}

func replaceParam(s, key, value string) string {
	target := "{{" + key + "}}"
	result := ""
	for i := 0; i < len(s); {
		if i+len(target) <= len(s) && s[i:i+len(target)] == target {
			result += value
			i += len(target)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}
