package repositories

import (
	"github.com/go-api/internal/domain/audit"
	"github.com/go-api/internal/domain/authorization/entities"
	"github.com/go-api/internal/domain/loginhistory"
	"github.com/go-api/internal/domain/session"
	domainToken "github.com/go-api/internal/domain/token"
	"github.com/go-api/internal/domain/user"
	db "github.com/go-api/internal/infrastructure/postgres/generated"
	"github.com/go-api/internal/infrastructure/postgres/types"
)

func toUserDomain(u db.User) *user.User {

	return &user.User{
		ID:         u.ID,
		Email:      u.Email,
		Password:   u.Password,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Phone:      types.FromNullablePGText(u.Phone),
		ImageURL:   types.FromNullablePGText(u.ImageUrl),
		IsVerified: u.IsVerified,
		CreatedAt:  types.FromPGTimestamp(u.CreatedAt),
		UpdatedAt:  types.FromPGTimestamp(u.UpdatedAt),
	}
}

func toTokenDomain(t db.Token) *domainToken.Token {
	return &domainToken.Token{
		ID:        t.ID,
		UserID:    t.UserID,
		Type:      domainToken.Type(t.Type),
		Token:     t.Token,
		ExpiresAt: types.FromPGTimestamp(t.ExpiresAt),
		UsedAt:    types.FromPGTimestampPtr(t.UsedAt),
		CreatedAt: types.FromPGTimestamp(t.CreatedAt),
		UpdatedAt: types.FromPGTimestamp(t.UpdatedAt),
	}
}

func toSessionDomain(s db.Session) *session.Session {
	return &session.Session{
		ID:            s.ID,
		UserID:        s.UserID,
		RefreshToken:  s.RefreshToken,
		UserAgent:     types.FromPGText(s.UserAgent),
		IPAddress:     types.FromPGInet(s.IpAddress),
		DeviceID:      types.FromPGText(s.DeviceID),
		Platform:      types.FromPGText(s.Platform),
		Browser:       types.FromPGText(s.Browser),
		LastSeenAt:    types.FromPGTimestamp(s.LastSeenAt),
		LastUsedAt:    types.FromPGTimestampPtr(s.LastUsedAt),
		ExpiresAt:     types.FromPGTimestamp(s.ExpiresAt),
		RevokedAt:     types.FromPGTimestampPtr(s.RevokedAt),
		RevokedReason: types.FromNullablePGText(s.RevokedReason),
		CreatedAt:     types.FromPGTimestamp(s.CreatedAt),
		UpdatedAt:     types.FromPGTimestamp(s.UpdatedAt),
	}
}

func toAuditLogDomain(log db.AuditLog) *audit.Log {
	return &audit.Log{
		ID:         log.ID,
		UserID:     types.FromNullableUUID(log.UserID),
		Action:     log.Action,
		EntityType: types.FromNullablePGText(log.EntityType),
		EntityID:   types.FromNullableUUID(log.EntityID),
		IPAddress:  types.FromPGInet(log.IpAddress),
		UserAgent:  types.FromPGText(log.UserAgent),
		Metadata:   types.FromJSONB(log.Metadata),
		CreatedAt:  types.FromPGTimestamp(log.CreatedAt),
		UpdatedAt:  types.FromPGTimestamp(log.UpdatedAt),
	}
}

func toLoginHistoryDomain(h db.LoginHistory) *loginhistory.LoginHistory {
	return &loginhistory.LoginHistory{
		ID:            h.ID,
		UserID:        types.FromNullableUUID(h.UserID),
		SessionID:     types.FromNullableUUID(h.SessionID),
		Status:        loginhistory.Status(h.Status),
		IPAddress:     types.FromPGInet(h.IpAddress),
		UserAgent:     types.FromPGText(h.UserAgent),
		DeviceName:    types.FromPGText(h.DeviceName),
		FailureReason: types.FromNullablePGText(h.FailureReason),
		LoginAt:       types.FromPGTimestamp(h.LoginAt),
		LogoutAt:      types.FromPGTimestampPtr(h.LogoutAt),
		CreatedAt:     types.FromPGTimestamp(h.CreatedAt),
		UpdatedAt:     types.FromPGTimestamp(h.UpdatedAt),
	}
}

func toRoleDomain(role db.Role) *entities.Role {
	return &entities.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: types.FromNullablePGText(role.Description),
		IsSystem:    role.IsSystem,
		CreatedAt:   types.FromPGTimestamp(role.CreatedAt),
		UpdatedAt:   types.FromPGTimestamp(role.UpdatedAt),
	}
}

func toPermissionDomain(permission db.Permission) *entities.Permission {
	return &entities.Permission{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: types.FromNullablePGText(permission.Description),
		Resource:    permission.Resource,
		Action:      permission.Action,
		CreatedAt:   types.FromPGTimestamp(permission.CreatedAt),
		UpdatedAt:   types.FromPGTimestamp(permission.UpdatedAt),
	}
}

func toRoleDomains(records []db.Role) []*entities.Role {
	roles := make([]*entities.Role, len(records))
	for i, record := range records {
		roles[i] = toRoleDomain(record)
	}
	return roles
}

func toPermissionDomains(records []db.Permission) []*entities.Permission {
	permissions := make([]*entities.Permission, len(records))
	for i, record := range records {
		permissions[i] = toPermissionDomain(record)
	}
	return permissions
}
