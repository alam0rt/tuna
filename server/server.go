package server

import (
	"log/slog"
	"net/http"
)

type Config struct {
	// Config fields go here
	Host string
	Port string
}

// NewServer creates a new HTTP server and returns it.
func NewServer(
	logger *slog.Logger,
	config *Config,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		logger,
		config,
	)
	var handler http.Handler = mux
	return handler
}

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	config *Config,
) {
	_ = logger
	_ = config
	mux.Handle("/", http.NotFoundHandler())
}
