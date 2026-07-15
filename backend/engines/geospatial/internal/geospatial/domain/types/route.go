package types

import "time"

type Route struct {
	Summary      string          `json:"summary"`
	Distance     Distance        `json:"distance"`
	Duration     Duration        `json:"duration"`
	Polyline     string          `json:"polyline"`
	Waypoints    []Coordinates   `json:"waypoints"`
	Steps        []RouteStep     `json:"steps"`
}

type RouteStep struct {
	Distance     Distance      `json:"distance"`
	Duration     Duration      `json:"duration"`
	Instruction  string        `json:"instruction"`
	Start        Coordinates   `json:"start"`
	End          Coordinates   `json:"end"`
	Polyline     string        `json:"polyline"`
}

type Distance struct {
	Meters    float64 `json:"meters"`
	Text      string  `json:"text"`
}

type Duration struct {
	Seconds   float64       `json:"seconds"`
	Text      string        `json:"text"`
}

type RouteRequest struct {
	Origin        Coordinates   `json:"origin"`
	Destination   Coordinates   `json:"destination"`
	Waypoints     []Coordinates `json:"waypoints,omitempty"`
	Mode          TravelMode    `json:"mode"`
	Alternatives  bool          `json:"alternatives"`
}

type TravelMode string

const (
	TravelModeDriving   TravelMode = "driving"
	TravelModeWalking   TravelMode = "walking"
	TravelModeBicycling TravelMode = "bicycling"
	TravelModeTransit   TravelMode = "transit"
)

type DistanceMatrix struct {
	OriginAddresses      []string           `json:"origin_addresses"`
	DestinationAddresses []string           `json:"destination_addresses"`
	Rows                 []DistanceMatrixRow `json:"rows"`
}

type DistanceMatrixRow struct {
	Elements []DistanceMatrixElement `json:"elements"`
}

type DistanceMatrixElement struct {
	Status   string   `json:"status"`
	Distance Distance `json:"distance"`
	Duration Duration `json:"duration"`
}

type DistanceMatrixRequest struct {
	Origins      []Coordinates `json:"origins"`
	Destinations []Coordinates `json:"destinations"`
	Mode         TravelMode    `json:"mode"`
}

type RouteBuilderInput struct {
	Origin      Coordinates
	Destination Coordinates
	Waypoints   []Coordinates
	Mode        TravelMode
}

type RouteBuilderResult struct {
	Route      *Route
	Duration   time.Duration
	DistanceKM float64
	Polyline   string
}
