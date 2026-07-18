package real

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/service"
	"github.com/sekaishopml/cytaxi/backend/engines/geospatial/internal/geospatial/domain/types"
)

type NominatimResult struct {
	PlaceID     int      `json:"place_id"`
	Licence     string   `json:"licence"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	Address     struct {
		Road        string `json:"road"`
		City        string `json:"city"`
		State       string `json:"state"`
		Country     string `json:"country"`
		Postcode    string `json:"postcode"`
	} `json:"address"`
}

type OSRMRouteResponse struct {
	Code    string `json:"code"`
	Routes  []OSRMRoute `json:"routes"`
}

type OSRMRoute struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
	Geometry string  `json:"geometry"`
	Legs     []OSRLeg `json:"legs"`
}

type OSRLeg struct {
	Distance float64        `json:"distance"`
	Duration float64        `json:"duration"`
	Steps    []interface{}  `json:"steps"`
}

type Provider struct {
	client    *http.Client
	userAgent string
}

func NewProvider(client *http.Client) service.GeospatialProvider {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &Provider{
		client:    client,
		userAgent: "CYTAXI/1.0 (beta; platform@cytaxi.app)",
	}
}

func (p *Provider) Name() string { return "openstreetmap" }

func (p *Provider) SearchPlaces(req types.PlaceSearchRequest) (*types.PlaceSearchResult, error) {
	u := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=%d&addressdetails=1",
		url.QueryEscape(req.Query), max(req.MaxResult, 5))
	if !req.Location.IsZero() {
		u += fmt.Sprintf("&viewbox=%f,%f,%f,%f&bounded=1",
			req.Location.Lng-0.1, req.Location.Lat-0.1,
			req.Location.Lng+0.1, req.Location.Lat+0.1)
	}

	var results []NominatimResult
	if err := p.doRequest(context.Background(), u, &results); err != nil {
		return nil, fmt.Errorf("nominatim search: %w", err)
	}

	places := make([]types.Place, 0, len(results))
	for _, r := range results {
		lat, _ := strconv.ParseFloat(r.Lat, 64)
		lng, _ := strconv.ParseFloat(r.Lon, 64)
		places = append(places, types.Place{
			ID:          fmt.Sprintf("osm_%d", r.PlaceID),
			Name:        strings.Split(r.DisplayName, ",")[0],
			Address:     r.DisplayName,
			Coordinates: types.Coordinates{Lat: lat, Lng: lng},
		})
	}

	return &types.PlaceSearchResult{
		Places:     places,
		TotalCount: len(places),
	}, nil
}

func (p *Provider) Autocomplete(req types.AutocompleteRequest) ([]types.AutocompletePrediction, error) {
	result, err := p.SearchPlaces(types.PlaceSearchRequest{
		Query:     req.Input,
		Location:  req.Location,
		MaxResult: 5,
	})
	if err != nil {
		return nil, err
	}

	predictions := make([]types.AutocompletePrediction, len(result.Places))
	for i, place := range result.Places {
		predictions[i] = types.AutocompletePrediction{
			ID:          place.ID,
			Description: place.Address,
		}
	}
	return predictions, nil
}

func (p *Provider) PlaceDetails(placeID string) (*types.Place, error) {
	u := fmt.Sprintf("https://nominatim.openstreetmap.org/details?format=json&place_id=%s", strings.TrimPrefix(placeID, "osm_"))
	var result NominatimResult
	if err := p.doRequest(context.Background(), u, &result); err != nil {
		return nil, fmt.Errorf("nominatim details: %w", err)
	}
	lat, _ := strconv.ParseFloat(result.Lat, 64)
	lng, _ := strconv.ParseFloat(result.Lon, 64)
	return &types.Place{
		ID:          placeID,
		Name:        strings.Split(result.DisplayName, ",")[0],
		Coordinates: types.Coordinates{Lat: lat, Lng: lng},
	}, nil
}

func (p *Provider) Geocode(address string) ([]types.Coordinates, error) {
	result, err := p.SearchPlaces(types.PlaceSearchRequest{Query: address, MaxResult: 1})
	if err != nil {
		return nil, err
	}
	if len(result.Places) == 0 {
		return nil, fmt.Errorf("no results for address: %s", address)
	}
	return []types.Coordinates{result.Places[0].Coordinates}, nil
}

func (p *Provider) ReverseGeocode(coords types.Coordinates) (*types.Address, error) {
	u := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f&addressdetails=1",
		coords.Lat, coords.Lng)

	var result NominatimResult
	if err := p.doRequest(context.Background(), u, &result); err != nil {
		return nil, fmt.Errorf("nominatim reverse: %w", err)
	}

	lat, _ := strconv.ParseFloat(result.Lat, 64)
	lng, _ := strconv.ParseFloat(result.Lon, 64)

	return &types.Address{
		FormattedAddress: result.DisplayName,
		Street:           result.Address.Road,
		City:             result.Address.City,
		State:            result.Address.State,
		Country:          result.Address.Country,
		PostalCode:       result.Address.Postcode,
		Coordinates:      types.Coordinates{Lat: lat, Lng: lng},
	}, nil
}

func (p *Provider) FindRoute(req types.RouteRequest) ([]types.Route, error) {
	coords := fmt.Sprintf("%f,%f;%f,%f",
		req.Origin.Lng, req.Origin.Lat,
		req.Destination.Lng, req.Destination.Lat)

	u := fmt.Sprintf("https://router.project-osrm.org/route/v1/driving/%s?overview=full&steps=true&geometries=polyline",
		coords)

	var resp OSRMRouteResponse
	if err := p.doRequest(context.Background(), u, &resp); err != nil {
		return nil, fmt.Errorf("osrm route: %w", err)
	}
	if resp.Code != "Ok" || len(resp.Routes) == 0 {
		return nil, fmt.Errorf("osrm: no route found")
	}

	r := resp.Routes[0]
	return []types.Route{{
		Distance: types.Distance{Meters: r.Distance, Text: fmt.Sprintf("%.1f km", r.Distance/1000)},
		Duration: types.Duration{Seconds: r.Duration, Text: fmt.Sprintf("%.0f min", r.Duration/60)},
		Polyline: r.Geometry,
		Waypoints: []types.Coordinates{req.Origin, req.Destination},
	}}, nil
}

func (p *Provider) GetDistanceMatrix(req types.DistanceMatrixRequest) (*types.DistanceMatrix, error) {
	matrix := &types.DistanceMatrix{}
	for _, origin := range req.Origins {
		row := types.DistanceMatrixRow{}
		for _, dest := range req.Destinations {
			dist := origin.DistanceTo(dest)
			eta := int(dist / 8.33)
			row.Elements = append(row.Elements, types.DistanceMatrixElement{
				Status: "OK",
				Distance: types.Distance{Meters: dist, Text: fmt.Sprintf("%.1f km", dist/1000)},
				Duration: types.Duration{Seconds: float64(eta), Text: fmt.Sprintf("%d min", eta/60)},
			})
		}
		matrix.Rows = append(matrix.Rows, row)
	}
	return matrix, nil
}

func (p *Provider) doRequest(ctx context.Context, urlStr string, result any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", p.userAgent)

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}
