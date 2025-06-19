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
	"github.com/amar2502/students-api/internal/http/handlers/student"
	"github.com/amar2502/students-api/internal/storage/sqlite"
)

func main() {

	// load config
	cfg := config.MustLoad()

	// database setup
	storage, er := sqlite.New(cfg)
	if er != nil {
		log.Fatal("cannot connect to database: ", er.Error())
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env))


	// setup router
	router := http.NewServeMux()


	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetbyId(storage))
	router.HandleFunc("GET /api/students", student.GetStudent(storage))


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