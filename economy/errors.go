package economy

import "errors"

var (
	ErrUnknownPlayer     = errors.New("unknown player")
	ErrNegativeAmount    = errors.New("amount must be positive")
	ErrInsufficientFunds = errors.New("insufficient funds")
)
