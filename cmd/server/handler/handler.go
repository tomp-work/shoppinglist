package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Item struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type Handler struct {
	Items map[string]*Item
}

func (h *Handler) GetItemList(c *echo.Context) error {
	items := []*Item{}
	for _, v := range h.Items {
		items = append(items, v)
	}
	return c.JSON(http.StatusOK, items)
}
