package model

type CurrencyConversionResponse struct {
	Success bool    `json:"success"`
	Terms   string  `json:"terms"`
	Privacy string  `json:"privacy"`
	Query   string  `json:"query"`
	Info    string  `json:"info"`
	Result  float64 `json:"result"`
}

type InfoResponse struct {
	Timestamp int64  `json:"timestamp"`
	Quote     string `json:"quote"`
}

type QueryResponse struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}