package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/channel"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/delivery"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/notification"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/template"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/valueobject"
	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/infrastructure/repository"
)

type NotificationService struct {
	notifRepo   repository.NotificationRepository
	templateRepo repository.TemplateRepository
	deliveryRepo repository.DeliveryRepository
	channels    *channel.ChannelRegistry
	logger      *slog.Logger
}

func NewNotificationService(
	notifRepo repository.NotificationRepository,
	templateRepo repository.TemplateRepository,
	deliveryRepo repository.DeliveryRepository,
	channels *channel.ChannelRegistry,
	logger *slog.Logger,
) *NotificationService {
	return &NotificationService{
		notifRepo:    notifRepo,
		templateRepo: templateRepo,
		deliveryRepo: deliveryRepo,
		channels:     channels,
		logger:       logger,
	}
}

func (s *NotificationService) Create(ctx context.Context, cmd command.CreateNotification) (*notification.Notification, error) {
	n := notification.NewNotification(cmd.RecipientID, cmd.Channel, cmd.TemplateID)
	n.Priority = cmd.Priority
	if cmd.Locale != "" {
		n.Locale = cmd.Locale
	}
	for k, v := range cmd.Parameters {
		n.SetParam(k, v)
	}
	if err := s.notifRepo.Save(ctx, n); err != nil {
		return nil, fmt.Errorf("save notification: %w", err)
	}
	return n, nil
}

func (s *NotificationService) Queue(ctx context.Context, cmd command.QueueNotification) error {
	n, err := s.notifRepo.FindByID(ctx, cmd.NotificationID)
	if err != nil {
		return fmt.Errorf("find notification: %w", err)
	}
	n.Queue()

	d := delivery.NewDelivery(n.ID, n.Channel)
	d.Queue()
	if err := s.notifRepo.Update(ctx, n); err != nil {
		return err
	}
	return s.deliveryRepo.Save(ctx, d)
}

func (s *NotificationService) Send(ctx context.Context, cmd command.SendNotification) (*delivery.Delivery, error) {
	n, err := s.notifRepo.FindByID(ctx, cmd.NotificationID)
	if err != nil {
		return nil, fmt.Errorf("find notification: %w", err)
	}

	tmpl, err := s.templateRepo.FindByID(ctx, n.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("find template: %w", err)
	}

	body := tmpl.Apply(n.Parameters)
	provider, err := s.channels.Get(n.Channel)
	if err != nil {
		return nil, err
	}

	d, err := s.deliveryRepo.FindByNotificationID(ctx, n.ID)
	if err != nil {
		return nil, fmt.Errorf("find delivery: %w", err)
	}

	n.Sending()
	s.notifRepo.Update(ctx, n)

	result, err := provider.Send(ctx, "", body)
	att := delivery.NewDeliveryAttempt(n.ID, d.Attempts+1)
	if err != nil || !result.Success {
		att.Status = valueobject.DeliveryStatus("failed")
		d.Fail()
		s.logger.Error("send failed", "notification", n.ID, "error", err)
	} else {
		att.Status = valueobject.DeliveryStatus("sent")
		d.Sent()
		n.Sent()
	}

	s.deliveryRepo.SaveAttempt(ctx, att)
	s.notifRepo.Update(ctx, n)
	s.deliveryRepo.Update(ctx, d)
	return d, nil
}

func (s *NotificationService) Retry(ctx context.Context, cmd command.RetryNotification) error {
	d, err := s.deliveryRepo.FindByNotificationID(ctx, cmd.NotificationID)
	if err != nil {
		return fmt.Errorf("find delivery: %w", err)
	}
	if !d.CanRetry() {
		return fmt.Errorf("cannot retry: attempts=%d, max=%d", d.Attempts, d.MaxRetries)
	}
	cmd2 := command.SendNotification{NotificationID: cmd.NotificationID}
	_, err = s.Send(ctx, cmd2)
	return err
}

func (s *NotificationService) Cancel(ctx context.Context, cmd command.CancelNotification) error {
	n, err := s.notifRepo.FindByID(ctx, cmd.NotificationID)
	if err != nil {
		return fmt.Errorf("find notification: %w", err)
	}
	n.Cancel()
	return s.notifRepo.Update(ctx, n)
}

func (s *NotificationService) UpdateDeliveryStatus(ctx context.Context, cmd command.UpdateDeliveryStatus) error {
	d, err := s.deliveryRepo.FindByNotificationID(ctx, cmd.NotificationID)
	if err != nil {
		return fmt.Errorf("find delivery: %w", err)
	}
	if cmd.Status == valueobject.DelStatusDelivered {
		d.Delivered()
		n, _ := s.notifRepo.FindByID(ctx, cmd.NotificationID)
		if n != nil {
			n.Delivered()
			s.notifRepo.Update(ctx, n)
		}
	}
	return s.deliveryRepo.Update(ctx, d)
}

func (s *NotificationService) Get(ctx context.Context, q query.GetNotification) (*notification.Notification, error) {
	return s.notifRepo.FindByID(ctx, q.NotificationID)
}

func (s *NotificationService) GetHistory(ctx context.Context, q query.GetNotificationHistory) ([]notification.Notification, error) {
	return s.notifRepo.FindByRecipientID(ctx, q.RecipientID)
}

func (s *NotificationService) GetDeliveryStatus(ctx context.Context, q query.GetDeliveryStatus) (*delivery.Delivery, error) {
	return s.deliveryRepo.FindByNotificationID(ctx, q.NotificationID)
}

func (s *NotificationService) GetPending(ctx context.Context, q query.GetPendingNotifications) ([]notification.Notification, error) {
	return s.notifRepo.FindPending(ctx, q.Channel)
}

func (s *NotificationService) GetTemplates(ctx context.Context, q query.GetTemplates) ([]template.NotificationTemplate, error) {
	return s.templateRepo.FindByChannel(ctx, q.Channel)
}
