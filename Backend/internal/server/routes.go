package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"nexus.local/internal/db"
)

type orderLine struct {
	ItemID   int `json:"item_id"`
	Quantity int `json:"quantity"`
}

type orderReq struct {
	UserID int         `json:"user_id"`
	Items  []orderLine `json:"items"`
}

type stockUpdateReq struct {
	ItemID int `json:"item_id"`
	Stock  int `json:"stock"`
}

type addItemReq struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// GET /items
func (s *Server) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	items, err := db.GetAllItems(s.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, items, http.StatusOK)
}

// POST /items/add
func (s *Server) addItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req addItemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Prevent duplicate by name (case‐insensitive)
	all, err := db.GetAllItems(s.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, it := range all {
		if strings.EqualFold(it.Name, req.Name) {
			http.Error(w, "item already exists", http.StatusConflict)
			return
		}
	}

	id, err := db.AddItem(s.DB, req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, map[string]int64{"item_id": id}, http.StatusCreated)
}

// POST /items/update
func (s *Server) updateStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req stockUpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := db.UpdateItemStock(s.DB, req.ItemID, req.Stock); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	items, err := db.GetAllItems(s.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, items, http.StatusOK)
}

// ----------------------------------------------------------------
// Orders: all operations live under the single /orders path.
// ----------------------------------------------------------------

// ordersHandler dispatches GET, POST, DELETE on /orders
func (s *Server) ordersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Get("order_id") != "" {
			s.getOrderHandler(w, r)
		} else {
			s.getOrdersHandler(w)
		}
	case http.MethodPost:
		s.placeOrderHandler(w, r)
	case http.MethodDelete:
		s.deleteOrderHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /orders
func (s *Server) getOrdersHandler(w http.ResponseWriter) {
	orders, err := db.GetAllOrders(s.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, orders, http.StatusOK)
}

// GET /orders?order_id=123
func (s *Server) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("order_id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid order_id", http.StatusBadRequest)
		return
	}
	order, lines, err := db.GetOrderByID(s.DB, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, map[string]any{
		"order":       order,
		"order_items": lines,
	}, http.StatusOK)
}

// POST /orders
func (s *Server) placeOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req orderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Accumulate quantities for duplicate itemIDs
	orderMap := make(map[int]int)
	for _, line := range req.Items {
		if line.Quantity <= 0 {
			http.Error(w, "quantity must be > 0", http.StatusBadRequest)
			return
		}
		orderMap[line.ItemID] += line.Quantity
	}

	// Verify each item exists
	for itemID := range orderMap {
		if _, err := db.GetItem(s.DB, itemID); err != nil {
			http.Error(w, fmt.Sprintf("item %d not found", itemID), http.StatusBadRequest)
			return
		}
	}

	orderID, err := db.PlaceOrder(s.DB, req.UserID, orderMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, map[string]int64{"order_id": orderID}, http.StatusCreated)
}

// DELETE /orders?order_id=123
func (s *Server) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("order_id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid order_id", http.StatusBadRequest)
		return
	}
	if err := db.DeleteOrder(s.DB, orderID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// helper to write JSON + status
func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
