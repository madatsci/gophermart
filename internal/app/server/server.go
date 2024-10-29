package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/madatsci/gophermart/internal/app/config"
	"go.uber.org/zap"
)

type Server struct {
	mux    http.Handler
	config *config.Config
	log    *zap.SugaredLogger
}

func New(config *config.Config, logger *zap.SugaredLogger) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/", tmpHandler)
	})

	server := &Server{
		mux:    r,
		config: config,
		log:    logger,
	}

	return server
}

// Start starts the server under the specified address.
func (s *Server) Start() error {
	s.log.Infof("starting server with config: %+v", s.config)

	return http.ListenAndServe(s.config.RunAddress, s.mux)
}

func tmpHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I can't do anything yet..."))
	w.WriteHeader(http.StatusOK)
}
