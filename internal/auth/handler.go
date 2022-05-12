package auth

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"in-backend/internal/handlers"
	"in-backend/internal/user"
	"in-backend/pkg/logging"
	"net/http"
)

const (
	loginUrl = "/login"
)

type handler struct {
	logger     *logging.Logger
	repository user.Repository
}

func NewHandler(logger *logging.Logger, repository user.Repository) handlers.Handler {
	return &handler{
		logger:     logger,
		repository: repository,
	}
}

func (h *handler) Register(router *gin.Engine) {
	//router.POST(loginUrl, h.Login)
}

func (h *handler) Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "" || password == "" {
		ctx.String(http.StatusNotFound, "empty string")
		return
	}
	usr, err := h.repository.GetByUsername(context.TODO(), username)
	if err != nil {
		ctx.String(http.StatusNotFound, "data not found")
		return
	}
	fmt.Println(usr.PasswordHash)
	if err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(password)); err != nil {
		ctx.String(http.StatusUnauthorized, "wrong password")
		return
	}
	//TODO return tokern
}
