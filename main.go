package main

import (
	"log"
	"net/http"

	"gostatus/internal/gateway"
	"gostatus/internal/handler"
	"gostatus/internal/store"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Token string
	Port  string
}

func main() {
	var conf Config

	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Fatal(err)
	}

	token := conf.Token
	if token == "" || token == "YOUR_TOKEN" {
		log.Fatal("Bot token is required in config.toml")
	}

	port := conf.Port
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
