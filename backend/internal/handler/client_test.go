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

type mockClientService struct {
	err error
}

func (m *mockClientService) CreateClient(ctx context.Context, input service.CreateClientInput) (*service.Client, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &service.Client{ID: uuid.New(), Name: "Test Client"}, nil
}

func (m *mockClientService) GetAllClients(ctx context.Context, grantWriterID uuid.UUID) ([]service.Client, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []service.Client{{ID: uuid.New(), Name: "Test Client"}}, nil
}

func (m *mockClientService) GetClientByID(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID) (*service.Client, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.Client{ID: uuid.New(), Name: "Test Client"}, nil
}

func (m *mockClientService) UpdateClient(ctx context.Context, input service.UpdateClientInput) (*service.Client, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &service.Client{ID: uuid.New(), Name: "Test Client"}, nil
}

func (m *mockClientService) DeleteClient(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

const validUserId = "0a7ea2a3-5a54-4686-92d5-c2a7e6709660"
const validClientId = "d468f6f9-b59b-4660-b270-647a13ab9273"
const invalidUserId = "Invalid_User_ID"
const invalidClientId = "Invalid_Client_ID"

func TestCreateClient(t *testing.T) {
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
			body:           `{"name": "Test Name"}`,
			err:            nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			body:           `{"name": "Test Name"}`,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid client name - 400",
			userId:         validUserId,
			body:           `{"email": "test@test.com"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			body:           `{"name": "Test Name"}`,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewClientHandler(&mockClientService{tt.err})

			req := httptest.NewRequest(http.MethodPost, "/clients", strings.NewReader(tt.body))

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)

			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.CreateClient(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetAllClients(t *testing.T) {
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
			handler := handler.NewClientHandler(&mockClientService{tt.err})

			req := httptest.NewRequest(http.MethodGet, "/clients", nil)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)

			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.GetAllClients(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetClientByID(t *testing.T) {
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
			name:           "client not found - 404",
			userId:         validUserId,
			clientId:       validClientId,
			err:            service.ErrClientNotFound,
			expectedStatus: http.StatusNotFound,
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
			handler := handler.NewClientHandler(&mockClientService{err: tt.err})

			req := httptest.NewRequest(http.MethodGet, "/clients", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.clientId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.GetClientByID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)

			}
		})
	}
}

func TestUpdateClient(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		clientId       string
		body           string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			userId:         validUserId,
			clientId:       validClientId,
			body:           `{"name": "Test Name"}`,
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			clientId:       validClientId,
			body:           `{"name": "Test Name"}`,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid client id - 400",
			userId:         validUserId,
			clientId:       invalidClientId,
			body:           `{"name": "Test Name"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid client name - 400",
			userId:         validUserId,
			clientId:       validClientId,
			body:           `{"email": "test@test.com"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "client not found - 404",
			userId:         validUserId,
			clientId:       validClientId,
			body:           `{"email": "test@test.com"}`,
			err:            service.ErrClientNotFound,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			clientId:       validClientId,
			body:           `{"name": "Test Name"}`,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewClientHandler(&mockClientService{tt.err})

			req := httptest.NewRequest(http.MethodPut, "/clients", strings.NewReader(tt.body))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.clientId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.UpdateClient(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestDeleteClient(t *testing.T) {
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
			expectedStatus: http.StatusNoContent,
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
			name:           "client not found - 404",
			userId:         validUserId,
			clientId:       validClientId,
			err:            service.ErrClientNotFound,
			expectedStatus: http.StatusNotFound,
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
			handler := handler.NewClientHandler(&mockClientService{tt.err})

			req := httptest.NewRequest(http.MethodDelete, "/clients", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.clientId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.DeleteClient(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
