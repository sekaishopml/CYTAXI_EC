package passenger

import "github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/domain/valueobject"

type Passenger struct {
	ID    valueobject.CustomerID `json:"id"`
	Phone string                 `json:"phone"`
	Name  string                 `json:"name"`
}
