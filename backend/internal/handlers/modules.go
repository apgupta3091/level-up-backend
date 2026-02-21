package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	dbgen "github.com/anujgupta/level-up-backend/generated/db"
	"github.com/anujgupta/level-up-backend/internal/middleware"
)

type ModulesHandler struct {
	queries *dbgen.Queries
}

func NewModulesHandler(q *dbgen.Queries) *ModulesHandler {
	return &ModulesHandler{queries: q}
}

func (h *ModulesHandler) ListModules(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())

	modules, err := h.queries.ListModules(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list modules")
		return
	}

	completedIDs, err := h.queries.GetCompletedLessonIDs(r.Context(), userID)
	if err != nil {
		completedIDs = nil
	}
	completedSet := make(map[string]bool, len(completedIDs))
	for _, id := range completedIDs {
		completedSet[id.String()] = true
	}

	type moduleItem struct {
		ID             string  `json:"id"`
		Title          string  `json:"title"`
		Slug           string  `json:"slug"`
		Description    string  `json:"description"`
		OrderIndex     int32   `json:"order_index"`
		EstimatedHours float64 `json:"estimated_hours"`
	}

	result := make([]moduleItem, len(modules))
	for i, m := range modules {
		hours, _ := m.EstimatedHours.Float64Value()
		result[i] = moduleItem{
			ID:             m.ID.String(),
			Title:          m.Title,
			Slug:           m.Slug,
			Description:    m.Description,
			OrderIndex:     m.OrderIndex,
			EstimatedHours: hours.Float64,
		}
	}

	respondOK(w, map[string]any{"modules": result})
}

func (h *ModulesHandler) GetModule(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	userID, _ := middleware.GetUserID(r.Context())

	module, err := h.queries.GetModuleBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "module not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get module")
		return
	}

	lessons, err := h.queries.GetLessonsByModule(r.Context(), module.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get lessons")
		return
	}

	completedCount, err := h.queries.GetCompletedLessonCountByModule(r.Context(), dbgen.GetCompletedLessonCountByModuleParams{
		UserID:   userID,
		ModuleID: module.ID,
	})
	if err != nil {
		completedCount = 0
	}

	type lessonItem struct {
		ID               string `json:"id"`
		Title            string `json:"title"`
		Slug             string `json:"slug"`
		OrderIndex       int32  `json:"order_index"`
		EstimatedMinutes int32  `json:"estimated_minutes"`
	}

	lessonList := make([]lessonItem, len(lessons))
	for i, l := range lessons {
		lessonList[i] = lessonItem{
			ID:               l.ID.String(),
			Title:            l.Title,
			Slug:             l.Slug,
			OrderIndex:       l.OrderIndex,
			EstimatedMinutes: l.EstimatedMinutes,
		}
	}

	hours, _ := module.EstimatedHours.Float64Value()
	respondOK(w, map[string]any{
		"id":               module.ID,
		"title":            module.Title,
		"slug":             module.Slug,
		"description":      module.Description,
		"order_index":      module.OrderIndex,
		"estimated_hours":  hours.Float64,
		"total_lessons":    len(lessons),
		"completed_lessons": completedCount,
		"lessons":          lessonList,
	})
}
