package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"in-backend/internal/composites"
	"in-backend/internal/config"
	"in-backend/pkg/logging"
	"in-backend/pkg/postgres"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Info("create router")
	r := gin.New()

	cfg := config.GetConfig()

	authorized := r.Group("/api")
	client := postgres.NewClient(context.Background(), 5, cfg.Postgres)

	itemComposite, _ := composites.NewItemComposite(client, &logger)
	itemComposite.Handler.RegisterAuth(authorized)

	userComposite, _ := composites.NewUserComposite(client, &logger)
	userComposite.Handler.RegisterAuth(authorized)

	authComposite, _ := composites.NewAuthComposite(&logger, userComposite.Repository)
	authComposite.Handler.Register(r)

	err := r.Run("0.0.0.0:8080")
	if err != nil {
		logger.Errorf("error startup application %t", err)
		return
	}
}
