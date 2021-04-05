package model

import "time"

type ListItem struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	ListID    uint       `json:"listId"`
	ItemID    uint       `json:"itemId"`
	Item      Item       `json:"item"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
