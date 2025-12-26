package service

import "github.com/rchhong/comiket-backend/internal/service/currency"

type CurrencyService struct {
	currencyConverter currency.CurrencyConverter
}

func NewCurrencyService(currencyConverter currency.CurrencyConverter) *CurrencyService {
	return &CurrencyService{
		currencyConverter: currencyConverter,
	}
}

func (currencyService *CurrencyService) Convert(fromCurrencyAmount float64) float64 {
	return currencyService.currencyConverter.Convert(fromCurrencyAmount)
}
