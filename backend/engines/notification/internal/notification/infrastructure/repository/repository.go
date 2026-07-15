package repository

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/delivery"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/notification"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/template"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/valueobject"
)

type NotificationRepository interface {
	FindByID(ctx context.Context, id valueobject.NotificationID) (*notification.Notification, error)
	FindByRecipientID(ctx context.Context, recipientID valueobject.RecipientID) ([]notification.Notification, error)
	FindPending(ctx context.Context, channel valueobject.ChannelType) ([]notification.Notification, error)
	Save(ctx context.Context, n *notification.Notification) error
	Update(ctx context.Context, n *notification.Notification) error
}

type TemplateRepository interface {
	FindByID(ctx context.Context, id valueobject.TemplateID) (*template.NotificationTemplate, error)
	FindByChannel(ctx context.Context, channel valueobject.ChannelType) ([]template.NotificationTemplate, error)
	Save(ctx context.Context, t *template.NotificationTemplate) error
}

type DeliveryRepository interface {
	FindByNotificationID(ctx context.Context, id valueobject.NotificationID) (*delivery.Delivery, error)
	Save(ctx context.Context, d *delivery.Delivery) error
	Update(ctx context.Context, d *delivery.Delivery) error
	SaveAttempt(ctx context.Context, a *delivery.DeliveryAttempt) error
}
