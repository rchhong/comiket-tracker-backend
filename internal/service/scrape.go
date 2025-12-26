package service

import (
	"net/http"

	"github.com/rchhong/comiket-backend/internal/models"
	"github.com/rchhong/comiket-backend/internal/service/currency"
	"github.com/rchhong/comiket-backend/internal/service/scrape"
)

type MelonbooksScraperService struct {
	melonbooksScraper *scrape.MelonbooksScraper
	currencyService   currency.CurrencyConverter
}

func NewMelonbooksScraperService(melonbooksScraper *scrape.MelonbooksScraper, currencyService currency.CurrencyConverter) *MelonbooksScraperService {
	return &MelonbooksScraperService{
		melonbooksScraper: melonbooksScraper,
		currencyService:   currencyService,
	}
}

func (melonbooksScraperService *MelonbooksScraperService) ScrapeMelonbooksProduct(melonbooksId int) (*models.Doujin, error) {
	melonbooksData, err := melonbooksScraperService.melonbooksScraper.ScrapeMelonbooksProduct(melonbooksId)
	if err != nil {
		return nil, models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	doujin := models.Doujin{
		MelonbooksId:    melonbooksData.MelonbooksId,
		Title:           melonbooksData.Title,
		PriceInYen:      melonbooksData.PriceInYen,
		PriceInUsd:      melonbooksScraperService.currencyService.Convert(float64(melonbooksData.PriceInYen)),
		IsR18:           melonbooksData.IsR18,
		ImagePreviewURL: melonbooksData.ImagePreviewURL,
		URL:             melonbooksData.URL,
		Circle:          melonbooksData.Circle,
		Authors:         melonbooksData.Authors,
		Genres:          melonbooksData.Genres,
		Events:          melonbooksData.Events,
	}

	return &doujin, nil
}
