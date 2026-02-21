package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	stripe "github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/webhook"

	dbgen "github.com/anujgupta/level-up-backend/generated/db"
	"github.com/anujgupta/level-up-backend/internal/config"
	stripehandler "github.com/anujgupta/level-up-backend/internal/stripe"
	"github.com/anujgupta/level-up-backend/internal/middleware"
	"github.com/anujgupta/level-up-backend/internal/mailer"
)

type PaymentsHandler struct {
	queries       *dbgen.Queries
	cfg           *config.Config
	webhookRouter *stripehandler.WebhookRouter
	logger        *slog.Logger
}

func NewPaymentsHandler(q *dbgen.Queries, cfg *config.Config, m *mailer.Mailer, logger *slog.Logger) *PaymentsHandler {
	stripe.Key = cfg.StripeSecretKey
	wr := stripehandler.NewWebhookRouter(q, m, logger)
	return &PaymentsHandler{queries: q, cfg: cfg, webhookRouter: wr, logger: logger}
}

func (h *PaymentsHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	user, err := h.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if user.SubscriptionStatus == dbgen.SubscriptionStatusActive {
		respondError(w, http.StatusConflict, "user already has an active subscription")
		return
	}

	successURL := fmt.Sprintf("http://localhost:3000/dashboard?checkout=success")
	cancelURL := fmt.Sprintf("http://localhost:3000/pricing?checkout=cancelled")

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(h.cfg.StripePriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL:        stripe.String(successURL),
		CancelURL:         stripe.String(cancelURL),
		ClientReferenceID: stripe.String(userID.String()),
		CustomerEmail:     stripe.String(user.Email),
	}

	if user.StripeCustomerID != nil {
		params.Customer = user.StripeCustomerID
		params.CustomerEmail = nil
	}

	sess, err := session.New(params)
	if err != nil {
		h.logger.Error("failed to create stripe checkout session", "err", err)
		respondError(w, http.StatusInternalServerError, "failed to create checkout session")
		return
	}

	respondOK(w, map[string]string{"url": sess.URL})
}

func (h *PaymentsHandler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	rawBody, ok := middleware.GetStripeRawBody(r.Context())
	if !ok {
		respondError(w, http.StatusBadRequest, "missing request body")
		return
	}

	sigHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(rawBody, sigHeader, h.cfg.StripeWebhookSecret)
	if err != nil {
		h.logger.Warn("stripe webhook signature verification failed", "err", err)
		respondError(w, http.StatusBadRequest, "invalid stripe signature")
		return
	}

	h.webhookRouter.Route(r.Context(), event)

	w.WriteHeader(http.StatusOK)
}

func (h *PaymentsHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	user, err := h.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	respondOK(w, map[string]any{
		"subscription_status":   user.SubscriptionStatus,
		"stripe_subscription_id": user.StripeSubscriptionID,
	})
}
