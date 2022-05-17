package user

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"in-backend/internal/config"
	"in-backend/internal/handlers"
	"in-backend/pkg/logging"
	"net/http"
	"strings"
)

const (
	userUrl    = "/user/:id"
	userOneUrl = "/user"
	usersUrl   = "/users"
)

type handler struct {
	logger     *logging.Logger
	repository Repository
	cfg        *config.Config
}

func (h *handler) RegisterAdmin(router *gin.RouterGroup) {
	router.PATCH(userUrl, h.EditUser)
	router.GET(usersUrl, h.getUsers)
	router.DELETE(userUrl, h.deleteUser)
}

func (h *handler) RegisterAuth(router *gin.RouterGroup) {
	router.GET(userOneUrl, h.getUser)
}

func NewHandler(repository Repository, logger *logging.Logger, cfg *config.Config) handlers.HandlerAuth {
	return &handler{
		logger:     logger,
		repository: repository,
		cfg:        cfg,
	}
}

func (h *handler) getUser(ctx *gin.Context) {
	id, _ := ctx.Get("id")

	usr, err := h.repository.GetOne(context.TODO(), fmt.Sprint(id))
	if err != nil {
		h.logger.Errorf("Error get user %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	ctx.JSON(http.StatusOK, usr)
}

func (h *handler) getUsers(ctx *gin.Context) {
	usrs, err := h.repository.GetAll(context.TODO())
	if err != nil {
		h.logger.Errorf("Error get user %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	ctx.JSON(http.StatusOK, usrs)
}

func (h *handler) EditUser(ctx *gin.Context) {
	var dto EditUserDTO
	id := ctx.Param("id")
	if err := ctx.ShouldBind(&dto); err != nil {
		ctx.String(http.StatusBadRequest, "missing vals")
		return
	}
	usr, err := h.repository.GetOne(context.TODO(), id)
	if err != nil {
		h.logger.Errorf("Error get user %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	if dto.Name != nil && strings.TrimSpace(*dto.Name) != "" {
		usr.Name = *dto.Name
	}

	if dto.Familia != nil && strings.TrimSpace(*dto.Familia) != "" {
		usr.Familia = *dto.Familia
	}

	if dto.Password != nil && strings.TrimSpace(*dto.Password) != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(*dto.Password)), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("error update user", err)
			return
		}
		usr.PasswordHash = string(passwordHash)
	}

	if err := h.repository.Update(context.TODO(), usr); err != nil {
		h.logger.Errorf("Error update user %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	ctx.JSON(http.StatusOK, usr)
}

func (h *handler) deleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.repository.Delete(context.TODO(), id)
	if err != nil {
		h.logger.Errorf("Error delete user %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "successful deleted",
	})
}
