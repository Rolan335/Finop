package finop

import (
	"context"
	"errors"
	"log"

	"github.com/Rolan335/Finop/internal/entity"
	"github.com/Rolan335/Finop/internal/repository/postgres"
	"github.com/google/uuid"
)

type Storage interface {
	AddUser(ctx context.Context, username string) error
	Deposit(ctx context.Context, txUUID uuid.UUID, username string, amount int) (int, error)
	Send(ctx context.Context, txUUID uuid.UUID, username string, receiver string, amount int) (int, error)
	GetTransactions(ctx context.Context, username string, limit int) ([]entity.Transaction, error)
}

type Finop struct {
	storage Storage
}

func NewFinop(storage Storage) *Finop {
	return &Finop{
		storage: storage,
	}
}

func (f *Finop) AddUser(ctx context.Context, username string) error {
	if len([]rune(username)) < 8 {
		return ErrShortUsername
	}
	if err := f.storage.AddUser(ctx, username); err != nil {
		if errors.Is(err, postgres.ErrAlreadyExists) {
			return ErrUserAlreadyExists
		}
		log.Println(err)
		return err
	}
	return nil
}

func (f *Finop) Deposit(ctx context.Context, username string, amount int) (int, error) {
	if amount < 1 {
		return 0, ErrBadAmount
	}
	txUUID, err := uuid.NewRandom()
	if err != nil {
		log.Println("failed to generate uuid", err.Error())
		return 0, err
	}
	newBalance, err := f.storage.Deposit(ctx, txUUID, username, amount)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return 0, ErrUserNotFound
		}
		log.Println(err)
		return 0, err
	}
	return newBalance, nil
}

func (f *Finop) Send(ctx context.Context, username string, receiver string, amount int) (int, error) {
	if amount < 1 {
		return 0, ErrBadAmount
	}
	txUUID, err := uuid.NewRandom()
	if err != nil {
		log.Println("failed to generate uuid", err.Error())
		return 0, err
	}
	newBalance, err := f.storage.Send(ctx, txUUID, username, receiver, amount)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return 0, ErrUserNotFound
		}
		if errors.Is(err, postgres.ErrReceiverNotFound) {
			return 0, ErrReceiverNotFound
		}
		if errors.Is(err, postgres.ErrInsufficientFunds) {
			return 0, ErrInsufficientFunds
		}
		log.Println(err)
		return 0, err
	}
	return newBalance, nil
}

func (f *Finop) GetUserTransactions(ctx context.Context, username string) ([]entity.Transaction, error) {
	limit := 10
	transactions, err := f.storage.GetTransactions(ctx, username, limit)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		log.Println(err)
		return nil, err
	}
	return transactions, nil
}
