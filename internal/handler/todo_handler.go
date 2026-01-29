package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"myapp/internal/service"
)

type TodoHandler struct {
	svc service.TodoService
}

func NewTodoHandler(svc service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// JSONレスポンスを返す共通処理
func (h *TodoHandler) respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			// ログ出力などをここに入れると良い
			http.Error(w, "JSON encode failed", http.StatusInternalServerError)
		}
	}
}

// エラーレスポンスもJSONで返す
func (h *TodoHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

// リクエストボディのデコード共通化
func (h *TodoHandler) decode(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// ID取得の共通化
func (h *TodoHandler) getID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

// GET /todos
func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListTodos(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, items)
}

// POST /todos
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := h.decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// ここで if req.Title == "" { ... } のようなバリデーションもしやすい

	item, err := h.svc.CreateTodo(r.Context(), req.Title)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.respondJSON(w, http.StatusCreated, item)
}

// GET /todos/{id}
func (h *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	item, err := h.svc.GetTodo(r.Context(), id)
	if err != nil {
		// Not Foundかどうかを厳密に判定したい場合はServiceのエラー型を見るなどの工夫が必要
		h.respondError(w, http.StatusNotFound, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, item)
}

// DELETE /todos/{id}
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.svc.DeleteTodo(r.Context(), id); err != nil {
		h.respondError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PATCH /todos/{id}/title
func (h *TodoHandler) UpdateTitle(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req updateTitleRequest
	if err := h.decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.svc.UpdateTitle(r.Context(), id, req.Title)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, item)
}

// PATCH /todos/{id}/completed
func (h *TodoHandler) UpdateCompleted(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req updateCompletedRequest
	if err := h.decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.svc.UpdateCompleted(r.Context(), id, req.Completed)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, item)
}