package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Item is a product in the shop.
type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// Connect opens & pings a MySQL database.
func Connect(user, pass, host, port, name string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, pass, host, port, name,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// GetAllItems returns every item in the `items` table.
func GetAllItems(db *sql.DB) ([]Item, error) {
	rows, err := db.Query(
		"SELECT id, name, description, price, stock FROM items",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(
			&it.ID, &it.Name, &it.Description, &it.Price, &it.Stock,
		); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}

// PlaceOrder creates an order + order_items, and deducts stock.
// orderItems maps itemID→quantity.
func PlaceOrder(db *sql.DB, userID int, orderItems map[int]int) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// 1) create the order
	res, err := tx.Exec(
		"INSERT INTO orders (user_id, created_at) VALUES (?, ?)",
		userID, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	orderID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// 2) insert line‐items & decrement stock
	for itemID, qty := range orderItems {
		if _, err := tx.Exec(
			"INSERT INTO order_items (order_id, item_id, quantity) VALUES (?, ?, ?)",
			orderID, itemID, qty,
		); err != nil {
			tx.Rollback()
			return 0, err
		}
		if _, err := tx.Exec(
			"UPDATE items SET stock = stock - ? WHERE id = ?",
			qty, itemID,
		); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return orderID, nil
}
