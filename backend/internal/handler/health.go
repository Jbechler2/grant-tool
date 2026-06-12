package handler

import (
	"database/sql"
	"net/http"
)

type HealthHandler struct {
	database *sql.DB
}

func NewHealthHandler(database *sql.DB) *HealthHandler {
	return &HealthHandler{database: database}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	err := h.database.Ping()

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}
