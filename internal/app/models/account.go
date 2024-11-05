package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID                 string          `bun:",pk,type:uuid" json:"-"`
	UserID             string          `bun:",unique,notnull" json:"-"`
	CurrentPointsTotal decimal.Decimal `bun:",notnull,default:0" json:"current"`
	WithdrawnTotal     decimal.Decimal `bun:",notnull,default:0" json:"withdrawn"`
	CreatedAt          time.Time       `bun:",notnull,default:current_timestamp" json:"-"`
	UpdatedAt          time.Time       `bun:",notnull,default:current_timestamp" json:"-"`

	User   User     `bun:"rel:belongs-to,join:user_id=id" json:"-"`
	Orders []*Order `bun:"rel:has-many,join:id=account_id" json:"-"`
}
