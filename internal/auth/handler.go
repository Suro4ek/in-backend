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

func (h *handler) RegisterAuth(router *gin.RouterGroup) {
	//TODO implement me
	panic("implement me")
}

func NewHandler(logger *logging.Logger, repository user.Repository) handlers.Handler {
	return &handler{
		logger:     logger,
		repository: repository,
	}
}

func (h *handler) Register(router *gin.Engine) {
	router.POST(loginUrl, h.Login)
}

func (h *handler) Login(ctx *gin.Context) {
	var credential LoginCredentials
	credential.Username = ctx.PostForm("username")
	credential.Password = ctx.PostForm("password")
	//TODO check empty username and password
	usr, err := h.repository.GetByUsername(context.TODO(), credential.Username)
	if err != nil {
		ctx.String(http.StatusNotFound, "data not found")
		return
	}
	fmt.Println(usr.PasswordHash)
	if err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(credential.Password)); err != nil {
		ctx.String(http.StatusUnauthorized, "wrong password")
		return
	}
	//TODO return tokern
}
