package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/service"
	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/types"
)

type RouteUseCase struct {
	provider service.GeospatialProvider
}

func NewRouteUseCase(provider service.GeospatialProvider) *RouteUseCase {
	return &RouteUseCase{provider: provider}
}

func (uc *RouteUseCase) FindRoute(ctx context.Context, input types.RouteBuilderInput) (*types.RouteBuilderResult, error) {
	routes, err := uc.provider.FindRoute(types.RouteRequest{
		Origin:      input.Origin,
		Destination: input.Destination,
		Waypoints:   input.Waypoints,
		Mode:        input.Mode,
	})
	if err != nil {
		return nil, fmt.Errorf("find route: %w", err)
	}
	if len(routes) == 0 {
		return nil, fmt.Errorf("find route: no routes found")
	}

	route := routes[0]
	return &types.RouteBuilderResult{
		Route:      &route,
		Duration:   time.Duration(route.Duration.Seconds) * time.Second,
		DistanceKM: route.Distance.Meters / 1000,
		Polyline:   route.Polyline,
	}, nil
}

func (uc *RouteUseCase) GetDistanceMatrix(ctx context.Context, req types.DistanceMatrixRequest) (*types.DistanceMatrix, error) {
	if len(req.Origins) == 0 || len(req.Destinations) == 0 {
		return nil, fmt.Errorf("distance matrix: origins and destinations required")
	}
	return uc.provider.GetDistanceMatrix(req)
}
