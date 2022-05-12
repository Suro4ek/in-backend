package main

import (
	"context"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"in-backend/internal/composites"
	"in-backend/internal/config"
	"in-backend/internal/middleware"
	"in-backend/pkg/logging"
	"in-backend/pkg/postgres"
	"log"
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
	r.NoRoute(authMiddleware.Auth.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	authcomposite, _ := composites.NewAuthComposite(&logger, userComposite.Repository, cfg)
	authcomposite.Handler.Register(r)
	authorized := r.Group("/api")
	authorized.GET("/refresh_token", authMiddleware.Auth.RefreshHandler)
	authorized.Use(authMiddleware.Auth.MiddlewareFunc())

	itemComposite.Handler.RegisterAuth(authorized)
	userComposite.Handler.RegisterAuth(authorized)

	err = r.Run(cfg.Listen.Host + ":" + cfg.Listen.Port)
	if err != nil {
		logger.Errorf("error startup application %t", err)
		return
	}
}
