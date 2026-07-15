package wallet

import (
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
)

type Wallet struct {
	ID         valueobject.WalletID `json:"id"`
	OwnerID    string               `json:"owner_id"`
	OwnerType  string               `json:"owner_type"` // customer, driver
	Balance    valueobject.Money    `json:"balance"`
	Currency   string               `json:"currency"`
	Active     bool                 `json:"active"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

func NewWallet(ownerID, ownerType string, currency string) *Wallet {
	now := time.Now()
	return &Wallet{
		ID:        valueobject.NewWalletID(),
		OwnerID:   ownerID,
		OwnerType: ownerType,
		Balance:   valueobject.Money{Currency: currency},
		Currency:  currency,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (w *Wallet) Credit(amount valueobject.Money) {
	w.Balance = w.Balance.Add(amount)
	w.UpdatedAt = time.Now()
}

func (w *Wallet) Debit(amount valueobject.Money) error {
	if w.Balance.Amount < amount.Amount {
		return fmt.Errorf("insufficient balance: have %.2f, need %.2f", w.Balance.Amount, amount.Amount)
	}
	w.Balance = w.Balance.Subtract(amount)
	w.UpdatedAt = time.Now()
	return nil
}
