package model

type CurrencyConversionResponse struct {
	Success bool          `json:"success"`
	Query   CurrencyQuery `json:"query"`
	Info    InfoResponse  `json:"info"`
	Result  float64       `json:"result"`
}

type CurrencyQuery struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type InfoResponse struct {
	Timestamp int64   `json:"timestamp"`
	Rate      float64 `json:"rate"`
}

type QueryResponse struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}