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

type ClientServicer interface {
	CreateClient(ctx context.Context, input service.CreateClientInput) (*service.Client, error)
	GetClientByID(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID) (*service.Client, error)
	GetAllClients(ctx context.Context, grantWriterID uuid.UUID) ([]service.Client, error)
	UpdateClient(ctx context.Context, input service.UpdateClientInput) (*service.Client, error)
	DeleteClient(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID) error
	GetAllTopics(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID) ([]service.Topic, error)
	AddTopic(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID, topicID uuid.UUID) error
	DeleteTopicFromClient(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID, topicID uuid.UUID) error
}

type ClientHandler struct {
	clientService ClientServicer
}

func NewClientHandler(clientService ClientServicer) *ClientHandler {
	return &ClientHandler{clientService: clientService}
}

type createClientRequest struct {
	Name         string `json:"name"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email"`
	Notes        string `json:"notes"`
}

type updateClientRequet struct {
	Name         *string `json:"name"`
	ContactName  *string `json:"contact_name"`
	ContactPhone *string `json:"contact_phone"`
	ContactEmail *string `json:"contact_email"`
	Notes        *string `json:"notes"`
}

type clientResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email"`
	Notes        string `json:"notes"`
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var req createClientRequest
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

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "client name is required")
		return
	}

	clientInput := service.CreateClientInput{
		GrantWriterID: grantWriterID,
		Name:          &req.Name,
		ContactPhone:  &req.ContactPhone,
		ContactName:   &req.ContactName,
		ContactEmail:  &req.ContactEmail,
		Notes:         &req.Notes,
	}

	result, err := h.clientService.CreateClient(r.Context(), clientInput)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create client")
		return
	}

	newClientResponse := toClientResponse(*result)

	writeJSON(w, http.StatusCreated, newClientResponse)
}

func (h *ClientHandler) GetClientByID(w http.ResponseWriter, r *http.Request) {
	clientIdString := chi.URLParam(r, "id")

	if clientIdString == "" {
		writeError(w, http.StatusBadRequest, "no client id provided")
		return
	}

	clientID, err := uuid.Parse(clientIdString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "unable to parse client id")
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

	result, err := h.clientService.GetClientByID(r.Context(), grantWriterID, clientID)
	if err != nil {
		if errors.Is(err, service.ErrClientNotFound) {
			writeError(w, http.StatusNotFound, "client not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to retrieve client")
		return
	}

	newClientResponse := toClientResponse(*result)

	writeJSON(w, http.StatusOK, newClientResponse)
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	var req updateClientRequet
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	clientIdString := chi.URLParam(r, "id")

	if clientIdString == "" {
		writeError(w, http.StatusBadRequest, "no client id provided")
		return
	}

	clientID, err := uuid.Parse(clientIdString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "unable to parse client id")
		return
	}

	if req.Name != nil && *req.Name == "" {
		writeError(w, http.StatusBadRequest, "client name cannot be empty")
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

	updateClientInput := service.UpdateClientInput{
		GrantWriterID: grantWriterID,
		ID:            clientID,
		Name:          req.Name,
		ContactPhone:  req.ContactPhone,
		ContactName:   req.ContactName,
		ContactEmail:  req.ContactEmail,
		Notes:         req.Notes,
	}

	result, err := h.clientService.UpdateClient(r.Context(), updateClientInput)

	if err != nil {
		if errors.Is(err, service.ErrClientNotFound) {
			writeError(w, http.StatusNotFound, "client not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update client")
		return
	}

	newClientResponse := toClientResponse(*result)

	writeJSON(w, http.StatusOK, newClientResponse)
}

func (h *ClientHandler) GetAllClients(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.clientService.GetAllClients(r.Context(), grantWriterID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve all clients")
		return
	}

	clients := make([]clientResponse, len(result))

	for i, record := range result {
		clients[i] = *toClientResponse(record)
	}

	writeJSON(w, http.StatusOK, clients)
}

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	clientIdString := chi.URLParam(r, "id")

	if clientIdString == "" {
		writeError(w, http.StatusBadRequest, "no client id provided")
		return
	}

	clientID, err := uuid.Parse(clientIdString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "unable to parse client id")
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

	err = h.clientService.DeleteClient(r.Context(), grantWriterID, clientID)
	if err != nil {
		if errors.Is(err, service.ErrClientNotFound) {
			writeError(w, http.StatusNotFound, "client not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete client")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ClientHandler) GetAllTopicsByClient(w http.ResponseWriter, r *http.Request) {
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

	results, err := h.clientService.GetAllTopics(r.Context(), grantWriterID, clientID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve topics")
		return
	}

	topics := make([]topicResponse, len(results))
	for i, t := range results {
		topics[i] = *toTopicResponse(&t)
	}

	writeJSON(w, http.StatusOK, topics)
}

func (h *ClientHandler) AddTopicToClient(w http.ResponseWriter, r *http.Request) {
	var req addTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

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

	topicID, err := uuid.Parse(req.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid topic id")
		return
	}

	err = h.clientService.AddTopic(r.Context(), grantWriterID, clientID, topicID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add topic to client")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ClientHandler) DeleteTopicFromClient(w http.ResponseWriter, r *http.Request) {
	clientIDString := chi.URLParam(r, "clientID")
	if clientIDString == "" {
		writeError(w, http.StatusBadRequest, "no client id provided")
		return
	}

	clientID, err := uuid.Parse(clientIDString)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid client id")
		return
	}

	topicIDString := chi.URLParam(r, "topicID")
	if topicIDString == "" {
		writeError(w, http.StatusBadRequest, "no client id provided")
		return
	}

	topicID, err := uuid.Parse(topicIDString)
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

	err = h.clientService.DeleteTopicFromClient(r.Context(), grantWriterID, clientID, topicID)
	if err != nil {
		if errors.Is(err, service.ErrClientNotFound) {
			writeError(w, http.StatusNotFound, "client not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete topic")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toClientResponse(serviceClient service.Client) *clientResponse {
	return &clientResponse{
		ID:           serviceClient.ID.String(),
		Name:         serviceClient.Name,
		ContactName:  serviceClient.ContactName,
		ContactPhone: serviceClient.ContactPhone,
		ContactEmail: serviceClient.ContactEmail,
		Notes:        serviceClient.Notes,
	}
}
