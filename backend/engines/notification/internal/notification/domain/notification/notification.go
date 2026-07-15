package notification

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/valueobject"
)

type Notification struct {
	ID          valueobject.NotificationID `json:"id"`
	RecipientID valueobject.RecipientID    `json:"recipient_id"`
	Channel     valueobject.ChannelType    `json:"channel"`
	TemplateID  valueobject.TemplateID     `json:"template_id"`
	Priority    valueobject.Priority       `json:"priority"`
	Status      valueobject.DeliveryStatus `json:"status"`
	Locale      valueobject.Locale         `json:"locale"`
	Parameters  map[string]string          `json:"parameters,omitempty"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
}

func NewNotification(recipientID valueobject.RecipientID, channel valueobject.ChannelType, templateID valueobject.TemplateID) *Notification {
	now := time.Now()
	return &Notification{
		ID:          valueobject.NewNotificationID(),
		RecipientID: recipientID,
		Channel:     channel,
		TemplateID:  templateID,
		Priority:    valueobject.PriorityNormal,
		Status:      valueobject.DelStatusPending,
		Locale:      "es",
		Parameters:  make(map[string]string),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (n *Notification) Queue()  { n.Status = valueobject.DelStatusQueued; n.UpdatedAt = time.Now() }
func (n *Notification) Sending() { n.Status = valueobject.DelStatusSending; n.UpdatedAt = time.Now() }
func (n *Notification) Sent()    { n.Status = valueobject.DelStatusSent; n.UpdatedAt = time.Now() }
func (n *Notification) Delivered() { n.Status = valueobject.DelStatusDelivered; n.UpdatedAt = time.Now() }
func (n *Notification) Fail()     { n.Status = valueobject.DelStatusFailed; n.UpdatedAt = time.Now() }
func (n *Notification) Cancel()   { n.Status = valueobject.DelStatusCancelled; n.UpdatedAt = time.Now() }

func (n *Notification) SetParam(key, value string) {
	n.Parameters[key] = value
	n.UpdatedAt = time.Now()
}
