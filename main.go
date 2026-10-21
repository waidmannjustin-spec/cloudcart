package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Product beschreibt die Datenstruktur für ein Produkt
type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PriceCents int    `json:"price_cents"`
	Stock      int    `json:"stock"`
}

func main() {
	// Factor XI (Logs): Strukturierter JSON-Logger auf stdout
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Statische Beispieldaten (Zustandslos bezüglich lokaler Schreibzugriffe)
	products := []Product{
		{ID: 1, Name: "Kopfhoerer", PriceCents: 7999, Stock: 10},
		{ID: 2, Name: "Tastatur", PriceCents: 4999, Stock: 15},
		{ID: 3, Name: "Maus", PriceCents: 2999, Stock: 20},
	}

	mux := http.NewServeMux()

	// 1. Root Route
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintln(w, "<h1>CloudCart</h1><p>Unsere erste Go-API läuft cloud-native!</p><a href='/products'>Produkte anzeigen</a>")
	})

	// 2. Health-Check Route (wichtig für Kubernetes Liveness/Readiness Probes)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	// 3. REST-API Endpunkt für Produkte (Faktor VII)
	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Abruf der Produktliste", "client_ip", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	// Factor III (Config): Lesen des Ports aus Umgebungsvariablen
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Server in eigener Goroutine starten, damit der Hauptprozess auf Signale lauschen kann
	go func() {
		logger.Info("CloudCart API startet", "port", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Unerwarteter Serverfehler", "error", err)
			os.Exit(1)
		}
	}()

	// Factor IX (Disposability): Graceful Shutdown abfangen (SIGTERM / SIGINT)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Warten auf Beendigungssignal von Docker/Kubernetes
	logger.Info("Herunterfahren der Anwendung eingeleitet...")

	// Dem Server 5 Sekunden Zeit geben, um laufende Anfragen zu beenden
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Fehler beim sauberen Stoppen des Servers", "error", err)
	} else {
		logger.Info("Server erfolgreich und sauber heruntergefahren")
	}
}