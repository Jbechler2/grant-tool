package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jbechler2/grant-tool/backend/internal/auth"
	"github.com/jbechler2/grant-tool/backend/internal/handler"
	"github.com/jbechler2/grant-tool/backend/internal/service"
)

type mockApplicationService struct {
	err error
}

func (m *mockApplicationService) CreateApplication(ctx context.Context, input service.CreateApplicationInput) (*service.Application, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.Application{ID: uuid.New(), Title: "New Application"}, nil
}

func (m *mockApplicationService) GetAllApplicationsByUserID(ctx context.Context, grantWriterID uuid.UUID) ([]service.Application, error) {
	if m.err != nil {
		return nil, m.err
	}

	return []service.Application{{ID: uuid.New(), Title: "New Application"}}, nil
}

func (m *mockApplicationService) GetAllApplicationsByClientID(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID) ([]service.Application, error) {
	if m.err != nil {
		return nil, m.err
	}

	return []service.Application{{ID: uuid.New(), Title: "New Application"}}, nil
}

func (m *mockApplicationService) GetApplicationByID(ctx context.Context, grantWriterID uuid.UUID, applicationID uuid.UUID) (*service.Application, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.Application{ID: uuid.New(), Title: "New Application"}, nil
}

func (m *mockApplicationService) UpdateApplication(ctx context.Context, input service.UpdateApplicationInput) (*service.Application, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.Application{ID: uuid.New(), Title: "New Application"}, nil
}

func (m *mockApplicationService) PublishApplication(ctx context.Context, grantWriterID uuid.UUID, applicationID uuid.UUID) (*service.Application, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.Application{ID: uuid.New(), Title: "New Application"}, nil
}

func (m *mockApplicationService) DeleteApplication(ctx context.Context, grantWriterID uuid.UUID, applicationID uuid.UUID) error {
	if m.err != nil {
		return m.err
	}

	return nil
}

const validApplicationId = "e2c5f9a4-3b6d-4f8e-9c1a-7d6b5e4f3a2c"
const invalidApplicationId = "Invalid_Application_ID"

func TestCreateApplication(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		body           string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 201",
			userId:         validUserId,
			body:           `{"grant_id": "` + validGrantId + `", "client_id": "` + validClientId + `", "title": "Test Title", "status": "draft"}`,
			err:            nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			body:           `{"grant_id": "` + validGrantId + `", "client_id": "` + validClientId + `", "title": "Test Title", "status": "draft"}`,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing grant id - 400",
			userId:         validUserId,
			body:           `{"client_id": "` + validClientId + `", "title": "Test Title", "status": "draft"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing client id - 400",
			userId:         validUserId,
			body:           `{"grant_id": "` + validGrantId + `", "title": "Test Title", "status": "draft"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing title - 400",
			userId:         validUserId,
			body:           `{"grant_id": "` + validGrantId + `", "client_id": "` + validClientId + `", "status": "draft"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing status - 400",
			userId:         validUserId,
			body:           `{"grant_id": "` + validGrantId + `", "client_id": "` + validClientId + `", "title": "Test Title"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			body:           `{"grant_id": "` + validGrantId + `", "client_id": "` + validClientId + `", "title": "Test Title", "status": "draft"}`,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewApplicationHandler(&mockApplicationService{tt.err})

			req := httptest.NewRequest(http.MethodPost, "/applications", strings.NewReader(tt.body))

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)

			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.CreateApplication(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetAllApplicationsByUserID(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid user id - 200",
			userId:         validUserId,
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewApplicationHandler(&mockApplicationService{tt.err})

			req := httptest.NewRequest(http.MethodGet, "/applications", nil)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)

			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.GetAllApplicationsByUserID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetAllApplicationsByClientID(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		clientId       string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			userId:         validUserId,
			clientId:       validClientId,
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			clientId:       validClientId,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid client id - 400",
			userId:         validUserId,
			clientId:       invalidClientId,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			clientId:       validClientId,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewApplicationHandler(&mockApplicationService{tt.err})

			req := httptest.NewRequest(http.MethodGet, "/clients", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.clientId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.GetAllApplicationsByClientID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetApplicationByID(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		applicationId  string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			userId:         validUserId,
			applicationId:  validApplicationId,
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			applicationId:  validApplicationId,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid application id - 400",
			userId:         validUserId,
			applicationId:  invalidApplicationId,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "application not found - 404",
			userId:         validUserId,
			applicationId:  validApplicationId,
			err:            service.ErrApplicationNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			applicationId:  validApplicationId,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewApplicationHandler(&mockApplicationService{err: tt.err})

			req := httptest.NewRequest(http.MethodGet, "/applications", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.applicationId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.GetApplicationByID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestUpdateApplication(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		applicationId  string
		body           string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			userId:         validUserId,
			applicationId:  validApplicationId,
			body:           `{"title": "Updated Title", "status": "submitted"}`,
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			applicationId:  validApplicationId,
			body:           `{"title": "Updated Title", "status": "submitted"}`,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid application id - 400",
			userId:         validUserId,
			applicationId:  invalidApplicationId,
			body:           `{"title": "Updated Title", "status": "submitted"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "application not found - 404",
			userId:         validUserId,
			applicationId:  validApplicationId,
			body:           `{"title": "Updated Title", "status": "submitted"}`,
			err:            service.ErrApplicationNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			applicationId:  validApplicationId,
			body:           `{"title": "Updated Title", "status": "submitted"}`,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewApplicationHandler(&mockApplicationService{tt.err})

			req := httptest.NewRequest(http.MethodPut, "/applications", strings.NewReader(tt.body))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.applicationId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.UpdateApplication(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestPublishApplication(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		applicationId  string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			userId:         validUserId,
			applicationId:  validApplicationId,
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			applicationId:  validApplicationId,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid application id - 400",
			userId:         validUserId,
			applicationId:  invalidApplicationId,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "application not found - 404",
			userId:         validUserId,
			applicationId:  validApplicationId,
			err:            service.ErrApplicationNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			applicationId:  validApplicationId,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewApplicationHandler(&mockApplicationService{tt.err})

			req := httptest.NewRequest(http.MethodPost, "/applications", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.applicationId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.PublishApplication(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestDeleteApplication(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		applicationId  string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 204",
			userId:         validUserId,
			applicationId:  validApplicationId,
			err:            nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			applicationId:  validApplicationId,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid application id - 400",
			userId:         validUserId,
			applicationId:  invalidApplicationId,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "application not found - 404",
			userId:         validUserId,
			applicationId:  validApplicationId,
			err:            service.ErrApplicationNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			applicationId:  validApplicationId,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewApplicationHandler(&mockApplicationService{tt.err})

			req := httptest.NewRequest(http.MethodDelete, "/applications", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.applicationId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.DeleteApplication(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
