package whatsapp

import "context"

type Client interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	SendMessage(ctx context.Context, to string, content string) (*SendResult, error)
	GetQRCode(ctx context.Context) (*QRCode, error)
	GetStatus(ctx context.Context) (ConnectionStatus, error)
}

type MessageHandler func(msg Message)

type ClientConfig struct {
	SessionID     string
	ProviderKind  ProviderKind
	WebhookURL    string
	Reconnect     bool
	AutoLoadQR    bool
}

type client struct {
	config    ClientConfig
	adapter   ProviderAdapter
	status    ConnectionStatus
	session   Session
	handler   MessageHandler
}

func NewClient(config ClientConfig, adapter ProviderAdapter) Client {
	return &client{
		config:  config,
		adapter: adapter,
		status:  StatusDisconnected,
		session: Session{
			ID:     config.SessionID,
			Status: StatusDisconnected,
		},
	}
}

func (c *client) Connect(ctx context.Context) error {
	c.status = StatusConnecting
	if err := c.adapter.Connect(ctx); err != nil {
		c.status = StatusDisconnected
		return err
	}
	c.status = StatusConnected
	return nil
}

func (c *client) Disconnect(ctx context.Context) error {
	if err := c.adapter.Disconnect(ctx); err != nil {
		return err
	}
	c.status = StatusDisconnected
	return nil
}

func (c *client) SendMessage(ctx context.Context, to string, content string) (*SendResult, error) {
	if c.status != StatusConnected {
		return nil, ErrNotConnected
	}
	return c.adapter.SendMessage(ctx, to, content)
}

func (c *client) GetQRCode(ctx context.Context) (*QRCode, error) {
	return c.adapter.GetQRCode(ctx)
}

func (c *client) GetStatus(ctx context.Context) (ConnectionStatus, error) {
	status, err := c.adapter.GetStatus(ctx)
	if err != nil {
		return c.status, err
	}
	c.status = status
	return status, nil
}

func (c *client) SetMessageHandler(handler MessageHandler) {
	c.handler = handler
}
