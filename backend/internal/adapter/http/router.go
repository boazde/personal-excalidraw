package http

import (
	"log/slog"
	"net/http"

	"github.com/personal-excalidraw/backend/internal/adapter/http/handler"
	"github.com/personal-excalidraw/backend/internal/adapter/http/middleware"
	"github.com/personal-excalidraw/backend/internal/infrastructure/config"
)

// NewRouter creates a new HTTP router with all routes and middleware
func NewRouter(
	cfg *config.Config,
	healthHandler *handler.HealthHandler,
	drawingHandler *handler.DrawingHandler,
	authHandler *handler.AuthHandler,
	logger *slog.Logger,
) http.Handler {
	// Create new ServeMux with Go 1.22+ routing
	mux := http.NewServeMux()

	// Health check endpoint (public)
	mux.HandleFunc("GET /health", healthHandler.Check)

	// Auth validation endpoint (protected by auth middleware)
	mux.HandleFunc("GET /auth/validate", authHandler.Validate)

	// Drawing API endpoints (nginx strips /api prefix)
	mux.HandleFunc("POST /drawings", drawingHandler.CreateDrawing)
	mux.HandleFunc("GET /drawings/{id}", drawingHandler.GetDrawing)
	mux.HandleFunc("GET /drawings", drawingHandler.ListDrawings)
	mux.HandleFunc("PUT /drawings/{id}", drawingHandler.UpdateDrawing)
	mux.HandleFunc("DELETE /drawings/{id}", drawingHandler.DeleteDrawing)

	// Share token endpoints (authenticated)
	mux.HandleFunc("POST /drawings/{id}/share", drawingHandler.GenerateShareToken)
	mux.HandleFunc("DELETE /drawings/{id}/share", drawingHandler.RevokeShareToken)

	// Public endpoints (no authentication required)
	mux.HandleFunc("GET /public/{token}", drawingHandler.GetPublicDrawing)

	// Apply middleware stack (in reverse order - outermost first)
	var handler http.Handler = mux
	handler = middleware.Auth(cfg, []string{"/health", "/public/{token}"})(handler)
	handler = middleware.CORS(cfg)(handler)
	handler = middleware.Logger(logger)(handler)
	handler = middleware.RequestID(handler)
	handler = middleware.Recover(logger)(handler)

	return handler
}
