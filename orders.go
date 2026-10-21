package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Order beschreibt eine Bestellung.
type Order struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
}

// Bestellungen vorübergehend im Arbeitsspeicher speichern.
var (
	orders      = make([]Order, 0)
	ordersMu    sync.Mutex
	nextOrderID = 1
)

// REST-Endpunkte für Bestellungen registrieren.
func registerOrderRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /orders", createOrder)
	mux.HandleFunc("GET /orders", getOrders)
	mux.HandleFunc("GET /orders/{id}", getOrderByID)
}

// POST /orders
func createOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Ungueltiges JSON", http.StatusBadRequest)
		return
	}

	if input.ProductID <= 0 || input.Quantity <= 0 {
		http.Error(w, "Produkt-ID und Menge muessen groesser als 0 sein",
			http.StatusBadRequest)
		return
	}

	ordersMu.Lock()

	order := Order{
		ID:        nextOrderID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
		Status:    "PENDING",
	}

	nextOrderID++
	orders = append(orders, order)

	ordersMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GET /orders
func getOrders(w http.ResponseWriter, r *http.Request) {
	ordersMu.Lock()
	defer ordersMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GET /orders/{id}
func getOrderByID(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimSpace(r.PathValue("id"))

	id, err := strconv.Atoi(idString)
	if err != nil || id <= 0 {
		http.Error(w, "Ungueltige Bestell-ID", http.StatusBadRequest)
		return
	}

	ordersMu.Lock()
	defer ordersMu.Unlock()

	for _, order := range orders {
		if order.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(order)
			return
		}
	}

	http.Error(w, "Bestellung nicht gefunden", http.StatusNotFound)
}
