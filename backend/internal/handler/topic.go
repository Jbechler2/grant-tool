package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jbechler2/grant-tool/backend/internal/auth"
	"github.com/jbechler2/grant-tool/backend/internal/service"
)

type TopicServicer interface {
	CreateTopic(ctx context.Context, input service.CreateTopicInput) (*service.Topic, error)
	GetAllTopics(ctx context.Context, grantWriterID uuid.UUID) ([]service.Topic, error)
	UpdateTopic(ctx context.Context, grantWriterID uuid.UUID, topicID uuid.UUID, newLabel string) (*service.Topic, error)
	DeleteTopic(ctx context.Context, grantWriterID uuid.UUID, topicID uuid.UUID) error
}

type TopicHandler struct {
	topicService TopicServicer
}

func NewTopicHandler(topicService TopicServicer) *TopicHandler {
	return &TopicHandler{topicService: topicService}
}

type createTopicRequest struct {
	Label string `json:"label"`
}

type updateTopicRequest struct {
	Label string `json:"label"`
}

func (h *TopicHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	var req createTopicRequest
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

	if req.Label == "" {
		writeError(w, http.StatusBadRequest, "label is required")
		return
	}

	result, err := h.topicService.CreateTopic(r.Context(), service.CreateTopicInput{
		GrantWriterID: grantWriterID,
		Label:         req.Label,
	})

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create topic")
		return
	}

	writeJSON(w, http.StatusCreated, toTopicResponse(result))
}

func (h *TopicHandler) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	var req updateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	topicIDString := chi.URLParam(r, "id")
	if topicIDString == "" {
		writeError(w, http.StatusBadRequest, "no topic id provided")
		return
	}

	topicID, err := uuid.Parse(topicIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid topic id")
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
		writeError(w, http.StatusBadRequest, "Label cannot be empty")
		return
	}

	result, err := h.topicService.UpdateTopic(r.Context(), grantWriterID, topicID, req.Label)
	if err != nil {
		if errors.Is(err, service.ErrTopicNotFound) {
			writeError(w, http.StatusNotFound, "topic not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update topic")
		return
	}

	writeJSON(w, http.StatusOK, toTopicResponse(result))

}

func (h *TopicHandler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	topicIDString := chi.URLParam(r, "id")
	if topicIDString == "" {
		writeError(w, http.StatusBadRequest, "no topic id provided")
		return
	}

	topicID, err := uuid.Parse(topicIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid topic id")
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

	err = h.topicService.DeleteTopic(r.Context(), grantWriterID, topicID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete topic")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TopicHandler) GetAllTopics(w http.ResponseWriter, r *http.Request) {
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

	results, err := h.topicService.GetAllTopics(r.Context(), grantWriterID)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create topic")
		return
	}

	writeJSON(w, http.StatusOK, results)
}
