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

	"github.com/amanguptak/students-api/internal/config"
	"github.com/amanguptak/students-api/internal/http/handlers/student"
	"github.com/amanguptak/students-api/internal/storage/sqlite"
)

func main() {
	// fmt.Println("Welcome to students api")
	//load config
	cfg := config.MustLoad()
	// fmt.Println(cfg,"hello")

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	//data base setup

	_, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env))
	// setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	fmt.Printf("server started at port  %s", cfg.HTTPServer.Addr)

	//channel
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// setup server
	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("failed to start server", err.Error())
		}
	}()

	<-done

	//structured log

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// giving time for shutting down in 5 second
	err = server.Shutdown(ctx) // shuting down server

	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
