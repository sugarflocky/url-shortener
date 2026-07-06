package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sugarflocky/url-shortener/internal/config"
	"github.com/sugarflocky/url-shortener/internal/httpapi"
	"github.com/sugarflocky/url-shortener/internal/shortener"
	"github.com/sugarflocky/url-shortener/internal/storage/memory"
)

// shutdownTimeout limits waiting active request while stopping.
const shutdownTimeout = 5 * time.Second

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	var st shortener.Storage
	switch cfg.Storage {
	case config.StorageMemory:
		st = memory.New()
	case config.StoragePostgres:
		// st = postgres.New()
		log.Fatal("postgres storage is not implemented yet")
	}
	svc := shortener.New(st)
	h := httpapi.New(svc)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{Addr: cfg.Addr, Handler: h.Router()}

	go func() {
		log.Printf("listening on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("server stopped")

}
