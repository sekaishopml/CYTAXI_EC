package definition

const (
	EventGeocodeRequested         = "geospatial.geocode_requested"
	EventGeocodeCompleted         = "geospatial.geocode_completed"
	EventReverseGeocodeRequested  = "geospatial.reverse_geocode_requested"
	EventReverseGeocodeCompleted  = "geospatial.reverse_geocode_completed"
	EventRouteRequested           = "geospatial.route_requested"
	EventRouteFound               = "geospatial.route_found"
	EventDistanceMatrixRequested  = "geospatial.distance_matrix_requested"
	EventDistanceMatrixCompleted  = "geospatial.distance_matrix_completed"
	EventPlaceSearchRequested     = "geospatial.place_search_requested"
	EventPlaceSearchCompleted     = "geospatial.place_search_completed"
	EventProviderError            = "geospatial.provider_error"
)

type GeocodeRequestedPayload struct {
	Address string `json:"address"`
}

type GeocodeCompletedPayload struct {
	Address     string  `json:"address"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Provider    string  `json:"provider"`
	LatencyMs   int64   `json:"latency_ms"`
}

type RouteRequestedPayload struct {
	OriginLat      float64 `json:"origin_lat"`
	OriginLng      float64 `json:"origin_lng"`
	DestinationLat float64 `json:"destination_lat"`
	DestinationLng float64 `json:"destination_lng"`
	Mode           string  `json:"mode"`
}

type RouteFoundPayload struct {
	DistanceKM  float64 `json:"distance_km"`
	DurationSec float64 `json:"duration_sec"`
	Provider    string  `json:"provider"`
}

type ProviderErrorPayload struct {
	Provider  string `json:"provider"`
	Operation string `json:"operation"`
	Error     string `json:"error"`
}
