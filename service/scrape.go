package service

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/rchhong/comiket-backend/models"
)

type MelonbooksScraper struct {
	collector *colly.Collector
}

func NewMelonbooksScraper() *MelonbooksScraper {
	return &MelonbooksScraper{
		collector: colly.NewCollector(),
	}
}

func (melonbooksScraper *MelonbooksScraper) ScrapeMelonbooksProduct(melonbooksProductId int) models.Doujin {
	var doujin models.Doujin

	melonbooksUrl := fmt.Sprintf("https://www.melonbooks.co.jp/detail/detail.php?product_id=%d&adult_view=1", melonbooksProductId)
	melonbooksScraper.collector.Visit(melonbooksUrl)

	return doujin
}
