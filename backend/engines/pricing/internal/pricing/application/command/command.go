package command

import "github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"

type CalculateFare struct {
	TripID   string
	DistanceKM float64
	DurationSec int
	WaitingSec  int
	IsNight     bool
	DemandLevel int // 1-5
	Region      string
}

type ApplyPromotion struct {
	TripID      string
	FareID      valueobject.FareID
	PromotionCode string
}

type ApplyCoupon struct {
	TripID string
	FareID valueobject.FareID
	CouponCode valueobject.CouponCode
}

type RemoveCoupon struct {
	TripID     string
	FareID     valueobject.FareID
	CouponCode valueobject.CouponCode
}

type CalculateTaxes struct {
	TripID  string
	FareID  valueobject.FareID
	Region  string
}

type CalculateCommission struct {
	TripID   string
	FareID   valueobject.FareID
}
