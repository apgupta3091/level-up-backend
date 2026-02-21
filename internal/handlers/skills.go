package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	dbgen "github.com/anujgupta/level-up-backend/generated/db"
	"github.com/anujgupta/level-up-backend/internal/middleware"
)

type SkillsHandler struct {
	queries *dbgen.Queries
}

func NewSkillsHandler(q *dbgen.Queries) *SkillsHandler {
	return &SkillsHandler{queries: q}
}

func (h *SkillsHandler) GetModuleSkills(w http.ResponseWriter, r *http.Request) {
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

	skills, err := h.queries.GetSkillsByModule(r.Context(), module.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get skills")
		return
	}

	type skillItem struct {
		ID         string `json:"id"`
		SkillName  string `json:"skill_name"`
		OrderIndex int32  `json:"order_index"`
	}

	result := make([]skillItem, len(skills))
	for i, s := range skills {
		result[i] = skillItem{
			ID:         s.ID.String(),
			SkillName:  s.SkillName,
			OrderIndex: s.OrderIndex,
		}
	}

	respondOK(w, map[string]any{"skills": result})
}

func (h *SkillsHandler) CompleteSkill(w http.ResponseWriter, r *http.Request) {
	skillIDStr := chi.URLParam(r, "id")
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	skillID, err := parseUUID(skillIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid skill id")
		return
	}

	if _, err := h.queries.GetSkillByID(r.Context(), skillID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "skill not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to verify skill")
		return
	}

	if err := h.queries.MarkSkillComplete(r.Context(), dbgen.MarkSkillCompleteParams{
		UserID:  userID,
		SkillID: skillID,
	}); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to mark skill complete")
		return
	}

	respondOK(w, map[string]string{"status": "completed"})
}
