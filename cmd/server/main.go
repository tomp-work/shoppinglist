package main

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/tomp-work/shoppinglist/cmd/server/handler"
)

func main() {
	e := echo.New()
	e.Use(middleware.RequestLogger())

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

	e.GET("/item", handler.GetItemList)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
