package handler

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v5"
)

type Item struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Picked   bool   `json:"picked"`
	SeqNum   int    `json:"seqnum"`
	Price    int    `json:"price"`
}

type ItemUpdate struct {
	Picked bool `json:"picked"`
}

type PriceReport struct {
	TotalPrice int `json:"totalprice"`
}

type Handler struct {
	Items map[string]*Item
}

func (h *Handler) generateID() string {
	// TODO: Consider using UUIDs.
	return fmt.Sprintf("%d", len(h.Items)+1)
}

// sortedItems returns a slice containing the items from the Items map sorted by sequence number.
func (h *Handler) sortedItems() []*Item {
	items := []*Item{}
	for _, v := range h.Items {
		items = append(items, v)
	}
	slices.SortFunc(items, func(a *Item, b *Item) int { return a.SeqNum - b.SeqNum })
	return items
}

func (h *Handler) GetItemList(c *echo.Context) error {
	return c.JSON(http.StatusOK, h.sortedItems())
}

// CreateItem will create the item, create a unique ID and set the sequence number so the item is at the end of the list.
func (h *Handler) CreateItem(c *echo.Context) error {
	item := Item{}
	if err := c.Bind(&item); err != nil {
		return fmt.Errorf("failed to Bind in CreateItem: %w", err)
	}
	item.Id = h.generateID()
	item.SeqNum = len(h.Items)
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

// UpdateItem updates the item with the ID from the path `items/:id`
func (h *Handler) UpdateItem(c *echo.Context) error {
	id := c.Param("id")
	if _, ok := h.Items[id]; !ok {
		return c.String(http.StatusNotFound, fmt.Sprintf("id (%s) not found", id))
	}
	update := ItemUpdate{}
	if err := c.Bind(&update); err != nil {
		return fmt.Errorf("failed to Bind in UpdateItem: %w", err)
	}
	h.Items[id].Picked = update.Picked
	return c.JSON(http.StatusOK, h.Items[id])
}

func (h *Handler) MoveItemUp(c *echo.Context) error {
	id := c.Param("id")
	if _, ok := h.Items[id]; !ok {
		return c.String(http.StatusNotFound, fmt.Sprintf("id (%s) not found", id))
	}
	seqNum := h.Items[id].SeqNum
	if seqNum == 0 {
		// Already at top of list.
		return c.NoContent(http.StatusOK)
	}
	// items contains a slice of all the items sorted by sequence number, which means that:
	// items[index].SeqNum == index
	items := h.sortedItems()
	// Swap the sequence numbers of our selected item and the one above (lower number is higher in the list).
	items[seqNum].SeqNum--
	items[seqNum-1].SeqNum++
	return c.NoContent(http.StatusOK)
}

func (h *Handler) MoveItemDown(c *echo.Context) error {
	id := c.Param("id")
	if _, ok := h.Items[id]; !ok {
		return c.String(http.StatusNotFound, fmt.Sprintf("id (%s) not found", id))
	}
	seqNum := h.Items[id].SeqNum
	if seqNum == len(h.Items)-1 {
		// Already at bottom of list.
		return c.NoContent(http.StatusOK)
	}
	// items contains a slice of all the items sorted by sequence number, which means that:
	// items[index].SeqNum == index
	items := h.sortedItems()
	// Swap the sequence numbers of our selected item and the one below (higher number is lower in the list).
	items[seqNum].SeqNum++
	items[seqNum+1].SeqNum--
	return c.NoContent(http.StatusOK)
}

func (h *Handler) CalcListTotalPrice(c *echo.Context) error {
	report := PriceReport{
		TotalPrice: 0,
	}
	for _, item := range h.Items {
		report.TotalPrice += item.Price
	}
	return c.JSON(http.StatusOK, &report)
}
