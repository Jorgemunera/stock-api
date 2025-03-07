package models

// Stock representa los datos de una acción
type Stock struct {
	Ticker     string `json:"ticker"`
	Company    string `json:"company"`
	Brokerage  string `json:"brokerage"`
	Action     string `json:"action"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
}

// StockScore representa la puntuación de una acción
type StockScore struct {
	Stock Stock
	Score float64
}

// APIResponse representa la respuesta de la API de acciones
type APIResponse struct {
	Items    []Stock `json:"items"`
	NextPage string  `json:"next_page"`
}
