package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jbechler2/grant-tool/backend/internal/handler"
	"github.com/jbechler2/grant-tool/backend/internal/service"
)

type mockAuthService struct {
	err error
}

func (m *mockAuthService) Register(ctx context.Context, input service.RegisterInput) (*service.AuthResult, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.AuthResult{Token: "test", RefreshToken: "test", RefreshExpiry: time.Now()}, nil
}

func (m *mockAuthService) Login(ctx context.Context, input service.LoginInput) (*service.AuthResult, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.AuthResult{Token: "test", RefreshToken: "test", RefreshExpiry: time.Now()}, nil
}

func (m *mockAuthService) RotateToken(ctx context.Context, tokenValue string, input service.RotateTokenInput) (*service.AuthResult, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.AuthResult{Token: "test", RefreshToken: "test", RefreshExpiry: time.Now()}, nil
}

func (m *mockAuthService) Logout(ctx context.Context, tokenValue string) error {
	if m.err != nil {
		return m.err
	}

	return nil
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 201",
			body:           `{"email": "test@test.com", "password": "password123"}`,
			err:            nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "missing email - 400",
			body:           `{"password": "password123"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid email - 400",
			body:           `{"email": "test", "password": "password123"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "mising password - 400",
			body:           `{"email": "test"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty password - 400",
			body:           `{"email": "test", "password": ""}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid password - 400",
			body:           `{"email": "test", "password": "pass"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "existing account - 409",
			body:           `{"email": "test@test.com", "password": "password123"}`,
			err:            service.ErrEmailTaken,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "service error - 500",
			body:           `{"email": "test@test.com", "password": "password123"}`,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewAuthHandler(&mockAuthService{tt.err}, false)

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.body))

			rr := httptest.NewRecorder()

			handler.Register(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			body:           `{"email": "test@test.com", "password": "password123"}`,
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing email - 400",
			body:           `{"password": "password123"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid email - 400",
			body:           `{"email": "", "password": "password123"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "mising password - 400",
			body:           `{"email": "test"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid password - 400",
			body:           `{"email": "test", "password": ""}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid creds - 401",
			body:           `{"email": "test@test.com", "password": "password123"}`,
			err:            service.ErrInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "too many sessions - 403",
			body:           `{"email": "test@test.com", "password": "password123"}`,
			err:            service.ErrTooManySessions,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "service error - 500",
			body:           `{"email": "test", "password": "password123"}`,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewAuthHandler(&mockAuthService{tt.err}, false)

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))

			rr := httptest.NewRecorder()

			handler.Login(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestRotateToken(t *testing.T) {
	tests := []struct {
		name           string
		cookieName     string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			cookieName:     "refresh_token",
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid cookie name - 401",
			cookieName:     "invalid",
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "refresh token not found - 401",
			cookieName:     "refresh_token",
			err:            service.ErrRefreshTokenNotFound,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "refresh token expired - 401",
			cookieName:     "refresh_token",
			err:            service.ErrRefreshTokenExpired,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "service error - 500",
			cookieName:     "refresh_token",
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewAuthHandler(&mockAuthService{tt.err}, false)

			cookie := &http.Cookie{
				Name:  tt.cookieName,
				Value: "test-refresh-token",
			}

			req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
			req.AddCookie(cookie)
			rr := httptest.NewRecorder()

			handler.Refresh(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name           string
		cookieName     string
		cookieValue    string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			cookieName:     "refresh_token",
			cookieValue:    "validCookie",
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid refresh token - 200",
			cookieName:     "invalid",
			cookieValue:    "validCookie",
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty refresh token - 200",
			cookieName:     "refresh_token",
			cookieValue:    "",
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "logout anyway - 200",
			cookieName:     "refresh_token",
			cookieValue:    "validCookie",
			err:            errors.New("cleaup failure"),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewAuthHandler(&mockAuthService{tt.err}, false)

			cookie := &http.Cookie{
				Name:  tt.cookieName,
				Value: tt.cookieValue,
			}
			req := httptest.NewRequest(http.MethodPost, "/logout", nil)
			req.AddCookie(cookie)
			rr := httptest.NewRecorder()

			handler.Logout(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
