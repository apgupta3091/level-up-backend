package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	dbgen "github.com/anujgupta/level-up-backend/generated/db"
	"github.com/anujgupta/level-up-backend/internal/middleware"
)

type LessonsHandler struct {
	queries *dbgen.Queries
}

func NewLessonsHandler(q *dbgen.Queries) *LessonsHandler {
	return &LessonsHandler{queries: q}
}

func (h *LessonsHandler) GetLesson(w http.ResponseWriter, r *http.Request) {
	moduleSlug := chi.URLParam(r, "slug")
	lessonSlug := chi.URLParam(r, "lessonSlug")

	module, err := h.queries.GetModuleBySlug(r.Context(), moduleSlug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "module not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get module")
		return
	}

	lesson, err := h.queries.GetLessonBySlug(r.Context(), dbgen.GetLessonBySlugParams{
		ModuleID: module.ID,
		Slug:     lessonSlug,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "lesson not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get lesson")
		return
	}

	respondOK(w, map[string]any{
		"id":                lesson.ID,
		"module_id":         lesson.ModuleID,
		"title":             lesson.Title,
		"slug":              lesson.Slug,
		"content":           lesson.Content,
		"order_index":       lesson.OrderIndex,
		"estimated_minutes": lesson.EstimatedMinutes,
	})
}

func (h *LessonsHandler) CompleteLesson(w http.ResponseWriter, r *http.Request) {
	lessonIDStr := chi.URLParam(r, "id")
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	lessonID, err := parseUUID(lessonIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid lesson id")
		return
	}

	// Verify lesson exists
	if _, err := h.queries.GetLessonByID(r.Context(), lessonID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "lesson not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to verify lesson")
		return
	}

	if err := h.queries.MarkLessonComplete(r.Context(), dbgen.MarkLessonCompleteParams{
		UserID:   userID,
		LessonID: lessonID,
	}); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to mark lesson complete")
		return
	}

	respondOK(w, map[string]string{"status": "completed"})
}
