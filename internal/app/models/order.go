package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	Order struct {
		ID        string          `bun:",pk,type:uuid"`
		AccountID string          `bun:",notnull"`
		Number    string          `bun:",unique,notnull"`
		Status    OrderStatus     `bun:",notnull"`
		Accrual   decimal.Decimal `bun:",nullzero"`
		CreatedAt time.Time       `bun:",notnull,default:current_timestamp"`
		UpdatedAt time.Time       `bun:",notnull,default:current_timestamp"`

		Account Account `bun:"rel:belongs-to,join:account_id=id"`
	}

	OrderStatus string
)

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)
