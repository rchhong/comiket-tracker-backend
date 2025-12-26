package ipgeoapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/rchhong/comiket-backend/internal/currency/ipgeoapi/dto"
)

type CurrencyConverterIpGeoAPI struct {
	currencyAPIURL string
	currencyAPIKey string
	fromCurrency   string
	toCurrency     string
	conversionRate float64
	updatedAt      time.Time
}

var IPGEO_API_CURRENCY_API_URL string = "https://api.getgeoapi.com/v2/currency/convert"

func NewCurrencyConverterIpGeoAPI(currencyAPIURL string, currencyAPIKey string, fromCurrency string, toCurrency string) (*CurrencyConverterIpGeoAPI, error) {
	currencyConverter := &CurrencyConverterIpGeoAPI{
		currencyAPIURL: currencyAPIURL,
		currencyAPIKey: currencyAPIKey,
		toCurrency:     toCurrency,
		fromCurrency:   fromCurrency,
		conversionRate: -1,
		updatedAt:      time.Now(),
	}

	err := currencyConverter.updateConversionRate()
	if err != nil {
		return nil, err
	}
	return currencyConverter, nil

}

func (currencyConverter *CurrencyConverterIpGeoAPI) Convert(fromCurrencyAmount float64) float64 {
	// At this point, the currency rate should be populated
	// Thus, we can just silently fail if we cannot update
	if time.Since(currencyConverter.updatedAt).Hours() >= 24 {
		err := currencyConverter.updateConversionRate()
		if err != nil {
			// TODO: better logging system
			log.Printf("[WARNING] unable to update currency conversion (%v)", err)
		}
	}

	return fromCurrencyAmount * currencyConverter.conversionRate
}

func (currencyConverter *CurrencyConverterIpGeoAPI) updateConversionRate() error {
	url, err := url.Parse(currencyConverter.currencyAPIURL)
	if err != nil {
		return err
	}

	queryParameters := url.Query()
	queryParameters.Set("api_key", currencyConverter.currencyAPIKey)
	queryParameters.Set("format", "json")
	queryParameters.Set("from", currencyConverter.fromCurrency)
	queryParameters.Set("to", currencyConverter.toCurrency)

	url.RawQuery = queryParameters.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned with non-OK status %d", resp.StatusCode)
	}

	var parsedResponse dto.CurrencyConverterAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&parsedResponse)
	if err != nil {
		return err
	}

	// Only works because there is exactly one item in the map
	for _, v := range parsedResponse.Rates {
		currencyConverter.conversionRate = v.Rate
	}
	currencyConverter.updatedAt = parsedResponse.UpdatedDate.Time

	return nil

}
