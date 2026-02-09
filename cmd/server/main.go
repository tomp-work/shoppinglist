package main

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/tomp-work/shoppinglist/cmd/server/handler"
)

func main() {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.CORS("http://localhost:1323", "http://localhost:5173"))

	handler := handler.Handler{
		ItemMaxID: 3,
		Items: map[string]*handler.Item{
			"1": {
				Id:       "1",
				Name:     "bread",
				Quantity: 1,
				Picked:   true,
				SeqNum:   0,
				Price:    5,
			},
			"2": {
				Id:       "2",
				Name:     "apple",
				Quantity: 3,
				Picked:   false,
				SeqNum:   1,
				Price:    10,
			},
			"3": {
				Id:       "3",
				Name:     "orange",
				Quantity: 4,
				Picked:   false,
				SeqNum:   2,
				Price:    15,
			},
		},
	}
	// Double check ItemMaxID matches number of items in map.
	if handler.ItemMaxID != len(handler.Items) {
		panic("ItemMaxID is invalid")
	}

	// List item routing.
	e.GET("/item", handler.GetItemList)
	e.POST("/item", handler.CreateItem)
	e.DELETE("/item/:id", handler.DeleteItem)
	e.PUT("/item/:id", handler.UpdateItem)
	e.POST("/item/:id/up", handler.MoveItemUp)
	e.POST("/item/:id/down", handler.MoveItemDown)
	// List details routing.
	e.GET("/list/total", handler.CalcListTotalPrice)
	e.GET("/list", handler.GetListDetails)
	e.PUT("/list", handler.UpdateListDetails)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
