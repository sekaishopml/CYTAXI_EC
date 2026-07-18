package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/coupon"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/fare"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/promotion"
	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/domain/valueobject"
)

var ErrFareNotFound = errors.New("fare not found")
var ErrPromoNotFound = errors.New("promotion not found")
var ErrCouponNotFound = errors.New("coupon not found")

// InMemoryFareRepository
type InMemoryFareRepository struct {
	mu     sync.RWMutex
	byID   map[string]*fare.Fare
	byTrip map[string][]fare.Fare
}

func NewInMemoryFareRepository() *InMemoryFareRepository {
	return &InMemoryFareRepository{
		byID:   make(map[string]*fare.Fare),
		byTrip: make(map[string][]fare.Fare),
	}
}

func (r *InMemoryFareRepository) FindByID(ctx context.Context, id valueobject.FareID) (*fare.Fare, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.byID[string(id)]
	if !ok {
		return nil, ErrFareNotFound
	}
	return f, nil
}

func (r *InMemoryFareRepository) FindByTripID(ctx context.Context, tripID string) ([]fare.Fare, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]fare.Fare, len(r.byTrip[tripID]))
	copy(out, r.byTrip[tripID])
	return out, nil
}

func (r *InMemoryFareRepository) Save(ctx context.Context, f *fare.Fare) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[string(f.ID)] = f
	r.byTrip[f.TripID] = append(r.byTrip[f.TripID], *f)
	return nil
}

func (r *InMemoryFareRepository) Update(ctx context.Context, f *fare.Fare) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[string(f.ID)] = f
	return nil
}

// InMemoryPromotionRepository
type InMemoryPromotionRepository struct {
	mu   sync.RWMutex
	data map[string]*promotion.Promotion
}

func NewInMemoryPromotionRepository() *InMemoryPromotionRepository {
	return &InMemoryPromotionRepository{data: make(map[string]*promotion.Promotion)}
}

func (r *InMemoryPromotionRepository) FindByCode(ctx context.Context, code string) (*promotion.Promotion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.data {
		if string(p.Code) == code {
			return p, nil
		}
	}
	return nil, ErrPromoNotFound
}

func (r *InMemoryPromotionRepository) FindAll(ctx context.Context) ([]promotion.Promotion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]promotion.Promotion, 0, len(r.data))
	for _, p := range r.data {
		out = append(out, *p)
	}
	return out, nil
}

func (r *InMemoryPromotionRepository) Save(ctx context.Context, p *promotion.Promotion) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[string(p.Code)] = p
	return nil
}

// InMemoryCouponRepository
type InMemoryCouponRepository struct {
	mu   sync.RWMutex
	data map[string]*coupon.Coupon
}

func NewInMemoryCouponRepository() *InMemoryCouponRepository {
	return &InMemoryCouponRepository{data: make(map[string]*coupon.Coupon)}
}

func (r *InMemoryCouponRepository) FindByCode(ctx context.Context, code valueobject.CouponCode) (*coupon.Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.data[string(code)]
	if !ok {
		return nil, ErrCouponNotFound
	}
	return c, nil
}

func (r *InMemoryCouponRepository) FindAll(ctx context.Context) ([]coupon.Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]coupon.Coupon, 0, len(r.data))
	for _, c := range r.data {
		out = append(out, *c)
	}
	return out, nil
}

func (r *InMemoryCouponRepository) Save(ctx context.Context, c *coupon.Coupon) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[string(c.Code)] = c
	return nil
}
