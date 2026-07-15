package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/context"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/customer"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/favorite"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/preference"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/domain/profile"
	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/infrastructure/repository"
)

type CustomerService struct {
	customerRepo  repository.CustomerRepository
	profileRepo   repository.ProfileRepository
	prefRepo      repository.PreferenceRepository
	favoriteRepo  repository.FavoritePlaceRepository
	contextRepo   repository.CustomerContextRepository
	logger        *slog.Logger
}

func NewCustomerService(
	customerRepo repository.CustomerRepository,
	profileRepo repository.ProfileRepository,
	prefRepo repository.PreferenceRepository,
	favoriteRepo repository.FavoritePlaceRepository,
	contextRepo repository.CustomerContextRepository,
	logger *slog.Logger,
) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		profileRepo:  profileRepo,
		prefRepo:     prefRepo,
		favoriteRepo: favoriteRepo,
		contextRepo:  contextRepo,
		logger:       logger,
	}
}

func (s *CustomerService) GetProfile(ctx context.Context, customerID string) (*profile.Profile, error) {
	return s.profileRepo.FindByCustomerID(ctx, customerID)
}

func (s *CustomerService) UpdateProfile(ctx context.Context, customerID string, p *profile.Profile) error {
	return s.profileRepo.Save(ctx, p)
}

func (s *CustomerService) GetPreferences(ctx context.Context, customerID string) (*preference.Preferences, error) {
	return s.prefRepo.FindByCustomerID(ctx, customerID)
}

func (s *CustomerService) UpdatePreferences(ctx context.Context, customerID string, update preference.PreferenceUpdate) error {
	prefs, err := s.prefRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return fmt.Errorf("find preferences: %w", err)
	}
	prefs.Apply(update)
	return s.prefRepo.Save(ctx, prefs)
}

func (s *CustomerService) GetFavorites(ctx context.Context, customerID string) ([]favorite.FavoritePlace, error) {
	return s.favoriteRepo.FindByCustomerID(ctx, customerID)
}

func (s *CustomerService) AddFavorite(ctx context.Context, customerID string, place *favorite.FavoritePlace) error {
	return s.favoriteRepo.Save(ctx, place)
}

func (s *CustomerService) RemoveFavorite(ctx context.Context, customerID string, placeID favorite.FavoritePlaceID) error {
	return s.favoriteRepo.Delete(ctx, placeID)
}

func (s *CustomerService) GetContext(ctx context.Context, customerID string) (*context.CustomerContext, error) {
	return s.contextRepo.FindByCustomerID(ctx, customerID)
}

func (s *CustomerService) UpdateContext(ctx context.Context, ctxData *context.CustomerContext) error {
	return s.contextRepo.Save(ctx, ctxData)
}
