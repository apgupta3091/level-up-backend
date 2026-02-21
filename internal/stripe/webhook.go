package stripehandler

import (
	"context"
	"log/slog"

	stripe "github.com/stripe/stripe-go/v82"

	dbgen "github.com/anujgupta/level-up-backend/generated/db"
	"github.com/anujgupta/level-up-backend/internal/mailer"
)

type WebhookRouter struct {
	queries *dbgen.Queries
	mailer  *mailer.Mailer
	logger  *slog.Logger
}

func NewWebhookRouter(q *dbgen.Queries, m *mailer.Mailer, l *slog.Logger) *WebhookRouter {
	return &WebhookRouter{queries: q, mailer: m, logger: l}
}

func (wr *WebhookRouter) Route(ctx context.Context, event stripe.Event) {
	switch event.Type {
	case "checkout.session.completed":
		wr.handleCheckoutCompleted(ctx, event)
	case "customer.subscription.created", "customer.subscription.updated":
		wr.handleSubscriptionActive(ctx, event)
	case "customer.subscription.deleted":
		wr.handleSubscriptionDeleted(ctx, event)
	case "invoice.payment_failed":
		wr.handlePaymentFailed(ctx, event)
	case "invoice.payment_succeeded":
		wr.handlePaymentSucceeded(ctx, event)
	default:
		wr.logger.Info("unhandled stripe event", "type", event.Type)
	}
}

func (wr *WebhookRouter) handleCheckoutCompleted(ctx context.Context, event stripe.Event) {
	session, ok := event.Data.Object["object"].(*stripe.CheckoutSession)
	if !ok {
		// Parse manually from raw data
		customerID, _ := event.Data.Object["customer"].(string)
		clientRefID, _ := event.Data.Object["client_reference_id"].(string)
		if customerID == "" || clientRefID == "" {
			wr.logger.Error("checkout.session.completed: missing customer or client_reference_id")
			return
		}
		session = &stripe.CheckoutSession{
			Customer:           &stripe.Customer{ID: customerID},
			ClientReferenceID:  clientRefID,
		}
	}
	_ = session

	customerID, _ := event.Data.Object["customer"].(string)
	clientRefID, _ := event.Data.Object["client_reference_id"].(string)

	if customerID == "" || clientRefID == "" {
		wr.logger.Error("checkout.session.completed: missing required fields")
		return
	}

	userID, err := parseUUID(clientRefID)
	if err != nil {
		wr.logger.Error("checkout.session.completed: invalid client_reference_id", "err", err)
		return
	}

	if _, err := wr.queries.UpdateUserStripeCustomerID(ctx, dbgen.UpdateUserStripeCustomerIDParams{
		ID:               userID,
		StripeCustomerID: &customerID,
	}); err != nil {
		wr.logger.Error("checkout.session.completed: failed to update stripe customer id", "err", err)
	}
}

func (wr *WebhookRouter) handleSubscriptionActive(ctx context.Context, event stripe.Event) {
	subID, _ := event.Data.Object["id"].(string)
	customerObj, _ := event.Data.Object["customer"].(string)

	if customerObj == "" || subID == "" {
		wr.logger.Error("subscription event: missing customer or subscription id")
		return
	}

	if _, err := wr.queries.UpdateUserSubscription(ctx, dbgen.UpdateUserSubscriptionParams{
		StripeCustomerID:     &customerObj,
		StripeSubscriptionID: &subID,
		SubscriptionStatus:   dbgen.SubscriptionStatusActive,
	}); err != nil {
		wr.logger.Error("subscription active: failed to update user", "err", err)
	}
}

func (wr *WebhookRouter) handleSubscriptionDeleted(ctx context.Context, event stripe.Event) {
	subID, _ := event.Data.Object["id"].(string)
	customerID, _ := event.Data.Object["customer"].(string)

	if customerID == "" {
		wr.logger.Error("subscription.deleted: missing customer id")
		return
	}

	if _, err := wr.queries.UpdateUserSubscription(ctx, dbgen.UpdateUserSubscriptionParams{
		StripeCustomerID:     &customerID,
		StripeSubscriptionID: &subID,
		SubscriptionStatus:   dbgen.SubscriptionStatusCancelled,
	}); err != nil {
		wr.logger.Error("subscription.deleted: failed to update user", "err", err)
	}
}

func (wr *WebhookRouter) handlePaymentFailed(ctx context.Context, event stripe.Event) {
	customerID, _ := event.Data.Object["customer"].(string)
	if customerID == "" {
		return
	}

	user, err := wr.queries.GetUserByStripeCustomerID(ctx, &customerID)
	if err != nil {
		wr.logger.Error("invoice.payment_failed: user not found", "customer", customerID)
		return
	}

	subID := user.StripeSubscriptionID
	if _, err := wr.queries.UpdateUserSubscription(ctx, dbgen.UpdateUserSubscriptionParams{
		StripeCustomerID:     &customerID,
		StripeSubscriptionID: subID,
		SubscriptionStatus:   dbgen.SubscriptionStatusPastDue,
	}); err != nil {
		wr.logger.Error("invoice.payment_failed: failed to update user", "err", err)
		return
	}

	wr.mailer.Send(mailer.EmailJob{
		To:       user.Email,
		Subject:  "Payment failed — update your billing info",
		Template: "payment_failed",
		Data:     map[string]string{"name": user.Name},
	})
}

func (wr *WebhookRouter) handlePaymentSucceeded(ctx context.Context, event stripe.Event) {
	customerID, _ := event.Data.Object["customer"].(string)
	if customerID == "" {
		return
	}

	user, err := wr.queries.GetUserByStripeCustomerID(ctx, &customerID)
	if err != nil {
		return
	}

	subID := user.StripeSubscriptionID
	if _, err := wr.queries.UpdateUserSubscription(ctx, dbgen.UpdateUserSubscriptionParams{
		StripeCustomerID:     &customerID,
		StripeSubscriptionID: subID,
		SubscriptionStatus:   dbgen.SubscriptionStatusActive,
	}); err != nil {
		wr.logger.Error("invoice.payment_succeeded: failed to update user", "err", err)
	}
}
