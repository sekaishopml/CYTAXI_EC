package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/types"
)

type GeocodeInputPort interface {
	Geocode(ctx context.Context, address string) (*types.Coordinates, error)
	ReverseGeocode(ctx context.Context, coords types.Coordinates) (*types.Address, error)
}

type RouteInputPort interface {
	FindRoute(ctx context.Context, req types.RouteBuilderInput) (*types.RouteBuilderResult, error)
	GetDistanceMatrix(ctx context.Context, req types.DistanceMatrixRequest) (*types.DistanceMatrix, error)
}

type PlaceInputPort interface {
	Search(ctx context.Context, req types.PlaceSearchRequest) (*types.PlaceSearchResult, error)
	Autocomplete(ctx context.Context, req types.AutocompleteRequest) ([]types.AutocompletePrediction, error)
	Details(ctx context.Context, placeID string) (*types.Place, error)
}

type ZoneInputPort interface {
	CreateZone(ctx context.Context, zone types.Zone) error
	GetZone(ctx context.Context, id types.ZoneID) (*types.Zone, error)
	ListZones(ctx context.Context) ([]types.Zone, error)
	FindZonesByPoint(ctx context.Context, point types.Coordinates) ([]types.Zone, error)
	DeleteZone(ctx context.Context, id types.ZoneID) error
}

type GeospatialOutputPort interface {
	GeocodeOutputPort
	RouteOutputPort
}

type GeocodeOutputPort interface {
	OnGeocodeResult(ctx context.Context, coords *types.Coordinates) error
	OnReverseGeocodeResult(ctx context.Context, address *types.Address) error
}

type RouteOutputPort interface {
	OnRouteFound(ctx context.Context, result *types.RouteBuilderResult) error
}
