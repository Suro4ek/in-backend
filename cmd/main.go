package main

import (
	"bufio"
	"context"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"in-backend/internal/composites"
	"in-backend/internal/config"
	"in-backend/internal/middleware"
	"in-backend/internal/user"
	"in-backend/pkg/logging"
	"in-backend/pkg/postgres"
	"log"
	"os"
	"strings"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Info("create router")
	r := gin.New()

	cfg := config.GetConfig()

	client := postgres.NewClient(context.Background(), 5, cfg.Postgres)

	userComposite, _ := composites.NewUserComposite(client, &logger, cfg)
	checkAdminUser(&logger, userComposite.Repository)
	itemComposite, _ := composites.NewItemComposite(client, userComposite.Repository, &logger)

	authMiddleware, err := middleware.NewAuthMiddleWare(userComposite.Repository, cfg)
	r.POST("/login", authMiddleware.Auth.LoginHandler)
	r.NoRoute(authMiddleware.Auth.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	authcomposite, _ := composites.NewAuthComposite(&logger, userComposite.Repository, cfg)
	authorized := r.Group("/api")
	authorized.GET("/refresh_token", authMiddleware.Auth.RefreshHandler)
	authorized.Use(authMiddleware.Auth.MiddlewareFunc())
	authorized1 := r.Group("/admin")
	authorized1.Use(authMiddleware.Auth.MiddlewareFunc())

	authcomposite.Handler.RegisterAdmin(authorized1)
	itemComposite.Handler.RegisterAuth(authorized)
	itemComposite.Handler.RegisterAdmin(authorized1)
	userComposite.Handler.RegisterAdmin(authorized1)
	userComposite.Handler.RegisterAuth(authorized)

	err = r.Run("0.0.0.0:8080")
	if err != nil {
		logger.Errorf("error startup application %t", err)
		return
	}
}

func checkAdminUser(logger *logging.Logger, repository user.Repository) {
	fmt.Println("Find admin user...")
	existsAdmin := repository.CheckAdmin(context.TODO())
	if !existsAdmin {
		fmt.Println("admin user not found")
		fmt.Println("create admin user")
		fmt.Print("Enter username:")
		reader := bufio.NewReader(os.Stdin)
		username, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		fmt.Print("Enter password:")
		reader = bufio.NewReader(os.Stdin)
		password, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(password)), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("error create user", err)
			return
		}
		usr := &user.User{
			Username:     strings.TrimSpace(username),
			PasswordHash: string(passwordHash),
			Role:         "admin",
			Name:         "Admin",
			Familia:      "Admin",
		}
		err = repository.Create(context.TODO(), usr)
		if err != nil {
			fmt.Println("error create user", err)
			return
		}
	} else {
		fmt.Println("admin user found")
	}
}
