package security

import (
	"context"
	"net"
	"time"

	"github.com/google/uuid"

	"github.com/go-api/internal/domain/audit"
	"github.com/go-api/internal/domain/loginhistory"
	"github.com/go-api/internal/domain/session"
	"github.com/go-api/internal/domain/user"
)

type Service struct {
	users        user.Repository
	sessions     session.Repository
	auditLogs    audit.Repository
	loginHistory loginhistory.Repository
}

func NewService(
	users user.Repository,
	sessions session.Repository,
	auditLogs audit.Repository,
	loginHistory loginhistory.Repository,
) *Service {
	return &Service{
		users:        users,
		sessions:     sessions,
		auditLogs:    auditLogs,
		loginHistory: loginHistory,
	}
}

func ptr[T any](v T) *T {
	return &v
}

func (s *Service) RecordSuccessfulLogin(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, ipAddress string, userAgent string, deviceName string) error {
	now := time.Now()

	history := &loginhistory.LoginHistory{
		ID:         uuid.New(),
		UserID:     &userID,
		SessionID:  &sessionID,
		Status:     loginhistory.StatusSuccess,
		UserAgent:  userAgent,
		DeviceName: deviceName,
		IPAddress:  net.ParseIP(ipAddress),
		LoginAt:    now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.loginHistory.Create(ctx, history); err != nil {
		return err
	}

	auditLog := &audit.Log{
		ID:         uuid.New(),
		UserID:     &userID,
		Action:     "USER_LOGIN",
		EntityType: ptr("user"),
		EntityID:   &userID,
		IPAddress:  net.ParseIP(ipAddress),
		UserAgent:  userAgent,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.auditLogs.Create(ctx, auditLog); err != nil {
		return err
	}

	if err := s.users.ResetFailedLoginAttempts(ctx, userID); err != nil {
		return err
	}

	return s.users.UpdateLastLogin(ctx, userID)
}

func (s *Service) RecordFailedLogin(ctx context.Context, userID *uuid.UUID, ipAddress string, userAgent string, reason string) error {
	now := time.Now()

	history := &loginhistory.LoginHistory{
		ID:            uuid.New(),
		UserID:        userID,
		Status:        loginhistory.StatusFailed,
		UserAgent:     userAgent,
		FailureReason: &reason,
		LoginAt:       now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.loginHistory.Create(ctx, history); err != nil {
		return err
	}

	var auditUser *uuid.UUID

	if userID != nil {
		auditUser = userID
	}

	log := &audit.Log{
		ID:         uuid.New(),
		UserID:     auditUser,
		Action:     "LOGIN_FAILED",
		EntityType: ptr("user"),
		EntityID:   auditUser,
		IPAddress:  net.ParseIP(ipAddress),
		UserAgent:  userAgent,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.auditLogs.Create(ctx, log); err != nil {
		return err
	}

	return nil
}

func (s *Service) RecordLogout(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	now := time.Now()
	if err := s.loginHistory.MarkLoginHistoryLogout(ctx, sessionID); err != nil {
		return err
	}

	log := &audit.Log{
		ID:        uuid.New(),
		UserID:    &userID,
		Action:    "USER_LOGOUT",
		CreatedAt: now,
	}

	if err := s.auditLogs.Create(ctx, log); err != nil {
		return err
	}

	return nil
}

func (s *Service) IncrementFailedLoginAttempt(ctx context.Context, userID uuid.UUID) error {

	return s.users.IncrementFailedLoginAttempts(ctx, userID)
}

func (s *Service) LockAccount(ctx context.Context, userID uuid.UUID) error {
	unlockDuration := time.Now().Add(15 * time.Minute)

	return s.users.LockUserAccount(ctx, userID, unlockDuration)
}
