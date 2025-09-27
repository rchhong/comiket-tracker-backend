package utils

import (
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

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
		if len(r.URL.Query()) != len(testQueryParameters) {
			t.Errorf("ERROR: testQueryParameters keys != URL request keys ([%s] != [%s])", strings.Join(slices.Collect(maps.Keys(testQueryParameters)), ", "), strings.Join(slices.Collect(maps.Keys(r.URL.Query())), ", "))
		}

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

	testCurrencyConverter, err := NewCurrencyConverter(server.URL, testQueryParameters["api_key"], testQueryParameters["from"], testQueryParameters["to"])
	assert.Nil(err, "No errors should occur when trying to update the currency converter")
	assert.Equal(testCurrencyConverter.conversionRate, 0.5, "The currency rate should be the expected value")

	convertedCurrency := testCurrencyConverter.Convert(10.0)
	assert.Equal(convertedCurrency, 5.0, "The currency conversion should have the expected result")

}
