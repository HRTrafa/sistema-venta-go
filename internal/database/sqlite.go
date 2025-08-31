package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // <--- CAMBIA ESTO
	"log"
)

var DB *sql.DB

func InitDB(dbPath string) {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		name TEXT,
		quantity INTEGER,
		price REAL
	);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS sales (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		client TEXT,
		product_id INTEGER,
		quantity INTEGER,
		price REAL,
		total REAL,
		status TEXT
	);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS cash_deliveries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		name TEXT,
		description TEXT,
		amount REAL
	);`)
	if err != nil {
		log.Fatal(err)
	}
}