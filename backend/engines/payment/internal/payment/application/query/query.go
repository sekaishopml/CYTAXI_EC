package query

import "github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"

type GetPayment struct {
	PaymentID valueobject.PaymentID
}

type GetPayments struct {
	TripID     string
	CustomerID string
	Limit      int
	Offset     int
}

type GetTransaction struct {
	TransactionID valueobject.TransactionID
}

type GetWallet struct {
	OwnerID string
}

type GetRefund struct {
	RefundID valueobject.RefundID
}

type GetSettlement struct {
	SettlementID valueobject.SettlementID
}

type GetInvoice struct {
	PaymentID valueobject.PaymentID
}

type GetReceipt struct {
	PaymentID valueobject.PaymentID
}
