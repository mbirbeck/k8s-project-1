package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Configure an HTTP server to listen on port 8080.
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/health", healthHandler)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	// Use an error channel to capture server errors.
	srvErrCh := make(chan error, 1)

	// Create a context to capture shutdown signals.
	shutdownTriggeredCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Start the HTTP server.
	go func() {
		log.Printf("Server is listening on %s\n", srv.Addr)
		srvErrCh <- srv.ListenAndServe()
	}()

	// Wait for a shutdown signal, or a server error.
	select {
	case <-shutdownTriggeredCtx.Done():
		log.Printf("Shutdown signal received")
		// Create a new 10 second timeout context for graceful shutdown.
		gracefulShutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		// Allow the server to shutdown gracefully within the timeout.
		if err := srv.Shutdown(gracefulShutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			if err := srv.Close(); err != nil {
				log.Fatalf("Forced close failed: %v", err)
			}
		}
		log.Printf("Server stopped")
	case err := <-srvErrCh:
		if !errors.Is(err, http.ErrServerClosed) {
			// Ensure a server startup or runtime failure doesn't go unnoticed.
			log.Fatalf("Server failed: %v", err)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}
	log.Printf("Received request for %s\n", name)
	fmt.Fprintf(w, "Hello, %s! v0.0.40\n", name)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
