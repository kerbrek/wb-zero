package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Wait for db connection in a loop for the given timeout
func connectLoop(driver string, DSN string, timeout time.Duration) (*sql.DB, error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("connectLoop: db connection failed after %s timeout", timeout)

		case <-ticker.C:
			database, _ := sql.Open("postgres", DSN)
			err := database.Ping()
			if err == nil {
				return database, nil
			}
		}
	}
}

func MakeDB(connStr string) (*sql.DB, error) {
	database, err := connectLoop("postgres", connStr, 15*time.Second)
	if err != nil {
		return nil, fmt.Errorf("MakeDB: %w", err)
	}

	db = database
	return database, nil
}

func SaveOrder(id string, orderJson []byte) error {
	_, err := db.Exec("INSERT INTO customer_order (id, json_data) VALUES ($1, $2)", id, orderJson)
	if err != nil {
		return fmt.Errorf("SaveOrder: %w", err)
	}

	return nil
}

func ReadAllOrders() (map[string][]byte, error) {
	rows, err := db.Query("SELECT id, json_data FROM customer_order")
	if err != nil {
		return nil, fmt.Errorf("ReadAllOrders: %w", err)
	}
	defer rows.Close()

	orderJsons := make(map[string][]byte)

	for rows.Next() {
		var id string
		var orderJson []byte
		if err := rows.Scan(&id, &orderJson); err != nil {
			return nil, fmt.Errorf("ReadAllOrders: %w", err)
		}

		orderJsons[id] = orderJson
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ReadAllOrders: %w", err)
	}

	return orderJsons, nil
}
