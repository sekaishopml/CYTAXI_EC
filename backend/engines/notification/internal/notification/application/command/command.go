package command

import "github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/valueobject"

type CreateNotification struct {
	RecipientID valueobject.RecipientID
	Channel     valueobject.ChannelType
	TemplateID  valueobject.TemplateID
	Priority    valueobject.Priority
	Locale      valueobject.Locale
	Parameters  map[string]string
}

type QueueNotification struct {
	NotificationID valueobject.NotificationID
}

type SendNotification struct {
	NotificationID valueobject.NotificationID
}

type RetryNotification struct {
	NotificationID valueobject.NotificationID
}

type CancelNotification struct {
	NotificationID valueobject.NotificationID
	Reason         string
}

type UpdateDeliveryStatus struct {
	NotificationID valueobject.NotificationID
	Status         valueobject.DeliveryStatus
}
