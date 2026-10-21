package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Product beschreibt die Datenstruktur eines Produkts.
type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PriceCents int    `json:"price_cents"`
	Stock      int    `json:"stock"`
}

func main() {
	// Factor XI: Strukturierte Logs auf stdout
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// HTTP-Router erstellen
	mux := http.NewServeMux()

	// Produkt- und Bestellrouten registrieren
	RegisterProductRoutes(mux, logger)
	registerOrderRoutes(mux)

	// Startseite: index.html anzeigen
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, "index.html")
	})

	// Health-Check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	// Factor III: Port über Umgebungsvariable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// HTTP-Server starten
	go func() {
		logger.Info("CloudCart API startet", "port", port)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error("Serverfehler", "error", err)
			os.Exit(1)
		}
	}()

	// Factor IX: Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info("Server wird heruntergefahren")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Shutdown fehlgeschlagen", "error", err)
	} else {
		logger.Info("Server erfolgreich gestoppt")
	}
}
