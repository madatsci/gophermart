package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID                 string          `bun:",pk,type:uuid"`
	UserID             string          `bun:",unique,notnull"`
	CurrentPointsTotal decimal.Decimal `bun:",notnull,default:0"`
	WithdrawnTotal     decimal.Decimal `bun:",notnull,default:0"`
	CreatedAt          time.Time       `bun:",notnull,default:current_timestamp"`
	UpdatedAt          time.Time       `bun:",notnull,default:current_timestamp"`

	User   User     `bun:"rel:belongs-to,join:user_id=id"`
	Orders []*Order `bun:"rel:has-many,join:id=account_id"`
}
