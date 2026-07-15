package invoice

import (
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
)

type Invoice struct {
	ID         valueobject.InvoiceNumber `json:"id"`
	PaymentID  valueobject.PaymentID     `json:"payment_id"`
	CustomerID string                    `json:"customer_id"`
	Amount     valueobject.Money         `json:"amount"`
	Items      []InvoiceItem             `json:"items"`
	Status     string                    `json:"status"`
	IssuedAt   time.Time                 `json:"issued_at"`
	PaidAt     *time.Time                `json:"paid_at,omitempty"`
}

type InvoiceItem struct {
	Description string           `json:"description"`
	Quantity    int              `json:"quantity"`
	UnitPrice   valueobject.Money `json:"unit_price"`
	Total       valueobject.Money `json:"total"`
}

func NewInvoice(paymentID valueobject.PaymentID, customerID string, amount valueobject.Money) *Invoice {
	return &Invoice{
		ID:         valueobject.InvoiceNumber(fmt.Sprintf("INV-%d", time.Now().Unix())),
		PaymentID:  paymentID,
		CustomerID: customerID,
		Amount:     amount,
		Status:     "issued",
		IssuedAt:   time.Now(),
	}
}
