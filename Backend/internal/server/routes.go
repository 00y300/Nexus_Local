// internal/server/routes.go
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"nexus.local/internal/db"
)

// orderLine is a single line‐item in the client’s payload.
type orderLine struct {
	ItemID   int `json:"item_id"`
	Quantity int `json:"quantity"`
}

// orderReq no longer has a UserID field.
type orderReq struct {
	Items []orderLine `json:"items"`
}

type stockUpdateReq struct {
	ItemID int `json:"item_id"`
	Stock  int `json:"stock"`
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

// POST /items/add (with image upload)
func (s *Server) addItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	desc := r.FormValue("description")
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	// 1) insert without image_url
	newID, err := db.AddItemWithImageURL(s.DB, name, desc, price, stock, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2) save uploaded image if present
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%d%s", newID, ext)
		dstPath := filepath.Join("uploads", filename)
		out, err := os.Create(dstPath)
		if err == nil {
			defer out.Close()
			io.Copy(out, file)
			imageURL := "/uploads/" + filename
			if err := db.UpdateItemImageURL(s.DB, int(newID), imageURL); err != nil {
				log.Println("failed to update image_url:", err)
			}
		}
	}

	jsonResponse(w, map[string]int64{"item_id": newID}, http.StatusCreated)
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

// ordersHandler dispatches GET, POST, DELETE on /orders
func (s *Server) ordersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Get("order_id") != "" {
			s.getOrderHandler(w, r)
		} else {
			s.getOrdersHandler(w, r)
		}
	case http.MethodPost:
		s.placeOrderHandler(w, r)
	case http.MethodDelete:
		s.deleteOrderHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// extractUserID verifies the id_token cookie and returns the Azure OID.
func (s *Server) extractUserID(r *http.Request) (string, error) {
	ck, err := r.Cookie("id_token")
	if err != nil {
		return "", err
	}
	idToken, err := s.AuthApp.Verifier.Verify(context.Background(), ck.Value)
	if err != nil {
		return "", err
	}
	var claims struct {
		OID string `json:"oid"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return "", err
	}
	return claims.OID, nil
}

// GET /orders — only the logged‑in user’s orders
func (s *Server) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := s.extractUserID(r)
	if err != nil {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	orders, err := db.GetOrdersByUser(s.DB, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, orders, http.StatusOK)
}

// GET /orders?order_id=123 — only if it belongs to the user
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
	userID, err := s.extractUserID(r)
	if err != nil {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	if order.UserID != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	jsonResponse(w, map[string]any{
		"order":       order,
		"order_items": lines,
	}, http.StatusOK)
}

// POST /orders — place under the extracted userID
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
	userID, err := s.extractUserID(r)
	if err != nil {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	orderMap := make(map[int]int)
	for _, line := range req.Items {
		if line.Quantity <= 0 {
			http.Error(w, "quantity must be > 0", http.StatusBadRequest)
			return
		}
		orderMap[line.ItemID] += line.Quantity
	}
	for itemID := range orderMap {
		if _, err := db.GetItem(s.DB, itemID); err != nil {
			http.Error(w, fmt.Sprintf("item %d not found", itemID), http.StatusBadRequest)
			return
		}
	}
	orderID, err := db.PlaceOrder(s.DB, userID, orderMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, map[string]int64{"order_id": orderID}, http.StatusCreated)
}

// DELETE /orders?order_id=123 — only if it belongs to the user
func (s *Server) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("order_id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid order_id", http.StatusBadRequest)
		return
	}
	userID, err := s.extractUserID(r)
	if err != nil {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	order, _, err := db.GetOrderByID(s.DB, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if order.UserID != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := db.DeleteOrder(s.DB, orderID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// jsonResponse is a helper for writing JSON + status code.
func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
