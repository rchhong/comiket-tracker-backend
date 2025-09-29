package models

import (
	"strings"
	"time"
)

type CurrencyConverterAPIResponse struct {
	BaseCurrencyCode string                     `json:"base_currency_code"`
	BaseCurrencyName string                     `json:"base_currency_name"`
	Amount           float64                    `json:"amount,string"`
	UpdatedDate      CurrencyConverterTime      `json:"updated_date"`
	Rates            map[string]RateInformation `json:"rates"`
	Status           string                     `json:"status"`
}

type RateInformation struct {
	CurrencyName  string  `json:"currency_name"`
	Rate          float64 `json:"rate,string"`
	RateForAmount float64 `json:"rate_for_amount,string"`
}

type CurrencyConverterTime struct {
	time.Time
}

func (currencyConverterTime *CurrencyConverterTime) UnmarshalJSON(b []byte) error {
	rawTime := strings.Trim(string(b), "\"")
	convertedTime, err := time.Parse(time.DateOnly, rawTime)

	if err != nil {
		return err
	}

	currencyConverterTime.Time = convertedTime
	return nil
}

type CurrencyConverter interface {
	Convert(fromCurrencyAmount float64) float64
}
