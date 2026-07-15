package fare

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"
)

type Fare struct {
	ID          valueobject.FareID  `json:"id"`
	TripID      string              `json:"trip_id"`
	Status      FareStatus          `json:"status"`
	Components  valueobject.FareComponents `json:"components"`
	Currency    string              `json:"currency"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type FareStatus string

const (
	FarePending    FareStatus = "pending"
	FareEstimated  FareStatus = "estimated"
	FareFinal      FareStatus = "final"
	FareDisputed   FareStatus = "disputed"
)

func NewFare(tripID string) *Fare {
	now := time.Now()
	return &Fare{
		ID:        valueobject.NewFareID(),
		TripID:    tripID,
		Status:    FarePending,
		Currency:  "USD",
		Components: valueobject.FareComponents{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (f *Fare) SetBaseFare(base valueobject.Money) {
	f.Components.BaseFare = base
	f.recalculate()
}

func (f *Fare) SetDistanceFare(dist valueobject.Money) {
	f.Components.DistanceFare = dist
	f.recalculate()
}

func (f *Fare) SetTimeFare(time valueobject.Money) {
	f.Components.TimeFare = time
	f.recalculate()
}

func (f *Fare) SetWaitingFare(waiting valueobject.Money) {
	f.Components.WaitingFare = waiting
	f.recalculate()
}

func (f *Fare) ApplyNightSurcharge(surcharge valueobject.Money) {
	f.Components.NightSurcharge = surcharge
	f.recalculate()
}

func (f *Fare) ApplyDemandSurcharge(surcharge valueobject.Money) {
	f.Components.DemandSurcharge = surcharge
	f.recalculate()
}

func (f *Fare) ApplyPromotion(amount valueobject.Money) {
	f.Components.Promotion = amount
	f.recalculate()
}

func (f *Fare) ApplyCoupon(amount valueobject.Money) {
	f.Components.Coupon = amount
	f.recalculate()
}

func (f *Fare) ApplyTax(rate valueobject.TaxRate) {
	taxable := f.Components.Subtotal
	f.Components.Discount = f.Components.Promotion
	f.Components.Discount = f.Components.Discount.Add(f.Components.Coupon)
	f.Components.Taxable = taxable
	f.Components.Tax = valueobject.Money{
		Amount:   rate.Apply(taxable.Amount),
		Currency: f.Currency,
	}
	f.recalculate()
}

func (f *Fare) ApplyCommission(rate valueobject.Percentage) {
	f.Components.Commission = valueobject.Money{
		Amount:   rate.Of(f.Components.Subtotal.Amount),
		Currency: f.Currency,
	}
	f.recalculate()
}

func (f *Fare) Finalize() {
	f.Status = FareFinal
	f.UpdatedAt = time.Now()
}

func (f *Fare) recalculate() {
	f.Components.Subtotal = valueobject.Money{Currency: f.Currency}
	f.Components.Subtotal = f.Components.Subtotal.Add(f.Components.BaseFare)
	f.Components.Subtotal = f.Components.Subtotal.Add(f.Components.DistanceFare)
	f.Components.Subtotal = f.Components.Subtotal.Add(f.Components.TimeFare)
	f.Components.Subtotal = f.Components.Subtotal.Add(f.Components.WaitingFare)
	f.Components.Subtotal = f.Components.Subtotal.Add(f.Components.NightSurcharge)
	f.Components.Subtotal = f.Components.Subtotal.Add(f.Components.DemandSurcharge)

	totalDisc := f.Components.Promotion.Add(f.Components.Coupon)
	afterDisc := f.Components.Subtotal
	afterDisc.Amount -= totalDisc.Amount
	if afterDisc.Amount < 0 {
		afterDisc.Amount = 0
	}

	f.Components.Total = afterDisc.Add(f.Components.Tax)
	f.Components.DriverEarnings = valueobject.Money{
		Amount:   afterDisc.Amount - f.Components.Commission.Amount,
		Currency: f.Currency,
	}
	if f.Components.DriverEarnings.Amount < 0 {
		f.Components.DriverEarnings.Amount = 0
	}

	f.UpdatedAt = time.Now()
}
