package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PriceCents int    `json:"price_cents"`
	Stock      int    `json:"stock"`
}

func main() {
	products := []Product{
		{ID: 1, Name: "Kopfhoerer", PriceCents: 7999, Stock: 10},
		{ID: 2, Name: "Tastatur", PriceCents: 4999, Stock: 15},
		{ID: 3, Name: "Maus", PriceCents: 2999, Stock: 20},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintln(w, "<h1>CloudCart</h1><p>Unsere erste Go-API läuft!</p><a href='/products'>Produkte anzeigen</a>")
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("CloudCart API startet auf Port", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
