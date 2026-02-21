package handlers

import (
	"net/http"

	dbgen "github.com/anujgupta/level-up-backend/generated/db"
	"github.com/anujgupta/level-up-backend/internal/middleware"
)

type ProgressHandler struct {
	queries *dbgen.Queries
}

func NewProgressHandler(q *dbgen.Queries) *ProgressHandler {
	return &ProgressHandler{queries: q}
}

func (h *ProgressHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	completedLessonIDs, err := h.queries.GetCompletedLessonIDs(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get progress")
		return
	}

	completedSkillIDs, err := h.queries.GetCompletedSkillIDs(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get skill progress")
		return
	}

	lessonIDs := make([]string, len(completedLessonIDs))
	for i, id := range completedLessonIDs {
		lessonIDs[i] = id.String()
	}

	skillIDs := make([]string, len(completedSkillIDs))
	for i, id := range completedSkillIDs {
		skillIDs[i] = id.String()
	}

	respondOK(w, map[string]any{
		"completed_lesson_ids": lessonIDs,
		"completed_skill_ids":  skillIDs,
	})
}
