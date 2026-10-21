package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

// Product beschreibt die Datenstruktur für ein Produkt
type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Handhabung von Produkten gemäß den 12-Factor-Prinzipien
type ProductHandler struct {
	logger *slog.Logger
}

func NewProductHandler(logger *slog.Logger) *ProductHandler {
	return &ProductHandler{
		logger: logger,
	}
}

// GetProducts verarbeitet synchrone REST-Anfragen (Faktor VII - Port Binding / REST API)
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	// Factor XI (Logs): Strukturierte Ausgabe als Event-Stream nach stdout
	h.logger.Info("Empfange GET-Anfrage fuer Produktkatalog", "method", r.Method, "path", r.URL.Path)

	// In einer vollen Cloud-Architektur würden die Daten aus einer externen DB geladen (Faktor VI: Statelessness).
	// Hier simulieren wir den Produktkatalog:
	products := []Product{
		{ID: "p1", Name: "Cloud Laptop", Price: 999.00},
		{ID: "p2", Name: "Kubernetes Buch", Price: 39.99},
		{ID: "p3", Name: "Docker T-Shirt", Price: 19.95},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		h.logger.Error("Fehler beim Serialisieren der Produkte", "error", err)
	}
}

// InitProductRoutes registriert die Endpunkte
func RegisterProductRoutes(mux *http.ServeMux, logger *slog.Logger) {
	handler := NewProductHandler(logger)
	mux.HandleFunc("/api/products", handler.GetProducts)
}
