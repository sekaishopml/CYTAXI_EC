package tenant

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("tenant: not found")

type Repository interface {
	GetByID(ctx context.Context, id ID) (*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)
	GetByDomain(ctx context.Context, domain string) (*Tenant, error)
	List(ctx context.Context) ([]*Tenant, error)
	Save(ctx context.Context, t *Tenant) error
	Update(ctx context.Context, t *Tenant) error
	Delete(ctx context.Context, id ID) error
}

type InMemoryRepository struct {
	tenants map[ID]*Tenant
}

func NewInMemoryRepository() *InMemoryRepository {
	repo := &InMemoryRepository{tenants: make(map[ID]*Tenant)}
	repo.seed()
	return repo
}

func (r *InMemoryRepository) seed() {
	r.tenants["tenant_cytaxi"] = &Tenant{
		ID: "tenant_cytaxi", Name: "CYTAXI Cooperativa", Slug: "cytaxi",
		Plan: PlanEnterprise, IsActive: true, MaxDrivers: 1000, MaxVehicles: 1200,
		Locale: "es", Timezone: "America/Guayaquil", Domain: "cytaxi.app",
		Branding: Branding{PrimaryColor: "#00a152", SecondaryColor: "#121212", AppName: "CYTAXI"},
		Features: []string{"all"}, CreatedAt: 1700000000, UpdatedAt: 1700000000,
	}
	r.tenants["tenant_demo"] = &Tenant{
		ID: "tenant_demo", Name: "Demo Cooperativa", Slug: "demo",
		Plan: PlanStarter, IsActive: true, MaxDrivers: 50, MaxVehicles: 60,
		Locale: "es", Timezone: "America/Guayaquil", Domain: "demo.cytaxi.app",
		Branding: Branding{PrimaryColor: "#2563eb", SecondaryColor: "#1e293b", AppName: "DemoTaxi"},
		Features: []string{"ride_hailing", "analytics"}, CreatedAt: 1700000000, UpdatedAt: 1700000000,
	}
}

func (r *InMemoryRepository) GetByID(_ context.Context, id ID) (*Tenant, error) {
	t, ok := r.tenants[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (r *InMemoryRepository) GetBySlug(_ context.Context, slug string) (*Tenant, error) {
	for _, t := range r.tenants {
		if t.Slug == slug {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

func (r *InMemoryRepository) GetByDomain(_ context.Context, domain string) (*Tenant, error) {
	for _, t := range r.tenants {
		if t.Domain == domain {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

func (r *InMemoryRepository) List(_ context.Context) ([]*Tenant, error) {
	list := make([]*Tenant, 0, len(r.tenants))
	for _, t := range r.tenants {
		list = append(list, t)
	}
	return list, nil
}

func (r *InMemoryRepository) Save(_ context.Context, t *Tenant) error {
	r.tenants[t.ID] = t
	return nil
}

func (r *InMemoryRepository) Update(_ context.Context, t *Tenant) error {
	r.tenants[t.ID] = t
	return nil
}

func (r *InMemoryRepository) Delete(_ context.Context, id ID) error {
	delete(r.tenants, id)
	return nil
}
