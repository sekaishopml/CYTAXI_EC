package query

import "github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/valueobject"

type GetNotification struct {
	NotificationID valueobject.NotificationID
}

type GetNotificationHistory struct {
	RecipientID valueobject.RecipientID
	Limit       int
	Offset      int
}

type GetDeliveryStatus struct {
	NotificationID valueobject.NotificationID
}

type GetPendingNotifications struct {
	Channel valueobject.ChannelType
	Limit   int
}

type GetTemplates struct {
	Channel valueobject.ChannelType
	Locale  valueobject.Locale
}
