package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v5"
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

	h := handler.Handler{
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
	if h.ItemMaxID != len(h.Items) {
		panic("ItemMaxID is invalid")
	}
	// Double check totalPrice equals sum of item prices.
	expectedTotalPrice := 0
	for _, item := range h.Items {
		expectedTotalPrice += item.Price
	}
	if h.ListDetails.TotalPrice != expectedTotalPrice {
		panic("TotalPrice is incorrect")
	}

	// Setup HTTP server.
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.CORS("http://localhost:1323", "http://localhost:5173"))

	// Login routing.
	e.POST("/login", h.Login)

	api := e.Group("/api")
	api.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: handler.AccessSecret,
	}))
	api.GET("/test", h.TestAuth)

	// List item routing.
	e.GET("/item", h.GetItemList)
	e.POST("/item", h.CreateItem)
	e.DELETE("/item/:id", h.DeleteItem)
	e.PUT("/item/:id", h.UpdateItem)
	e.POST("/item/:id/up", h.MoveItemUp)
	e.POST("/item/:id/down", h.MoveItemDown)

	// List details routing.
	e.GET("/list", h.GetListDetails)
	e.PUT("/list", h.UpdateListDetails)
	e.POST("/list/send", h.SendListEmail)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
