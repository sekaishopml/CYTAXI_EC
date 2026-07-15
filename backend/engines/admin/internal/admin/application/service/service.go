package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/announcement"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/audit"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/configuration"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/featureflag"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/maintenance"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/permission"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/role"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/valueobject"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/infrastructure/repository"
)

type AdminService struct {
	roleRepo   repository.RoleRepository
	permRepo   repository.PermissionRepository
	configRepo repository.ConfigurationRepository
	flagRepo   repository.FeatureFlagRepository
	auditRepo  repository.AuditRepository
	maintRepo  repository.MaintenanceRepository
	annRepo    repository.AnnouncementRepository
	logger     *slog.Logger
}

func NewAdminService(
	roleRepo repository.RoleRepository,
	permRepo repository.PermissionRepository,
	configRepo repository.ConfigurationRepository,
	flagRepo repository.FeatureFlagRepository,
	auditRepo repository.AuditRepository,
	maintRepo repository.MaintenanceRepository,
	annRepo repository.AnnouncementRepository,
	logger *slog.Logger,
) *AdminService {
	return &AdminService{
		roleRepo: roleRepo, permRepo: permRepo, configRepo: configRepo,
		flagRepo: flagRepo, auditRepo: auditRepo, maintRepo: maintRepo,
		annRepo: annRepo, logger: logger,
	}
}

func (s *AdminService) CreateRole(ctx context.Context, cmd command.CreateRole) (*role.Role, error) {
	r := role.NewRole(cmd.Name, cmd.Description)
	if err := s.roleRepo.Save(ctx, r); err != nil {
		return nil, fmt.Errorf("save role: %w", err)
	}
	return r, nil
}

func (s *AdminService) UpdateRole(ctx context.Context, cmd command.UpdateRole) (*role.Role, error) {
	r, err := s.roleRepo.FindByID(ctx, cmd.RoleID)
	if err != nil {
		return nil, fmt.Errorf("find role: %w", err)
	}
	r.Name = cmd.Name
	r.Description = cmd.Description
	r.UpdatedAt = time.Now()
	s.roleRepo.Update(ctx, r)
	return r, nil
}

func (s *AdminService) DeleteRole(ctx context.Context, cmd command.DeleteRole) error {
	return s.roleRepo.Delete(ctx, cmd.RoleID)
}

func (s *AdminService) AssignPermission(ctx context.Context, cmd command.AssignPermission) (*role.Role, error) {
	r, err := s.roleRepo.FindByID(ctx, cmd.RoleID)
	if err != nil {
		return nil, fmt.Errorf("find role: %w", err)
	}
	r.AddPermission(cmd.PermissionID)
	s.roleRepo.Update(ctx, r)
	return r, nil
}

func (s *AdminService) UpdateConfiguration(ctx context.Context, cmd command.UpdateConfiguration) (*configuration.Configuration, error) {
	cfg, err := s.configRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		cfg = configuration.NewConfiguration(cmd.Key, cmd.Value, configuration.CfgString)
		s.configRepo.Save(ctx, cfg)
		return cfg, nil
	}
	cfg.Update(cmd.Value)
	s.configRepo.Save(ctx, cfg)
	return cfg, nil
}

func (s *AdminService) EnableFeature(ctx context.Context, cmd command.EnableFeature) error {
	f, err := s.flagRepo.FindByID(ctx, cmd.FeatureID)
	if err != nil {
		return fmt.Errorf("find feature: %w", err)
	}
	f.Enable()
	return s.flagRepo.Update(ctx, f)
}

func (s *AdminService) DisableFeature(ctx context.Context, cmd command.DisableFeature) error {
	f, err := s.flagRepo.FindByID(ctx, cmd.FeatureID)
	if err != nil {
		return fmt.Errorf("find feature: %w", err)
	}
	f.Disable()
	return s.flagRepo.Update(ctx, f)
}

func (s *AdminService) CreateAnnouncement(ctx context.Context, cmd command.CreateAnnouncement) (*announcement.Announcement, error) {
	a := announcement.NewAnnouncement(cmd.Title, cmd.Body, cmd.Priority, time.Unix(cmd.PublishAt, 0))
	s.annRepo.Save(ctx, a)
	return a, nil
}

func (s *AdminService) ScheduleMaintenance(ctx context.Context, cmd command.ScheduleMaintenance) (*maintenance.MaintenanceWindow, error) {
	m := maintenance.NewMaintenanceWindow(cmd.Title, time.Unix(cmd.StartAt, 0), time.Unix(cmd.EndAt, 0))
	m.Description = cmd.Description
	m.Affects = cmd.Affects
	s.maintRepo.Save(ctx, m)
	return m, nil
}

func (s *AdminService) RegisterAudit(ctx context.Context, cmd command.RegisterAudit) (*audit.AuditLog, error) {
	a := audit.NewAuditLog(cmd.OperatorID, cmd.Action, cmd.Resource, cmd.ResourceID)
	a.Changes = cmd.Changes
	s.auditRepo.Save(ctx, a)
	return a, nil
}

func (s *AdminService) GetConfiguration(ctx context.Context, q query.GetConfiguration) (*configuration.Configuration, error) {
	return s.configRepo.FindByKey(ctx, q.Key)
}

func (s *AdminService) GetRoles(ctx context.Context, q query.GetRoles) ([]role.Role, error) {
	return s.roleRepo.FindAll(ctx)
}

func (s *AdminService) GetPermissions(ctx context.Context, q query.GetPermissions) ([]permission.Permission, error) {
	if q.RoleID != "" {
		r, _ := s.roleRepo.FindByID(ctx, q.RoleID)
		if r != nil {
			var perms []permission.Permission
			for _, pid := range r.Permissions {
				p, _ := s.permRepo.FindByID(ctx, pid)
				if p != nil {
					perms = append(perms, *p)
				}
			}
			return perms, nil
		}
	}
	return s.permRepo.FindAll(ctx)
}

func (s *AdminService) GetAuditLogs(ctx context.Context, q query.GetAuditLogs) ([]audit.AuditLog, error) {
	return s.auditRepo.FindAll(ctx)
}

func (s *AdminService) GetMaintenance(ctx context.Context, q query.GetMaintenance) ([]maintenance.MaintenanceWindow, error) {
	return s.maintRepo.FindAll(ctx)
}

func (s *AdminService) GetAnnouncements(ctx context.Context, q query.GetAnnouncements) ([]announcement.Announcement, error) {
	return s.annRepo.FindAll(ctx)
}

func (s *AdminService) GetFeatureFlags(ctx context.Context, q query.GetFeatureFlags) ([]featureflag.FeatureFlag, error) {
	return s.flagRepo.FindAll(ctx)
}
