package scrape

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rchhong/comiket-backend/models"
)

type MelonbooksScraper struct {
	currencyConverter models.CurrencyConverter
}

func NewMelonbooksScraper(currencyConverter models.CurrencyConverter) *MelonbooksScraper {
	return &MelonbooksScraper{
		currencyConverter: currencyConverter,
	}
}

func (melonbooksScraper *MelonbooksScraper) ScrapeMelonbooksProduct(melonbooksProductId int) (*models.Doujin, error) {
	var doujin models.Doujin
	var scrapeError error

	collector := colly.NewCollector()

	melonbooksUrl := fmt.Sprintf("https://www.melonbooks.co.jp/detail/detail.php?product_id=%d&adult_view=1", melonbooksProductId)

	doujin.MelonbooksId = melonbooksProductId
	doujin.URL = melonbooksUrl
	// Retrieve title
	collector.OnHTML("div.item-header > h1", func(e *colly.HTMLElement) {
		doujin.Title = e.Text
	})

	// Retrieve cost in yen (+ convert to USD)
	collector.OnHTML("span.price--value", func(e *colly.HTMLElement) {
		parsedText := strings.TrimSpace(e.Text)
		re := regexp.MustCompile(`[^\d]`)
		parsedText = re.ReplaceAllString(parsedText, "")
		priceInYen, err := strconv.ParseInt(parsedText, 10, 64)
		if err != nil {
			scrapeError = err
			return
		}
		doujin.PriceInYen = int(priceInYen)
		doujin.PriceInUsd = melonbooksScraper.currencyConverter.Convert(float64(doujin.PriceInYen))
	})

	// Retrieve image preview URL
	collector.OnHTML("div.item-img", func(e *colly.HTMLElement) {
		doujin.ImagePreviewURL = fmt.Sprintf("https:%s", e.ChildAttr("img", "src"))
	})

	// Retrieve all other metadata from table at bottom
	collector.OnHTML("div.table-wrapper > table > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, c *colly.HTMLElement) {
			metadataItem := c.ChildText("th")

			// initialize all arrays just in case we don't have data for them
			switch metadataItem {
			case "サークル名":
				rawText := c.ChildText("td > a:nth-child(1)")
				doujin.Circle = strings.Split(strings.TrimSpace(rawText), "\u00a0")[0]
			case "作家名":
				c.ForEach("a:not(.fa-heart)", func(_ int, c2 *colly.HTMLElement) {
					doujin.Authors = append(doujin.Authors, strings.TrimSpace(c2.Text))
				})
			case "ジャンル":
				c.ForEach("a:not(.fa-heart)", func(_ int, c2 *colly.HTMLElement) {
					doujin.Genres = append(doujin.Genres, strings.TrimSpace(c2.Text))
				})
			case "イベント":
				c.ForEach("a:not(.fa-heart)", func(_ int, c2 *colly.HTMLElement) {
					doujin.Events = append(doujin.Events, strings.TrimSpace(c2.Text))
				})
			case "作品種別":
				doujin.IsR18 = c.ChildText("td") == "18禁"
			default:
				log.Printf("[WARNING] Ignoring metadataItem %s for doujin %s", metadataItem, melonbooksUrl)
			}
		})
	})

	collector.Visit(melonbooksUrl)

	if scrapeError != nil {
		return nil, scrapeError
	}

	return &doujin, nil
}
