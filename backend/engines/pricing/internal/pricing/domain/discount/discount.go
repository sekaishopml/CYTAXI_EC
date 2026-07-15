package discount

import (
	"fmt"
	"time"
)

type Discount struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        DiscountType `json:"type"`
	Value       float64 `json:"value"`
	Enabled     bool    `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
}

type DiscountType string

const (
	DiscountPercentage DiscountType = "percentage"
	DiscountFixed      DiscountType = "fixed"
)

func NewDiscount(name string, dType DiscountType, value float64) *Discount {
	return &Discount{
		ID:        fmt.Sprintf("disc_%d", time.Now().UnixNano()),
		Name:      name,
		Type:      dType,
		Value:     value,
		Enabled:   true,
		CreatedAt: time.Now(),
	}
}
