package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	dbgen "github.com/anujgupta/level-up-backend/generated/db"
)

type AdminHandler struct {
	queries *dbgen.Queries
}

func NewAdminHandler(q *dbgen.Queries) *AdminHandler {
	return &AdminHandler{queries: q}
}

func (h *AdminHandler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
	submissions, err := h.queries.ListPendingSubmissions(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list submissions")
		return
	}

	respondOK(w, map[string]any{"submissions": submissions})
}

type reviewSubmissionRequest struct {
	Status   string `json:"status"`
	Feedback string `json:"feedback"`
}

func (h *AdminHandler) ReviewSubmission(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := parseUUID(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid submission id")
		return
	}

	var req reviewSubmissionRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	status := dbgen.SubmissionStatus(req.Status)
	if !status.Valid() {
		respondError(w, http.StatusBadRequest, "invalid status: must be reviewed, approved, or needs_revision")
		return
	}

	submission, err := h.queries.ReviewSubmission(r.Context(), dbgen.ReviewSubmissionParams{
		ID:       id,
		Status:   status,
		Feedback: &req.Feedback,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "submission not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to review submission")
		return
	}

	respondOK(w, submission)
}
