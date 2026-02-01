package handler

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *TodoHandler) {

	mux.HandleFunc("GET /todos", h.List)
	mux.HandleFunc("POST /todos", h.Create)

	mux.HandleFunc("GET /todos/{id}", h.Get)
	mux.HandleFunc("DELETE /todos/{id}", h.Delete)

	mux.HandleFunc("PATCH /todos/{id}/title", h.UpdateTitle)
	mux.HandleFunc("PATCH /todos/{id}/completed", h.UpdateCompleted)
}
