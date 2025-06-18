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

	"github.com/amar2502/students-api/internal/config"
)

func main() {

	// load config
	cfg := config.MustLoad()

	// database setup


	// setup router
	router := http.NewServeMux()


	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Students API"))
	})


	// setup server
	server := http.Server {
		Addr: cfg.HttpServer.Addr,
		Handler: router,
	}

	slog.Info("server started ", slog.String("address", cfg.HttpServer.Addr))
	fmt.Println("server starting on ", cfg.HttpServer.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

	if err != nil { 
		log.Fatal("cannot start server: ", err.Error())
	}
	} ()

	<- done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")

}