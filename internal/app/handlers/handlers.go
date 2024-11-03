package handlers

import (
	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/store"
	"go.uber.org/zap"
)

type Handlers struct {
	s   store.Store
	c   *config.Config
	log *zap.SugaredLogger
}

// New creates new Handlers.
func New(config *config.Config, logger *zap.SugaredLogger, store store.Store) *Handlers {
	return &Handlers{c: config, s: store, log: logger}
}

func (h *Handlers) handleError(method string, err error) {
	h.log.With("method", method, "err", err).Errorln("error handling request")
}
