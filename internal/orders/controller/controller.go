package controller

import (
	"net/http"

	"github.com/Jhnvlglmlbrt/wb-order/internal/cache"
	"github.com/labstack/echo/v4"
)

type CacheHandler struct {
	c *cache.Cache
}

func NewController(c *cache.Cache) *CacheHandler {
	return &CacheHandler{
		c: c,
	}
}

// func (ch *CacheHandler) GetOrder(ctx echo.Context) error {
// 	order := ch.c.GetOrderByUid(ctx.Param("order"))

// 	or, err := json.MarshalIndent(order, "", "\t")
// 	if err != nil {
// 		return fmt.Errorf("error at marshaling respond: %v", err)
// 	}

// 	return ctx.JSONBlob(http.StatusOK, or)
// }

// func (ch *CacheHandler) GetAllOrders(ctx echo.Context) error {
// 	order := ch.c.GetOrders()

// 	or, err := json.MarshalIndent(order, "", "\t")
// 	if err != nil {
// 		return fmt.Errorf("error at marshaling respond: %v", err)
// 	}

// 	return ctx.JSONBlob(http.StatusOK, or)
// }

func (ch *CacheHandler) HandlePage(ctx echo.Context) error {
	return ctx.File("static/order.html")
}

func (ch *CacheHandler) HandleGetData(ctx echo.Context) error {
	id := ctx.QueryParam("id")

	if id != "" {
		order := ch.c.GetOrderByUid(id)

		return ctx.JSON(http.StatusOK, order)
	}

	orders := ch.c.GetOrders()

	return ctx.JSON(http.StatusOK, orders)
}
