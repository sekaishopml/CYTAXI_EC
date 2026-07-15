package repository

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/candidate"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/matching"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/ranking"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/valueobject"
)

type MatchingRepository interface {
	FindByID(ctx context.Context, id valueobject.MatchingID) (*matching.Matching, error)
	FindByTripID(ctx context.Context, tripID valueobject.TripID) ([]matching.Matching, error)
	Save(ctx context.Context, m *matching.Matching) error
	Update(ctx context.Context, m *matching.Matching) error
}

type CandidateRepository interface {
	FindByMatchingID(ctx context.Context, id valueobject.MatchingID) (*candidate.CandidateSet, error)
	FindRanking(ctx context.Context, id valueobject.MatchingID) (*ranking.CandidateRanking, error)
	Save(ctx context.Context, cs *candidate.CandidateSet) error
	SaveAttempt(ctx context.Context, a *candidate.AssignmentAttempt) error
}
