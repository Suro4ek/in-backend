package items

import (
	"context"
	"github.com/gin-gonic/gin"
	"in-backend/internal/handlers"
	"in-backend/internal/user"
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
	logger         *logging.Logger
	repository     Repository
	userRepository user.Repository
}

func (h *handler) RegisterAdmin(router *gin.RouterGroup) {
	router.GET(itemsUrl, h.GetItems)
	router.GET(itemUrl, h.GetItem)
	router.POST(itemOneUrl, h.CreateItem)
	router.DELETE(itemUrl, h.Delete)
	router.PATCH(itemUrl, h.EditAdmin)
}

func (h *handler) RegisterAuth(router *gin.RouterGroup) {
	router.GET(itemsUrl, h.GetItems)
	router.PATCH(itemUrl, h.Edit)
}

func NewHandler(repository Repository, userRepository user.Repository, logger *logging.Logger) handlers.HandlerAuth {
	return &handler{
		logger:         logger,
		repository:     repository,
		userRepository: userRepository,
	}
}

func (h *handler) GetItems(ctx *gin.Context) {
	all, err := h.repository.GetAll(context.TODO())
	if err != nil {
		ctx.Status(http.StatusNotFound)
	}
	items := make([]Item, 0)
	for _, itm := range all {
		if itm.OwnerID != nil {
			usr, err := h.userRepository.GetOne(context.TODO(), strconv.Itoa(*itm.OwnerID))
			if err != nil {
				h.logger.Errorf("Error get user by item %t", err)
				itm.Owner = nil
			} else {
				itm.Owner = &usr
			}
		}
		items = append(items, itm)
	}
	ctx.JSON(http.StatusOK, items)
}

func (h *handler) GetItem(ctx *gin.Context) {
	id := ctx.Param("id")
	one, err := h.repository.GetOne(context.TODO(), id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
	}

	ctx.JSON(http.StatusOK, one)
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
	id := ctx.Param("id")
	if err := ctx.ShouldBind(&dto); err != nil {
		ctx.String(http.StatusBadRequest, "missing vals")
		return
	}

	itm, err := h.repository.GetOne(context.TODO(), id)
	if err != nil {
		h.logger.Errorf("Error get item %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	if dto.OwnerID != nil {
		if *dto.OwnerID < 0 {
			itm.Owner = nil
		} else {
			usr, err := h.userRepository.GetOne(context.TODO(), strconv.Itoa(*dto.OwnerID))
			if err != nil {
				h.logger.Errorf("Error get user by item %t", err)
				ctx.String(http.StatusInternalServerError, "error %t", err)
			}
			itm.Owner = &usr
		}
	}

	if err := h.repository.Update(context.TODO(), itm); err != nil {
		h.logger.Errorf("Error update item %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	ctx.JSON(http.StatusOK, itm)
}

func (h *handler) EditAdmin(ctx *gin.Context) {
	var dto EditItemDTO
	id := ctx.Param("id")
	if err := ctx.ShouldBind(&dto); err != nil {
		ctx.String(http.StatusBadRequest, "missing vals")
		return
	}

	itm, err := h.repository.GetOne(context.TODO(), id)
	if err != nil {
		h.logger.Errorf("Error get item %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	if dto.OwnerID != nil {
		if *dto.OwnerID < 0 {
			itm.Owner = nil
		} else {
			usr, err := h.userRepository.GetOne(context.TODO(), strconv.Itoa(*dto.OwnerID))
			if err != nil {
				h.logger.Errorf("Error get user by item %t", err)
				ctx.String(http.StatusInternalServerError, "error %t", err)
			}
			itm.Owner = &usr
		}
	}
	if dto.Name != nil && strings.TrimSpace(*dto.Name) != "" {
		itm.Name = *dto.Name
	}

	if dto.ProductName != nil && strings.TrimSpace(*dto.ProductName) != "" {
		itm.ProductName = *dto.ProductName
	}

	if dto.SerialNumber != nil && strings.TrimSpace(*dto.SerialNumber) != "" {
		itm.SerialNumber = *dto.SerialNumber
	}

	if err := h.repository.Update(context.TODO(), itm); err != nil {
		h.logger.Errorf("Error update item %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
	}
	ctx.JSON(http.StatusOK, itm)
}
