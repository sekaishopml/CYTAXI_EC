package transaction

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
)

type Transaction struct {
	ID        valueobject.TransactionID `json:"id"`
	PaymentID valueobject.PaymentID    `json:"payment_id"`
	Type      TransactionType          `json:"type"`
	Amount    valueobject.Money        `json:"amount"`
	GatewayRef string                  `json:"gateway_ref,omitempty"`
	CreatedAt time.Time                `json:"created_at"`
}

type TransactionType string

const (
	TxnAuthorization TransactionType = "authorization"
	TxnCapture       TransactionType = "capture"
	TxnRefund        TransactionType = "refund"
	TxnSettlement    TransactionType = "settlement"
	TxnPayout        TransactionType = "payout"
)

func NewTransaction(paymentID valueobject.PaymentID, txnType TransactionType, amount valueobject.Money) *Transaction {
	return &Transaction{
		ID:        valueobject.NewTransactionID(),
		PaymentID: paymentID,
		Type:      txnType,
		Amount:    amount,
		CreatedAt: time.Now(),
	}
}
