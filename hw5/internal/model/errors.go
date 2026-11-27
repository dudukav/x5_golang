package model

import "errors"

var (
	ErrAlreadyExists = errors.New("order already exists")
	ErrNotFound      = errors.New("order not found")
)
