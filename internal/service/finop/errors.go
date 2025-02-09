package finop

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrShortUsername = errors.New("username should be at least 8 characters")
var ErrUserNotFound = errors.New("user not found")
var ErrBadAmount = errors.New("amount should be greater than 0")
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrReceiverNotFound = errors.New("receiver not found")