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
		Items: map[string]*handler.Item{
			"1": {
				Id:       "1",
				Name:     "bread",
				Quantity: 1,
			},
			"2": {
				Id:       "2",
				Name:     "apple",
				Quantity: 3,
			},
			"3": {
				Id:       "3",
				Name:     "orange",
				Quantity: 4,
			},
		},
	}

	// Routing.
	e.GET("/item", handler.GetItemList)
	e.POST("/item", handler.CreateItem)
	e.DELETE("/item/:id", handler.DeleteItem)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
