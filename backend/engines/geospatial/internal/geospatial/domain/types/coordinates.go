package types

import "math"

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Bounds struct {
	SouthWest Coordinates `json:"south_west"`
	NorthEast Coordinates `json:"north_east"`
}

func (c Coordinates) DistanceTo(other Coordinates) float64 {
	const R = 6371000
	dLat := (other.Lat - c.Lat) * math.Pi / 180
	dLng := (other.Lng - c.Lng) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(c.Lat*math.Pi/180)*math.Cos(other.Lat*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func (c Coordinates) IsZero() bool {
	return c.Lat == 0 && c.Lng == 0
}
