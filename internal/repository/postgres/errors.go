package postgres

import "errors"

var ErrAlreadyExists = errors.New("already exists")
var ErrNotFound = errors.New("not found")
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrReceiverNotFound = errors.New("receiver not found")