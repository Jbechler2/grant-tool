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

type GrantHandler struct {
	grantService *service.GrantService
}

func NewGrantHandler(grantService *service.GrantService) *GrantHandler {
	return &GrantHandler{grantService: grantService}
}

type createGrantRequest struct {
	Title                     string   `json:"title"`
	FunderName                string   `json:"funder_name"`
	FunderWebsite             *string  `json:"funder_website"`
	Description               *string  `json:"description"`
	AwardAmountMin            *float64 `json:"award_amount_min"`
	AwardAmountMax            *float64 `json:"award_amount_max"`
	EligibilityNotes          *string  `json:"eligibility_notes"`
	EstimatedApplicationHours *float64 `json:"estimated_application_hours"`
}

type updateGrantRequest struct {
	Title                     *string  `json:"title"`
	FunderName                *string  `json:"funder_name"`
	FunderWebsite             *string  `json:"funder_website"`
	Description               *string  `json:"description"`
	AwardAmountMin            *float64 `json:"award_amount_min"`
	AwardAmountMax            *float64 `json:"award_amount_max"`
	EligibilityNotes          *string  `json:"eligibility_notes"`
	EstimatedApplicationHours *float64 `json:"estimated_application_hours"`
	Visibility                *string  `json:"visibility"`
}

type addDeadlineRequest struct {
	Label       string  `json:"label"`
	Date        string  `json:"date"`
	Description *string `json:"description"`
}

type grantResponse struct {
	ID                        string    `json:"id"`
	GrantWriterID             string    `json:"grant_writer_id"`
	Title                     string    `json:"title"`
	FunderName                string    `json:"funder_name"`
	FunderWebsite             string    `json:"funder_website"`
	Description               string    `json:"description"`
	AwardAmountMin            *float64  `json:"award_amount_min"`
	AwardAmountMax            *float64  `json:"award_amount_max"`
	EligibilityNotes          string    `json:"eligibility_notes"`
	EstimatedApplicationHours *float64  `json:"estimated_application_hours"`
	Visibility                string    `json:"visibility"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type deadlineResponse struct {
	ID          string    `json:"id"`
	GrantID     string    `json:"grant_id"`
	Label       string    `json:"label"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (h *GrantHandler) CreateGrant(w http.ResponseWriter, r *http.Request) {
	var req createGrantRequest
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

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	if req.FunderName == "" {
		writeError(w, http.StatusBadRequest, "funder name is required")
		return
	}

	result, err := h.grantService.CreateGrant(r.Context(), service.CreateGrantInput{
		GrantWriterID:             grantWriterID,
		Title:                     req.Title,
		FunderName:                req.FunderName,
		FunderWebsite:             req.FunderWebsite,
		Description:               req.Description,
		AwardAmountMin:            req.AwardAmountMin,
		AwardAmountMax:            req.AwardAmountMax,
		EligibilityNotes:          req.EligibilityNotes,
		EstimatedApplicationHours: req.EstimatedApplicationHours,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create grant")
		return
	}

	writeJSON(w, http.StatusCreated, toGrantResponse(result))
}

func (h *GrantHandler) GetGrantByID(w http.ResponseWriter, r *http.Request) {
	grantIDString := chi.URLParam(r, "id")
	if grantIDString == "" {
		writeError(w, http.StatusBadRequest, "no grant id provided")
		return
	}

	grantID, err := uuid.Parse(grantIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid grant id")
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

	result, err := h.grantService.GetGrantByID(r.Context(), grantWriterID, grantID)
	if err != nil {
		if errors.Is(err, service.ErrGrantNotFound) {
			writeError(w, http.StatusNotFound, "grant not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to retrieve grant")
		return
	}

	writeJSON(w, http.StatusOK, toGrantResponse(result))
}

func (h *GrantHandler) GetAllGrants(w http.ResponseWriter, r *http.Request) {
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

	results, err := h.grantService.GetAllGrants(r.Context(), grantWriterID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve grants")
		return
	}

	grants := make([]grantResponse, len(results))
	for i, g := range results {
		grants[i] = *toGrantResponse(&g)
	}

	writeJSON(w, http.StatusOK, grants)
}

func (h *GrantHandler) UpdateGrant(w http.ResponseWriter, r *http.Request) {
	var req updateGrantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	grantIDString := chi.URLParam(r, "id")
	if grantIDString == "" {
		writeError(w, http.StatusBadRequest, "no grant id provided")
		return
	}

	grantID, err := uuid.Parse(grantIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid grant id")
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

	if req.Title != nil && *req.Title == "" {
		writeError(w, http.StatusBadRequest, "title cannot be empty")
		return
	}

	if req.FunderName != nil && *req.FunderName == "" {
		writeError(w, http.StatusBadRequest, "funder name cannot be empty")
		return
	}

	result, err := h.grantService.UpdateGrant(r.Context(), service.UpdateGrantInput{
		ID:                        grantID,
		GrantWriterID:             grantWriterID,
		Title:                     req.Title,
		FunderName:                req.FunderName,
		FunderWebsite:             req.FunderWebsite,
		Description:               req.Description,
		AwardAmountMin:            req.AwardAmountMin,
		AwardAmountMax:            req.AwardAmountMax,
		EligibilityNotes:          req.EligibilityNotes,
		EstimatedApplicationHours: req.EstimatedApplicationHours,
		Visibility:                req.Visibility,
	})
	if err != nil {
		if errors.Is(err, service.ErrGrantNotFound) {
			writeError(w, http.StatusNotFound, "grant not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update grant")
		return
	}

	writeJSON(w, http.StatusOK, toGrantResponse(result))
}

func (h *GrantHandler) DeleteGrant(w http.ResponseWriter, r *http.Request) {
	grantIDString := chi.URLParam(r, "id")
	if grantIDString == "" {
		writeError(w, http.StatusBadRequest, "no grant id provided")
		return
	}

	grantID, err := uuid.Parse(grantIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid grant id")
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

	err = h.grantService.DeleteGrant(r.Context(), grantWriterID, grantID)
	if err != nil {
		if errors.Is(err, service.ErrGrantNotFound) {
			writeError(w, http.StatusNotFound, "grant not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete grant")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GrantHandler) AddDeadline(w http.ResponseWriter, r *http.Request) {
	var req addDeadlineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	grantIDString := chi.URLParam(r, "id")
	if grantIDString == "" {
		writeError(w, http.StatusBadRequest, "no grant id provided")
		return
	}

	grantID, err := uuid.Parse(grantIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid grant id")
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

	if req.Label == "" {
		writeError(w, http.StatusBadRequest, "label is required")
		return
	}

	if req.Date == "" {
		writeError(w, http.StatusBadRequest, "date is required")
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
		return
	}

	result, err := h.grantService.AddDeadline(r.Context(), service.AddDeadlineInput{
		GrantWriterID: grantWriterID,
		GrantID:       grantID,
		Label:         req.Label,
		Date:          date,
		Description:   req.Description,
	})
	if err != nil {
		if errors.Is(err, service.ErrGrantNotFound) {
			writeError(w, http.StatusNotFound, "grant not found")
			return
		}
		if errors.Is(err, service.ErrInvalidDeadlineLabel) {
			writeError(w, http.StatusBadRequest, "invalid deadline label")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to add deadline")
		return
	}

	writeJSON(w, http.StatusCreated, toDeadlineResponse(result))
}

func (h *GrantHandler) GetDeadlinesByGrantID(w http.ResponseWriter, r *http.Request) {
	grantIDString := chi.URLParam(r, "id")
	if grantIDString == "" {
		writeError(w, http.StatusBadRequest, "no grant id provided")
		return
	}

	grantID, err := uuid.Parse(grantIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid grant id")
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

	results, err := h.grantService.GetDeadlinesByGrantID(r.Context(), grantWriterID, grantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve grants")
		return
	}

	grants := make([]deadlineResponse, len(results))
	for i, g := range results {
		grants[i] = *toDeadlineResponse(&g)
	}

	writeJSON(w, http.StatusOK, grants)

}

func (h *GrantHandler) DeleteDeadline(w http.ResponseWriter, r *http.Request) {
	grantIDString := chi.URLParam(r, "id")
	if grantIDString == "" {
		writeError(w, http.StatusBadRequest, "no grant id provided")
		return
	}

	grantID, err := uuid.Parse(grantIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid grant id")
		return
	}

	deadlineIDString := chi.URLParam(r, "deadlineID")
	if deadlineIDString == "" {
		writeError(w, http.StatusBadRequest, "no deadline id provided")
		return
	}

	deadlineID, err := uuid.Parse(deadlineIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid deadline id")
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

	err = h.grantService.DeleteDeadline(r.Context(), grantWriterID, grantID, deadlineID)
	if err != nil {
		if errors.Is(err, service.ErrGrantNotFound) {
			writeError(w, http.StatusNotFound, "grant not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete deadline")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toGrantResponse(g *service.Grant) *grantResponse {
	return &grantResponse{
		ID:                        g.ID.String(),
		GrantWriterID:             g.GrantWriterID.String(),
		Title:                     g.Title,
		FunderName:                g.FunderName,
		FunderWebsite:             g.FunderWebsite,
		Description:               g.Description,
		AwardAmountMin:            g.AwardAmountMin,
		AwardAmountMax:            g.AwardAmountMax,
		EligibilityNotes:          g.EligibilityNotes,
		EstimatedApplicationHours: g.EstimatedApplicationHours,
		Visibility:                g.Visibility,
		CreatedAt:                 g.CreatedAt,
		UpdatedAt:                 g.UpdatedAt,
	}
}

func toDeadlineResponse(d *service.Deadline) *deadlineResponse {
	return &deadlineResponse{
		ID:          d.ID.String(),
		GrantID:     d.GrantID.String(),
		Label:       d.Label,
		Date:        d.Date,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
	}
}
