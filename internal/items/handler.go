package items

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"in-backend/internal/handlers"
	"in-backend/pkg/logging"
	"net/http"
)

const (
	itemsUrl   = "/items"
	itemUrl    = "/item/:id"
	itemOneUrl = "/item/"
)

type handler struct {
	logger     *logging.Logger
	repository Repository
}

func (h *handler) RegisterAuth(router *gin.RouterGroup) {
	router.GET(itemsUrl, h.GetItems)
	router.GET(itemUrl, h.GetItem)
	router.POST(itemOneUrl, h.CreateItem)
	router.PATCH(itemUrl, h.Edit)
}

func NewHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handler{
		logger:     logger,
		repository: repository,
	}
}

func (h *handler) Register(router *gin.Engine) {

}

func (h *handler) GetItems(ctx *gin.Context) {
	all, err := h.repository.GetAll(context.TODO())
	if err != nil {
		ctx.Status(http.StatusNotFound)
	}

	//allBytes, err := json.Marshal(all)
	//if err != nil {
	//	ctx.String(http.StatusInternalServerError, "error %t", err)
	//}

	ctx.JSON(http.StatusOK, all)
}

func (h *handler) GetItem(ctx *gin.Context) {
	id := ctx.Param("id")
	one, err := h.repository.GetOne(context.TODO(), id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
	}

	allBytes, err := json.Marshal(one)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}

	ctx.JSONP(http.StatusOK, allBytes)
}

func (h *handler) CreateItem(ctx *gin.Context) {
	var itm Item
	if err := json.NewDecoder(ctx.Request.Body).Decode(&itm); err != nil {
		h.logger.Errorf("error json decode %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}

	if err := h.repository.Create(context.TODO(), &itm); err != nil {
		h.logger.Errorf("Error create item %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	//TODO redirect location ctx.Redirect(code, location)
	ctx.String(http.StatusOK, "item successful created")
}

func (h *handler) Edit(ctx *gin.Context) {
	var itm Item
	if err := json.NewDecoder(ctx.Request.Body).Decode(&itm); err != nil {
		h.logger.Errorf("error json decode %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}

	if err := h.repository.Update(context.TODO(), itm); err != nil {
		h.logger.Errorf("Error create user %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	//TODO redirect location ctx.Redirect(code, location)
	ctx.String(http.StatusOK, "user updated")
}
