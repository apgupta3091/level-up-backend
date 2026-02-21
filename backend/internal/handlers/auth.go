package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/anujgupta/level-up-backend/internal/auth"
	dbgen "github.com/anujgupta/level-up-backend/generated/db"
	"github.com/anujgupta/level-up-backend/internal/mailer"
)

type AuthHandler struct {
	queries *dbgen.Queries
	auth    *auth.Service
	mailer  *mailer.Mailer
}

func NewAuthHandler(q *dbgen.Queries, a *auth.Service, m *mailer.Mailer) *AuthHandler {
	return &AuthHandler{queries: q, auth: a, mailer: m}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	if req.Email == "" || req.Password == "" || req.Name == "" {
		respondError(w, http.StatusBadRequest, "email, password, and name are required")
		return
	}
	if len(req.Password) < 8 {
		respondError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user, err := h.queries.CreateUser(r.Context(), dbgen.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
	})
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			respondError(w, http.StatusConflict, "email already registered")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	tokens, err := h.auth.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	h.mailer.Send(mailer.EmailJob{
		To:       user.Email,
		Subject:  "Welcome to Level Up Backend",
		Template: "welcome",
		Data:     map[string]string{"name": user.Name},
	})

	respondCreated(w, map[string]any{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"user": map[string]any{
			"id":                  user.ID,
			"email":               user.Email,
			"name":                user.Name,
			"subscription_status": user.SubscriptionStatus,
		},
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, err := h.queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	tokens, err := h.auth.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	respondOK(w, map[string]any{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"user": map[string]any{
			"id":                  user.ID,
			"email":               user.Email,
			"name":                user.Name,
			"subscription_status": user.SubscriptionStatus,
		},
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claims, err := h.auth.ParseToken(req.RefreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}
	if claims.TokenType != auth.TokenTypeRefresh {
		respondError(w, http.StatusUnauthorized, "refresh token required")
		return
	}

	// Fetch current role in case it changed
	user, err := h.queries.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "user not found")
		return
	}

	tokens, err := h.auth.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	respondOK(w, map[string]string{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}
