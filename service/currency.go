package service

import (
	"github.com/rchhong/comiket-backend/utils"
)

type CurrencyConverterService struct {
	currencyConverter utils.CurrencyConverter
}

func NewCurrencyConverterService(currencyConverter *utils.CurrencyConverter) *CurrencyConverterService {
	return &CurrencyConverterService{
		currencyConverter: *currencyConverter,
	}
}

func (currencyConverterService *CurrencyConverterService) Convert(fromCurrencyAmount float64) float64 {
	return currencyConverterService.currencyConverter.Convert(fromCurrencyAmount)
}
