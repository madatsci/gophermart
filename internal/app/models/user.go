package models

import "time"

type User struct {
	ID        string    `bun:",pk,type:uuid"`
	Login     string    `bun:",unique,notnull"`
	Password  string    `bun:",notnull"`
	CreatedAt time.Time `bun:",notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp"`

	Account *Account `bun:"rel:has-one,join:id=user_id"`
}
