package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/invoice"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/payment"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/receipt"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/refund"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/settlement"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/transaction"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/wallet"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, cmd command.CreatePayment) (*payment.Payment, error)
	Authorize(ctx context.Context, cmd command.AuthorizePayment) (*payment.Payment, error)
	Capture(ctx context.Context, cmd command.CapturePayment) (*payment.Payment, error)
	Cancel(ctx context.Context, cmd command.CancelPayment) error
	Refund(ctx context.Context, cmd command.RefundPayment) (*refund.Refund, error)
	CreateSettlement(ctx context.Context, cmd command.CreateSettlement) (*settlement.Settlement, error)
	PayDriver(ctx context.Context, cmd command.PayDriver) error
	GenerateInvoice(ctx context.Context, cmd command.GenerateInvoice) (*invoice.Invoice, error)
	GenerateReceipt(ctx context.Context, cmd command.GenerateReceipt) (*receipt.Receipt, error)
	GetPayment(ctx context.Context, q query.GetPayment) (*payment.Payment, error)
	GetPayments(ctx context.Context, q query.GetPayments) ([]payment.Payment, error)
	GetTransaction(ctx context.Context, q query.GetTransaction) (*transaction.Transaction, error)
	GetWallet(ctx context.Context, q query.GetWallet) (*wallet.Wallet, error)
	GetRefund(ctx context.Context, q query.GetRefund) (*refund.Refund, error)
	GetSettlement(ctx context.Context, q query.GetSettlement) (*settlement.Settlement, error)
	GetInvoice(ctx context.Context, q query.GetInvoice) (*invoice.Invoice, error)
	GetReceipt(ctx context.Context, q query.GetReceipt) (*receipt.Receipt, error)
}
