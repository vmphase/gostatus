package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gostatus/internal/gateway"
	"gostatus/internal/handler"
	"gostatus/internal/store"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Token string
}

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	healthCheck := flag.Bool("health", false, "Run internal healthcheck")
	configPath := flag.String("config", "config.toml", "Path to config.toml")
	flag.Parse()

	if *healthCheck {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/badge/status/0", *port))
		if err != nil || resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	}

	var conf Config
	if _, err := toml.DecodeFile(*configPath, &conf); err != nil {
		log.Fatalf("config: %v", err)
	}

	if conf.Token == "" || conf.Token == "YOUR_TOKEN" {
		log.Fatal("Bot token is required in config.toml")
	}

	s := store.New()
	go gateway.Connect(conf.Token, s)

	mux := http.NewServeMux()
	handler.New(s).Register(mux)

	log.Printf("Listening on :%s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, mux))
}
