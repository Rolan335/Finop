// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// PostUsersJSONBody defines parameters for PostUsers.
type PostUsersJSONBody struct {
	Username *string `json:"username,omitempty"`
}

// PostUsersUsernameBalanceJSONBody defines parameters for PostUsersUsernameBalance.
type PostUsersUsernameBalanceJSONBody struct {
	Amount *float32 `json:"amount,omitempty"`
}

// PostUsersUsernameTransferJSONBody defines parameters for PostUsersUsernameTransfer.
type PostUsersUsernameTransferJSONBody struct {
	Amount     *float32 `json:"amount,omitempty"`
	ReceiverId *string  `json:"receiverId,omitempty"`
}

// PostUsersJSONRequestBody defines body for PostUsers for application/json ContentType.
type PostUsersJSONRequestBody PostUsersJSONBody

// PostUsersUsernameBalanceJSONRequestBody defines body for PostUsersUsernameBalance for application/json ContentType.
type PostUsersUsernameBalanceJSONRequestBody PostUsersUsernameBalanceJSONBody

// PostUsersUsernameTransferJSONRequestBody defines body for PostUsersUsernameTransfer for application/json ContentType.
type PostUsersUsernameTransferJSONRequestBody PostUsersUsernameTransferJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Добавление нового пользователя
	// (POST /users)
	PostUsers(c *gin.Context)
	// Пополнение баланса пользователя
	// (POST /users/{username}/balance)
	PostUsersUsernameBalance(c *gin.Context, username string)
	// Просмотр 10 последних операций пользователя
	// (GET /users/{username}/transactions)
	GetUsersUsernameTransactions(c *gin.Context, username string)
	// Перевод денег от одного пользователя к другому
	// (POST /users/{username}/transfer)
	PostUsersUsernameTransfer(c *gin.Context, username string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// PostUsers operation middleware
func (siw *ServerInterfaceWrapper) PostUsers(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostUsers(c)
}

// PostUsersUsernameBalance operation middleware
func (siw *ServerInterfaceWrapper) PostUsersUsernameBalance(c *gin.Context) {

	var err error

	// ------------- Path parameter "username" -------------
	var username string

	err = runtime.BindStyledParameterWithOptions("simple", "username", c.Param("username"), &username, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter username: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostUsersUsernameBalance(c, username)
}

// GetUsersUsernameTransactions operation middleware
func (siw *ServerInterfaceWrapper) GetUsersUsernameTransactions(c *gin.Context) {

	var err error

	// ------------- Path parameter "username" -------------
	var username string

	err = runtime.BindStyledParameterWithOptions("simple", "username", c.Param("username"), &username, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter username: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetUsersUsernameTransactions(c, username)
}

// PostUsersUsernameTransfer operation middleware
func (siw *ServerInterfaceWrapper) PostUsersUsernameTransfer(c *gin.Context) {

	var err error

	// ------------- Path parameter "username" -------------
	var username string

	err = runtime.BindStyledParameterWithOptions("simple", "username", c.Param("username"), &username, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter username: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostUsersUsernameTransfer(c, username)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.POST(options.BaseURL+"/users", wrapper.PostUsers)
	router.POST(options.BaseURL+"/users/:username/balance", wrapper.PostUsersUsernameBalance)
	router.GET(options.BaseURL+"/users/:username/transactions", wrapper.GetUsersUsernameTransactions)
	router.POST(options.BaseURL+"/users/:username/transfer", wrapper.PostUsersUsernameTransfer)
}
