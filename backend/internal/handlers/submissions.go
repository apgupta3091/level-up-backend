package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	dbgen "github.com/anujgupta/level-up-backend/generated/db"
	"github.com/anujgupta/level-up-backend/internal/middleware"
)

type SubmissionsHandler struct {
	queries *dbgen.Queries
}

func NewSubmissionsHandler(q *dbgen.Queries) *SubmissionsHandler {
	return &SubmissionsHandler{queries: q}
}

func (h *SubmissionsHandler) GetAssignment(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	module, err := h.queries.GetModuleBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "module not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get module")
		return
	}

	assignment, err := h.queries.GetAssignmentByModuleID(r.Context(), module.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "assignment not found for this module")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get assignment")
		return
	}

	hours, _ := assignment.EstimatedHours.Float64Value()
	respondOK(w, map[string]any{
		"id":              assignment.ID,
		"module_id":       assignment.ModuleID,
		"title":           assignment.Title,
		"description":     assignment.Description,
		"rubric":          assignment.Rubric,
		"estimated_hours": hours.Float64,
	})
}

type createSubmissionRequest struct {
	AssignmentID   string `json:"assignment_id"`
	GithubURL      string `json:"github_url"`
	WrittenAnswers string `json:"written_answers"`
}

func (h *SubmissionsHandler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req createSubmissionRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.AssignmentID == "" || req.GithubURL == "" {
		respondError(w, http.StatusBadRequest, "assignment_id and github_url are required")
		return
	}

	assignmentID, err := parseUUID(req.AssignmentID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid assignment_id")
		return
	}

	if _, err := h.queries.GetAssignmentByID(r.Context(), assignmentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "assignment not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to verify assignment")
		return
	}

	submission, err := h.queries.CreateSubmission(r.Context(), dbgen.CreateSubmissionParams{
		AssignmentID:   assignmentID,
		UserID:         userID,
		GithubUrl:      strings.TrimSpace(req.GithubURL),
		WrittenAnswers: req.WrittenAnswers,
	})
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			respondError(w, http.StatusConflict, "submission already exists for this assignment")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to create submission")
		return
	}

	respondCreated(w, submission)
}

func (h *SubmissionsHandler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	submissions, err := h.queries.GetSubmissionsByUser(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list submissions")
		return
	}

	respondOK(w, map[string]any{"submissions": submissions})
}

func (h *SubmissionsHandler) GetSubmission(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := parseUUID(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid submission id")
		return
	}

	submission, err := h.queries.GetSubmissionByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "submission not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get submission")
		return
	}

	role, _ := middleware.GetRole(r.Context())
	if submission.UserID != userID && role != "admin" {
		respondError(w, http.StatusForbidden, "access denied")
		return
	}

	respondOK(w, submission)
}
