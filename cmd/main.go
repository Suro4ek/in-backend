package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"in-backend/internal/composites"
	"in-backend/internal/config"
	"in-backend/internal/middleware"
	"in-backend/pkg/logging"
	"in-backend/pkg/postgres"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Info("create router")
	r := gin.New()

	cfg := config.GetConfig()

	client := postgres.NewClient(context.Background(), 5, cfg.Postgres)

	itemComposite, _ := composites.NewItemComposite(client, &logger)

	userComposite, _ := composites.NewUserComposite(client, &logger, cfg)

	authMiddleware, err := middleware.NewAuthMiddleWare(userComposite.Repository, cfg)
	r.POST("/login", authMiddleware.Auth.LoginHandler)

	authorized := r.Group("/api")
	authorized.Use(authMiddleware.Auth.MiddlewareFunc())

	itemComposite.Handler.RegisterAuth(authorized)
	userComposite.Handler.Register(r)

	err = r.Run(cfg.Listen.Host + ":" + cfg.Listen.Port)
	if err != nil {
		logger.Errorf("error startup application %t", err)
		return
	}
}
