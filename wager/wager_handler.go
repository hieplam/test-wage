package wager

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"test-wage/domain"

	"github.com/gin-gonic/gin"
)

// Handler handler routing of wager
type Handler struct {
	service domain.WagerService
}

func NewWagerHandler(e *gin.Engine, service domain.WagerService) {
	handler := &Handler{
		service: service,
	}
	e.GET("/wagers", handler.ListWagers)
	e.POST("/wagers", handler.PlaceWager)
	e.POST("buy/:wager_id", handler.BuyWager)
}
func (h *Handler) ListWagers(c *gin.Context) {
	pageInfo := GetPagingInfo(c)
	if pageInfo.OrderBy != "asc" && pageInfo.OrderBy != "desc" {
		pageInfo.OrderBy = "asc"
	}

	wagers, err := h.service.ListWagers(pageInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, wagers)

}
func GetPagingInfo(c *gin.Context) domain.PageInfo {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Printf("err when getting page from context, set page to 1, err: %+v", err)
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		log.Printf("err when getting limit from context, set limit to 10, err: %+v", err)
		limit = 10
	}

	sortBy := strings.ToLower(c.DefaultQuery("sort_by", "id"))
	orderBy := strings.ToLower(c.DefaultQuery("order_by", "asc"))

	return domain.PageInfo{Page: page, Limit: limit, SortBy: sortBy, OrderBy: orderBy}
}
func (h *Handler) PlaceWager(c *gin.Context) {
	var req domain.PlaceWagerReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	wager, err := h.service.PlaceWager(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, wager)
}

func (h *Handler) BuyWager(c *gin.Context) {
	var req domain.BuyWagerReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	idStr := c.Param("wager_id")
	wagerID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	po, err := h.service.BuyWager(uint(wagerID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, po)
}
