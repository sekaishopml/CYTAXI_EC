package usecase

import (
	"context"
	"fmt"

	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/service"
	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/types"
)

type GeocodeUseCase struct {
	provider service.GeospatialProvider
}

func NewGeocodeUseCase(provider service.GeospatialProvider) *GeocodeUseCase {
	return &GeocodeUseCase{provider: provider}
}

func (uc *GeocodeUseCase) Geocode(ctx context.Context, address string) (*types.Coordinates, error) {
	results, err := uc.provider.Geocode(address)
	if err != nil {
		return nil, fmt.Errorf("geocode: %w", err)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("geocode: no results for address: %s", address)
	}
	return &results[0], nil
}

func (uc *GeocodeUseCase) ReverseGeocode(ctx context.Context, coords types.Coordinates) (*types.Address, error) {
	addr, err := uc.provider.ReverseGeocode(coords)
	if err != nil {
		return nil, fmt.Errorf("reverse geocode: %w", err)
	}
	return addr, nil
}
