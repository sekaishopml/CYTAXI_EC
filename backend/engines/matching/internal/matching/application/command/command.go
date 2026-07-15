package command

import "github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/valueobject"

type StartMatching struct {
	TripID    valueobject.TripID
	PickupLat float64
	PickupLng float64
	Strategy  string
}

type EvaluateCandidates struct {
	MatchingID valueobject.MatchingID
}

type RankCandidates struct {
	MatchingID valueobject.MatchingID
	Strategy   string
}

type SelectCandidate struct {
	MatchingID valueobject.MatchingID
}

type RetryMatching struct {
	MatchingID valueobject.MatchingID
}

type CancelMatching struct {
	MatchingID valueobject.MatchingID
	Reason     string
}
