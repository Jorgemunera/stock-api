package main

import (
	"fmt"
	"net/http"
	"stock-api/database"
	"stock-api/middlewares"
	api "stock-api/routes"
	"stock-api/services"
	"time"
)

func main() {
	fmt.Println("Starting server...")

	// Esperar 5 segundos para que CockroachDB esté listo
	fmt.Println("Waiting for CockroachDB to start...")
	time.Sleep(5 * time.Second)

	// Conectar a la base de datos
	db, err := database.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Crear la tabla si no existe
	if err := database.CreateTable(db); err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	// Insertar datos de las acciones en la tabla
	if err := database.InsertStocks(db); err != nil {
		fmt.Println("Error inserting stocks:", err)
		return
	}

	// Crear el servicio de acciones
	stockService := services.NewStockService(db)

	// Configurar las rutas de la API
	api.SetupRoutes(stockService)

	// Aplicar Middleware CORS a TODAS las rutas
	handler := middlewares.CORSHandler(http.DefaultServeMux)

	// Iniciar el servidor con el middleware
	fmt.Println("Server is running on http://localhost:3000")
	err = http.ListenAndServe(":3000", handler)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
