package model

import (
	"errors"
	"strings"
)

type List struct {
	BaseModel
	OwnerID string     `json:"ownerId"`
	Name    string     `json:"name"`
	Default bool       `json:"default"`
	Items   []ListItem `json:"items"`
}

var (
	ErrListNameInvalid = errors.New("name is empty")
	ErrListNameIsInUse = errors.New("name is used by another list")
)

type AddList struct {
	Name string `json:"name"`
}

func (list AddList) Validation(existingLists []List) error {
	if len(strings.TrimSpace(list.Name)) == 0 {
		return ErrListNameInvalid
	}

	for _, otherList := range existingLists {
		if list.Name == otherList.Name {
			return ErrListNameIsInUse
		}
	}

	return nil
}

type UpdateList struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

func (list UpdateList) Validation(existingLists []List) error {
	if len(strings.TrimSpace(list.Name)) == 0 {
		return ErrListNameInvalid
	}

	for _, otherList := range existingLists {
		if list.Name == otherList.Name && list.ID != otherList.ID {
			return ErrListNameIsInUse
		}
	}

	return nil
}
