package repository

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/invoice"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/payment"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/receipt"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/refund"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/settlement"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/transaction"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/wallet"
)

type PaymentRepository interface {
	FindByID(ctx context.Context, id valueobject.PaymentID) (*payment.Payment, error)
	FindByTripID(ctx context.Context, tripID string) ([]payment.Payment, error)
	Save(ctx context.Context, p *payment.Payment) error
	Update(ctx context.Context, p *payment.Payment) error
}

type TransactionRepository interface {
	FindByID(ctx context.Context, id valueobject.TransactionID) (*transaction.Transaction, error)
	Save(ctx context.Context, t *transaction.Transaction) error
}

type SettlementRepository interface {
	FindByID(ctx context.Context, id valueobject.SettlementID) (*settlement.Settlement, error)
	Save(ctx context.Context, s *settlement.Settlement) error
	Update(ctx context.Context, s *settlement.Settlement) error
}

type RefundRepository interface {
	FindByID(ctx context.Context, id valueobject.RefundID) (*refund.Refund, error)
	Save(ctx context.Context, r *refund.Refund) error
	Update(ctx context.Context, r *refund.Refund) error
}

type WalletRepository interface {
	FindByOwnerID(ctx context.Context, ownerID string) (*wallet.Wallet, error)
	Save(ctx context.Context, w *wallet.Wallet) error
	Update(ctx context.Context, w *wallet.Wallet) error
}

type InvoiceRepository interface {
	FindByPaymentID(ctx context.Context, paymentID valueobject.PaymentID) (*invoice.Invoice, error)
	Save(ctx context.Context, i *invoice.Invoice) error
}

type ReceiptRepository interface {
	FindByPaymentID(ctx context.Context, paymentID valueobject.PaymentID) (*receipt.Receipt, error)
	Save(ctx context.Context, r *receipt.Receipt) error
}
