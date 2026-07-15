package method

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/payment/internal/payment/domain/valueobject"
)

type PaymentMethod struct {
	ID      valueobject.PaymentMethodID   `json:"id"`
	OwnerID string                        `json:"owner_id"`
	Type    valueobject.PaymentMethodType `json:"type"`
	Details PaymentMethodDetails          `json:"details"`
	Default bool                          `json:"default"`
	Active  bool                          `json:"active"`
	CreatedAt time.Time                   `json:"created_at"`
}

type PaymentMethodDetails struct {
	Last4     string `json:"last_4,omitempty"`
	Brand     string `json:"brand,omitempty"`
	ExpMonth  int    `json:"exp_month,omitempty"`
	ExpYear   int    `json:"exp_year,omitempty"`
	BankName  string `json:"bank_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
}
