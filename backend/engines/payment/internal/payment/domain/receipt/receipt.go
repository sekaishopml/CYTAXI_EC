package receipt

import (
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
)

type Receipt struct {
	ID          valueobject.ReceiptNumber `json:"id"`
	PaymentID   valueobject.PaymentID     `json:"payment_id"`
	CustomerID  string                    `json:"customer_id"`
	Amount      valueobject.Money         `json:"amount"`
	Method      valueobject.PaymentMethodType `json:"method"`
	GatewayRef  string                    `json:"gateway_ref,omitempty"`
	GeneratedAt time.Time                 `json:"generated_at"`
}

func NewReceipt(paymentID valueobject.PaymentID, customerID string, amount valueobject.Money, method valueobject.PaymentMethodType) *Receipt {
	return &Receipt{
		ID:          valueobject.ReceiptNumber(fmt.Sprintf("RCP-%d", time.Now().Unix())),
		PaymentID:   paymentID,
		CustomerID:  customerID,
		Amount:      amount,
		Method:      method,
		GeneratedAt: time.Now(),
	}
}
