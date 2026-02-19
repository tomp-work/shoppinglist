package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v2"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env.local file")
	}
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{"tomp.work.email@gmail.com"},
		Subject: "Hello World",
		Html:    "<p>Congrats on sending your <strong>first email</strong>!</p>",
	}

	_, err = client.Emails.Send(params)
	if err != nil {
		panic(err)
	}
}
