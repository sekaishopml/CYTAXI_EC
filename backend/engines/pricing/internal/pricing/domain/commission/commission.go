package commission

import "github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"

type Commission struct {
	Name      string             `json:"name"`
	Rate      valueobject.Percentage `json:"rate"`
	MinAmount float64            `json:"min_amount"`
	Enabled   bool               `json:"enabled"`
}

func NewCommission(name string, rate valueobject.Percentage, minAmount float64) *Commission {
	return &Commission{
		Name:      name,
		Rate:      rate,
		MinAmount: minAmount,
		Enabled:   true,
	}
}

func (c *Commission) Calculate(amount float64) float64 {
	result := c.Rate.Of(amount)
	if result < c.MinAmount {
		return c.MinAmount
	}
	return result
}
