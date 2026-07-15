package valueobject

import (
	"fmt"
	"time"
)

type NotificationID string
type RecipientID string
type TemplateID string
type AttemptID string

type ChannelType string

const (
	ChannelWhatsApp  ChannelType = "whatsapp"
	ChannelPush      ChannelType = "push"
	ChannelEmail     ChannelType = "email"
	ChannelSMS       ChannelType = "sms"
	ChannelWebSocket ChannelType = "websocket"
	ChannelInApp     ChannelType = "in_app"
)

type Priority int

const (
	PriorityLow    Priority = 0
	PriorityNormal Priority = 5
	PriorityHigh   Priority = 7
	PriorityUrgent Priority = 10
)

type DeliveryStatus string

const (
	DelStatusPending   DeliveryStatus = "pending"
	DelStatusQueued    DeliveryStatus = "queued"
	DelStatusSending   DeliveryStatus = "sending"
	DelStatusSent      DeliveryStatus = "sent"
	DelStatusDelivered DeliveryStatus = "delivered"
	DelStatusFailed    DeliveryStatus = "failed"
	DelStatusCancelled DeliveryStatus = "cancelled"
)

type Locale string

func NewNotificationID() NotificationID {
	return NotificationID(fmt.Sprintf("notif_%d", time.Now().UnixNano()))
}

func NewAttemptID() AttemptID {
	return AttemptID(fmt.Sprintf("att_%d", time.Now().UnixNano()))
}

type RecipientContact struct {
	Phone       string
	Email       string
	DeviceToken string
	UserID      string
}
