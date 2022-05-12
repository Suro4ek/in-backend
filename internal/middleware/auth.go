package middleware

import (
	"context"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"in-backend/internal/config"
	"in-backend/internal/user"
	"time"
)

var identityKey = "id"

type authMiddleware struct {
	Auth *jwt.GinJWTMiddleware
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func NewAuthMiddleWare(Repository user.Repository, cfg *config.Config) (*authMiddleware, error) {
	auth, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "auth",
		Key:         []byte(cfg.Secret),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(uint); ok {
				return jwt.MapClaims{
					identityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(g *gin.Context) interface{} {
			claims := jwt.ExtractClaims(g)
			return claims["id"].(float64)
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password
			usr, err := Repository.GetByUsername(context.TODO(), userID)
			if err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			if err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return usr.ID, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(float64); ok {
				_, err := Repository.GetOne(context.TODO(), fmt.Sprint(v))
				if err != nil {
					return false
				}
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	return &authMiddleware{
		Auth: auth,
	}, err
}
