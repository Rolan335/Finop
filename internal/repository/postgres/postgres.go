package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Rolan335/Finop/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	depositID = 1
	sendID    = 2
)

type Config struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Name     string `env:"POSTGRES_NAME"`
}

type Storage struct {
	db *pgxpool.Pool
}

func MustNewStorage(cfg *Config) *Storage {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		panic("failed to create pool: " + err.Error())
	}

	if err := conn.Ping(context.Background()); err != nil {
		panic("can't connect to postgres: " + err.Error())
	}
	return &Storage{
		db: conn,
	}
}

func (s *Storage) AddUser(ctx context.Context, username string) error {
	initBalance := 0
	query := `INSERT INTO users (name, balance) VALUES ($1, $2) RETURNING id`
	if err := s.db.QueryRow(ctx, query, username, initBalance).Scan(nil); err != nil {
		//unique constraint violation check
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ErrAlreadyExists
			}
			return pgErr
		}
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}
func (s *Storage) Deposit(ctx context.Context, txUUID uuid.UUID, username string, amount int) (newBalance int, err error) {
	var userID int
	if err := s.db.QueryRow(ctx, "SELECT id FROM users WHERE name = $1 LIMIT 1", username).Scan(&userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("failed to get user: %w", err)
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		if err := tx.Commit(ctx); err != nil {
			tx.Rollback(ctx)
			return
		}
	}()
	query := "INSERT INTO transactions (id, user_id, operation_id, amount, created_at) VALUES ($1, $2, $3, $4, $5)"
	if _, err := tx.Exec(ctx, query, txUUID, userID, depositID, amount, time.Now()); err != nil {
		return 0, fmt.Errorf("failed to insert transaction: %w", err)
	}
	var balance int
	if err := tx.QueryRow(ctx, "UPDATE users SET balance = balance + $1 WHERE name = $2 RETURNING balance", amount, username).Scan(&balance); err != nil {
		return 0, fmt.Errorf("failed to update balance: %w", err)
	}
	return balance, nil
}
func (s *Storage) Send(ctx context.Context, txUUID uuid.UUID, username string, receiver string, amount int) (newBalance int, err error) {
	var userID int
	if err := s.db.QueryRow(ctx, "SELECT id FROM users WHERE name = $1 LIMIT 1", username).Scan(&userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("failed to get user: %w", err)
	}
	var receiverID int
	if err := s.db.QueryRow(ctx, "SELECT id FROM users WHERE name = $1 LIMIT 1", receiver).Scan(&receiverID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrReceiverNotFound
		}
		return 0, fmt.Errorf("failed to get receiver: %w", err)
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		if err := tx.Commit(ctx); err != nil {
			tx.Rollback(ctx)
			return
		}
	}()
	var balance int
	if err := tx.QueryRow(ctx, "UPDATE users SET balance = balance - $1 WHERE id = $2 RETURNING balance", amount, userID).Scan(&balance); err != nil {
		return 0, fmt.Errorf("failed to update balance: %w", err)
	}
	if balance < 0 {
		return 0, ErrInsufficientFunds
	}
	if err := tx.QueryRow(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2 RETURNING balance", amount, receiverID).Scan(nil); err != nil {
		return 0, fmt.Errorf("failed to update balance: %w", err)
	}
	query := "INSERT INTO transactions (id, user_id, receiver_id, operation_id, amount, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	if _, err := tx.Exec(ctx, query, txUUID, userID, receiverID, sendID, amount, time.Now()); err != nil {
		return 0, fmt.Errorf("failed to insert transaction: %w", err)
	}
	return balance, nil
}
func (s *Storage) GetTransactions(ctx context.Context, username string, limit int) ([]entity.Transaction, error) {
	var userID int
	if err := s.db.QueryRow(ctx, "SELECT id FROM users WHERE name = $1 LIMIT 1", username).Scan(&userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	query := `SELECT 
	t.id,
    u.name as receiver,
	o.type,
	t.amount,
	t.created_at
	FROM transactions t
	INNER JOIN operations o ON t.operation_id = o.id
    LEFT JOIN users u ON t.receiver_id = u.id
	WHERE t.user_id = $1
	ORDER BY 
		t.created_at DESC
	LIMIT $2;
	`
	rows, err := s.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find transactions: %w", err)
	}
	defer rows.Close()
	transactions := make([]entity.Transaction, 0, limit)
	for rows.Next() {
		var transaction entity.Transaction
		var receiver sql.NullString
		err := rows.Scan(&transaction.ID, &receiver, &transaction.Operation, &transaction.Amount, &transaction.Time)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if receiver.Valid {
			transaction.Receiver = receiver.String
		} else {
			transaction.Receiver = ""
		}
		transactions = append(transactions, transaction)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error in row: %w", rows.Err())
	}
	return transactions, nil
}
