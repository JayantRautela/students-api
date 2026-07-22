package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JayantRautela/students-api/internal/config"
)

// main entry point
func main() {
	// load config
	cfg := config.MustLoad()

	// router setup
	router := http.NewServeMux() // gives a new Server Mux (router)

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api"))
	})

	// server setup
	server := http.Server {
		Addr: cfg.Addr,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// go routine
	go func()  {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("Failed to start server")
		}
	} ()

	<- done

	// graceful shutdown
	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}