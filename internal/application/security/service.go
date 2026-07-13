package security

import (
	"context"
	"net"
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

func (s *Service) RecordSuccessfulLogin(ctx context.Context, userID, sessionID uuid.UUID, ipAddress, userAgent, deviceName string) error {
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

	metadata := map[string]any{
		"session_id": sessionID,
		"device":     deviceName,
	}

	auditLog := &audit.Log{
		ID:         uuid.New(),
		UserID:     &userID,
		Action:     "USER_LOGIN",
		EntityType: ptr("user"),
		EntityID:   &userID,
		IPAddress:  net.ParseIP(ipAddress),
		UserAgent:  userAgent,
		Metadata:   metadata,
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

func (s *Service) RecordFailedLogin(ctx context.Context, userID *uuid.UUID, ipAddress, userAgent, reason string) error {
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

func (s *Service) RecordLogout(ctx context.Context, userID, sessionID uuid.UUID, ipAddress, userAgent, deviceName string) error {
	now := time.Now()
	if err := s.loginHistory.MarkLoginHistoryLogout(ctx, sessionID); err != nil {
		return err
	}

	metadata := map[string]any{
		"session_id": sessionID,
		"device":     deviceName,
	}

	log := &audit.Log{
		ID:         uuid.New(),
		UserID:     &userID,
		Action:     "USER_LOGOUT",
		EntityType: ptr("user"),
		EntityID:   &userID,
		IPAddress:  net.ParseIP(ipAddress),
		UserAgent:  userAgent,
		Metadata:   metadata,
		CreatedAt:  now,
		UpdatedAt:  now,
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

func (s *Service) GetAuditLogs(ctx context.Context, userID uuid.UUID, req common.PaginationRequest) ([]AuditLogResponse, common.PaginationResponse, error) {

	req = req.Normalize()

	logs, total, err := s.auditLogs.GetByUserPaginated(ctx, userID, req.PageSize, req.Offset())
	if err != nil {
		return nil, common.PaginationResponse{}, err
	}

	var response []AuditLogResponse
	for _, l := range logs {
		response = append(response, AuditLogResponse{
			ID:         l.ID,
			Action:     l.Action,
			EntityType: l.EntityType,
			CreatedAt:  l.CreatedAt,
		})
	}

	pagination := common.PaginationResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
	}

	return response, pagination, nil
}

func (s *Service) GetLoginHistory(ctx context.Context, userID uuid.UUID, req common.PaginationRequest) ([]LoginHistoryResponse, common.PaginationResponse, error) {

	req = req.Normalize()

	data, total, err := s.loginHistory.GetByUserPaginated(ctx, userID, req.PageSize, req.Offset())
	if err != nil {
		return nil, common.PaginationResponse{}, err
	}

	var response []LoginHistoryResponse
	for _, d := range data {
		response = append(response, LoginHistoryResponse{
			ID:         d.ID,
			Status:     string(d.Status),
			IPAddress:  d.IPAddress.String(),
			DeviceName: d.DeviceName,
			LoginAt:    d.LoginAt,
			LogoutAt:   d.LogoutAt,
		})
	}

	pagination := common.PaginationResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
	}

	return response, pagination, nil
}

func (s *Service) GetUserSessions(ctx context.Context, userID uuid.UUID, req common.PaginationRequest) ([]SessionResponse, common.PaginationResponse, error) {

	req = req.Normalize()

	data, total, err := s.sessions.GetSessionsByUserPaginated(ctx, userID, req.PageSize, req.Offset())
	if err != nil {
		return nil, common.PaginationResponse{}, err
	}

	var response []SessionResponse
	for _, d := range data {
		response = append(response, SessionResponse{
			ID:         d.ID,
			Platform:   d.Platform,
			Browser:    d.Browser,
			LastUsedAt: d.LastUsedAt,
			ExpiresAt:  d.ExpiresAt,
		})
	}

	pagination := common.PaginationResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
	}

	return response, pagination, nil
}

func (s *Service) GetCurrentSessions(ctx context.Context, userID uuid.UUID) ([]SessionResponse, error) {

	sessions, err := s.sessions.GetCurrentSessions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var response []SessionResponse
	for _, s := range sessions {
		response = append(response, SessionResponse{
			ID:         s.ID,
			Platform:   s.Platform,
			Browser:    s.Browser,
			LastUsedAt: s.LastUsedAt,
			ExpiresAt:  s.ExpiresAt,
		})
	}

	return response, nil
}
