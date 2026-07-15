package refund

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
)

type Refund struct {
	ID         valueobject.RefundID     `json:"id"`
	PaymentID  valueobject.PaymentID    `json:"payment_id"`
	Amount     valueobject.Money        `json:"amount"`
	Reason     string                   `json:"reason"`
	Status     valueobject.RefundStatus `json:"status"`
	GatewayRef string                   `json:"gateway_ref,omitempty"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
}

func NewRefund(paymentID valueobject.PaymentID, amount valueobject.Money, reason string) *Refund {
	now := time.Now()
	return &Refund{
		ID:        valueobject.NewRefundID(),
		PaymentID: paymentID,
		Amount:    amount,
		Reason:    reason,
		Status:    valueobject.RefStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (r *Refund) Process()  { r.Status = valueobject.RefStatusProcessing; r.UpdatedAt = time.Now() }
func (r *Refund) Complete(ref string) { r.Status = valueobject.RefStatusCompleted; r.GatewayRef = ref; r.UpdatedAt = time.Now() }
func (r *Refund) Reject()   { r.Status = valueobject.RefStatusRejected; r.UpdatedAt = time.Now() }
