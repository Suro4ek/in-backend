package items

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"in-backend/internal/handlers"
	"in-backend/pkg/logging"
	"net/http"
	"strconv"
	"strings"
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
	router.PATCH(itemOneUrl, h.Edit)
	router.DELETE(itemUrl, h.Delete)
}

func NewHandler(repository Repository, logger *logging.Logger) handlers.HandlerAuth {
	return &handler{
		logger:     logger,
		repository: repository,
	}
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
	var dto CreateItemDTO
	if err := ctx.ShouldBind(&dto); err != nil || strings.TrimSpace(dto.Name) == "" || strings.TrimSpace(dto.ProductName) == "" || strings.TrimSpace(dto.SerialNumber) == "" {
		ctx.String(http.StatusBadRequest, "missing vals")
		return
	}

	itm := &Item{Name: dto.Name, ProductName: dto.ProductName, SerialNumber: dto.SerialNumber}

	if err := h.repository.Create(context.TODO(), itm); err != nil {
		h.logger.Errorf("Error create item %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}

	ctx.JSON(http.StatusOK, itm)
}

func (h *handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.repository.Delete(context.TODO(), id)
	if err != nil {
		h.logger.Errorf("Error get item %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":   200,
		"status": "Успешно удалено",
	})
}

func (h *handler) Edit(ctx *gin.Context) {
	var dto EditItemDTO
	itm, err := h.repository.GetOne(context.TODO(), strconv.Itoa(int(dto.ID)))
	if err != nil {
		h.logger.Errorf("Error get item %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	itm.Name = dto.Name
	itm.ProductName = dto.ProductName
	itm.SerialNumber = dto.SerialNumber
	itm.OwnerID = dto.OwnerID

	if err := h.repository.Update(context.TODO(), itm); err != nil {
		h.logger.Errorf("Error create user %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	ctx.JSON(http.StatusOK, itm)
}
