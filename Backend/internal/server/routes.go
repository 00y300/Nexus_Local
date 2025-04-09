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

func (s *Server) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

func (s *Server) placeOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req orderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	orderMap := make(map[int]int)
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
