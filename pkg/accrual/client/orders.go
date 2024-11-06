package client

import "github.com/shopspring/decimal"

// OrderResponse represents order object that is received from accrual service.
type OrderResponse struct {
	Order   string          `json:"order"`
	Status  OrderStatus     `json:"status"`
	Accrual decimal.Decimal `json:"accrual"`
}

// IsFinalStatus returns true if the order is in its final status.
func (o OrderResponse) IsFinalStatus() bool {
	return o.Status == OrderStatusProcessed || o.Status == OrderStatusInvalid
}

type OrderStatus string

const (
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

// GetOrder returns order status from accrual system.
func (c *Client) GetOrder(number string) (OrderResponse, error) {
	var res OrderResponse

	req := RequestOptions{
		Name:   "get_order",
		Path:   "/api/orders/" + number,
		Result: &res,
	}

	_, err := c.get(req)

	return res, err
}
