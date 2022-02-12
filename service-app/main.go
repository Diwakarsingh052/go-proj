package main

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service-app/auth"
	"service-app/data/user"
	"service-app/database"
	"service-app/handlers"
	"syscall"
	"time"
)

func main() {
	l := log.New(os.Stdout, "Users : ", log.LstdFlags)
	startApp(l)
}

func startApp(log *log.Logger) error {
	// =========================================================================
	// Start Database
	db, err := database.Open()
	if err != nil {
		return fmt.Errorf("connecting to db %w", err)
	}
	uDB := &user.DbService{DB: db}

	defer func() {
		log.Printf("main: Database Stopping : ")
		db.Close()
	}()

	// =========================================================================
	// Initialize authentication support
	log.Println("main : Started : Initializing authentication support")

	privatePEM, err := os.ReadFile("private.pem")
	if err != nil {
		return fmt.Errorf("reading auth private key %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("parsing auth private key %w", err)
	}
	a, err := auth.NewAuth(privateKey, "RS256")

	if err != nil {
		return fmt.Errorf("constructing auth %w", err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := http.Server{
		Addr:         ":8080",
		Handler:      handlers.API(shutdown, log, a, uDB),
		ReadTimeout:  8000 * time.Second,
		WriteTimeout: 800 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error %w", err)
	case sig := <-shutdown:
		log.Printf("main: %v : Start shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil { // first trying to cleanly shutdown
			api.Close() // forcing shutdown
			return fmt.Errorf("could not stop server gracefully %w", err)
		}
	}
	return nil

}
