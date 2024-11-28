package handlers

import (
	"net/http"

	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/server/middleware"
	"github.com/madatsci/gophermart/internal/app/store"
	"github.com/madatsci/gophermart/pkg/jwt"
	"github.com/pkg/errors"
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

func New(opts Options) *Handlers {
	return &Handlers{c: opts.Config, s: opts.Store, jwt: opts.JWT, log: opts.Logger}
}

func ensureUserID(r *http.Request) (string, error) {
	userIDCtx := r.Context().Value(middleware.AuthenticatedUserKey)
	userID, ok := userIDCtx.(string)
	if !ok {
		return "", errors.New("authenticated user is required")
	}

	return userID, nil
}

func (h *Handlers) handleError(method string, err error) {
	h.log.With("method", method, "err", err).Errorln("error handling request")
}
