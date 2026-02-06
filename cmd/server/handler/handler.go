package handler

import (
	"fmt"
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

func (h *Handler) generateID() string {
	// TODO: Consider using UUIDs.
	return fmt.Sprintf("%d", len(h.Items)+1)
}

func (h *Handler) GetItemList(c *echo.Context) error {
	items := []*Item{}
	for _, v := range h.Items {
		items = append(items, v)
	}
	return c.JSON(http.StatusOK, items)
}

func (h *Handler) CreateItem(c *echo.Context) error {
	item := Item{}
	if err := c.Bind(&item); err != nil {
		return fmt.Errorf("failed to Bind in CreateItem: %w", err)
	}
	item.Id = h.generateID()
	h.Items[item.Id] = &item
	return c.JSON(http.StatusCreated, &item)
}

// DeleteItem deletes the item with the given ID from the path `items/:id`
func (h *Handler) DeleteItem(c *echo.Context) error {
	id := c.Param("id")
	if _, ok := h.Items[id]; !ok {
		return c.String(http.StatusNotFound, fmt.Sprintf("id (%s) not found", id))
	}
	delete(h.Items, id)
	return c.NoContent(http.StatusOK)
}
