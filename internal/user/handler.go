package user

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"in-backend/internal/config"
	"in-backend/internal/handlers"
	"in-backend/pkg/logging"
	"net/http"
)

const (
	userUrl    = "/user/:id"
	userOneUrl = "/user"
)

type handler struct {
	logger     *logging.Logger
	repository Repository
	cfg        *config.Config
}

func (h *handler) RegisterAuth(router *gin.RouterGroup) {
	router.GET(userOneUrl, h.getUser)
}

func (h *handler) Register(router *gin.Engine) {

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
		return
	}
	ctx.JSON(http.StatusOK, usr)
}
