package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/customer"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/context"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/favorite"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/preference"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/profile"
)

type ProfileInputPort interface {
	GetProfile(ctx context.Context, customerID string) (*profile.Profile, error)
	UpdateProfile(ctx context.Context, customerID string, p *profile.Profile) error
}

type PreferenceInputPort interface {
	GetPreferences(ctx context.Context, customerID string) (*preference.Preferences, error)
	UpdatePreferences(ctx context.Context, customerID string, update preference.PreferenceUpdate) error
}

type FavoritePlaceInputPort interface {
	GetFavorites(ctx context.Context, customerID string) ([]favorite.FavoritePlace, error)
	AddFavorite(ctx context.Context, customerID string, place *favorite.FavoritePlace) error
	RemoveFavorite(ctx context.Context, customerID string, placeID favorite.FavoritePlaceID) error
}

type CustomerContextInputPort interface {
	GetContext(ctx context.Context, customerID string) (*context.CustomerContext, error)
	UpdateContext(ctx context.Context, ctxData *context.CustomerContext) error
}

type CustomerProfileOutputPort interface {
	OnProfileChanged(ctx context.Context, customerID string) error
	OnPreferencesChanged(ctx context.Context, customerID string) error
	OnFavoritePlaceAdded(ctx context.Context, customerID string, placeID favorite.FavoritePlaceID) error
}

type CustomerService interface {
	ProfileInputPort
	PreferenceInputPort
	FavoritePlaceInputPort
	CustomerContextInputPort
}
