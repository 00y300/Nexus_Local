// internal/db/db.go
package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url,omitempty"`
}

type Order struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItem struct {
	OrderID  int64 `json:"order_id"`
	ItemID   int   `json:"item_id"`
	Quantity int   `json:"quantity"`
}

// Connect opens & verifies a MySQL database connection.
func Connect(user, pass, host, port, name string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
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

// GetAllItems returns every item in the items table, including image_url.
func GetAllItems(db *sql.DB) ([]Item, error) {
	rows, err := db.Query(
		"SELECT id, name, description, price, stock, image_url FROM items",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var it Item
		var img sql.NullString
		if err := rows.Scan(
			&it.ID,
			&it.Name,
			&it.Description,
			&it.Price,
			&it.Stock,
			&img,
		); err != nil {
			return nil, err
		}
		if img.Valid {
			it.ImageURL = img.String
		}
		items = append(items, it)
	}
	return items, nil
}

// GetItem fetches a single item by its ID, including image_url.
func GetItem(db *sql.DB, itemID int) (*Item, error) {
	var it Item
	var img sql.NullString
	err := db.QueryRow(
		"SELECT id, name, description, price, stock, image_url FROM items WHERE id = ?",
		itemID,
	).Scan(
		&it.ID,
		&it.Name,
		&it.Description,
		&it.Price,
		&it.Stock,
		&img,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item %d not found", itemID)
	}
	if err != nil {
		return nil, err
	}
	if img.Valid {
		it.ImageURL = img.String
	}
	return &it, nil
}

// AddItem inserts a new product into the items table.
// It returns the newly created item's ID.
func AddItem(db *sql.DB, name, description string, price float64, stock int) (int64, error) {
	res, err := db.Exec(
		"INSERT INTO items (name, description, price, stock) VALUES (?, ?, ?, ?)",
		name, description, price, stock,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateItemStock sets the stock for a given item.
func UpdateItemStock(db *sql.DB, itemID, newStock int) error {
	res, err := db.Exec(
		"UPDATE items SET stock = ? WHERE id = ?",
		newStock, itemID,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no item with id %d", itemID)
	}
	return nil
}

// UpdateItem updates all modifiable fields of an item.
func UpdateItem(db *sql.DB, item Item) error {
	res, err := db.Exec(
		"UPDATE items SET name = ?, description = ?, price = ?, stock = ? WHERE id = ?",
		item.Name, item.Description, item.Price, item.Stock, item.ID,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no item with id %d", item.ID)
	}
	return nil
}

// DeleteItem removes an item from the database.
func DeleteItem(db *sql.DB, itemID int) error {
	res, err := db.Exec("DELETE FROM items WHERE id = ?", itemID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no item with id %d", itemID)
	}
	return nil
}

// PlaceOrder creates an order + order_items, and deducts stock.
func PlaceOrder(db *sql.DB, userID string, orderItems map[int]int) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// 1) create the order header
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

	// 3) commit transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return orderID, nil
}

// GetAllOrders returns every order header.
func GetAllOrders(db *sql.DB) ([]Order, error) {
	rows, err := db.Query("SELECT id, user_id, created_at FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// GetOrderByID fetches one order and its line‐items.
func GetOrderByID(db *sql.DB, orderID int64) (*Order, []OrderItem, error) {
	var o Order
	err := db.QueryRow(
		"SELECT id, user_id, created_at FROM orders WHERE id = ?",
		orderID,
	).Scan(&o.ID, &o.UserID, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil, fmt.Errorf("order %d not found", orderID)
	}
	if err != nil {
		return nil, nil, err
	}

	rows, err := db.Query(
		"SELECT order_id, item_id, quantity FROM order_items WHERE order_id = ?",
		orderID,
	)
	if err != nil {
		return &o, nil, err
	}
	defer rows.Close()

	var lines []OrderItem
	for rows.Next() {
		var li OrderItem
		if err := rows.Scan(&li.OrderID, &li.ItemID, &li.Quantity); err != nil {
			return &o, nil, err
		}
		lines = append(lines, li)
	}
	return &o, lines, nil
}

// DeleteOrder deletes an order (and cascades to order_items).
func DeleteOrder(db *sql.DB, orderID int64) error {
	res, err := db.Exec("DELETE FROM orders WHERE id = ?", orderID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no order with id %d", orderID)
	}
	return nil
}

// AddItemWithImageURL inserts a new item and allows setting image_url.
func AddItemWithImageURL(db *sql.DB, name, desc string, price float64, stock int, imageURL string) (int64, error) {
	res, err := db.Exec(
		"INSERT INTO items (name, description, price, stock, image_url) VALUES (?, ?, ?, ?, ?)",
		name, desc, price, stock, imageURL,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateItemImageURL updates only the image_url column for an item.
func UpdateItemImageURL(db *sql.DB, itemID int, imageURL string) error {
	_, err := db.Exec(
		"UPDATE items SET image_url = ? WHERE id = ?",
		imageURL, itemID,
	)
	return err
}
