package accrual

import (
	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/pkg/accrual/client"
	"go.uber.org/zap"
)

type (
	AccrualService struct {
		client *client.Client
		logger *zap.SugaredLogger
		config *config.Config
	}

	AccrualProvider interface {
		GetOrder(number string) (client.OrderResponse, error)
	}
)

// New creates new accrual service.
func New(config *config.Config, logger *zap.SugaredLogger) *AccrualService {
	return &AccrualService{
		client: client.New(config),
		logger: logger,
		config: config,
	}
}
