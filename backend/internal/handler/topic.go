package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jbechler2/grant-tool/backend/internal/auth"
	"github.com/jbechler2/grant-tool/backend/internal/service"
)

type TopicHandler struct {
	topicService *service.TopicService
}

func NewTopicHandler(topicService *service.TopicService) *TopicHandler {
	return &TopicHandler{topicService: topicService}
}

type createTopicRequest struct {
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
