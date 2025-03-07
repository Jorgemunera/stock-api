package services

import (
	"database/sql"
	"stock-api/models"
	"strconv"
	"strings"
)

// StockService maneja la lógica de negocio relacionada con las acciones
type StockService struct {
	db *sql.DB
}

// NewStockService crea una nueva instancia de StockService
func NewStockService(db *sql.DB) *StockService {
	return &StockService{db: db}
}

// GetAllStocks obtiene todas las acciones de la base de datos
func (s *StockService) GetAllStocks() ([]models.Stock, error) {
	rows, err := s.db.Query("SELECT * FROM stocks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var s models.Stock
		err := rows.Scan(&s.Ticker, &s.Company, &s.Brokerage, &s.Action, &s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}

	return stocks, nil
}

// GetRecommendations obtiene las mejores acciones para invertir
func (s *StockService) GetRecommendations() ([]models.StockScore, error) {
	stocks, err := s.GetAllStocks()
	if err != nil {
		return nil, err
	}

	var scoredStocks []models.StockScore
	for _, stock := range stocks {
		score := calculateStockScore(stock)
		scoredStocks = append(scoredStocks, models.StockScore{Stock: stock, Score: score})
	}

	sortStocksByScore(scoredStocks)

	if len(scoredStocks) > 5 {
		return scoredStocks[:5], nil
	}
	return scoredStocks, nil
}

// calculateStockScore calcula la puntuación de una acción
func calculateStockScore(stock models.Stock) float64 {
	ratingChangeScore := calculateRatingChangeScore(stock)
	growthScore := calculateGrowthScore(stock)
	actionScore := calculateActionScore(stock)

	// Ponderación de los factores
	return (0.5 * growthScore) + (0.3 * ratingChangeScore) + (0.2 * actionScore)
}

// calculateGrowthScore calcula el potencial de crecimiento basado en el target price
func calculateGrowthScore(stock models.Stock) float64 {
	targetFrom, _ := strconv.ParseFloat(strings.TrimPrefix(stock.TargetFrom, "$"), 64)
	targetTo, _ := strconv.ParseFloat(strings.TrimPrefix(stock.TargetTo, "$"), 64)
	if targetFrom == 0 {
		return 0.0
	}
	return ((targetTo - targetFrom) / targetFrom) * 10
}

// calculateRatingChangeScore evalúa el cambio de rating
func calculateRatingChangeScore(stock models.Stock) float64 {
	ratingChanges := map[string]float64{
		"SellNeutral": 2, "NeutralBuy": 4, "SellBuy": 6,
		"BuyNeutral": -4, "NeutralSell": -2, "BuySell": -6,
		"BuyBuy": 1, "NeutralNeutral": 0, "SellSell": -1,
	}
	key := stock.RatingFrom + stock.RatingTo
	return ratingChanges[key]
}

// calculateActionScore evalúa el tipo de acción tomada
func calculateActionScore(stock models.Stock) float64 {
	switch stock.Action {
	case "upgraded by":
		return 2.0
	case "target raised by":
		return 1.5
	case "reiterated by":
		return 0.5
	case "target lowered by":
		return -1.5
	case "downgraded by":
		return -2.0
	default:
		return 0.0
	}
}

// sortStocksByScore ordena las acciones por puntuación (de mayor a menor)
func sortStocksByScore(stocks []models.StockScore) {
	for i := 0; i < len(stocks); i++ {
		for j := i + 1; j < len(stocks); j++ {
			if stocks[i].Score < stocks[j].Score {
				stocks[i], stocks[j] = stocks[j], stocks[i]
			}
		}
	}
}
