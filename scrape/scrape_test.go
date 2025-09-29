package scrape

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/rchhong/comiket-backend/models"
	"github.com/stretchr/testify/assert"
)

type CurrencyConverterMock struct {
	mockRate float64
}

func (currencyConverterMock CurrencyConverterMock) Convert(fromCurrencyAmount float64) float64 {
	return fromCurrencyAmount * currencyConverterMock.mockRate
}

func NewMockCurrencyConverter(mockRate float64) *CurrencyConverterMock {
	return &CurrencyConverterMock{
		mockRate: mockRate,
	}
}

func TestScrape(t *testing.T) {
	assert := assert.New(t)

	mockRate := 0.5

	mockCurrencyConverter := NewMockCurrencyConverter(mockRate)
	melonbooksScraper := NewMelonbooksScraper(mockCurrencyConverter)

	testFile, err := os.Open("./testdata/scrape.json")
	assert.Nil(err)

	var doujinsToTest []models.Doujin
	json.NewDecoder(testFile).Decode(&doujinsToTest)

	for _, test := range doujinsToTest {
		doujin, err := melonbooksScraper.ScrapeMelonbooksProduct(test.MelonbooksId)
		assert.Nil(err)

		assert.Equal(test.MelonbooksId, doujin.MelonbooksId, "MelonbooksId is set to the expected value")
		assert.Equal(test.Title, doujin.Title, "Title is set to the expected value")
		assert.Equal(test.PriceInYen, doujin.PriceInYen, "PriceInYen is set to the expected value")
		assert.Equal(test.PriceInUsd, doujin.PriceInUsd, "PriceInUsd is set to the expected value")
		assert.Equal(test.IsR18, doujin.IsR18, "IsR18 is set to the expected value")
		assert.Equal(test.ImagePreviewURL, doujin.ImagePreviewURL, "ImagePreviewURL is set to the expected value")
		assert.Equal(test.URL, doujin.URL, "URL is set to the expected value")
		assert.Equal(test.Circle, doujin.Circle, "Circle is set to the expected value")
		assert.ElementsMatch(test.Authors, doujin.Authors, "Authors is set to the expected value")
		assert.ElementsMatch(test.Genres, doujin.Genres, "Genres is set to the expected value")
		assert.ElementsMatch(test.Events, doujin.Events, "Events is set to the expected value")

	}
}
