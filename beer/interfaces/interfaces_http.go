package interfaces

import (
	"go_trainning/beer-api/beer/external/model"
)

type CurrencyLayer interface {
	GetExchangeCurrency(currency, currencyBeer string, amountBeer float64) (model.CurrencyConversionResponse, error)
}