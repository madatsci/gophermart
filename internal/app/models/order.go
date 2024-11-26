package models

import (
	"time"
)

type (
	Order struct {
		ID        string      `bun:",pk,type:uuid" json:"-"`
		AccountID string      `bun:",notnull" json:"-"`
		Number    string      `bun:",unique,notnull" json:"number"`
		Status    OrderStatus `bun:",notnull" json:"status"`
		Accrual   float32     `bun:",nullzero" json:"accrual"`
		CreatedAt time.Time   `bun:",notnull,default:current_timestamp" json:"uploaded_at"`
		UpdatedAt time.Time   `bun:",notnull,default:current_timestamp" json:"-"`

		Account Account `bun:"rel:belongs-to,join:account_id=id" json:"-"`
	}

	OrderStatus string
)

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)
