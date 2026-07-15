package intent

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
)

type PaymentIntent struct {
	ID          valueobject.PaymentID           `json:"id"`
	TripID      string                          `json:"trip_id"`
	CustomerID  string                          `json:"customer_id"`
	Amount      valueobject.Money               `json:"amount"`
	Method      valueobject.PaymentMethodType   `json:"method"`
	Status      valueobject.PaymentStatus       `json:"status"`
	ProviderRef string                           `json:"provider_ref,omitempty"`
	ClientSecret string                          `json:"client_secret,omitempty"`
	CreatedAt   time.Time                       `json:"created_at"`
	ExpiresAt   time.Time                       `json:"expires_at"`
}

func NewPaymentIntent(tripID, customerID string, amount valueobject.Money, method valueobject.PaymentMethodType) *PaymentIntent {
	now := time.Now()
	return &PaymentIntent{
		ID:         valueobject.NewPaymentID(),
		TripID:     tripID,
		CustomerID: customerID,
		Amount:     amount,
		Method:     method,
		Status:     valueobject.PayStatusPending,
		CreatedAt:  now,
		ExpiresAt:  now.Add(30 * time.Minute),
	}
}
