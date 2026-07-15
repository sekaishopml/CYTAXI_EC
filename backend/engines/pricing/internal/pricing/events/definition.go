package events

const (
	EventFareCalculated     = "pricing.fare_calculated"
	EventFareUpdated        = "pricing.fare_updated"
	EventPromotionApplied   = "pricing.promotion_applied"
	EventPromotionRemoved   = "pricing.promotion_removed"
	EventCouponApplied      = "pricing.coupon_applied"
	EventCouponRejected     = "pricing.coupon_rejected"
	EventTaxesCalculated    = "pricing.taxes_calculated"
	EventCommissionCalculated = "pricing.commission_calculated"
)

type FareCalculatedPayload struct {
	TripID    string  `json:"trip_id"`
	FareID    string  `json:"fare_id"`
	BaseFare  float64 `json:"base_fare"`
	Total     float64 `json:"total"`
	Currency  string  `json:"currency"`
}

type PromotionAppliedPayload struct {
	TripID    string  `json:"trip_id"`
	FareID    string  `json:"fare_id"`
	PromoCode string  `json:"promo_code"`
	Amount    float64 `json:"amount"`
}

type CouponAppliedPayload struct {
	TripID     string  `json:"trip_id"`
	FareID     string  `json:"fare_id"`
	CouponCode string  `json:"coupon_code"`
	Amount     float64 `json:"amount"`
}

type CouponRejectedPayload struct {
	TripID     string `json:"trip_id"`
	CouponCode string `json:"coupon_code"`
	Reason     string `json:"reason"`
}

type TaxesCalculatedPayload struct {
	TripID     string  `json:"trip_id"`
	FareID     string  `json:"fare_id"`
	Taxable    float64 `json:"taxable_amount"`
	TaxAmount  float64 `json:"tax_amount"`
	Region     string  `json:"region"`
}

type CommissionCalculatedPayload struct {
	TripID   string  `json:"trip_id"`
	FareID   string  `json:"fare_id"`
	Amount   float64 `json:"amount"`
	Rate     float64 `json:"rate"`
}
