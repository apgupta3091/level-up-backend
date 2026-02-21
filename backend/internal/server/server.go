package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/httprate"

	dbgen "github.com/anujgupta/level-up-backend/generated/db"
	"github.com/anujgupta/level-up-backend/internal/auth"
	"github.com/anujgupta/level-up-backend/internal/config"
	"github.com/anujgupta/level-up-backend/internal/handlers"
	"github.com/anujgupta/level-up-backend/internal/mailer"
	appmiddleware "github.com/anujgupta/level-up-backend/internal/middleware"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(
	cfg *config.Config,
	queries *dbgen.Queries,
	authSvc *auth.Service,
	mailerSvc *mailer.Mailer,
	logger *slog.Logger,
) *Server {
	r := chi.NewRouter()

	// Build an httplog.Logger that wraps our slog.Logger
	httpLogger := &httplog.Logger{
		Logger: logger,
		Options: httplog.Options{
			Concise:  cfg.IsDevelopment(),
			LogLevel: slog.LevelInfo,
		},
	}

	// ── Global middleware ─────────────────────────────────────────────────────
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(httplog.RequestLogger(httpLogger))
	r.Use(chimiddleware.Recoverer)

	// ── Handlers ─────────────────────────────────────────────────────────────
	authHandler := handlers.NewAuthHandler(queries, authSvc, mailerSvc)
	paymentsHandler := handlers.NewPaymentsHandler(queries, cfg, mailerSvc, logger)
	modulesHandler := handlers.NewModulesHandler(queries)
	lessonsHandler := handlers.NewLessonsHandler(queries)
	progressHandler := handlers.NewProgressHandler(queries)
	skillsHandler := handlers.NewSkillsHandler(queries)
	submissionsHandler := handlers.NewSubmissionsHandler(queries)
	adminHandler := handlers.NewAdminHandler(queries)

	// ── Routes ───────────────────────────────────────────────────────────────

	r.Get("/health", handlers.Health)

	// Auth (rate limited)
	r.With(httprate.LimitByIP(10, 60)).Post("/auth/register", authHandler.Register)
	r.With(httprate.LimitByIP(10, 60)).Post("/auth/login", authHandler.Login)
	r.With(httprate.LimitByIP(20, 60)).Post("/auth/refresh", authHandler.RefreshToken)

	// Stripe webhook — raw body must be captured before any body parsing
	r.With(appmiddleware.StripeRawBody).Post("/payments/webhook", paymentsHandler.StripeWebhook)

	// JWT-protected routes
	authenticate := appmiddleware.Authenticate(authSvc)
	requireActive := appmiddleware.RequireActive(queries)

	r.Group(func(r chi.Router) {
		r.Use(authenticate)

		// Payments
		r.Post("/payments/checkout", paymentsHandler.CreateCheckoutSession)
		r.Get("/payments/subscription", paymentsHandler.GetSubscription)

		// Progress + submissions (JWT only, no sub gate)
		r.Get("/progress", progressHandler.GetProgress)
		r.Get("/submissions", submissionsHandler.ListSubmissions)
		r.Get("/submissions/{id}", submissionsHandler.GetSubmission)

		// Subscription-gated content
		r.Group(func(r chi.Router) {
			r.Use(requireActive)

			r.Get("/modules", modulesHandler.ListModules)
			r.Get("/modules/{slug}", modulesHandler.GetModule)
			r.Get("/modules/{slug}/lessons/{lessonSlug}", lessonsHandler.GetLesson)
			r.Post("/lessons/{id}/complete", lessonsHandler.CompleteLesson)
			r.Get("/modules/{slug}/skills", skillsHandler.GetModuleSkills)
			r.Post("/skills/{id}/complete", skillsHandler.CompleteSkill)
			r.Get("/modules/{slug}/assignment", submissionsHandler.GetAssignment)
			r.Post("/submissions", submissionsHandler.CreateSubmission)
		})

		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.RequireAdmin)

			r.Get("/admin/submissions", adminHandler.ListSubmissions)
			r.Put("/admin/submissions/{id}/review", adminHandler.ReviewSubmission)
		})
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("server starting", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) HTTPServer() *http.Server {
	return s.httpServer
}
