package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/handlers"
	"github.com/madatsci/gophermart/internal/app/store"
	"github.com/madatsci/gophermart/pkg/jwt"
	"go.uber.org/zap"

	mw "github.com/madatsci/gophermart/internal/app/server/middleware"
)

type Server struct {
	mux    http.Handler
	config *config.Config
	h      *handlers.Handlers
	log    *zap.SugaredLogger
}

func New(config *config.Config, store store.Store, logger *zap.SugaredLogger) *Server {
	jwt := jwt.New(jwt.Options{
		Secret:   config.TokenSecret,
		Duration: config.TokenDuration,
		Issuer:   config.TokenIssuer,
	})

	h := handlers.New(handlers.Options{
		Store:  store,
		Config: config,
		JWT:    jwt,
		Logger: logger,
	})

	r := chi.NewRouter()
	loggerMiddleware := mw.NewLogger(logger)
	r.Use(loggerMiddleware.Logger)
	r.Use(middleware.Recoverer)

	authMiddleware := mw.NewAuth(mw.Options{
		Config: config,
		JWT:    jwt,
		Log:    logger,
	})

	r.Route("/", func(r chi.Router) {
		// Public API
		r.Post("/api/user/register", h.RegisterUser)
		r.Post("/api/user/login", h.LoginUser)

		// Private API
		r.Route("/api/user/orders", func(r chi.Router) {
			r.Use(authMiddleware.PrivateAPIAuth)
			r.Post("/", h.CreateOrder)
		})
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
