package query

import "github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/valueobject"

type GetMatching struct {
	MatchingID valueobject.MatchingID
}

type GetCandidates struct {
	MatchingID valueobject.MatchingID
}

type GetAssignmentHistory struct {
	TripID valueobject.TripID
}

type GetRanking struct {
	MatchingID valueobject.MatchingID
}

type PreviewCandidates struct {
	PickupLat   float64
	PickupLng   float64
	RadiusMeters float64
	VehicleType string
}
