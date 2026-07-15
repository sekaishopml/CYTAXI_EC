package settlement

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
)

type Settlement struct {
	ID           valueobject.SettlementID `json:"id"`
	PaymentID    valueobject.PaymentID    `json:"payment_id"`
	DriverID     string                   `json:"driver_id"`
	Amount       valueobject.Money        `json:"amount"`
	Commission   valueobject.Money        `json:"commission"`
	NetAmount    valueobject.Money        `json:"net_amount"`
	Status       valueobject.SettlementStatus `json:"status"`
	PayableAt    time.Time                `json:"payable_at,omitempty"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

func NewSettlement(paymentID valueobject.PaymentID, driverID string, amount, commission valueobject.Money) *Settlement {
	now := time.Now()
	net := amount.Subtract(commission)
	return &Settlement{
		ID:         valueobject.NewSettlementID(),
		PaymentID:  paymentID,
		DriverID:   driverID,
		Amount:     amount,
		Commission: commission,
		NetAmount:  net,
		Status:     valueobject.SettStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (s *Settlement) Process()  { s.Status = valueobject.SettStatusProcessing; s.UpdatedAt = time.Now() }
func (s *Settlement) Complete() { s.Status = valueobject.SettStatusCompleted; s.UpdatedAt = time.Now() }
func (s *Settlement) Fail()     { s.Status = valueobject.SettStatusFailed; s.UpdatedAt = time.Now() }
