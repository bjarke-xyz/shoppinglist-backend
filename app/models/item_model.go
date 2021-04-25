package models

import (
	"time"
)

type Item struct {
	ID        uint      `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt time.Time `db:"deleted_at" json:"deletedAt"`

	Name    string `db:"name" json:"name"`
	OwnerID string `db:"owner_id" json:"ownerId"`
}
