package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gostatus/internal/gateway"
	"gostatus/internal/handler"
	"gostatus/internal/store"

	"github.com/BurntSushi/toml"
)

// Config holds the runtime configuration loaded from config.toml.
type Config struct {
	Token string
}

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	healthCheck := flag.Bool("health", false, "Run internal healthcheck")
	configPath := flag.String("config", "config.toml", "Path to config.toml")
	flag.Parse()

	if *healthCheck {
		if err := checkHealth(*port); err != nil {
			log.Println("health check failed:", err)
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

	srv := &http.Server{
		Addr:              ":" + *port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Printf("Listening on :%s", *port)
	log.Fatal(srv.ListenAndServe())
}

func checkHealth(port string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://localhost:%s/badge/status/0", port), http.NoBody)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %s", resp.Status)
	}
	return nil
}
