package services

import (
	"testing"

	"stock-api/models"

	"github.com/stretchr/testify/assert"
)

// Mock para la base de datos
type MockStockService struct{}

func (m *MockStockService) GetAllStocks() ([]models.Stock, error) {
	return []models.Stock{
		{Ticker: "AAPL", RatingFrom: "Sell", RatingTo: "Buy", Action: "upgraded by", TargetFrom: "$100", TargetTo: "$120"},
		{Ticker: "GOOGL", RatingFrom: "Neutral", RatingTo: "Buy", Action: "reiterated by", TargetFrom: "$1500", TargetTo: "$1800"},
		{Ticker: "TSLA", RatingFrom: "Buy", RatingTo: "Sell", Action: "downgraded by", TargetFrom: "$800", TargetTo: "$700"},
	}, nil
}

func (m *MockStockService) GetRecommendations() ([]models.StockScore, error) {
	stocks, _ := m.GetAllStocks()
	var scoredStocks []models.StockScore
	for _, stock := range stocks {
		score := calculateStockScore(stock)
		scoredStocks = append(scoredStocks, models.StockScore{Stock: stock, Score: score})
	}

	// Simular ordenamiento de mayor a menor
	if scoredStocks[0].Score < scoredStocks[1].Score {
		scoredStocks[0], scoredStocks[1] = scoredStocks[1], scoredStocks[0]
	}

	return scoredStocks, nil
}

func TestGetAllStocks(t *testing.T) {
	mockService := &MockStockService{}
	stocks, err := mockService.GetAllStocks()
	assert.NoError(t, err)
	assert.NotEmpty(t, stocks)
	assert.Equal(t, 3, len(stocks))
	assert.Equal(t, "AAPL", stocks[0].Ticker)
}

func TestGetRecommendations(t *testing.T) {
	mockService := &MockStockService{}

	recommendations, err := mockService.GetRecommendations()
	assert.NoError(t, err)
	assert.NotEmpty(t, recommendations)

	// Verificar que la lista tiene al menos un elemento
	assert.Greater(t, len(recommendations), 0)

	// Verificar que las acciones están ordenadas por Score
	assert.GreaterOrEqual(t, recommendations[0].Score, recommendations[1].Score)
}

// Pruebas de funciones auxiliares

func TestCalculateRatingChangeScore(t *testing.T) {
	tests := []struct {
		name  string
		stock models.Stock
		want  float64
	}{
		{"Sell to Buy", models.Stock{RatingFrom: "Sell", RatingTo: "Buy"}, 10},
		{"Neutral to Buy", models.Stock{RatingFrom: "Neutral", RatingTo: "Buy"}, 5},
		{"Buy to Sell", models.Stock{RatingFrom: "Buy", RatingTo: "Sell"}, -10},
		{"Unknown transition", models.Stock{RatingFrom: "Hold", RatingTo: "Buy"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateRatingChangeScore(tt.stock)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalculateActionScore(t *testing.T) {
	tests := []struct {
		name  string
		stock models.Stock
		want  float64
	}{
		{"Upgraded by", models.Stock{Action: "upgraded by"}, 4},
		{"Downgraded by", models.Stock{Action: "downgraded by"}, -4},
		{"Target raised by", models.Stock{Action: "target raised by"}, 3},
		{"Unknown action", models.Stock{Action: "modified by"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateActionScore(tt.stock)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalculateTargetPriceChangeScore(t *testing.T) {
	tests := []struct {
		name  string
		stock models.Stock
		want  float64
	}{
		{"Increase from 100 to 120", models.Stock{TargetFrom: "$100", TargetTo: "$120"}, 2},
		{"Decrease from 200 to 180", models.Stock{TargetFrom: "$200", TargetTo: "$180"}, -1},
		{"Invalid format", models.Stock{TargetFrom: "invalid", TargetTo: "$150"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateTargetPriceChangeScore(tt.stock)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalculateStockScore(t *testing.T) {
	stock := models.Stock{
		RatingFrom: "Sell", RatingTo: "Buy",
		Action:     "upgraded by",
		TargetFrom: "$100", TargetTo: "$120",
	}

	got := calculateStockScore(stock)
	expected := 10*0.4 + 4*0.3 + 2*0.3
	assert.Equal(t, expected, got)
}
