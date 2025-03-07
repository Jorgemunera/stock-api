package database

import (
	"database/sql"
	"fmt"
	"stock-api/models"
	"time"

	"github.com/go-resty/resty/v2"
	_ "github.com/lib/pq"
)

// Conectar a CockroachDB con reintentos
func ConnectDB() (*sql.DB, error) {
	connStr := "postgresql://root@localhost:26257/?sslmode=disable"
	var db *sql.DB
	var err error

	// Intentar conectarse varias veces
	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			// Verificar si la base de datos está lista
			err = db.Ping()
			if err == nil {
				break
			}
		}

		fmt.Printf("Attempt %d: Error connecting to database: %v. Retrying in %v...\n", i+1, err, retryDelay)
		time.Sleep(retryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("error connecting to database after %d attempts: %w", maxRetries, err)
	}

	// Crear la base de datos si no existe
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS stockdb")
	if err != nil {
		return nil, fmt.Errorf("error creating database: %w", err)
	}

	// Conéctate a la base de datos stockdb
	connStr = "postgresql://root@localhost:26257/stockdb?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}

// Crear la tabla en CockroachDB
func CreateTable(db *sql.DB) error {
	fmt.Println("Creating table 'stocks'...")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS stocks (
			ticker TEXT PRIMARY KEY,
			company TEXT,
			brokerage TEXT,
			action TEXT,
			rating_from TEXT,
			rating_to TEXT,
			target_from TEXT,
			target_to TEXT
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	fmt.Println("Table 'stocks' created successfully!")
	return nil
}

// Insertar los datos iniciales en la tabla
func InsertStocks(db *sql.DB) error {
	fmt.Println("Inserting initial stocks into the database...")

	// Obtener los datos de las acciones desde la API
	stocks, err := fetchStocks()
	if err != nil {
		return fmt.Errorf("error fetching stocks: %w", err)
	}

	// Insertar cada acción en la tabla
	for _, stock := range stocks {
		err := InsertStock(db, stock)
		if err != nil {
			return fmt.Errorf("error inserting stock %s: %w", stock.Ticker, err)
		}
	}

	fmt.Println("Stocks inserted successfully!")
	return nil
}

// InsertStock inserta una acción en la tabla, evitando duplicados
func InsertStock(db *sql.DB, stock models.Stock) error {
	_, err := db.Exec(`
		INSERT INTO stocks (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (ticker) DO NOTHING;
	`, stock.Ticker, stock.Company, stock.Brokerage, stock.Action, stock.RatingFrom, stock.RatingTo, stock.TargetFrom, stock.TargetTo)

	if err != nil {
		return fmt.Errorf("error inserting stock: %w", err)
	}
	return nil
}

// Función para obtener los datos de las acciones desde la API
func fetchStocks() ([]models.Stock, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdHRlbXB0cyI6MzQsImVtYWlsIjoibSIsImV4cCI6MTc0MTQ0Nzc5OCwiaWQiOiIwIiwicGFzc3dvcmQiOiInIE9SICcxJz0nMSJ9.HDP5iWr7UvRdWc4l99nIuUPkZM4Uw5eAD1X5QH-P3Ow").
		SetResult(&models.APIResponse{}).
		Get("https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list")

	if err != nil {
		return nil, err
	}

	apiResponse := resp.Result().(*models.APIResponse)
	return apiResponse.Items, nil
}
