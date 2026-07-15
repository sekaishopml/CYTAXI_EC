package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/candidate"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/matching"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/ranking"
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/infrastructure/repository"
)

type MatchingService struct {
	matchRepo repository.MatchingRepository
	candidateRepo repository.CandidateRepository
	logger   *slog.Logger
}

func NewMatchingService(
	matchRepo repository.MatchingRepository,
	candidateRepo repository.CandidateRepository,
	logger *slog.Logger,
) *MatchingService {
	return &MatchingService{
		matchRepo:     matchRepo,
		candidateRepo: candidateRepo,
		logger:        logger,
	}
}

func (s *MatchingService) StartMatching(ctx context.Context, cmd command.StartMatching) (*matching.Matching, error) {
	session := matching.MatchingSession{
		PickupLat: cmd.PickupLat,
		PickupLng: cmd.PickupLng,
	}
	m := matching.NewMatching(cmd.TripID, session)
	m.Strategy = cmd.Strategy
	m.StartSearching()

	if err := s.matchRepo.Save(ctx, m); err != nil {
		return nil, fmt.Errorf("save matching: %w", err)
	}
	return m, nil
}

func (s *MatchingService) EvaluateCandidates(ctx context.Context, cmd command.EvaluateCandidates) (*candidate.CandidateSet, error) {
	m, err := s.matchRepo.FindByID(ctx, cmd.MatchingID)
	if err != nil {
		return nil, fmt.Errorf("find matching: %w", err)
	}
	m.StartEvaluating()
	s.matchRepo.Update(ctx, m)

	set := candidate.NewCandidateSet(cmd.MatchingID)
	s.logger.Info("candidates evaluated", "matching_id", cmd.MatchingID, "count", 0)
	return set, nil
}

func (s *MatchingService) RankCandidates(ctx context.Context, cmd command.RankCandidates) (*ranking.CandidateRanking, error) {
	cr := ranking.NewCandidateRanking(cmd.MatchingID, cmd.Strategy)
	s.logger.Info("candidates ranked", "matching_id", cmd.MatchingID, "strategy", cmd.Strategy)
	return cr, nil
}

func (s *MatchingService) SelectCandidate(ctx context.Context, cmd command.SelectCandidate) (*candidate.AssignmentResult, error) {
	result := &candidate.AssignmentResult{
		MatchingID: string(cmd.MatchingID),
		Strategy:   "balanced",
		Success:    false,
	}
	s.logger.Info("candidate selection", "matching_id", cmd.MatchingID)
	return result, nil
}

func (s *MatchingService) RetryMatching(ctx context.Context, cmd command.RetryMatching) (*matching.Matching, error) {
	m, err := s.matchRepo.FindByID(ctx, cmd.MatchingID)
	if err != nil {
		return nil, fmt.Errorf("find matching: %w", err)
	}
	if err := m.Retry(); err != nil {
		return nil, err
	}
	if err := s.matchRepo.Update(ctx, m); err != nil {
		return nil, fmt.Errorf("update matching: %w", err)
	}
	return m, nil
}

func (s *MatchingService) CancelMatching(ctx context.Context, cmd command.CancelMatching) error {
	m, err := s.matchRepo.FindByID(ctx, cmd.MatchingID)
	if err != nil {
		return fmt.Errorf("find matching: %w", err)
	}
	m.Cancel()
	return s.matchRepo.Update(ctx, m)
}

func (s *MatchingService) GetMatching(ctx context.Context, q query.GetMatching) (*matching.Matching, error) {
	return s.matchRepo.FindByID(ctx, q.MatchingID)
}

func (s *MatchingService) GetCandidates(ctx context.Context, q query.GetCandidates) (*candidate.CandidateSet, error) {
	return s.candidateRepo.FindByMatchingID(ctx, q.MatchingID)
}

func (s *MatchingService) GetRanking(ctx context.Context, q query.GetRanking) (*ranking.CandidateRanking, error) {
	return s.candidateRepo.FindRanking(ctx, q.MatchingID)
}
