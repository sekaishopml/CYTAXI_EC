package coupon

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"
)

type Coupon struct {
	Code       valueobject.CouponCode `json:"code"`
	Name       string                 `json:"name"`
	Type       CouponType             `json:"type"`
	Value      float64                `json:"value"`
	Enabled    bool                   `json:"enabled"`
	MaxUses    int                    `json:"max_uses"`
	UsedCount  int                    `json:"used_count"`
	ExpiresAt  *time.Time             `json:"expires_at,omitempty"`
	MinFare    float64                `json:"min_fare,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type CouponType string

const (
	CouponTypePercentage CouponType = "percentage"
	CouponTypeFixed      CouponType = "fixed"
)

func NewCoupon(code string, name string, cType CouponType, value float64) *Coupon {
	return &Coupon{
		Code:      valueobject.CouponCode(code),
		Name:      name,
		Type:      cType,
		Value:     value,
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (c *Coupon) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

func (c *Coupon) CanUse() bool {
	return c.Enabled && !c.IsExpired() && (c.MaxUses == 0 || c.UsedCount < c.MaxUses)
}

func (c *Coupon) Use() {
	c.UsedCount++
	c.UpdatedAt = time.Now()
}
