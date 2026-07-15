package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/candidate"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/matching"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/ranking"
)

type MatchingService interface {
	StartMatching(ctx context.Context, cmd command.StartMatching) (*matching.Matching, error)
	EvaluateCandidates(ctx context.Context, cmd command.EvaluateCandidates) (*candidate.CandidateSet, error)
	RankCandidates(ctx context.Context, cmd command.RankCandidates) (*ranking.CandidateRanking, error)
	SelectCandidate(ctx context.Context, cmd command.SelectCandidate) (*candidate.AssignmentResult, error)
	RetryMatching(ctx context.Context, cmd command.RetryMatching) (*matching.Matching, error)
	CancelMatching(ctx context.Context, cmd command.CancelMatching) error
	GetMatching(ctx context.Context, q query.GetMatching) (*matching.Matching, error)
	GetCandidates(ctx context.Context, q query.GetCandidates) (*candidate.CandidateSet, error)
	GetRanking(ctx context.Context, q query.GetRanking) (*ranking.CandidateRanking, error)
}
