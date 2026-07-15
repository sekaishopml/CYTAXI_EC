package whatsapp

import (
	"context"
	"errors"
)

var ErrNotConnected = errors.New("whatsapp: not connected")
var ErrInvalidQR = errors.New("whatsapp: invalid qr code")
var ErrSendFailed = errors.New("whatsapp: send failed")

type ProviderKind string

const (
	ProviderWhatsMeow   ProviderKind = "whatsmeow"
	ProviderWAWebJS     ProviderKind = "wawebjs"
	ProviderBusinessAPI ProviderKind = "business_api"
)

type ProviderAdapter interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	SendMessage(ctx context.Context, to string, content string) (*SendResult, error)
	GetQRCode(ctx context.Context) (*QRCode, error)
	GetStatus(ctx context.Context) (ConnectionStatus, error)
}

type adapterFactory struct{}

func NewAdapter(kind ProviderKind) (ProviderAdapter, error) {
	switch kind {
	case ProviderWhatsMeow:
		return newWhatsMeowAdapter(), nil
	case ProviderWAWebJS:
		return nil, errors.New("whatsapp: wawebjs adapter not implemented")
	case ProviderBusinessAPI:
		return nil, errors.New("whatsapp: business api adapter not implemented")
	default:
		return nil, errors.New("whatsapp: unknown provider kind")
	}
}

type whatsMeowAdapter struct {
	connected bool
	qrCode    *QRCode
}

func newWhatsMeowAdapter() *whatsMeowAdapter {
	return &whatsMeowAdapter{}
}

func (a *whatsMeowAdapter) Connect(ctx context.Context) error {
	a.connected = true
	return nil
}

func (a *whatsMeowAdapter) Disconnect(ctx context.Context) error {
	a.connected = false
	return nil
}

func (a *whatsMeowAdapter) SendMessage(ctx context.Context, to string, content string) (*SendResult, error) {
	if !a.connected {
		return nil, ErrNotConnected
	}
	return &SendResult{
		MessageID: MessageID("msg_" + to),
		Status:    "sent",
	}, nil
}

func (a *whatsMeowAdapter) GetQRCode(ctx context.Context) (*QRCode, error) {
	if a.qrCode != nil {
		return a.qrCode, nil
	}
	return nil, ErrInvalidQR
}

func (a *whatsMeowAdapter) GetStatus(ctx context.Context) (ConnectionStatus, error) {
	if a.connected {
		return StatusConnected, nil
	}
	return StatusDisconnected, nil
}
