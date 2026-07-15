package tax

import "github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"

type Tax struct {
	Name      string           `json:"name"`
	Rate      valueobject.TaxRate `json:"rate"`
	Country   string           `json:"country"`
	Region    string           `json:"region,omitempty"`
	Enabled   bool             `json:"enabled"`
}

func NewTax(name string, rate valueobject.TaxRate, country string) *Tax {
	return &Tax{
		Name:    name,
		Rate:    rate,
		Country: country,
		Enabled: true,
	}
}
