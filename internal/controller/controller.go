package controller

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Rolan335/Finop/internal/entity"
	"github.com/Rolan335/Finop/internal/service/finop"
	"github.com/gin-gonic/gin"
)

type Server struct {
	timeout time.Duration
	service *finop.Finop
}

func New(timeout time.Duration, service *finop.Finop) *Server {
	return &Server{
		timeout: timeout,
		service: service,
	}
}

// Добавление нового пользователя
// (POST /users)
func (s *Server) PostUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), s.timeout)
	defer cancel()
	var user entity.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBadRequest.Error()})
		return
	}
	if user == (entity.User{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBadRequest.Error()})
		return
	}
	if err := s.service.AddUser(ctx, user.Username); err != nil {
		if errors.Is(err, finop.ErrShortUsername) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrShortUsername.Error()})
			return
		}
		if errors.Is(err, finop.ErrUserAlreadyExists) {
			c.JSON(http.StatusAlreadyReported, gin.H{"error": ErrAlreadyExists.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternalServer.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Пополнение баланса пользователя
// (POST /users/{username}/balance)
func (s *Server) PostUsersUsernameBalance(c *gin.Context, username string) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), s.timeout)
	defer cancel()
	var amount *entity.Amount
	if err := c.BindJSON(&amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBadRequest.Error()})
		return
	}
	if amount == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBadRequest.Error()})
		return
	}
	newBalance, err := s.service.Deposit(ctx, username, amount.Amount)
	if err != nil {
		if errors.Is(err, finop.ErrBadAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrBadAmount.Error()})
			return
		}
		if errors.Is(err, finop.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternalServer.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": newBalance})
}

// Просмотр 10 последних операций пользователя
// (GET /users/{username}/transactions)
func (s *Server) GetUsersUsernameTransactions(c *gin.Context, username string) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), s.timeout)
	defer cancel()
	transactions, err := s.service.GetUserTransactions(ctx, username)
	if err != nil {
		if errors.Is(err, finop.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternalServer.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

// Перевод денег от одного пользователя к другому
// (POST /users/{username}/transfer)
func (s *Server) PostUsersUsernameTransfer(c *gin.Context, username string) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), s.timeout)
	defer cancel()
	var send *entity.Send
	if err := c.BindJSON(&send); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBadRequest.Error()})
		return
	}
	if send == nil || send.Receiver == "" || send.Amount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrBadRequest.Error()})
		return
	}
	newBalance, err := s.service.Send(ctx, username, send.Receiver, send.Amount)
	if err != nil {
		if errors.Is(err, finop.ErrBadAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrBadAmount.Error()})
			return
		}
		if errors.Is(err, finop.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error_code": "user_not_found", "message": ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, finop.ErrReceiverNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error_code": "user_not_found", "message": ErrReceiverNotFound.Error()})
			return
		}
		if errors.Is(err, finop.ErrInsufficientFunds) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrInsufficientFunds.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternalServer.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": newBalance})
}
