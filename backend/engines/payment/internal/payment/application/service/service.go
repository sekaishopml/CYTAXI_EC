package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/invoice"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/payment"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/receipt"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/refund"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/settlement"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/transaction"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/wallet"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/infrastructure/gateway"
	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/infrastructure/repository"
)

type PaymentService struct {
	payRepo      repository.PaymentRepository
	txnRepo      repository.TransactionRepository
	stlRepo      repository.SettlementRepository
	refRepo      repository.RefundRepository
	walRepo      repository.WalletRepository
	invRepo      repository.InvoiceRepository
	rcpRepo      repository.ReceiptRepository
	gateways     *gateway.GatewayRegistry
	logger       *slog.Logger
}

func NewPaymentService(
	payRepo repository.PaymentRepository,
	txnRepo repository.TransactionRepository,
	stlRepo repository.SettlementRepository,
	refRepo repository.RefundRepository,
	walRepo repository.WalletRepository,
	invRepo repository.InvoiceRepository,
	rcpRepo repository.ReceiptRepository,
	gateways *gateway.GatewayRegistry,
	logger *slog.Logger,
) *PaymentService {
	return &PaymentService{
		payRepo: payRepo, txnRepo: txnRepo, stlRepo: stlRepo,
		refRepo: refRepo, walRepo: walRepo, invRepo: invRepo,
		rcpRepo: rcpRepo, gateways: gateways, logger: logger,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, cmd command.CreatePayment) (*payment.Payment, error) {
	p := payment.NewPayment(cmd.TripID, cmd.CustomerID, cmd.Amount, cmd.Method)
	if err := s.payRepo.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("save payment: %w", err)
	}
	txn := transaction.NewTransaction(p.ID, transaction.TxnAuthorization, p.Amount)
	s.txnRepo.Save(ctx, txn)
	return p, nil
}

func (s *PaymentService) Authorize(ctx context.Context, cmd command.AuthorizePayment) (*payment.Payment, error) {
	p, err := s.payRepo.FindByID(ctx, cmd.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("find payment: %w", err)
	}
	gw, err := s.gateways.Get(p.Method)
	if err != nil {
		p.Fail()
		s.payRepo.Update(ctx, p)
		return p, err
	}
	result, err := gw.Authorize(ctx, p.Amount, p.Method, string(p.ID))
	if err != nil || !result.Success {
		p.Fail()
		s.payRepo.Update(ctx, p)
		return p, fmt.Errorf("authorize: %w", err)
	}
	p.Authorize(result.GatewayRef)
	s.payRepo.Update(ctx, p)
	txn := transaction.NewTransaction(p.ID, transaction.TxnAuthorization, p.Amount)
	txn.GatewayRef = result.GatewayRef
	s.txnRepo.Save(ctx, txn)
	return p, nil
}

func (s *PaymentService) Capture(ctx context.Context, cmd command.CapturePayment) (*payment.Payment, error) {
	p, err := s.payRepo.FindByID(ctx, cmd.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("find payment: %w", err)
	}
	gw, _ := s.gateways.Get(p.Method)
	if gw != nil {
		gw.Capture(ctx, p.GatewayRef, p.Amount)
	}
	p.Capture()
	s.payRepo.Update(ctx, p)
	txn := transaction.NewTransaction(p.ID, transaction.TxnCapture, p.Amount)
	s.txnRepo.Save(ctx, txn)
	return p, nil
}

func (s *PaymentService) Cancel(ctx context.Context, cmd command.CancelPayment) error {
	p, err := s.payRepo.FindByID(ctx, cmd.PaymentID)
	if err != nil {
		return fmt.Errorf("find payment: %w", err)
	}
	p.Cancel()
	return s.payRepo.Update(ctx, p)
}

func (s *PaymentService) Refund(ctx context.Context, cmd command.RefundPayment) (*refund.Refund, error) {
	r := refund.NewRefund(cmd.PaymentID, cmd.Amount, cmd.Reason)
	r.Process()
	p, _ := s.payRepo.FindByID(ctx, cmd.PaymentID)
	if p != nil {
		p.Refund()
		s.payRepo.Update(ctx, p)
	}
	r.Complete("ref_gateway_placeholder")
	s.refRepo.Save(ctx, r)
	return r, nil
}

func (s *PaymentService) CreateSettlement(ctx context.Context, cmd command.CreateSettlement) (*settlement.Settlement, error) {
	stl := settlement.NewSettlement(cmd.PaymentID, cmd.DriverID, cmd.Amount, cmd.Commission)
	stl.Process()
	if err := s.stlRepo.Save(ctx, stl); err != nil {
		return nil, fmt.Errorf("save settlement: %w", err)
	}
	return stl, nil
}

func (s *PaymentService) PayDriver(ctx context.Context, cmd command.PayDriver) error {
	stl, err := s.stlRepo.FindByID(ctx, cmd.SettlementID)
	if err != nil {
		return fmt.Errorf("find settlement: %w", err)
	}
	stl.Complete()
	return s.stlRepo.Update(ctx, stl)
}

func (s *PaymentService) GenerateInvoice(ctx context.Context, cmd command.GenerateInvoice) (*invoice.Invoice, error) {
	p, err := s.payRepo.FindByID(ctx, cmd.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("find payment: %w", err)
	}
	inv := invoice.NewInvoice(p.ID, p.CustomerID, p.Amount)
	s.invRepo.Save(ctx, inv)
	return inv, nil
}

func (s *PaymentService) GenerateReceipt(ctx context.Context, cmd command.GenerateReceipt) (*receipt.Receipt, error) {
	p, err := s.payRepo.FindByID(ctx, cmd.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("find payment: %w", err)
	}
	r := receipt.NewReceipt(p.ID, p.CustomerID, p.Amount, p.Method)
	s.rcpRepo.Save(ctx, r)
	return r, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, q query.GetPayment) (*payment.Payment, error) {
	return s.payRepo.FindByID(ctx, q.PaymentID)
}

func (s *PaymentService) GetPayments(ctx context.Context, q query.GetPayments) ([]payment.Payment, error) {
	if q.TripID != "" {
		return s.payRepo.FindByTripID(ctx, q.TripID)
	}
	return nil, nil
}

func (s *PaymentService) GetTransaction(ctx context.Context, q query.GetTransaction) (*transaction.Transaction, error) {
	return s.txnRepo.FindByID(ctx, q.TransactionID)
}

func (s *PaymentService) GetWallet(ctx context.Context, q query.GetWallet) (*wallet.Wallet, error) {
	return s.walRepo.FindByOwnerID(ctx, q.OwnerID)
}

func (s *PaymentService) GetRefund(ctx context.Context, q query.GetRefund) (*refund.Refund, error) {
	return s.refRepo.FindByID(ctx, q.RefundID)
}

func (s *PaymentService) GetSettlement(ctx context.Context, q query.GetSettlement) (*settlement.Settlement, error) {
	return s.stlRepo.FindByID(ctx, q.SettlementID)
}

func (s *PaymentService) GetInvoice(ctx context.Context, q query.GetInvoice) (*invoice.Invoice, error) {
	return s.invRepo.FindByPaymentID(ctx, q.PaymentID)
}

func (s *PaymentService) GetReceipt(ctx context.Context, q query.GetReceipt) (*receipt.Receipt, error) {
	return s.rcpRepo.FindByPaymentID(ctx, q.PaymentID)
}
