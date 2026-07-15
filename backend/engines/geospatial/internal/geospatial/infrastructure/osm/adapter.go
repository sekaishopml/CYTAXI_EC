package osm

import (
	"fmt"

	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/types"
)

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) Name() string {
	return "openstreetmap"
}

func (a *Adapter) Geocode(address string) ([]types.Coordinates, error) {
	return nil, fmt.Errorf("osm: geocode not implemented")
}

func (a *Adapter) ReverseGeocode(coords types.Coordinates) (*types.Address, error) {
	return nil, fmt.Errorf("osm: reverse geocode not implemented")
}

func (a *Adapter) FindRoute(req types.RouteRequest) ([]types.Route, error) {
	return nil, fmt.Errorf("osm: find route not implemented")
}

func (a *Adapter) GetDistanceMatrix(req types.DistanceMatrixRequest) (*types.DistanceMatrix, error) {
	return nil, fmt.Errorf("osm: distance matrix not implemented")
}

func (a *Adapter) SearchPlaces(req types.PlaceSearchRequest) (*types.PlaceSearchResult, error) {
	return nil, fmt.Errorf("osm: places search not implemented")
}

func (a *Adapter) Autocomplete(req types.AutocompleteRequest) ([]types.AutocompletePrediction, error) {
	return nil, fmt.Errorf("osm: autocomplete not implemented")
}

func (a *Adapter) PlaceDetails(placeID string) (*types.Place, error) {
	return nil, fmt.Errorf("osm: place details not implemented")
}
