package destination

import "github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/valueobject"

type Destination struct {
	Address  string                  `json:"address"`
	Location valueobject.Coordinates `json:"location"`
	PlaceID  string                  `json:"place_id,omitempty"`
}
