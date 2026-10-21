
package main

import (
    "encoding/json"
    "log/slog"
    "net/http"
)

func RegisterProductRoutes(mux *http.ServeMux, logger *slog.Logger) {

    mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {

        logger.Info("Produktliste wird abgerufen")

        products := []Product{
            {ID: 1, Name: "Cloud Laptop", PriceCents: 99900, Stock: 10},
            {ID: 2, Name: "Kubernetes Buch", PriceCents: 3999, Stock: 20},
            {ID: 3, Name: "Docker T-Shirt", PriceCents: 1995, Stock: 30},
        }

        w.Header().Set("Content-Type", "application/json")

        if err := json.NewEncoder(w).Encode(products); err != nil {
            logger.Error("Fehler bei JSON-Ausgabe", "error", err)
        }
    })
}
