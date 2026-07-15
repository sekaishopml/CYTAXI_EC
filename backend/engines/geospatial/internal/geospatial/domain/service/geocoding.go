package service

import (
	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/types"
)

type Geocoder interface {
	Geocode(address string) ([]types.Coordinates, error)
	ReverseGeocode(coords types.Coordinates) (*types.Address, error)
}

type DirectionFinder interface {
	FindRoute(req types.RouteRequest) ([]types.Route, error)
	GetDistanceMatrix(req types.DistanceMatrixRequest) (*types.DistanceMatrix, error)
}

type PlaceSearcher interface {
	SearchPlaces(req types.PlaceSearchRequest) (*types.PlaceSearchResult, error)
	Autocomplete(req types.AutocompleteRequest) ([]types.AutocompletePrediction, error)
	PlaceDetails(placeID string) (*types.Place, error)
}

type GeospatialProvider interface {
	Geocoder
	DirectionFinder
	PlaceSearcher
	Name() string
}
