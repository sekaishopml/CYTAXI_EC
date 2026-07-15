package valueobject

import (
	"fmt"
	"time"
)

type PaymentID string
type TransactionID string
type WalletID string
type RefundID string
type SettlementID string
type InvoiceNumber string
type ReceiptNumber string
type PaymentMethodID string

type Money struct {
	Amount   float64
	Currency string
}

type Percentage float64

type PaymentStatus string

const (
	PayStatusPending    PaymentStatus = "pending"
	PayStatusAuthorized PaymentStatus = "authorized"
	PayStatusCaptured   PaymentStatus = "captured"
	PayStatusFailed     PaymentStatus = "failed"
	PayStatusCancelled  PaymentStatus = "cancelled"
	PayStatusRefunded   PaymentStatus = "refunded"
	PayStatusPartialRefund PaymentStatus = "partial_refund"
)

type PaymentMethodType string

const (
	PMCard       PaymentMethodType = "card"
	PMWallet     PaymentMethodType = "wallet"
	PMCash       PaymentMethodType = "cash"
	PMTransfer   PaymentMethodType = "transfer"
	PMBNPL       PaymentMethodType = "bnpl"
)

type SettlementStatus string

const (
	SettStatusPending   SettlementStatus = "pending"
	SettStatusProcessing SettlementStatus = "processing"
	SettStatusCompleted  SettlementStatus = "completed"
	SettStatusFailed     SettlementStatus = "failed"
)

type RefundStatus string

const (
	RefStatusPending    RefundStatus = "pending"
	RefStatusProcessing RefundStatus = "processing"
	RefStatusCompleted  RefundStatus = "completed"
	RefStatusRejected   RefundStatus = "rejected"
)

func NewPaymentID() PaymentID       { return PaymentID(fmt.Sprintf("pay_%d", time.Now().UnixNano())) }
func NewTransactionID() TransactionID { return TransactionID(fmt.Sprintf("txn_%d", time.Now().UnixNano())) }
func NewWalletID() WalletID         { return WalletID(fmt.Sprintf("wal_%d", time.Now().UnixNano())) }
func NewRefundID() RefundID         { return RefundID(fmt.Sprintf("ref_%d", time.Now().UnixNano())) }
func NewSettlementID() SettlementID { return SettlementID(fmt.Sprintf("stl_%d", time.Now().UnixNano())) }

func (m Money) Add(o Money) Money           { return Money{Amount: m.Amount + o.Amount, Currency: m.Currency} }
func (m Money) Subtract(o Money) Money      { return Money{Amount: m.Amount - o.Amount, Currency: m.Currency} }
func (m Money) Multiply(f float64) Money    { return Money{Amount: m.Amount * f, Currency: m.Currency} }
func (m Money) IsZero() bool                { return m.Amount == 0 }
func (m Money) IsNegative() bool            { return m.Amount < 0 }

func (p Percentage) Of(amount float64) float64 { return amount * float64(p) / 100 }
