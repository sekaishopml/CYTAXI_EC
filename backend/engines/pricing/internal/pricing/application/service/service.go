package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/fare"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/promotion"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/coupon"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/infrastructure/repository"
)

type PricingService struct {
	fareRepo       repository.FareRepository
	promotionRepo  repository.PromotionRepository
	couponRepo     repository.CouponRepository
	logger         *slog.Logger
}

func NewPricingService(
	fareRepo repository.FareRepository,
	promotionRepo repository.PromotionRepository,
	couponRepo repository.CouponRepository,
	logger *slog.Logger,
) *PricingService {
	return &PricingService{
		fareRepo:      fareRepo,
		promotionRepo: promotionRepo,
		couponRepo:    couponRepo,
		logger:        logger,
	}
}

func (s *PricingService) CalculateFare(ctx context.Context, cmd command.CalculateFare) (*fare.Fare, error) {
	f := fare.NewFare(cmd.TripID)

	f.SetBaseFare(valueobject.Money{Amount: 1.00, Currency: "USD"})
	f.SetDistanceFare(valueobject.Money{Amount: cmd.DistanceKM * 0.50, Currency: "USD"})
	f.SetTimeFare(valueobject.Money{Amount: float64(cmd.DurationSec) * 0.02, Currency: "USD"})

	if cmd.WaitingSec > 0 {
		f.SetWaitingFare(valueobject.Money{Amount: float64(cmd.WaitingSec) * 0.015, Currency: "USD"})
	}
	if cmd.IsNight {
		f.ApplyNightSurcharge(valueobject.Money{Amount: 1.50, Currency: "USD"})
	}
	if cmd.DemandLevel >= 4 {
		f.ApplyDemandSurcharge(valueobject.Money{Amount: 2.0, Currency: "USD"})
	}

	if err := s.fareRepo.Save(ctx, f); err != nil {
		return nil, fmt.Errorf("save fare: %w", err)
	}
	return f, nil
}

func (s *PricingService) ApplyPromotion(ctx context.Context, cmd command.ApplyPromotion) (*fare.Fare, error) {
	f, err := s.fareRepo.FindByID(ctx, cmd.FareID)
	if err != nil {
		return nil, fmt.Errorf("find fare: %w", err)
	}

	promo, err := s.promotionRepo.FindByCode(ctx, cmd.PromotionCode)
	if err != nil {
		return nil, fmt.Errorf("find promotion: %w", err)
	}

	var amount valueobject.Money
	switch promo.Type {
	case promotion.TypePercentage:
		amount = valueobject.Money{Amount: valueobject.Percentage(promo.Value).Of(f.Components.Subtotal.Amount), Currency: f.Currency}
	case promotion.TypeFixed:
		amount = valueobject.Money{Amount: promo.Value, Currency: f.Currency}
	}

	f.ApplyPromotion(amount)
	if err := s.fareRepo.Update(ctx, f); err != nil {
		return nil, fmt.Errorf("update fare: %w", err)
	}
	return f, nil
}

func (s *PricingService) ApplyCoupon(ctx context.Context, cmd command.ApplyCoupon) (*fare.Fare, error) {
	c, err := s.couponRepo.FindByCode(ctx, cmd.CouponCode)
	if err != nil {
		return nil, fmt.Errorf("find coupon: %w", err)
	}
	if !c.CanUse() {
		return nil, fmt.Errorf("coupon %s cannot be used", cmd.CouponCode)
	}

	f, err := s.fareRepo.FindByID(ctx, cmd.FareID)
	if err != nil {
		return nil, fmt.Errorf("find fare: %w", err)
	}

	var amount valueobject.Money
	switch c.Type {
	case coupon.CouponTypePercentage:
		amount = valueobject.Money{Amount: valueobject.Percentage(c.Value).Of(f.Components.Subtotal.Amount), Currency: f.Currency}
	case coupon.CouponTypeFixed:
		amount = valueobject.Money{Amount: c.Value, Currency: f.Currency}
	}

	f.ApplyCoupon(amount)
	c.Use()

	if err := s.fareRepo.Update(ctx, f); err != nil {
		return nil, fmt.Errorf("update fare: %w", err)
	}
	return f, nil
}

func (s *PricingService) RemoveCoupon(ctx context.Context, cmd command.RemoveCoupon) (*fare.Fare, error) {
	f, err := s.fareRepo.FindByID(ctx, cmd.FareID)
	if err != nil {
		return nil, fmt.Errorf("find fare: %w", err)
	}
	f.ApplyCoupon(valueobject.Money{Currency: f.Currency})
	if err := s.fareRepo.Update(ctx, f); err != nil {
		return nil, fmt.Errorf("update fare: %w", err)
	}
	return f, nil
}

func (s *PricingService) CalculateTaxes(ctx context.Context, cmd command.CalculateTaxes) (*fare.Fare, error) {
	f, err := s.fareRepo.FindByID(ctx, cmd.FareID)
	if err != nil {
		return nil, fmt.Errorf("find fare: %w", err)
	}
	f.ApplyTax(12.0)
	if err := s.fareRepo.Update(ctx, f); err != nil {
		return nil, fmt.Errorf("update fare: %w", err)
	}
	return f, nil
}

func (s *PricingService) CalculateCommission(ctx context.Context, cmd command.CalculateCommission) (*fare.Fare, error) {
	f, err := s.fareRepo.FindByID(ctx, cmd.FareID)
	if err != nil {
		return nil, fmt.Errorf("find fare: %w", err)
	}
	f.ApplyCommission(15.0)
	if err := s.fareRepo.Update(ctx, f); err != nil {
		return nil, fmt.Errorf("update fare: %w", err)
	}
	return f, nil
}

func (s *PricingService) GetFare(ctx context.Context, q query.GetFare) (*fare.Fare, error) {
	return s.fareRepo.FindByID(ctx, q.FareID)
}

func (s *PricingService) GetFareHistory(ctx context.Context, q query.GetFareHistory) ([]fare.Fare, error) {
	return s.fareRepo.FindByTripID(ctx, q.TripID)
}

func (s *PricingService) GetPromotions(ctx context.Context, q query.GetPromotions) ([]promotion.Promotion, error) {
	return s.promotionRepo.FindAll(ctx)
}

func (s *PricingService) GetCoupons(ctx context.Context, q query.GetCoupons) ([]coupon.Coupon, error) {
	return s.couponRepo.FindAll(ctx)
}

func (s *PricingService) PreviewFare(ctx context.Context, q query.PreviewFare) (*fare.Fare, error) {
	cmd := command.CalculateFare{
		TripID:      "preview",
		DistanceKM:  q.DistanceKM,
		DurationSec: q.DurationSec,
		Region:      q.Region,
	}
	return s.CalculateFare(ctx, cmd)
}
