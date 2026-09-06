package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gostatus/internal/config"
	"gostatus/internal/gateway"
	"gostatus/internal/handler"
	"gostatus/internal/store"
)

func main() {
	cfg := config.Load()

	if cfg.Healthcheck {
		if err := checkHealth(cfg.Port); err != nil {
			log.Println("health check failed:", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	s := store.New()
	go gateway.Connect(cfg.Token, s)

	mux := http.NewServeMux()
	handler.New(s).Register(mux)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Printf("Listening on :%s", cfg.Port)
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
