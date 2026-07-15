package usecase

import (
	"context"
	"fmt"

	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/service"
	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/types"
)

type PlaceUseCase struct {
	provider service.GeospatialProvider
}

func NewPlaceUseCase(provider service.GeospatialProvider) *PlaceUseCase {
	return &PlaceUseCase{provider: provider}
}

func (uc *PlaceUseCase) Search(ctx context.Context, req types.PlaceSearchRequest) (*types.PlaceSearchResult, error) {
	if req.Query == "" {
		return nil, fmt.Errorf("place search: query required")
	}
	return uc.provider.SearchPlaces(req)
}

func (uc *PlaceUseCase) Autocomplete(ctx context.Context, req types.AutocompleteRequest) ([]types.AutocompletePrediction, error) {
	if req.Input == "" {
		return nil, fmt.Errorf("autocomplete: input required")
	}
	return uc.provider.Autocomplete(req)
}

func (uc *PlaceUseCase) Details(ctx context.Context, placeID string) (*types.Place, error) {
	if placeID == "" {
		return nil, fmt.Errorf("place details: id required")
	}
	return uc.provider.PlaceDetails(placeID)
}
