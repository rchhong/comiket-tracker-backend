package currency

type CurrencyConverter interface {
	Convert(fromCurrencyAmount float64) float64
}
