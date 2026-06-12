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

type mockGrantService struct {
	err error
}

func (m *mockGrantService) CreateGrant(ctx context.Context, input service.CreateGrantInput) (*service.Grant, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.Grant{ID: uuid.New(), Title: "New Grant"}, nil
}

func (m *mockGrantService) GetGrantByID(ctx context.Context, grantWriterID uuid.UUID, grantID uuid.UUID) (*service.Grant, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.Grant{ID: uuid.New(), Title: "New Grant"}, nil
}

func (m *mockGrantService) GetAllGrants(ctx context.Context, grantWriterID uuid.UUID) ([]service.Grant, error) {
	if m.err != nil {
		return nil, m.err
	}

	return []service.Grant{{ID: uuid.New(), Title: "New Grant"}}, nil
}

func (m *mockGrantService) UpdateGrant(ctx context.Context, input service.UpdateGrantInput) (*service.Grant, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.Grant{ID: uuid.New(), Title: "New Grant"}, nil
}

func (m *mockGrantService) DeleteGrant(ctx context.Context, grantWriterID uuid.UUID, grantID uuid.UUID) error {
	if m.err != nil {
		return m.err
	}

	return nil
}

func (m *mockGrantService) AddDeadline(ctx context.Context, input service.AddDeadlineInput) (*service.Deadline, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.Deadline{ID: uuid.New(), Label: "New Deadline"}, nil
}

func (m *mockGrantService) GetDeadlinesByGrantID(ctx context.Context, grantWriterID uuid.UUID, grantID uuid.UUID) ([]service.Deadline, error) {
	if m.err != nil {
		return nil, m.err
	}

	return []service.Deadline{{ID: uuid.New(), Label: "New Deadline"}}, nil
}

func (m *mockGrantService) DeleteDeadline(ctx context.Context, grantWriterID uuid.UUID, grantID uuid.UUID, deadlineID uuid.UUID) error {
	if m.err != nil {
		return m.err
	}

	return nil
}

const validGrantId = "d468f6f9-b59b-4660-b270-647a13ab9273"
const invalidGrantId = "Invalid_Grant_ID"
const validDeadlineId = "7131d5d3-83c4-4ea2-b1ea-2f5f90833e57"
const invalidDeadlineId = "Invalid_Deadline_ID"

func TestCreateGrant(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		body           string
		err            error
		expectedStatus int
	}{{
		name:           "valid - 201",
		userId:         validUserId,
		body:           `{"title": "Test Title", "funder_name": "Test Funder"}`,
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
			name:           "invalid grant title - 400",
			userId:         validUserId,
			body:           `{"email": "test@test.com"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			body:           `{"title": "Test Title", "funder_name": "Test Funder"}`,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewGrantHandler(&mockGrantService{tt.err})

			req := httptest.NewRequest(http.MethodPost, "/grants", strings.NewReader(tt.body))

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)

			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.CreateGrant(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetGrantByID(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		grantId        string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			userId:         validUserId,
			grantId:        validGrantId,
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			grantId:        validGrantId,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid grant id - 400",
			userId:         validUserId,
			grantId:        invalidGrantId,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "grant not found - 404",
			userId:         validUserId,
			grantId:        validGrantId,
			err:            service.ErrGrantNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			grantId:        validGrantId,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewGrantHandler(&mockGrantService{err: tt.err})

			req := httptest.NewRequest(http.MethodGet, "/grants", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.grantId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.GetGrantByID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)

			}
		})
	}
}

func TestGetAllGrants(t *testing.T) {
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
			handler := handler.NewGrantHandler(&mockGrantService{tt.err})

			req := httptest.NewRequest(http.MethodGet, "/grants", nil)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)

			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.GetAllGrants(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestUpdateGrant(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		grantId        string
		body           string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			userId:         validUserId,
			grantId:        validGrantId,
			body:           `{"name": "Test Name"}`,
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			grantId:        validGrantId,
			body:           `{"name": "Test Name"}`,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid grant id - 400",
			userId:         validUserId,
			grantId:        invalidGrantId,
			body:           `{"name": "Test Name"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid grant name - 400",
			userId:         validUserId,
			grantId:        validGrantId,
			body:           `{"title": ""}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "grant not found - 404",
			userId:         validUserId,
			grantId:        validGrantId,
			body:           `{"email": "test@test.com"}`,
			err:            service.ErrGrantNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			grantId:        validGrantId,
			body:           `{"name": "Test Name"}`,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewGrantHandler(&mockGrantService{tt.err})

			req := httptest.NewRequest(http.MethodPut, "/grants", strings.NewReader(tt.body))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.grantId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.UpdateGrant(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestDeleteGrant(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		grantId        string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			userId:         validUserId,
			grantId:        validGrantId,
			err:            nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			grantId:        validGrantId,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid grant id - 400",
			userId:         validUserId,
			grantId:        invalidGrantId,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "grant not found - 404",
			userId:         validUserId,
			grantId:        validGrantId,
			err:            service.ErrGrantNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			grantId:        validGrantId,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewGrantHandler(&mockGrantService{tt.err})

			req := httptest.NewRequest(http.MethodDelete, "/grants", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.grantId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.DeleteGrant(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestAddDeadline(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		grantId        string
		body           string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 201",
			userId:         validUserId,
			grantId:        validGrantId,
			body:           `{"label": "test", "date": "2026-01-02"}`,
			err:            nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			grantId:        validGrantId,
			body:           `{"label": "test", "date": "2026-01-02"}`,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid grant id - 400",
			userId:         validUserId,
			grantId:        invalidGrantId,
			body:           `{"label": "test", "date": "2026-01-02"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing deadline label - 400",
			userId:         validUserId,
			grantId:        validGrantId,
			body:           `{"date": "2026-01-02"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid deadline label - 400",
			userId:         validUserId,
			grantId:        validGrantId,
			body:           `{"label": "", "date": "2026-01-02"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty deadline date - 400",
			userId:         validUserId,
			grantId:        validGrantId,
			body:           `{"label": "test", "date": ""}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "malformed deadline date - 400",
			userId:         validUserId,
			grantId:        validGrantId,
			body:           `{"label": "test", "date": "10/10/10"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "grant not found - 404",
			userId:         validUserId,
			grantId:        validGrantId,
			body:           `{"label": "test", "date": "2026-01-02"}`,
			err:            service.ErrGrantNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			grantId:        validGrantId,
			body:           `{"label": "test", "date": "2026-01-02"}`,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewGrantHandler(&mockGrantService{tt.err})

			req := httptest.NewRequest(http.MethodPost, "/grants", strings.NewReader(tt.body))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.grantId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.AddDeadline(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetAllDeadlinesByGrantID(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		grantId        string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid user id - 200",
			userId:         validUserId,
			grantId:        validGrantId,
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			grantId:        validGrantId,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid grant id - 400",
			userId:         validUserId,
			grantId:        invalidGrantId,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			grantId:        validGrantId,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewGrantHandler(&mockGrantService{tt.err})

			req := httptest.NewRequest(http.MethodGet, "/grants", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.grantId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.GetDeadlinesByGrantID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestDeleteDeadline(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		grantId        string
		deadlineId     string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			userId:         validUserId,
			grantId:        validGrantId,
			deadlineId:     validDeadlineId,
			err:            nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			grantId:        validGrantId,
			deadlineId:     validDeadlineId,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid grant id - 400",
			userId:         validUserId,
			grantId:        invalidGrantId,
			deadlineId:     validDeadlineId,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid deadline id - 400",
			userId:         validUserId,
			grantId:        validGrantId,
			deadlineId:     invalidDeadlineId,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			grantId:        validGrantId,
			deadlineId:     validDeadlineId,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewGrantHandler(&mockGrantService{tt.err})

			req := httptest.NewRequest(http.MethodDelete, "/grants", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.grantId)
			rctx.URLParams.Add("deadlineID", tt.deadlineId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.DeleteDeadline(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
