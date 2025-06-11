package errors

import "errors"

var (
	ErrUnknownPlayer     = errors.New("unknown player")
	ErrNegativeAmount    = errors.New("amount must be positive")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrPageNotFound      = errors.New("page is not found")
)
