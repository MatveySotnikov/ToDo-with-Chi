package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repo
}

const ( //константы валидации длины title
	minTitleLen = 3
	maxTitleLen = 100
)

const (
	defaultPage  = 1
	defaultLimit = 10
	maxLimit     = 50
)

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)          // GET /tasks
	r.Post("/", h.create)       // POST /tasks
	r.Get("/{id}", h.get)       // GET /tasks/{id}
	r.Put("/{id}", h.update)    // PUT /tasks/{id}
	r.Delete("/{id}", h.delete) // DELETE /tasks/{id}
	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// 1. Парсинг пагинации (page)
	page := defaultPage
	if pageStr := q.Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// 2. Парсинг пагинации (limit)
	limit := defaultLimit
	if limitStr := q.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			if l > maxLimit {
				limit = maxLimit
			} else {
				limit = l
			}
		}
	}

	// 3. Парсинг фильтра (done)
	var doneFilter *bool // nil = не фильтровать
	if doneStr := q.Get("done"); doneStr != "" {
		if d, err := strconv.ParseBool(doneStr); err == nil {
			doneFilter = &d // Устанавливаем указатель на true или false
		} else {
			// Ошибка, если 'done' не 'true' и не 'false'
			httpError(w, http.StatusBadRequest, "invalid 'done' filter: must be 'true' or 'false'")
			return
		}
	}

	// 4. Вызов репозитория с новыми параметрами
	tasks := h.repo.List(page, limit, doneFilter)
	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	t, err := h.repo.Get(id)
	if err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

type createReq struct {
	Title string `json:"title"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		httpError(w, http.StatusBadRequest, "invalid json: require non-empty title")
		return
	}
	titleLen := len(req.Title)
	if titleLen < minTitleLen || titleLen > maxTitleLen {
		httpError(w, http.StatusBadRequest, "title length must be between 3 and 100 characters")
		return
	}
	t := h.repo.Create(req.Title)
	writeJSON(w, http.StatusCreated, t)
}

type updateReq struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	var req updateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		httpError(w, http.StatusBadRequest, "invalid json: require non-empty title")
		return
	}

	titleLen := len(req.Title)
	if titleLen < minTitleLen || titleLen > maxTitleLen {
		httpError(w, http.StatusBadRequest, "title length must be between 3 and 100 characters")
		return
	}

	t, err := h.repo.Update(id, req.Title, req.Done)
	if err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	if err := h.repo.Delete(id); err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// helpers

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	raw := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		httpError(w, http.StatusBadRequest, "invalid id")
		return 0, true
	}
	return id, false
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
