package controller

import "errors"

var ErrBadRequest = errors.New("bad request")
var ErrInternalServer = errors.New("internal server error")
var ErrAlreadyExists = errors.New("entity is already exists")
var ErrShortUsername = errors.New("username should be at least 8 characters")
var ErrBadAmount = errors.New("amount should be greater than 0")
var ErrUserNotFound = errors.New("user not found")
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrReceiverNotFound = errors.New("receiver not found")