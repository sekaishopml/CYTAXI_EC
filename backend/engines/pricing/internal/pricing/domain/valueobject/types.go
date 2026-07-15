package valueobject

import (
	"fmt"
	"time"
)

type FareID string
type PromotionID string
type CouponCode string

type Money struct {
	Amount   float64
	Currency string
}

type Percentage float64

type TaxRate float64

func NewFareID() FareID { return FareID(fmt.Sprintf("fare_%d", time.Now().UnixNano())) }

func (m Money) Add(other Money) Money {
	return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}
}

func (m Money) Multiply(factor float64) Money {
	return Money{Amount: m.Amount * factor, Currency: m.Currency}
}

func (m Money) IsZero() bool { return m.Amount == 0 }

func (m Money) Round(decimals int) Money {
	pow := 1.0
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return Money{Amount: float64(int(m.Amount*pow)) / pow, Currency: m.Currency}
}

func (p Percentage) Of(amount float64) float64 {
	return amount * float64(p) / 100
}

func (p Percentage) Decimal() float64 {
	return float64(p) / 100
}

func (t TaxRate) Apply(amount float64) float64 {
	return amount * float64(t) / 100
}

type FareComponents struct {
	BaseFare       Money
	DistanceFare   Money
	TimeFare       Money
	WaitingFare    Money
	NightSurcharge Money
	DemandSurcharge Money
	Subtotal       Money
	Discount       Money
	Promotion      Money
	Coupon         Money
	Taxable        Money
	Tax            Money
	Commission     Money
	Total          Money
	DriverEarnings Money
}
