package security

import "time"

type AuditLogFilter struct {
	Action     *string
	EntityType *string
	Search     *string
	FromDate   *time.Time
	ToDate     *time.Time
}

type SessionFilter struct {
	Platform *string
	Browser  *string
	IsActive *bool
}

type LoginHistoryFilter struct {
	Status   *string
	FromDate *time.Time
	ToDate   *time.Time
}
