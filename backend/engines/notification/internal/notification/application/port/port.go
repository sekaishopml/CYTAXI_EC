package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/delivery"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/notification"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/template"
)

type NotificationService interface {
	Create(ctx context.Context, cmd command.CreateNotification) (*notification.Notification, error)
	Queue(ctx context.Context, cmd command.QueueNotification) error
	Send(ctx context.Context, cmd command.SendNotification) (*delivery.Delivery, error)
	Retry(ctx context.Context, cmd command.RetryNotification) error
	Cancel(ctx context.Context, cmd command.CancelNotification) error
	UpdateDeliveryStatus(ctx context.Context, cmd command.UpdateDeliveryStatus) error
	Get(ctx context.Context, q query.GetNotification) (*notification.Notification, error)
	GetHistory(ctx context.Context, q query.GetNotificationHistory) ([]notification.Notification, error)
	GetDeliveryStatus(ctx context.Context, q query.GetDeliveryStatus) (*delivery.Delivery, error)
	GetPending(ctx context.Context, q query.GetPendingNotifications) ([]notification.Notification, error)
	GetTemplates(ctx context.Context, q query.GetTemplates) ([]template.NotificationTemplate, error)
}
