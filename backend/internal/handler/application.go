package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jbechler2/grant-tool/backend/internal/auth"
	"github.com/jbechler2/grant-tool/backend/internal/service"
)

type ApplicationHandler struct {
	applicationService *service.ApplicationService
}

func NewApplicationHandler(applicationService *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{applicationService: applicationService}
}

type createApplicationRequest struct {
	GrantID     uuid.UUID `json:"grant_id"`
	ClientID    uuid.UUID `json:"client_id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	IsExclusive bool      `json:"is_exclusive"`
	Notes       *string   `json:"notes"`
}

type updateApplicationRequest struct {
	GrantWriterID uuid.UUID `json:"grant_writer_id"`
	GrantID       uuid.UUID `json:"grant_id"`
	Title         string    `json:"title"`
	Status        string    `json:"status"`
	IsExclusive   bool      `json:"is_exclusive"`
	Notes         *string   `json:"notes"`
}

type applicationResponse struct {
	ID            string    `json:"id"`
	GrantWriterID string    `json:"grant_writer_id"`
	GrantID       string    `json:"grant_id"`
	ClientID      string    `json:"client_id"`
	Title         string    `json:"title"`
	Status        string    `json:"status"`
	IsExclusive   bool      `json:"is_exclusive"`
	PublishedAt   time.Time `json:"published_at"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}

func (h *ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	var req createApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	grantWriterID, err := uuid.Parse(userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if req.GrantID == uuid.Nil {
		writeError(w, http.StatusBadRequest, "grant is required")
		return
	}

	if req.ClientID == uuid.Nil {
		writeError(w, http.StatusBadRequest, "client is required")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	if req.Status == "" {
		writeError(w, http.StatusBadRequest, "status is required")
		return
	}

	result, err := h.applicationService.CreateApplication(r.Context(), service.CreateApplicationInput{
		GrantWriterID: grantWriterID,
		GrantID:       req.GrantID,
		ClientID:      req.ClientID,
		Title:         req.Title,
		Status:        req.Status,
		IsExclusive:   req.IsExclusive,
		Notes:         req.Notes,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create application")
		return
	}

	writeJSON(w, http.StatusCreated, toApplicationResponse(result))
}

func (h *ApplicationHandler) GetAllApplicationsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	grantWriterID, err := uuid.Parse(userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	results, err := h.applicationService.GetAllApplicationsByUserID(r.Context(), grantWriterID)
	if err != nil {
		if errors.Is(err, service.ErrApplicationNotFound) {
			writeError(w, http.StatusNotFound, "application not found")
			return
		}
		writeError(w, http.StatusNotFound, "failed to retrieve applications")
	}

	applications := make([]applicationResponse, len(results))
	for i, a := range results {
		applications[i] = *toApplicationResponse(&a)
	}

	writeJSON(w, http.StatusOK, applications)
}

func (h *ApplicationHandler) GetAllApplicationsByClientID(w http.ResponseWriter, r *http.Request) {
	clientIDString := chi.URLParam(r, "id")
	if clientIDString == "" {
		writeError(w, http.StatusBadRequest, "no client id provided")
		return
	}

	clientID, err := uuid.Parse(clientIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid client id")
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	grantWriterID, err := uuid.Parse(userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	results, err := h.applicationService.GetAllApplicationsByClientID(r.Context(), grantWriterID, clientID)
	if err != nil {
		if errors.Is(err, service.ErrApplicationNotFound) {
			writeError(w, http.StatusNotFound, "application not found")
			return
		}
		writeError(w, http.StatusNotFound, "failed to retrieve applications")
	}

	applications := make([]applicationResponse, len(results))
	for i, a := range results {
		applications[i] = *toApplicationResponse(&a)
	}

	writeJSON(w, http.StatusOK, applications)
}

func (h *ApplicationHandler) GetApplicationByID(w http.ResponseWriter, r *http.Request) {
	applicationIDString := chi.URLParam(r, "id")
	if applicationIDString == "" {
		writeError(w, http.StatusBadRequest, "no client id provided")
		return
	}

	applicationID, err := uuid.Parse(applicationIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid client id")
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	grantWriterID, err := uuid.Parse(userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.applicationService.GetApplicationByID(r.Context(), grantWriterID, applicationID)
	if err != nil {
		if errors.Is(err, service.ErrApplicationNotFound) {
			writeError(w, http.StatusNotFound, "application not found")
			return
		}
		writeError(w, http.StatusNotFound, "failed to retrieve applications")
	}

	writeJSON(w, http.StatusOK, toApplicationResponse(result))
}

func (h *ApplicationHandler) UpdateApplication(w http.ResponseWriter, r *http.Request) {
	var req updateApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	applicationIDString := chi.URLParam(r, "id")
	if applicationIDString == "" {
		writeError(w, http.StatusBadRequest, "no application id provided")
		return
	}

	applicationID, err := uuid.Parse(applicationIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid application id")
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	grantWriterID, err := uuid.Parse(userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.applicationService.UpdateApplication(r.Context(), service.UpdateApplicationInput{
		ApplicationID: applicationID,
		GrantWriterID: grantWriterID,
		Title:         req.Title,
		Status:        req.Status,
		IsExclusive:   req.IsExclusive,
		Notes:         req.Notes,
	})
	if err != nil {
		if errors.Is(err, service.ErrApplicationNotFound) {
			writeError(w, http.StatusNotFound, "application not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update application")
		return
	}

	writeJSON(w, http.StatusOK, toApplicationResponse(result))
}

func (h *ApplicationHandler) PublishApplication(w http.ResponseWriter, r *http.Request) {
	applicationIDString := chi.URLParam(r, "id")
	if applicationIDString == "" {
		writeError(w, http.StatusBadRequest, "no application id provided")
		return
	}

	applicationID, err := uuid.Parse(applicationIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid application id")
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	grantWriterID, err := uuid.Parse(userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.applicationService.PublishApplication(r.Context(), grantWriterID, applicationID)
	if err != nil {
		if errors.Is(err, service.ErrApplicationNotFound) {
			writeError(w, http.StatusNotFound, "application not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to publish application")
		return
	}

	writeJSON(w, http.StatusOK, toApplicationResponse(result))
}

func (h *ApplicationHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	applicationIDString := chi.URLParam(r, "id")
	if applicationIDString == "" {
		writeError(w, http.StatusBadRequest, "no application id provided")
		return
	}

	applicationID, err := uuid.Parse(applicationIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid application id")
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	grantWriterID, err := uuid.Parse(userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err = h.applicationService.DeleteApplication(r.Context(), grantWriterID, applicationID)
	if err != nil {
		if errors.Is(err, service.ErrApplicationNotFound) {
			writeError(w, http.StatusNotFound, "application not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete application")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toApplicationResponse(a *service.Application) *applicationResponse {
	return &applicationResponse{
		ID:            a.ID.String(),
		GrantWriterID: a.GrantWriterID.String(),
		GrantID:       a.GrantID.String(),
		ClientID:      a.ClientID.String(),
		Title:         a.Title,
		Status:        a.Status,
		IsExclusive:   a.IsExclusive,
		PublishedAt:   derefTime(a.PublishedAt),
		Notes:         derefString(a.Notes),
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
		DeletedAt:     a.DeletedAt,
	}
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
