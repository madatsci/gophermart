package models

import "time"

type Order struct {
	ID        string    `bun:",pk,type:uuid"`
	AccountID string    `bun:",notnull"`
	Number    string    `bun:",unique,notnull"`
	CreatedAt time.Time `bun:",notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp"`

	Account Account `bun:"rel:belongs-to,join:account_id=id"`
}
