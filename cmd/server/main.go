package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/tomp-work/shoppinglist/cmd/server/emails"
	"github.com/tomp-work/shoppinglist/cmd/server/handler"
	"github.com/tomp-work/shoppinglist/cmd/server/models"
)

func main() {
	// Read config.
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env.local file")
	}
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the Resend API key from the environment and add to sender implementation.
	emailSender := emails.ResendSender{
		ApiKey: os.Getenv("RESEND_API_KEY"),
	}

	handler := handler.Handler{
		EmailSender: &emailSender,
		ItemMaxID:   3,
		Items: map[string]*models.Item{
			"1": {
				Id:     "1",
				Name:   "bread",
				Picked: true,
				SeqNum: 0,
				Price:  5,
			},
			"2": {
				Id:     "2",
				Name:   "red wine",
				Picked: false,
				SeqNum: 1,
				Price:  10,
			},
			"3": {
				Id:     "3",
				Name:   "cheese",
				Picked: false,
				SeqNum: 2,
				Price:  15,
			},
		},
		ListDetails: models.ListDetails{
			TotalPrice:    30,
			SpendingLimit: 100,
		},
	}
	// Double check ItemMaxID matches number of items in map.
	if handler.ItemMaxID != len(handler.Items) {
		panic("ItemMaxID is invalid")
	}
	// Double check totalPrice equals sum of item prices.
	expectedTotalPrice := 0
	for _, item := range handler.Items {
		expectedTotalPrice += item.Price
	}
	if handler.ListDetails.TotalPrice != expectedTotalPrice {
		panic("TotalPrice is incorrect")
	}

	// Setup HTTP server.
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.CORS("http://localhost:1323", "http://localhost:5173"))

	// List item routing.
	e.GET("/item", handler.GetItemList)
	e.POST("/item", handler.CreateItem)
	e.DELETE("/item/:id", handler.DeleteItem)
	e.PUT("/item/:id", handler.UpdateItem)
	e.POST("/item/:id/up", handler.MoveItemUp)
	e.POST("/item/:id/down", handler.MoveItemDown)
	// List details routing.
	e.GET("/list", handler.GetListDetails)
	e.PUT("/list", handler.UpdateListDetails)
	e.POST("/list/send", handler.SendListEmail)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
