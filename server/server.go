package server

import (
	"log/slog"
	"net/http"
	"path"

	"github.com/alam0rt/tuna/vtuner"
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
	_ = config
	mux.Handle("/setupapp/", handleSetupApp(logger))
	mux.Handle("/", handleLandingPage(logger))
}

func handleSetupApp(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling setupapp request", "path", r.URL.Path, "method", r.Method, "remote_addr", r.RemoteAddr)

		if !r.URL.Query().Has("token") {
			_, err := w.Write(vtuner.EncryptedToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		switch _, file := path.Split(r.URL.Path); file {
		case "loginXML.asp":
			handleLandingPage(logger)
		}

		http.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

func handleLandingPage(logger *slog.Logger) http.Handler {
	radiobrowser := &vtuner.Directory{
		Title:          "Radiobrowser",
		DestinationURL: "http://localhost:8080/radiobrowser",
		Count:          4,
	}

	notImplemented := &vtuner.Display{
		Display: "Not implemented",
	}

	page := vtuner.NewPage([]vtuner.Item{
		radiobrowser,
		notImplemented,
	}, false)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling request", "path", r.URL.Path, "method", r.Method, "remote_addr", r.RemoteAddr)

		w.Header().Add("Content-Type", "application/xml")
		if err := page.Write(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
