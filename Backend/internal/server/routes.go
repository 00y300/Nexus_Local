package server

import (
	"encoding/json"
	"net/http"

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// POST /orders
// { "user_id":1, "items":[{"item_id":1,"quantity":2},…] }
func (s *Server) placeOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req orderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	orderMap := make(map[int]int, len(req.Items))
	for _, line := range req.Items {
		orderMap[line.ItemID] = line.Quantity
	}

	orderID, err := db.PlaceOrder(s.DB, req.UserID, orderMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"order_id": orderID})
}

// POST /items/update
// { "item_id":1, "stock":42 }
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

	// Return updated list
	items, err := db.GetAllItems(s.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// POST /items/add
// { "name":"Foo","description":"Bar","price":9.99,"stock":100 }
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

	id, err := db.AddItem(s.DB, req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"item_id": id})
}
