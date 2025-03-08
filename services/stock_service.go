package services

import (
	"database/sql"
	"sort"
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

	// Ordenar de mayor a menor por Score
	sort.Slice(scoredStocks, func(i, j int) bool {
		return scoredStocks[i].Score > scoredStocks[j].Score
	})

	// Devolver solo las 5 mejores recomendaciones
	if len(scoredStocks) > 5 {
		return scoredStocks[:5], nil
	}
	return scoredStocks, nil
}

// calculateStockScore calcula la puntuación de una acción con los nuevos pesos
func calculateStockScore(stock models.Stock) float64 {
	ratingChangeScore := calculateRatingChangeScore(stock) * 0.4
	actionScore := calculateActionScore(stock) * 0.3
	targetPriceScore := calculateTargetPriceChangeScore(stock) * 0.3

	return ratingChangeScore + actionScore + targetPriceScore
}

// calculateRatingChangeScore evalúa el cambio de rating
func calculateRatingChangeScore(stock models.Stock) float64 {
	ratingChanges := map[string]float64{
		"SellBuy": 10, "SellOverweight": 8, "SellNeutral": 5,
		"NeutralBuy": 5, "NeutralOverweight": 4, "Equal WeightOverweight": 4,
		"Market PerformBuy": 5, "In-LineOverweight": 4, "UnderweightOverweight": 9,
		"UnderweightNeutral": 4, "BuyNeutral": -5, "OverweightNeutral": -4,
		"BuySell": -10, "OverweightSell": -8, "NeutralSell": -5,
		"UnderweightSell": -7, "OverweightUnderweight": -9, "BuyOverweight": 3,
		"NeutralNeutral": 0, "BuyBuy": 3,
	}

	key := stock.RatingFrom + stock.RatingTo
	if score, exists := ratingChanges[key]; exists {
		return score
	}
	return 0
}

// calculateActionScore evalúa el tipo de acción tomada
func calculateActionScore(stock models.Stock) float64 {
	actionScores := map[string]float64{
		"upgraded by":       4,
		"downgraded by":     -4,
		"target raised by":  3,
		"target lowered by": -3,
		"reiterated by":     1,
	}

	if score, exists := actionScores[stock.Action]; exists {
		return score
	}
	return 0
}

// calculateTargetPriceChangeScore evalúa el impacto del cambio de target price
func calculateTargetPriceChangeScore(stock models.Stock) float64 {
	targetFrom, err1 := strconv.ParseFloat(strings.TrimPrefix(stock.TargetFrom, "$"), 64)
	targetTo, err2 := strconv.ParseFloat(strings.TrimPrefix(stock.TargetTo, "$"), 64)
	if err1 != nil || err2 != nil || targetFrom == 0 {
		return 0
	}

	percentageChange := ((targetTo - targetFrom) / targetFrom) * 10
	return percentageChange
}
