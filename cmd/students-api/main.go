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
	"github.com/Rajit-Dutta/StudentAPI/internal/http/handler/students"
	"github.com/Rajit-Dutta/StudentAPI/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal("Error connecting DB")
	}
	slog.Info("storage initialised", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	router := http.NewServeMux()

	router.HandleFunc("POST /", students.New(storage))

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
