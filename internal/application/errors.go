package application

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrNotEnoughFunds  = errors.New("not enough funds")
)
