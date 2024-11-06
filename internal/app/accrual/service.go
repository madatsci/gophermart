package accrual

import (
	"context"
	"fmt"
	"time"

	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/models"
	"github.com/madatsci/gophermart/internal/app/store"
	"github.com/madatsci/gophermart/pkg/accrual/client"
	"go.uber.org/zap"
)

type (
	AccrualService struct {
		client *client.Client
		store  store.Store
		logger *zap.SugaredLogger
	}

	AccrualProvider interface {
		GetOrder(number string) (client.OrderResponse, error)
	}
)

const syncOrdersLimit = 10

// New creates new accrual service.
func New(config *config.Config, store store.Store, logger *zap.SugaredLogger) *AccrualService {
	return &AccrualService{
		client: client.New(config),
		store:  store,
		logger: logger,
	}
}

// UpdateOrders fetches orders from accrual system and updates their status and accrual.
func (a *AccrualService) SyncOrders(ctx context.Context) error {
	orders, err := a.store.ListOrdersByStatus(ctx, []models.OrderStatus{models.OrderStatusNew, models.OrderStatusProcessing}, syncOrdersLimit)
	if err != nil {
		return err
	}

	for _, o := range orders {
		if ctx.Err() != nil {
			return nil
		}

		or, err := a.client.GetOrder(o.Number)
		if err != nil {
			a.logError(o.Number, err)
			continue
		}

		newStatus, err := mapOrderStatus(or.Status)
		if err != nil {
			a.logError(o.Number, err)
			continue
		}

		prevStatus := o.Status
		if newStatus != prevStatus {
			o.Status = newStatus
			o.Accrual = or.Accrual
			o.UpdatedAt = time.Now()

			_, err := a.store.UpdateOrder(ctx, o, prevStatus)
			if err != nil {
				a.logError(o.Number, err)
				continue
			}

			a.logger.With(
				"number", o.Number,
				"prev_status", prevStatus,
				"new_status", newStatus,
				"accrual", o.Accrual,
			).Info("updated order")
		}
	}

	return nil
}

func (a *AccrualService) logError(orderNumber string, err error) {
	a.logger.With("number", orderNumber, "err", err).Errorln("could not sync order")
}

func mapOrderStatus(accrualOrderStatus client.OrderStatus) (models.OrderStatus, error) {
	switch accrualOrderStatus {
	case client.OrderStatusRegistered:
		return models.OrderStatusNew, nil
	case client.OrderStatusProcessing:
		return models.OrderStatusProcessing, nil
	case client.OrderStatusProcessed:
		return models.OrderStatusProcessed, nil
	case client.OrderStatusInvalid:
		return models.OrderStatusInvalid, nil
	default:
		return "", fmt.Errorf("unknown order status received from accrual system: %s", accrualOrderStatus)
	}
}
