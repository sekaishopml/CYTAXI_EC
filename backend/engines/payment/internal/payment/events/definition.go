package events

const (
	EventPaymentCreated      = "payment.created"
	EventPaymentAuthorized   = "payment.authorized"
	EventPaymentCaptured     = "payment.captured"
	EventPaymentCancelled    = "payment.cancelled"
	EventPaymentFailed       = "payment.failed"
	EventRefundCreated       = "payment.refund_created"
	EventRefundCompleted     = "payment.refund_completed"
	EventSettlementCreated   = "payment.settlement_created"
	EventDriverPaid          = "payment.driver_paid"
	EventInvoiceGenerated    = "payment.invoice_generated"
	EventReceiptGenerated    = "payment.receipt_generated"
)

type PaymentCreatedPayload struct {
	PaymentID  string  `json:"payment_id"`
	TripID     string  `json:"trip_id"`
	CustomerID string  `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	Method     string  `json:"method"`
}

type PaymentAuthorizedPayload struct {
	PaymentID  string `json:"payment_id"`
	GatewayRef string `json:"gateway_ref"`
}

type PaymentCapturedPayload struct {
	PaymentID  string  `json:"payment_id"`
	Amount     float64 `json:"amount"`
}

type RefundCreatedPayload struct {
	RefundID  string  `json:"refund_id"`
	PaymentID string  `json:"payment_id"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason"`
}

type RefundCompletedPayload struct {
	RefundID   string `json:"refund_id"`
	GatewayRef string `json:"gateway_ref"`
}

type SettlementCreatedPayload struct {
	SettlementID string  `json:"settlement_id"`
	PaymentID    string  `json:"payment_id"`
	DriverID     string  `json:"driver_id"`
	NetAmount    float64 `json:"net_amount"`
}

type DriverPaidPayload struct {
	SettlementID string  `json:"settlement_id"`
	DriverID     string  `json:"driver_id"`
	Amount       float64 `json:"amount"`
}

type InvoiceGeneratedPayload struct {
	InvoiceID  string  `json:"invoice_id"`
	PaymentID  string  `json:"payment_id"`
	Amount     float64 `json:"amount"`
}

type ReceiptGeneratedPayload struct {
	ReceiptID  string  `json:"receipt_id"`
	PaymentID  string  `json:"payment_id"`
	Amount     float64 `json:"amount"`
	Method     string  `json:"method"`
}
