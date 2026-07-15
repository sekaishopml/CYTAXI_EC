package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/blacklist"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/document"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/fraud"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/identity"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/trustscore"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/verification"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/whitelist"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/infrastructure/repository"
)

type TrustService struct {
	identityRepo repository.IdentityRepository
	verifRepo    repository.VerificationRepository
	docRepo      repository.DocumentRepository
	fraudRepo    repository.FraudRepository
	trustRepo    repository.TrustRepository
	blacklistRepo repository.BlacklistRepository
	logger       *slog.Logger
}

func NewTrustService(
	identityRepo repository.IdentityRepository,
	verifRepo repository.VerificationRepository,
	docRepo repository.DocumentRepository,
	fraudRepo repository.FraudRepository,
	trustRepo repository.TrustRepository,
	blacklistRepo repository.BlacklistRepository,
	logger *slog.Logger,
) *TrustService {
	return &TrustService{
		identityRepo: identityRepo, verifRepo: verifRepo, docRepo: docRepo,
		fraudRepo: fraudRepo, trustRepo: trustRepo, blacklistRepo: blacklistRepo,
		logger: logger,
	}
}

func (s *TrustService) CreateIdentity(ctx context.Context, cmd command.CreateIdentity) (*identity.Identity, error) {
	id := identity.NewIdentity(cmd.OwnerID, cmd.Type, cmd.Phone)
	id.Email = cmd.Email
	if err := s.identityRepo.Save(ctx, id); err != nil {
		return nil, fmt.Errorf("save identity: %w", err)
	}
	return id, nil
}

func (s *TrustService) VerifyIdentity(ctx context.Context, cmd command.VerifyIdentity) (*verification.Verification, error) {
	v := verification.NewVerification(cmd.IdentityID, verification.VerificationType(cmd.Type))
	v.StartReview()
	if err := s.verifRepo.Save(ctx, v); err != nil {
		return nil, fmt.Errorf("save verification: %w", err)
	}
	return v, nil
}

func (s *TrustService) UploadDocument(ctx context.Context, cmd command.UploadDocument) (*document.Document, error) {
	d := document.NewDocument(cmd.IdentityID, cmd.DocType, cmd.URL)
	if err := s.docRepo.Save(ctx, d); err != nil {
		return nil, fmt.Errorf("save document: %w", err)
	}
	return d, nil
}

func (s *TrustService) ApproveVerification(ctx context.Context, cmd command.ApproveVerification) error {
	v, err := s.verifRepo.FindByID(ctx, cmd.VerificationID)
	if err != nil {
		return fmt.Errorf("find verification: %w", err)
	}
	v.Approve(verification.VerificationResult{Passed: true, Score: cmd.Score})
	s.verifRepo.Update(ctx, v)

	id, _ := s.identityRepo.FindByID(ctx, v.IdentityID)
	if id != nil {
		id.Verify()
		s.identityRepo.Update(ctx, id)
	}
	return nil
}

func (s *TrustService) RejectVerification(ctx context.Context, cmd command.RejectVerification) error {
	v, err := s.verifRepo.FindByID(ctx, cmd.VerificationID)
	if err != nil {
		return fmt.Errorf("find verification: %w", err)
	}
	v.Reject(cmd.Reason)
	return s.verifRepo.Update(ctx, v)
}

func (s *TrustService) CalculateTrustScore(ctx context.Context, cmd command.CalculateTrustScore) (*trustscore.TrustProfile, error) {
	tp := trustscore.NewTrustProfile(cmd.IdentityID)
	tp.Components.VerificationScore = cmd.VerifyScore
	tp.Components.ActivityScore = cmd.Activity
	tp.Components.CommunityScore = cmd.Community
	tp.Components.ComplianceScore = cmd.Compliance
	tp.Calculate()
	if err := s.trustRepo.Save(ctx, tp); err != nil {
		return nil, fmt.Errorf("save trust: %w", err)
	}
	return tp, nil
}

func (s *TrustService) RunFraudCheck(ctx context.Context, cmd command.RunFraudCheck) (*fraud.FraudAssessment, error) {
	f := fraud.NewFraudAssessment(cmd.IdentityID)
	if err := s.fraudRepo.Save(ctx, f); err != nil {
		return nil, fmt.Errorf("save fraud check: %w", err)
	}
	return f, nil
}

func (s *TrustService) Blacklist(ctx context.Context, cmd command.BlacklistIdentity) (*blacklist.BlacklistEntry, error) {
	e := blacklist.NewBlacklistEntry(cmd.IdentityID, cmd.Reason, cmd.Severity)
	if err := s.blacklistRepo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("save blacklist: %w", err)
	}
	return e, nil
}

func (s *TrustService) Whitelist(ctx context.Context, cmd command.WhitelistIdentity) (*whitelist.WhitelistEntry, error) {
	e := whitelist.NewWhitelistEntry(cmd.IdentityID, cmd.Reason)
	return e, nil
}

func (s *TrustService) GetIdentity(ctx context.Context, q query.GetIdentity) (*identity.Identity, error) {
	return s.identityRepo.FindByID(ctx, q.IdentityID)
}

func (s *TrustService) GetTrustScore(ctx context.Context, q query.GetTrustScore) (*trustscore.TrustProfile, error) {
	return s.trustRepo.FindByID(ctx, q.IdentityID)
}

func (s *TrustService) GetDocuments(ctx context.Context, q query.GetDocuments) ([]document.Document, error) {
	return s.docRepo.FindByIdentityID(ctx, q.IdentityID)
}

func (s *TrustService) GetFraudHistory(ctx context.Context, q query.GetFraudHistory) ([]fraud.FraudAssessment, error) {
	return s.fraudRepo.FindByIdentityID(ctx, q.IdentityID)
}
