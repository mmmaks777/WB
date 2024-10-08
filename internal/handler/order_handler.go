package handler

import (
	"net/http"
	"sync"
	"wb_l0/internal/cache"
	"wb_l0/internal/domain"

	"github.com/gin-gonic/gin"
)

type ServiceHandler interface {
	GetOrder(ctx *gin.Context, cache map[string]domain.Order, mu *sync.RWMutex) gin.HandlerFunc
}

type OrderHandler struct{}

func (h *OrderHandler) GetOrder(cache *cache.OrderCache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		order, ok := cache.Get(id)
		if ok {
			ctx.HTML(http.StatusOK, "model.html", gin.H{"data": order})
		} else {
			ctx.HTML(http.StatusNotFound, "model.html", gin.H{
				"error":   http.StatusNotFound,
				"errMsg":  "Order not found",
				"orderId": id,
			})
		}
	}
}
