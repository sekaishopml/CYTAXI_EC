package repository

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/coupon"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/fare"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/promotion"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"
)

type FareRepository interface {
	FindByID(ctx context.Context, id valueobject.FareID) (*fare.Fare, error)
	FindByTripID(ctx context.Context, tripID string) ([]fare.Fare, error)
	Save(ctx context.Context, f *fare.Fare) error
	Update(ctx context.Context, f *fare.Fare) error
}

type PromotionRepository interface {
	FindByCode(ctx context.Context, code string) (*promotion.Promotion, error)
	FindAll(ctx context.Context) ([]promotion.Promotion, error)
	Save(ctx context.Context, p *promotion.Promotion) error
}

type CouponRepository interface {
	FindByCode(ctx context.Context, code valueobject.CouponCode) (*coupon.Coupon, error)
	FindAll(ctx context.Context) ([]coupon.Coupon, error)
	Save(ctx context.Context, c *coupon.Coupon) error
}
