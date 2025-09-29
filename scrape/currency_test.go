package scrape

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCurrencyConverter(t *testing.T) {
	assert := assert.New(t)

	testQueryParameters := make(map[string]string)
	testQueryParameters["api_key"] = "test_api_key"
	testQueryParameters["format"] = "json"
	testQueryParameters["from"] = "JPY"
	testQueryParameters["to"] = "USD"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Contains correct query parameters
		assert.Equal(len(testQueryParameters), len(r.URL.Query()), "The number of URL parameters should be the same")

		for k, v := range r.URL.Query() {
			actualValue, exists := testQueryParameters[k]
			assert.True(exists, fmt.Sprintf("%s should exist in testQueryParameters", k))
			assert.Equal(v[0], actualValue, fmt.Sprintf("%s should be set to the expected value", k))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
		  "base_currency_code": "JPY",
		  "base_currency_name": "Japanese yen",
		  "amount": "1.0000",
		  "updated_date": "2025-09-27",
		  "rates": {
			"USD": {
			  "currency_name": "United States dollar",
			  "rate": "0.5",
			  "rate_for_amount": "0.5"
			}
		  },
		  "status": "success"
		}`))
	}))

	defer server.Close()

	testCurrencyConverter, err := NewCurrencyConverterImpl(server.URL, testQueryParameters["api_key"], testQueryParameters["from"], testQueryParameters["to"])
	assert.Nil(err, "No errors should occur when trying to update the currency converter")
	assert.Equal(0.5, testCurrencyConverter.conversionRate, "The currency rate should be the expected value")
	expectedUpdatedAtTimestamp, err := time.Parse(time.DateOnly, "2025-09-27")
	assert.Equal(expectedUpdatedAtTimestamp, testCurrencyConverter.updatedAt, "The updatedAt timestamp should have been updated")
	convertedCurrency := testCurrencyConverter.Convert(10.0)
	assert.Equal(5.0, convertedCurrency, "The currency conversion should have the expected result")

}

func TestCurrencyConverterWithUpdate(t *testing.T) {
	assert := assert.New(t)

	testQueryParameters := make(map[string]string)
	testQueryParameters["api_key"] = "test_api_key"
	testQueryParameters["format"] = "json"
	testQueryParameters["from"] = "JPY"
	testQueryParameters["to"] = "USD"

	firstCall := true
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Contains correct query parameters
		assert.Equal(len(testQueryParameters), len(r.URL.Query()), "The number of URL parameters should be the same")

		for k, v := range r.URL.Query() {
			actualValue, exists := testQueryParameters[k]
			assert.True(exists, fmt.Sprintf("%s should exist in testQueryParameters", k))
			assert.Equal(v[0], actualValue, fmt.Sprintf("%s should be set to the expected value", k))
		}

		if firstCall {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
			  "base_currency_code": "JPY",
			  "base_currency_name": "Japanese yen",
			  "amount": "1.0000",
			  "updated_date": "1970-01-01",
			  "rates": {
				"USD": {
				  "currency_name": "United States dollar",
				  "rate": "0.5",
				  "rate_for_amount": "0.5"
				}
			  },
			  "status": "success"
			}`))
			firstCall = false
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
			  "base_currency_code": "JPY",
			  "base_currency_name": "Japanese yen",
			  "amount": "1.0000",
			  "updated_date": "2025-09-27",
			  "rates": {
				"USD": {
				  "currency_name": "United States dollar",
				  "rate": "0.75",
				  "rate_for_amount": "0.75"
				}
			  },
			  "status": "success"
			}`))
		}

	}))

	defer server.Close()

	testCurrencyConverter, err := NewCurrencyConverterImpl(server.URL, testQueryParameters["api_key"], testQueryParameters["from"], testQueryParameters["to"])
	assert.Nil(err, "No errors should occur when trying to update the currency converter")
	assert.Equal(0.5, testCurrencyConverter.conversionRate, "The currency rate should be the expected value")
	expectedUpdatedAtTimestamp, err := time.Parse(time.DateOnly, "1970-01-01")
	assert.Equal(expectedUpdatedAtTimestamp, testCurrencyConverter.updatedAt, "The updatedAt timestamp should have been updated")

	convertedCurrency := testCurrencyConverter.Convert(10.0)
	assert.Equal(0.75, testCurrencyConverter.conversionRate, "The currency rate should have been updated")
	expectedUpdatedAtTimestamp, err = time.Parse(time.DateOnly, "2025-09-27")
	assert.Equal(expectedUpdatedAtTimestamp, testCurrencyConverter.updatedAt, "The updatedAt timestamp should have been updated")
	assert.Equal(7.5, convertedCurrency, "The currency conversion should have the expected result")

}
