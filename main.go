package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AdityaHarindar/Richpanel-Assignment/store"
	"github.com/AdityaHarindar/Richpanel-Assignment/transport"
)

func main() {
	// Initialize datastore, cache
	ds := store.NewStore()
	c := store.NewCache(30 * time.Second)

	r := transport.NewRouter(ds, c)

	//Auth middleware
	authorizedKey := "someSecretApiKey"
	r.Use(transport.AuthMiddleware(authorizedKey))

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Start server
	go func() {
		log.Println("Server is running on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
}
