package service

import (
	"github.com/rchhong/comiket-backend/models"
	"github.com/rchhong/comiket-backend/utils"
)

type MelonbooksScraperService struct {
	melonbooksScraper utils.MelonbooksScraper
}

func NewMelonbooksScraperService(melonbooksScraper *utils.MelonbooksScraper) *MelonbooksScraperService {
	return &MelonbooksScraperService{
		melonbooksScraper: *melonbooksScraper,
	}
}

func (melonbooksScraperService *MelonbooksScraperService) ScrapeMelonbooksProduct(melonbooksId int) (*models.Doujin, error) {
	return melonbooksScraperService.melonbooksScraper.ScrapeMelonbooksProduct(melonbooksId)
}
