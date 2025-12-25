package service

import (
	"github.com/rchhong/comiket-backend/models"
	"github.com/rchhong/comiket-backend/scrape"
)

type MelonbooksScraperService struct {
	melonbooksScraper scrape.MelonbooksScraper
}

func NewMelonbooksScraperService(melonbooksScraper *scrape.MelonbooksScraper) *MelonbooksScraperService {
	return &MelonbooksScraperService{
		melonbooksScraper: *melonbooksScraper,
	}
}

func (melonbooksScraperService *MelonbooksScraperService) ScrapeMelonbooksProduct(melonbooksId int) (*models.Doujin, error) {
	return melonbooksScraperService.melonbooksScraper.ScrapeMelonbooksProduct(melonbooksId)
}
