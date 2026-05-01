package scrape

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rchhong/comiket-backend/internal/logging"
	"github.com/rchhong/comiket-backend/internal/service/scrape/dto"
)

type MelonbooksScraper struct {
}

func NewMelonbooksScraper() *MelonbooksScraper {
	return &MelonbooksScraper{}
}

func (melonbooksScraper *MelonbooksScraper) ScrapeMelonbooksProduct(melonbooksProductId int) (*dto.MelonbooksData, error) {
	var melonbooksData dto.MelonbooksData
	var scrapeError error

	collector := colly.NewCollector()

	melonbooksUrl := fmt.Sprintf("https://www.melonbooks.co.jp/detail/detail.php?product_id=%d&adult_view=1", melonbooksProductId)

	melonbooksData.MelonbooksId = melonbooksProductId
	melonbooksData.URL = melonbooksUrl
	// Retrieve title
	collector.OnHTML("div.item-header > h1", func(e *colly.HTMLElement) {
		melonbooksData.Title = e.Text
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
		melonbooksData.PriceInYen = int(priceInYen)
	})

	// Retrieve image preview URL
	collector.OnHTML("div.item-img", func(e *colly.HTMLElement) {
		melonbooksData.ImagePreviewURL = fmt.Sprintf("https:%s", e.ChildAttr("img", "src"))
	})

	// Retrieve all other metadata from table at bottom
	collector.OnHTML("div.table-wrapper > table > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, c *colly.HTMLElement) {
			metadataItem := c.ChildText("th")

			// initialize all arrays just in case we don't have data for them
			switch metadataItem {
			case "サークル名":
				rawText := c.ChildText("td > a:nth-child(1)")
				melonbooksData.Circle = strings.Split(strings.TrimSpace(rawText), "\u00a0")[0]
			case "作家名":
				c.ForEach("a:not(.fa-heart)", func(_ int, c2 *colly.HTMLElement) {
					melonbooksData.Authors = append(melonbooksData.Authors, strings.TrimSpace(c2.Text))
				})
			case "ジャンル":
				c.ForEach("a:not(.fa-heart)", func(_ int, c2 *colly.HTMLElement) {
					melonbooksData.Genres = append(melonbooksData.Genres, strings.TrimSpace(c2.Text))
				})
			case "イベント":
				c.ForEach("a:not(.fa-heart)", func(_ int, c2 *colly.HTMLElement) {
					melonbooksData.Events = append(melonbooksData.Events, strings.TrimSpace(c2.Text))
				})
			case "作品種別":
				melonbooksData.IsR18 = c.ChildText("td") == "18禁"
			default:
				logging.Logger.Debug(fmt.Sprintf("Ignoring metadataItem %s for melonbooksData %s", metadataItem, melonbooksUrl))
			}
		})
	})

	collector.Visit(melonbooksUrl)

	if scrapeError != nil {
		return nil, scrapeError
	}

	return &melonbooksData, nil
}
