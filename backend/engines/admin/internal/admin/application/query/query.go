package query

import "github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/valueobject"

type GetConfiguration struct {
	Key valueobject.ConfigurationKey
}

type GetRoles struct{}

type GetPermissions struct {
	RoleID valueobject.RoleID
}

type GetAuditLogs struct {
	OperatorID valueobject.OperatorID
	Resource   string
	Limit      int
}

type GetMaintenance struct {
	Active bool
}

type GetAnnouncements struct {
	Active bool
}

type GetFeatureFlags struct {
	Scope string
}
