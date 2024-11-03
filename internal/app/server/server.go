package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/handlers"
	"github.com/madatsci/gophermart/internal/app/store"
	"go.uber.org/zap"
)

type Server struct {
	mux    http.Handler
	config *config.Config
	h      *handlers.Handlers
	log    *zap.SugaredLogger
}

func New(config *config.Config, store store.Store, logger *zap.SugaredLogger) *Server {
	h := handlers.New(config, logger, store)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/api/user/register", h.RegisterUser)
	})

	server := &Server{
		mux:    r,
		config: config,
		h:      h,
		log:    logger,
	}

	return server
}

// Start starts the server under the specified address.
func (s *Server) Start() error {
	s.log.Infof("starting server with config: %+v", s.config)

	return http.ListenAndServe(s.config.RunAddress, s.mux)
}
