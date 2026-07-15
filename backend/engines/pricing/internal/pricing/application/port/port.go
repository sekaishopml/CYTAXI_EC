package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/fare"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/promotion"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/coupon"
)

type PricingService interface {
	CalculateFare(ctx context.Context, cmd command.CalculateFare) (*fare.Fare, error)
	ApplyPromotion(ctx context.Context, cmd command.ApplyPromotion) (*fare.Fare, error)
	ApplyCoupon(ctx context.Context, cmd command.ApplyCoupon) (*fare.Fare, error)
	RemoveCoupon(ctx context.Context, cmd command.RemoveCoupon) (*fare.Fare, error)
	CalculateTaxes(ctx context.Context, cmd command.CalculateTaxes) (*fare.Fare, error)
	CalculateCommission(ctx context.Context, cmd command.CalculateCommission) (*fare.Fare, error)
	GetFare(ctx context.Context, q query.GetFare) (*fare.Fare, error)
	GetFareHistory(ctx context.Context, q query.GetFareHistory) ([]fare.Fare, error)
	GetPromotions(ctx context.Context, q query.GetPromotions) ([]promotion.Promotion, error)
	GetCoupons(ctx context.Context, q query.GetCoupons) ([]coupon.Coupon, error)
	PreviewFare(ctx context.Context, q query.PreviewFare) (*fare.Fare, error)
}
