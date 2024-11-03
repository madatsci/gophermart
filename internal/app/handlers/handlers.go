package handlers

import (
	"net/http"

	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/store"
	"github.com/madatsci/gophermart/pkg/jwt"
	"go.uber.org/zap"
)

type (
	Handlers struct {
		s   store.Store
		c   *config.Config
		jwt *jwt.JWT
		log *zap.SugaredLogger
	}

	Options struct {
		Store  store.Store
		Config *config.Config
		JWT    *jwt.JWT
		Logger *zap.SugaredLogger
	}
)

// New creates new Handlers.
func New(opts Options) *Handlers {
	return &Handlers{c: opts.Config, s: opts.Store, jwt: opts.JWT, log: opts.Logger}
}

// TODO delete
func (h *Handlers) PrivateAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) handleError(method string, err error) {
	h.log.With("method", method, "err", err).Errorln("error handling request")
}
