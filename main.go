package main

import (
	"log"
	"net/http"
	"os"

	"gostatus/internal/gateway"
	"gostatus/internal/handler"
	"gostatus/internal/store"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN env var is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := store.New()
	go gateway.Connect(token, s)

	mux := http.NewServeMux()
	handler.New(s).Register(mux)

	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
