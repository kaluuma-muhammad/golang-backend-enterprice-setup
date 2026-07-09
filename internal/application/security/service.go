package security

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/go-api/internal/application/common"
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
	tx           common.TransactionManager
}

func NewService(
	users user.Repository,
	sessions session.Repository,
	auditLogs audit.Repository,
	loginHistory loginhistory.Repository,
	tx common.TransactionManager,
) *Service {
	return &Service{
		users:        users,
		sessions:     sessions,
		auditLogs:    auditLogs,
		loginHistory: loginHistory,
		tx:           tx,
	}
}

func (s *Service) RecordSuccessfulLogin(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, ipAddress string, userAgent string, deviceName string) error {
	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
		now := time.Now()

		history := &loginhistory.LoginHistory{
			ID:         uuid.New(),
			UserID:     &userID,
			SessionID:  &sessionID,
			Status:     loginhistory.StatusSuccess,
			UserAgent:  userAgent,
			DeviceName: deviceName,
			LoginAt:    now,
			CreatedAt:  now,
		}

		if err := s.loginHistory.Create(txCtx, history); err != nil {
			return err
		}

		auditLog := &audit.Log{
			ID:        uuid.New(),
			UserID:    &userID,
			Action:    "USER_LOGIN",
			UserAgent: userAgent,
			CreatedAt: now,
		}

		if err := s.auditLogs.Create(txCtx, auditLog); err != nil {
			return err
		}

		if err := s.users.ResetFailedLoginAttempts(txCtx, userID); err != nil {
			return err
		}

		return s.users.UpdateLastLogin(txCtx, userID)
	})

	return err
}

func (s *Service) RecordFailedLogin(ctx context.Context, userID *uuid.UUID, ipAddress string, userAgent string, reason string) error {
	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
		now := time.Now()

		history := &loginhistory.LoginHistory{
			ID:            uuid.New(),
			UserID:        userID,
			Status:        loginhistory.StatusFailed,
			UserAgent:     userAgent,
			FailureReason: &reason,
			LoginAt:       now,
			CreatedAt:     now,
		}

		if err := s.loginHistory.Create(ctx, history); err != nil {
			return err
		}

		var auditUser *uuid.UUID

		if userID != nil {
			auditUser = userID
		}

		log := &audit.Log{
			ID:        uuid.New(),
			UserID:    auditUser,
			Action:    "LOGIN_FAILED",
			UserAgent: userAgent,
			CreatedAt: now,
		}

		if err := s.auditLogs.Create(txCtx, log); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *Service) RecordLogout(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	err := s.tx.Execute(ctx, func(txCtx context.Context) error {
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

		if err := s.auditLogs.Create(txCtx, log); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *Service) IncrementFailedLoginAttempt(ctx context.Context, userID uuid.UUID) error {

	return s.users.IncrementFailedLoginAttempts(ctx, userID)
}

func (s *Service) LockAccount(ctx context.Context, userID uuid.UUID) error {
	unlockDuration := time.Now().Add(15 * time.Minute)

	return s.users.LockUserAccount(ctx, userID, unlockDuration)
}
