package stop

import (
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/valueobject"
)

type Stop struct {
	ID        valueobject.StopID       `json:"id"`
	Address   string                   `json:"address"`
	Location  valueobject.Coordinates  `json:"location"`
	Notes     string                   `json:"notes,omitempty"`
	ReachedAt *time.Time               `json:"reached_at,omitempty"`
}

func NewStop(address string, loc valueobject.Coordinates) Stop {
	return Stop{
		ID:       valueobject.StopID(fmt.Sprintf("stop_%d", time.Now().UnixNano())),
		Address:  address,
		Location: loc,
	}
}
