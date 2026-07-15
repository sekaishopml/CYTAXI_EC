package command

import "github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"

type CreatePayment struct {
	TripID     string
	CustomerID string
	Amount     valueobject.Money
	Method     valueobject.PaymentMethodType
}

type AuthorizePayment struct {
	PaymentID valueobject.PaymentID
}

type CapturePayment struct {
	PaymentID valueobject.PaymentID
}

type CancelPayment struct {
	PaymentID valueobject.PaymentID
	Reason    string
}

type RefundPayment struct {
	PaymentID valueobject.PaymentID
	Amount    valueobject.Money
	Reason    string
}

type CreateSettlement struct {
	PaymentID  valueobject.PaymentID
	DriverID   string
	Amount     valueobject.Money
	Commission valueobject.Money
}

type PayDriver struct {
	SettlementID valueobject.SettlementID
}

type GenerateInvoice struct {
	PaymentID valueobject.PaymentID
}

type GenerateReceipt struct {
	PaymentID valueobject.PaymentID
}
