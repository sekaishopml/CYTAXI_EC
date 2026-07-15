package promotion

import (
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"
)

type Promotion struct {
	ID          valueobject.PromotionID `json:"id"`
	Code        string                  `json:"code"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Type        PromotionType           `json:"type"`
	Value       float64                 `json:"value"`
	Enabled     bool                    `json:"enabled"`
	MinFare     float64                 `json:"min_fare,omitempty"`
	MaxDiscount float64                 `json:"max_discount,omitempty"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

type PromotionType string

const (
	TypePercentage PromotionType = "percentage"
	TypeFixed      PromotionType = "fixed"
	TypeFreeRide   PromotionType = "free_ride"
)

func NewPromotion(code, name string, promType PromotionType, value float64) *Promotion {
	now := time.Now()
	return &Promotion{
		ID:        valueobject.PromotionID(fmt.Sprintf("promo_%d", now.UnixNano())),
		Code:      code,
		Name:      name,
		Type:      promType,
		Value:     value,
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
