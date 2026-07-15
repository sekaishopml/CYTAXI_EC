package whatsapp

import "time"

type MessageID string

type Message struct {
	ID        MessageID
	From      string
	To        string
	Content   string
	Timestamp time.Time
}

type QRCode struct {
	Code      string
	ExpiresAt time.Time
}

type ConnectionStatus string

const (
	StatusConnected    ConnectionStatus = "connected"
	StatusDisconnected ConnectionStatus = "disconnected"
	StatusConnecting   ConnectionStatus = "connecting"
	StatusQRReady      ConnectionStatus = "qr_ready"
)

type Session struct {
	ID        string
	Status    ConnectionStatus
	QRCode    *QRCode
	CreatedAt time.Time
}

type SendResult struct {
	MessageID MessageID
	Status    string
	Error     error
}
