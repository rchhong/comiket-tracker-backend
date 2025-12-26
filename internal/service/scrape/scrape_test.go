package scrape

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/rchhong/comiket-backend/internal/service/scrape/dto"
	"github.com/stretchr/testify/assert"
)

func TestScrape(t *testing.T) {
	assert := assert.New(t)

	melonbooksScraper := NewMelonbooksScraper()

	testFile, err := os.Open("./testdata/scrape.json")
	assert.Nil(err)

	var doujinsToTest []dto.MelonbooksData
	json.NewDecoder(testFile).Decode(&doujinsToTest)

	for _, test := range doujinsToTest {
		doujin, err := melonbooksScraper.ScrapeMelonbooksProduct(test.MelonbooksId)
		assert.Nil(err)

		assert.Equal(test.MelonbooksId, doujin.MelonbooksId, "MelonbooksId is set to the expected value")
		assert.Equal(test.Title, doujin.Title, "Title is set to the expected value")
		assert.Equal(test.PriceInYen, doujin.PriceInYen, "PriceInYen is set to the expected value")
		assert.Equal(test.IsR18, doujin.IsR18, "IsR18 is set to the expected value")
		assert.Equal(test.ImagePreviewURL, doujin.ImagePreviewURL, "ImagePreviewURL is set to the expected value")
		assert.Equal(test.URL, doujin.URL, "URL is set to the expected value")
		assert.Equal(test.Circle, doujin.Circle, "Circle is set to the expected value")
		assert.ElementsMatch(test.Authors, doujin.Authors, "Authors is set to the expected value")
		assert.ElementsMatch(test.Genres, doujin.Genres, "Genres is set to the expected value")
		assert.ElementsMatch(test.Events, doujin.Events, "Events is set to the expected value")

	}
}
