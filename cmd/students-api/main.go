package main

import (
	"fmt"
	"log"
	"net/http"

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

	fmt.Printf("server started %s", cfg.HTTpServer.Addr)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Failed to start server")
	}
}