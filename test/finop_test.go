package test

import (
	"context"
	"log"
	"testing"

	"github.com/Rolan335/Finop/internal/config"
	"github.com/Rolan335/Finop/internal/entity"
	"github.com/Rolan335/Finop/internal/repository/postgres"
	"github.com/Rolan335/Finop/internal/service/finop"
	"github.com/stretchr/testify/assert"
)

var (
	user1        = "useruser1"
	user1Balance = 0
	user2        = "useruser2"
	user2Balance = 0
)

func TestBusinessLogicWithRealDB(t *testing.T) {
	a := assert.New(t)
	envPath := ".env.test"
	cfg := config.MustNewConfig(envPath)
	cfg.Migration.Action = "down"
	if err := postgres.Migrate(&cfg.Migration); err != nil {
		log.Println("migrations doesn't initialized yet")
	}
	cfg.Migration.Action = "up"
	if err := postgres.Migrate(&cfg.Migration); err != nil {
		panic("failed to migrate: " + err.Error())
	}
	storage := postgres.MustNewStorage(&cfg.DB)

	service := finop.NewFinop(storage)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout)
	defer cancel()
	t.Run("AddUser", func(t *testing.T) {
		users := []string{user1, user2, "useruser2", "user4"}
		err := service.AddUser(ctx, users[0])
		a.NoError(err, "failed to add user")
		err = service.AddUser(ctx, users[1])
		a.NoError(err, "failed to add user")
		err = service.AddUser(ctx, users[2])
		a.ErrorIs(err, finop.ErrUserAlreadyExists, "should fail to add user")
		err = service.AddUser(ctx, users[3])
		a.ErrorIs(err, finop.ErrShortUsername, "should fail to add user")
	})
	t.Run("Deposit", func(t *testing.T) {
		deposits := []struct {
			username    string
			amount      int
			expectedErr error
			expected    int
		}{
			{
				username:    "hello",
				amount:      100,
				expectedErr: finop.ErrUserNotFound,
			},
			{
				username: user1,
				amount:   500,
				expected: 500,
			},
			{
				username: user1,
				amount:   1000,
				expected: 1500,
			},
			{
				username:    user2,
				amount:      0,
				expectedErr: finop.ErrBadAmount,
			},
			{
				username:    user1,
				amount:      -1000,
				expectedErr: finop.ErrBadAmount,
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout)
		defer cancel()
		for _, test := range deposits {
			newBalance, err := service.Deposit(ctx, test.username, test.amount)
			if test.expectedErr != nil {
				a.ErrorIs(err, test.expectedErr, "should return proper error")
				continue
			}
			if user1 == test.username {
				user1Balance = newBalance
				continue
			}
			if user2 == test.username {
				user1Balance = newBalance
				continue
			}
			a.Equal(newBalance, test.expected)
		}
		user1Balance = 1500
	})
	t.Run("Send", func(t *testing.T) {
		sends := []struct {
			username    string
			receiver    string
			amount      int
			expectedErr error
			expected    int
		}{
			{
				username:    "hello",
				receiver:    user1,
				amount:      100,
				expectedErr: finop.ErrUserNotFound,
			},
			{
				username:    user1,
				receiver:    user2,
				amount:      500_000_000,
				expectedErr: finop.ErrInsufficientFunds,
			},
			{
				username:    user1,
				receiver:    user2,
				amount:      -10000,
				expectedErr: finop.ErrBadAmount,
			},
			{
				username:    user2,
				receiver:    "hellohello",
				amount:      1000,
				expectedErr: finop.ErrReceiverNotFound,
			},
			{
				username: user1,
				receiver: user2,
				amount:   50,
				expected: user1Balance - 50,
			},
			{
				username: user2,
				receiver: user1,
				amount:   50,
				expected: user2Balance - 50,
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout)
		defer cancel()
		for _, test := range sends {
			newBalance, err := service.Send(ctx, test.username, test.receiver, test.amount)
			if test.expectedErr != nil {
				a.ErrorIs(err, test.expectedErr, "should return proper error")
				continue
			}
			if user1 == test.username {
				user1Balance = newBalance
				user2Balance += test.amount
				continue
			}
			if user2 == test.username {
				user2Balance = newBalance
				user1Balance += test.amount
				continue
			}
			a.Equal(newBalance, test.expected)
		}
	})
	t.Run("GetTransactions", func(t *testing.T) {
		getTransactions := []struct {
			username    string
			expected    []entity.Transaction
			expectedErr error
		}{
			{
				username:    "hellow",
				expectedErr: finop.ErrUserNotFound,
			},
			{
				username: user1,
				expected: []entity.Transaction{
					{
						Operation: "send",
						Receiver:  user2,
						Amount:    50,
					},
					{
						Operation: "deposit",
						Receiver:  "",
						Amount:    1000,
					},
					{
						Operation: "deposit",
						Receiver:  "",
						Amount:    500,
					},
				},
			},
			{
				username: user2,
				expected: []entity.Transaction{
					{
						Operation: "send",
						Receiver:  user1,
						Amount:    50,
					},
				},
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout)
		defer cancel()
		for _, test := range getTransactions {
			transactions, err := service.GetUserTransactions(ctx, test.username)
			if test.expectedErr != nil {
				a.ErrorIs(err, test.expectedErr, "should return proper error")
				continue
			}
			a.Equal(len(test.expected), len(transactions))
			for i := range test.expected {
				a.Equal(test.expected[i].Amount, transactions[i].Amount, "should be equal amount")
				a.Equal(test.expected[i].Operation, transactions[i].Operation, "should be same operation")
				a.Equal(test.expected[i].Receiver, transactions[i].Receiver, "should be same receiver")
			}
		}
	})
}
