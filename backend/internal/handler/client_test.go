package handler_test

import (
	"context"
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

type mockClientService struct{}

func (m *mockClientService) CreateClient(ctx context.Context, input service.CreateClientInput) (*service.Client, error) {
	return &service.Client{ID: uuid.New(), Name: "Test Client"}, nil
}

func (m *mockClientService) GetAllClients(ctx context.Context, grantWriterID uuid.UUID) ([]service.Client, error) {
	return []service.Client{{ID: uuid.New(), Name: "Test Client"}}, nil
}

func (m *mockClientService) GetClientByID(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID) (*service.Client, error) {
	return &service.Client{ID: uuid.New(), Name: "Test Client"}, nil
}

func (m *mockClientService) UpdateClient(ctx context.Context, input service.UpdateClientInput) (*service.Client, error) {
	return &service.Client{ID: uuid.New(), Name: "Test Client"}, nil
}

func (m *mockClientService) DeleteClient(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID) error {
	return nil
}

func TestGetAllClients_ReturnOK(t *testing.T) {
	handler := handler.NewClientHandler(&mockClientService{})

	req := httptest.NewRequest(http.MethodGet, "/clients", nil)

	ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, uuid.New().String())

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.GetAllClients(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

}

func TestGetAllClients_InvalidUserID(t *testing.T) {
	handler := handler.NewClientHandler(&mockClientService{})

	req := httptest.NewRequest(http.MethodGet, "/clients", nil)

	ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, "Invalid_Credential")

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.GetAllClients(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestCreateClient_ReturnOK(t *testing.T) {
	handler := handler.NewClientHandler(&mockClientService{})

	body := `{"Name": "Jason"}`

	req := httptest.NewRequest(http.MethodPost, "/clients", strings.NewReader(body))

	ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, uuid.New().String())

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.CreateClient(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
}

func TestCreateClient_NoName(t *testing.T) {
	handler := handler.NewClientHandler(&mockClientService{})

	body := `{"Email": "jason@test.com"}`

	req := httptest.NewRequest(http.MethodPost, "/clients", strings.NewReader(body))

	ctx := context.WithValue(req.Context(), auth.ContextKeyUserID, uuid.New().String())

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.CreateClient(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestGetClientByID_ReturnsOK(t *testing.T) {
	handler := handler.NewClientHandler(&mockClientService{})

	req := httptest.NewRequest(http.MethodGet, "/clients", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", uuid.New().String())

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, auth.ContextKeyUserID, uuid.New().String())

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.GetClientByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetClientByID_InvalidClientID(t *testing.T) {
	handler := handler.NewClientHandler(&mockClientService{})

	req := httptest.NewRequest(http.MethodGet, "/clients", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "Invalid_Client_ID")

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, auth.ContextKeyUserID, uuid.New().String())

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.GetClientByID(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestGetClientID_InvalidUserID(t *testing.T) {
	handler := handler.NewClientHandler(&mockClientService{})

	req := httptest.NewRequest(http.MethodGet, "/clients", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", uuid.New().String())

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, auth.ContextKeyUserID, "Invalid_User_ID")

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.GetClientByID(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestUpdateClient_ReturnsOK(t *testing.T) {
	handler := handler.NewClientHandler(&mockClientService{})

	body := `{"Name": "client_name"}`

	req := httptest.NewRequest(http.MethodPost, "/clients", strings.NewReader(body))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", uuid.New().String())

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, auth.ContextKeyUserID, uuid.New().String())

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.UpdateClient(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestUpdateClient_NoName(t *testing.T) {
	handler := handler.NewClientHandler(&mockClientService{})

	body := `{"email": "test@test.com"}`

	req := httptest.NewRequest(http.MethodPut, "/clients", strings.NewReader(body))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", uuid.New().String())

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, auth.ContextKeyUserID, uuid.New().String())

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.UpdateClient(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestDeleteClient_ReturnsOK(t *testing.T) {
	handler := handler.NewClientHandler(&mockClientService{})

	req := httptest.NewRequest(http.MethodDelete, "/clients", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", uuid.New().String())

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, auth.ContextKeyUserID, uuid.New().String())

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.DeleteClient(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rr.Code)
	}
}

func TestDeleteClient_InvalidClientID(t *testing.T) {
	handler := handler.NewClientHandler(&mockClientService{})

	req := httptest.NewRequest(http.MethodDelete, "/clients", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "Invalid_Client_ID")

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, auth.ContextKeyUserID, uuid.New().String())

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.DeleteClient(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestDeleteClient_InvalidUserID(t *testing.T) {
	handler := handler.NewClientHandler(&mockClientService{})

	req := httptest.NewRequest(http.MethodDelete, "/clients", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", uuid.New().String())

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, auth.ContextKeyUserID, "Invalid_User_ID")

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.DeleteClient(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}
