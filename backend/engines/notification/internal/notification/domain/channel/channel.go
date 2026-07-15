package channel

import (
	"context"
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/valueobject"
)

type ChannelProvider interface {
	Name() string
	Kind() valueobject.ChannelType
	Send(ctx context.Context, to string, body string) (*SendResult, error)
	IsAvailable(ctx context.Context) bool
}

type SendResult struct {
	Success    bool
	MessageID  string
	ProviderID string
	Error      string
}

type ChannelRegistry struct {
	providers map[valueobject.ChannelType]ChannelProvider
}

func NewChannelRegistry() *ChannelRegistry {
	return &ChannelRegistry{
		providers: make(map[valueobject.ChannelType]ChannelProvider),
	}
}

func (r *ChannelRegistry) Register(p ChannelProvider) {
	r.providers[p.Kind()] = p
}

func (r *ChannelRegistry) Get(kind valueobject.ChannelType) (ChannelProvider, error) {
	p, ok := r.providers[kind]
	if !ok {
		return nil, fmt.Errorf("channel provider %s not registered", kind)
	}
	return p, nil
}

func (r *ChannelRegistry) AvailableChannels() []valueobject.ChannelType {
	var channels []valueobject.ChannelType
	for _, p := range r.providers {
		channels = append(channels, p.Kind())
	}
	return channels
}

type LogChannel struct {
	kind valueobject.ChannelType
}

func NewLogChannel(kind valueobject.ChannelType) ChannelProvider {
	return &LogChannel{kind: kind}
}

func (c *LogChannel) Name() string            { return "log_" + string(c.kind) }
func (c *LogChannel) Kind() valueobject.ChannelType { return c.kind }

func (c *LogChannel) Send(ctx context.Context, to string, body string) (*SendResult, error) {
	return &SendResult{Success: true, MessageID: fmt.Sprintf("log_%d", time.Now().UnixNano())}, nil
}

func (c *LogChannel) IsAvailable(ctx context.Context) bool { return true }
