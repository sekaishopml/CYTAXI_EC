package events

const (
	EventNotificationCreated   = "notification.created"
	EventNotificationQueued    = "notification.queued"
	EventNotificationSent      = "notification.sent"
	EventNotificationDelivered = "notification.delivered"
	EventNotificationFailed    = "notification.failed"
	EventNotificationRetried   = "notification.retried"
	EventNotificationCancelled = "notification.cancelled"
	EventTemplateUpdated       = "notification.template_updated"
)

type NotificationCreatedPayload struct {
	NotificationID string `json:"notification_id"`
	RecipientID    string `json:"recipient_id"`
	Channel        string `json:"channel"`
	TemplateID     string `json:"template_id"`
	Priority       int    `json:"priority"`
}

type NotificationQueuedPayload struct {
	NotificationID string `json:"notification_id"`
	Channel        string `json:"channel"`
}

type NotificationSentPayload struct {
	NotificationID string `json:"notification_id"`
	Channel        string `json:"channel"`
	ProviderID     string `json:"provider_id,omitempty"`
}

type NotificationDeliveredPayload struct {
	NotificationID string `json:"notification_id"`
	Channel        string `json:"channel"`
	DeliveredAt    string `json:"delivered_at"`
}

type NotificationFailedPayload struct {
	NotificationID string `json:"notification_id"`
	Channel        string `json:"channel"`
	Error          string `json:"error"`
	Attempts       int    `json:"attempts"`
}

type NotificationRetriedPayload struct {
	NotificationID string `json:"notification_id"`
	Attempt        int    `json:"attempt"`
}

type NotificationCancelledPayload struct {
	NotificationID string `json:"notification_id"`
	Reason         string `json:"reason"`
}

type TemplateUpdatedPayload struct {
	TemplateID string `json:"template_id"`
	Version    int    `json:"version"`
}
