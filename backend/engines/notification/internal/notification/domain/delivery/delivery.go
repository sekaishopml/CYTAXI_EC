package delivery

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/valueobject"
)

type Delivery struct {
	ID             valueobject.NotificationID  `json:"id"`
	NotificationID valueobject.NotificationID  `json:"notification_id"`
	Channel        valueobject.ChannelType     `json:"channel"`
	Status         valueobject.DeliveryStatus  `json:"status"`
	Attempts       int                         `json:"attempts"`
	MaxRetries     int                         `json:"max_retries"`
	LastAttemptAt  *time.Time                  `json:"last_attempt_at,omitempty"`
	CompletedAt    *time.Time                  `json:"completed_at,omitempty"`
	CreatedAt      time.Time                   `json:"created_at"`
	UpdatedAt      time.Time                   `json:"updated_at"`
}

func NewDelivery(notificationID valueobject.NotificationID, channel valueobject.ChannelType) *Delivery {
	now := time.Now()
	return &Delivery{
		ID:             notificationID,
		NotificationID: notificationID,
		Channel:        channel,
		Status:         valueobject.DelStatusPending,
		MaxRetries:     3,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (d *Delivery) Queue() {
	d.Status = valueobject.DelStatusQueued
	d.UpdatedAt = time.Now()
}

func (d *Delivery) Sent() {
	now := time.Now()
	d.Status = valueobject.DelStatusSent
	d.LastAttemptAt = &now
	d.Attempts++
	d.UpdatedAt = now
}

func (d *Delivery) Delivered() {
	now := time.Now()
	d.Status = valueobject.DelStatusDelivered
	d.CompletedAt = &now
	d.UpdatedAt = now
}

func (d *Delivery) Fail() {
	d.Status = valueobject.DelStatusFailed
	d.UpdatedAt = time.Now()
}

func (d *Delivery) Cancel() {
	d.Status = valueobject.DelStatusCancelled
	d.UpdatedAt = time.Now()
}

func (d *Delivery) CanRetry() bool {
	return d.Attempts < d.MaxRetries && d.Status == valueobject.DelStatusFailed
}

type DeliveryAttempt struct {
	ID        valueobject.AttemptID     `json:"id"`
	DeliveryID valueobject.NotificationID `json:"delivery_id"`
	Attempt   int                       `json:"attempt"`
	Status    valueobject.DeliveryStatus `json:"status"`
	Error     string                    `json:"error,omitempty"`
	CreatedAt time.Time                 `json:"created_at"`
}

func NewDeliveryAttempt(deliveryID valueobject.NotificationID, attempt int) *DeliveryAttempt {
	return &DeliveryAttempt{
		ID:         valueobject.NewAttemptID(),
		DeliveryID: deliveryID,
		Attempt:    attempt,
		Status:     valueobject.DelStatusSending,
		CreatedAt:  time.Now(),
	}
}
