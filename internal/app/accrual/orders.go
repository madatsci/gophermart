package accrual

import "github.com/madatsci/gophermart/pkg/accrual/client"

// GetOrder uses accrual system client to get order status.
func (a *AccrualService) GetOrder(number string) (client.OrderResponse, error) {
	return a.client.GetOrder(number)
}
