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

type mockTopicService struct {
	err error
}

const validTopicId = "d94c0c2e-bf2b-44df-9e52-b2d83671766c"
const invalidTopicId = "Invalid_Topic_ID"

func (m *mockTopicService) CreateTopic(ctx context.Context, input service.CreateTopicInput) (*service.Topic, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.Topic{ID: uuid.New(), Label: "New Topic"}, nil
}

func (m *mockTopicService) GetAllTopics(ctx context.Context, grantWriterID uuid.UUID) ([]service.Topic, error) {
	if m.err != nil {
		return nil, m.err
	}

	return []service.Topic{{ID: uuid.New(), Label: "New T"}}, nil
}

func (m *mockTopicService) UpdateTopic(ctx context.Context, grantWriterID uuid.UUID, topicID uuid.UUID, newLabel string) (*service.Topic, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &service.Topic{ID: uuid.New(), Label: newLabel}, nil
}

func (m *mockTopicService) DeleteTopic(ctx context.Context, grantWriterID uuid.UUID, topicID uuid.UUID) error {
	if m.err != nil {
		return m.err
	}

	return nil
}

func TestCreateTopic(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		body           string
		err            error
		expectedStatus int
	}{{
		name:           "valid - 201",
		userId:         validUserId,
		body:           `{"label": "Test Label"}`,
		err:            nil,
		expectedStatus: http.StatusCreated,
	},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			body:           `{"label": "Test Label"}`,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing label - 400",
			userId:         validUserId,
			body:           `{}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid label - 400",
			userId:         validUserId,
			body:           `{"label": ""}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			body:           `{"label": "Test Label"}`,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewTopicHandler(&mockTopicService{tt.err})

			req := httptest.NewRequest(http.MethodPost, "/topics", strings.NewReader(tt.body))

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)

			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.CreateTopic(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestGetAllTopics(t *testing.T) {
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
			handler := handler.NewTopicHandler(&mockTopicService{tt.err})

			req := httptest.NewRequest(http.MethodGet, "/topics", nil)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)

			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.GetAllTopics(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestUpdateTopic(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		topicId        string
		body           string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 200",
			userId:         validUserId,
			topicId:        validTopicId,
			body:           `{"label": "Test Label"}`,
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			topicId:        validTopicId,
			body:           `{"label": "Test Label"}`,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid topic id - 400",
			userId:         validUserId,
			topicId:        invalidTopicId,
			body:           `{"label": "Test Label"}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid topic label - 400",
			userId:         validUserId,
			topicId:        validTopicId,
			body:           `{"label": ""}`,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "topic not found - 404",
			userId:         validUserId,
			topicId:        validTopicId,
			body:           `{"label": "Test Label"}`,
			err:            service.ErrTopicNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			topicId:        validTopicId,
			body:           `{"label": "Test Label"}`,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewTopicHandler(&mockTopicService{tt.err})

			req := httptest.NewRequest(http.MethodPut, "/topics", strings.NewReader(tt.body))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.topicId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.UpdateTopic(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestDeleteTopic(t *testing.T) {
	tests := []struct {
		name           string
		userId         string
		topicId        string
		err            error
		expectedStatus int
	}{
		{
			name:           "valid - 204",
			userId:         validUserId,
			topicId:        validTopicId,
			err:            nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "invalid user id - 401",
			userId:         invalidUserId,
			topicId:        validTopicId,
			err:            nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid topic id - 400",
			userId:         validUserId,
			topicId:        invalidTopicId,
			err:            nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "topic not found - 500",
			userId:         validUserId,
			topicId:        validTopicId,
			err:            service.ErrTopicNotFound,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "service error - 500",
			userId:         validUserId,
			topicId:        validTopicId,
			err:            errors.New("db connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler.NewTopicHandler(&mockTopicService{tt.err})

			req := httptest.NewRequest(http.MethodDelete, "/topics", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.topicId)

			ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, tt.userId)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.DeleteTopic(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
