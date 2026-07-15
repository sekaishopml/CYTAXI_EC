package publisher

import "context"

type EventPublisher interface {
	Publish(ctx context.Context, event string, payload any) error
}

type LogPublisher struct{}

func NewLogPublisher() *LogPublisher {
	return &LogPublisher{}
}

func (p *LogPublisher) Publish(ctx context.Context, event string, payload any) error {
	return nil
}
