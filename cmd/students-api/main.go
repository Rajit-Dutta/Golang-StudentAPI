package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rajit-Dutta/StudentAPI/internal/config"
)

func main() {
	cfg := config.MustLoad()
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api"))
	})

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	fmt.Println("Server started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("Cannot connect to the server")
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown unsuccessful", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown successful")
}
