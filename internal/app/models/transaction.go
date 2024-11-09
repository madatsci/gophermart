package models

import (
	"time"
)

type (
	Transaction struct {
		ID          string      `bun:",pk,type:uuid" json:"-"`
		AccountID   string      `bun:",notnull" json:"-"`
		Amount      float32     `bun:",notnull" json:"sum"`
		OrderNumber string      `bun:",notnull" json:"order"`
		Direction   TxDirection `bun:",notnull" json:"-"`
		CreatedAt   time.Time   `bun:",notnull,default:current_timestamp" json:"processed_at"`
		UpdatedAt   time.Time   `bun:",notnull,default:current_timestamp" json:"-"`

		Account Account `bun:"rel:belongs-to,join:account_id=id" json:"-"`
	}

	TxDirection string
)

const (
	TxDirectionAccrual    TxDirection = "accrual"
	TxDirectionWithdrawal TxDirection = "withdrawal"
)
