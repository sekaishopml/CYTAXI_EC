package query

import "github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"

type GetFare struct {
	FareID valueobject.FareID
}

type GetFareHistory struct {
	TripID string
}

type GetPromotions struct {
	Active bool
}

type GetCoupons struct {
	Active bool
}

type PreviewFare struct {
	DistanceKM float64
	DurationSec int
	Region     string
}
