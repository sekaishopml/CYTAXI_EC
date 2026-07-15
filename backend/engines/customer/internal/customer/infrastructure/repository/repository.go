package repository

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/context"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/customer"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/favorite"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/preference"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/profile"
)

type CustomerRepository interface {
	FindByID(ctx context.Context, id customer.CustomerID) (*customer.Customer, error)
	FindByPhone(ctx context.Context, phone string) (*customer.Customer, error)
	Save(ctx context.Context, c *customer.Customer) error
	Update(ctx context.Context, c *customer.Customer) error
}

type ProfileRepository interface {
	FindByCustomerID(ctx context.Context, customerID string) (*profile.Profile, error)
	Save(ctx context.Context, p *profile.Profile) error
}

type PreferenceRepository interface {
	FindByCustomerID(ctx context.Context, customerID string) (*preference.Preferences, error)
	Save(ctx context.Context, p *preference.Preferences) error
}

type FavoritePlaceRepository interface {
	FindByCustomerID(ctx context.Context, customerID string) ([]favorite.FavoritePlace, error)
	FindByID(ctx context.Context, id favorite.FavoritePlaceID) (*favorite.FavoritePlace, error)
	Save(ctx context.Context, f *favorite.FavoritePlace) error
	Delete(ctx context.Context, id favorite.FavoritePlaceID) error
}

type CustomerContextRepository interface {
	FindByCustomerID(ctx context.Context, customerID string) (*context.CustomerContext, error)
	Save(ctx context.Context, c *context.CustomerContext) error
}
