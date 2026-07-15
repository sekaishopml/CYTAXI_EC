package types

import "time"

type ZoneID string

type Zone struct {
	ID        ZoneID       `json:"id"`
	Name      string       `json:"name"`
	Area      GeoFence     `json:"area"`
	Active    bool         `json:"active"`
	CreatedAt time.Time    `json:"created_at"`
}

type GeoFence struct {
	Type        string        `json:"type"` // polygon, circle, rectangle
	Coordinates []Coordinates `json:"coordinates"`
	Center      Coordinates   `json:"center,omitempty"`
	Radius      float64       `json:"radius,omitempty"` // meters (for circle)
}

type ZoneService interface {
	CreateZone(ctx Zone) error
	GetZone(id ZoneID) (*Zone, error)
	ListZones() ([]Zone, error)
	UpdateZone(zone Zone) error
	DeleteZone(id ZoneID) error
	FindZonesByPoint(point Coordinates) ([]Zone, error)
}

func PointInGeoFence(point Coordinates, fence GeoFence) bool {
	switch fence.Type {
	case "circle":
		return point.DistanceTo(fence.Center) <= fence.Radius
	case "polygon":
		return pointInPolygon(point, fence.Coordinates)
	default:
		return false
	}
}

func pointInPolygon(point Coordinates, polygon []Coordinates) bool {
	inside := false
	j := len(polygon) - 1
	for i := 0; i < len(polygon); i++ {
		if (polygon[i].Lat > point.Lat) != (polygon[j].Lat > point.Lat) &&
			point.Lng < (polygon[j].Lng-polygon[i].Lng)*(point.Lat-polygon[i].Lat)/(polygon[j].Lat-polygon[i].Lat)+polygon[i].Lng {
			inside = !inside
		}
		j = i
	}
	return inside
}
