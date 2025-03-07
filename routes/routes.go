package api

import (
	"encoding/json"
	"net/http"
	"stock-api/services"
)

// Configurar las rutas de la API
func SetupRoutes(stockService *services.StockService) {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/stocks", StocksHandler(stockService))
	http.HandleFunc("/recommendations", RecommendationsHandler(stockService))
}

// Manejador para la ruta principal
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Stock API!"))
}

// Manejador para listar todas las acciones
func StocksHandler(stockService *services.StockService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stocks, err := stockService.GetAllStocks()
		if err != nil {
			http.Error(w, "Error fetching stocks", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stocks)
	}
}

// Manejador para recomendar las mejores acciones
func RecommendationsHandler(stockService *services.StockService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recommendations, err := stockService.GetRecommendations()
		if err != nil {
			http.Error(w, "Error fetching recommendations", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recommendations)
	}
}
