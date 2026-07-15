package payment

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
)

type Payment struct {
	ID            valueobject.PaymentID    `json:"id"`
	TripID        string                   `json:"trip_id"`
	CustomerID    string                   `json:"customer_id"`
	DriverID      string                   `json:"driver_id,omitempty"`
	Amount        valueobject.Money        `json:"amount"`
	Method        valueobject.PaymentMethodType `json:"method"`
	Status        valueobject.PaymentStatus `json:"status"`
	GatewayRef    string                   `json:"gateway_ref,omitempty"`
	Description   string                   `json:"description"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
	CompletedAt   *time.Time               `json:"completed_at,omitempty"`
}

func NewPayment(tripID, customerID string, amount valueobject.Money, method valueobject.PaymentMethodType) *Payment {
	now := time.Now()
	return &Payment{
		ID:          valueobject.NewPaymentID(),
		TripID:      tripID,
		CustomerID:  customerID,
		Amount:      amount,
		Method:      method,
		Status:      valueobject.PayStatusPending,
		Description: "Trip payment",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Payment) Authorize(gatewayRef string) {
	p.Status = valueobject.PayStatusAuthorized
	p.GatewayRef = gatewayRef
	p.UpdatedAt = time.Now()
}

func (p *Payment) Capture() {
	now := time.Now()
	p.Status = valueobject.PayStatusCaptured
	p.CompletedAt = &now
	p.UpdatedAt = now
}

func (p *Payment) Fail() {
	p.Status = valueobject.PayStatusFailed
	p.UpdatedAt = time.Now()
}

func (p *Payment) Cancel() {
	p.Status = valueobject.PayStatusCancelled
	p.UpdatedAt = time.Now()
}

func (p *Payment) Refund() {
	p.Status = valueobject.PayStatusRefunded
	p.UpdatedAt = time.Now()
}
