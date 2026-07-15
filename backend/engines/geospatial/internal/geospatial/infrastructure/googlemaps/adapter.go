package googlemaps

import (
	"fmt"

	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/types"
)

type Adapter struct {
	apiKey string
}

func NewAdapter(apiKey string) *Adapter {
	return &Adapter{apiKey: apiKey}
}

func (a *Adapter) Name() string {
	return "google_maps"
}

func (a *Adapter) Geocode(address string) ([]types.Coordinates, error) {
	return nil, fmt.Errorf("google_maps: geocode not implemented")
}

func (a *Adapter) ReverseGeocode(coords types.Coordinates) (*types.Address, error) {
	return nil, fmt.Errorf("google_maps: reverse geocode not implemented")
}

func (a *Adapter) FindRoute(req types.RouteRequest) ([]types.Route, error) {
	return nil, fmt.Errorf("google_maps: find route not implemented")
}

func (a *Adapter) GetDistanceMatrix(req types.DistanceMatrixRequest) (*types.DistanceMatrix, error) {
	return nil, fmt.Errorf("google_maps: distance matrix not implemented")
}

func (a *Adapter) SearchPlaces(req types.PlaceSearchRequest) (*types.PlaceSearchResult, error) {
	return nil, fmt.Errorf("google_maps: places search not implemented")
}

func (a *Adapter) Autocomplete(req types.AutocompleteRequest) ([]types.AutocompletePrediction, error) {
	return nil, fmt.Errorf("google_maps: autocomplete not implemented")
}

func (a *Adapter) PlaceDetails(placeID string) (*types.Place, error) {
	return nil, fmt.Errorf("google_maps: place details not implemented")
}
