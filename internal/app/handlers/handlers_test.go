package handlers

import (
	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/store/database/mocks"
	"github.com/madatsci/gophermart/pkg/jwt"
	"go.uber.org/zap"
)

func newTestHandlers(m *mocks.MockStore) *Handlers {
	return &Handlers{
		s: m,
		c: &config.Config{
			AuthCookieName: "auth_token",
		},
		jwt: &jwt.JWT{},
		log: zap.NewNop().Sugar(),
	}
}
