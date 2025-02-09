package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID        uuid.UUID `json:"transactionId"`
	Operation string    `json:"operation"`
	Receiver  string    `json:"receiver"`
	Amount    int       `json:"amount"`
	Time      time.Time `json:"time"`
}

type User struct {
	Username string `json:"username"`
}

type Send struct{
	Receiver string `json:"receiver"`
	Amount int `json:"amount"`
}

type Amount struct {
	Amount int `json:"amount"`
}
