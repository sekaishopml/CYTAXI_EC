package mapbox

import (
	"fmt"

	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/types"
)

type Adapter struct {
	accessToken string
}

func NewAdapter(accessToken string) *Adapter {
	return &Adapter{accessToken: accessToken}
}

func (a *Adapter) Name() string {
	return "mapbox"
}

func (a *Adapter) Geocode(address string) ([]types.Coordinates, error) {
	return nil, fmt.Errorf("mapbox: geocode not implemented")
}

func (a *Adapter) ReverseGeocode(coords types.Coordinates) (*types.Address, error) {
	return nil, fmt.Errorf("mapbox: reverse geocode not implemented")
}

func (a *Adapter) FindRoute(req types.RouteRequest) ([]types.Route, error) {
	return nil, fmt.Errorf("mapbox: find route not implemented")
}

func (a *Adapter) GetDistanceMatrix(req types.DistanceMatrixRequest) (*types.DistanceMatrix, error) {
	return nil, fmt.Errorf("mapbox: distance matrix not implemented")
}

func (a *Adapter) SearchPlaces(req types.PlaceSearchRequest) (*types.PlaceSearchResult, error) {
	return nil, fmt.Errorf("mapbox: places search not implemented")
}

func (a *Adapter) Autocomplete(req types.AutocompleteRequest) ([]types.AutocompletePrediction, error) {
	return nil, fmt.Errorf("mapbox: autocomplete not implemented")
}

func (a *Adapter) PlaceDetails(placeID string) (*types.Place, error) {
	return nil, fmt.Errorf("mapbox: place details not implemented")
}
