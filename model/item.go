package model

import (
	"errors"
	"strings"
)

type Item struct {
	BaseModel
	OwnerID string `json:"ownerId"`
	Name    string `json:"name"`
}

var (
	ErrItemNameInvalid = errors.New("name is empty")
	ErrItemNameInUse   = errors.New("name is used by another item")
)

type AddItem struct {
	Name string `json:"name"`
}

func (item AddItem) Validation(existingItems []Item) error {
	if len(strings.TrimSpace(item.Name)) == 0 {
		return ErrListNameInvalid
	}

	for _, otherItem := range existingItems {
		if item.Name == otherItem.Name {
			return ErrListNameIsInUse
		}
	}

	return nil
}

type UpdateItem struct {
	Name string `json:"name"`
}

func (item UpdateItem) Validation(existingItems []Item) error {
	if len(strings.TrimSpace(item.Name)) == 0 {
		return ErrListNameInvalid
	}

	for _, otherItem := range existingItems {
		if item.Name == otherItem.Name {
			return ErrListNameIsInUse
		}
	}

	return nil
}
