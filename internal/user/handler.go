package user

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"in-backend/internal/handlers"
	"in-backend/pkg/logging"
	"net/http"
)

const (
	userUrl    = "/user/:id"
	userOneUrl = "/user"
	loginUrl   = "/login"
)

type handler struct {
	logger     *logging.Logger
	repository Repository
}

func (h *handler) RegisterAuth(router *gin.RouterGroup) {
	router.POST(userOneUrl, h.CreateUser)
}

func NewHandler(repository Repository, logger *logging.Logger) handlers.HandlerAuth {
	return &handler{
		logger:     logger,
		repository: repository,
	}
}

func (h *handler) CreateUser(ctx *gin.Context) {
	var usr User
	if err := json.NewDecoder(ctx.Request.Body).Decode(&usr); err != nil {
		h.logger.Errorf("error json decode %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	password, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Errorf("error to bcrypt password %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	usr.PasswordHash = string(password)
	if err := h.repository.Create(context.TODO(), &usr); err != nil {
		h.logger.Errorf("Error create user %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	//TODO redirect location ctx.Redirect(code, location)
	ctx.String(http.StatusOK, "user successful created %s", usr.ID)
}
