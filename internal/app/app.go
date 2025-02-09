package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Rolan335/Finop/internal/config"
	"github.com/Rolan335/Finop/internal/controller"
	"github.com/Rolan335/Finop/pkg/api"
	"github.com/gin-gonic/gin"
)

type App struct {
	server *http.Server
}

// for graceful shutdown of services (at our case postgres), they should have method Close
type Close interface {
	Close()
}

func NewService(config *config.Config, server *controller.Server) *App {
	gin.SetMode(config.GinMode)
	r := gin.Default()

	api.RegisterHandlers(r, server)

	return &App{
		server: &http.Server{
			Addr:    config.Port,
			Handler: r,
		},
	}
}

func (a *App) Start() {
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			log.Println("server shut ", err.Error())
		}
	}()
}

func (a *App) GracefulStop(services ...interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Println("Failed to graceful shutdown", err.Error())
	}
	log.Println("gracefully shut")
	for _, service := range services {
		if asserted, ok := service.(Close); ok {
			asserted.Close()
		}
	}
}
