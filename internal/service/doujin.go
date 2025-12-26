package service

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/rchhong/comiket-backend/internal/models"
	"github.com/rchhong/comiket-backend/internal/repositories"
	"github.com/rchhong/comiket-backend/internal/scrape"
)

type DoujinService struct {
	doujinRepository  repositories.DoujinRepository
	melonbooksScraper *scrape.MelonbooksScraper
}

func NewDoujinService(doujinRepository repositories.DoujinRepository, melonbooksScraper *scrape.MelonbooksScraper) *DoujinService {
	return &DoujinService{
		doujinRepository:  doujinRepository,
		melonbooksScraper: melonbooksScraper,
	}
}

func (doujinService DoujinService) GetDoujinByMelonbooksId(melonbooksId int) (*models.DoujinWithMetadata, *models.ComiketBackendError) {
	doujin, err := doujinService.doujinRepository.GetDoujinByMelonbooksId(melonbooksId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusNotFound}
		} else {
			return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
		}
	}
	return doujin, nil
}

func (doujinService DoujinService) CreateDoujin(melonbooksId int) (*models.DoujinWithMetadata, *models.ComiketBackendError) {
	scrapedData, err := doujinService.melonbooksScraper.ScrapeMelonbooksProduct(melonbooksId)
	if err != nil {
		return nil, &models.ComiketBackendError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	doujin, err := doujinService.doujinRepository.CreateDoujin(*scrapedData)
	if err != nil {
		return nil, &models.ComiketBackendError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return doujin, nil
}

func (doujinService DoujinService) UpdateDoujin(melonbooksId int) (*models.DoujinWithMetadata, *models.ComiketBackendError) {
	scrapedData, err := doujinService.melonbooksScraper.ScrapeMelonbooksProduct(melonbooksId)
	if err != nil {
		return nil, &models.ComiketBackendError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	updatedDoujin, err := doujinService.doujinRepository.UpdateDoujin(melonbooksId, *scrapedData)
	if err != nil {
		return nil, &models.ComiketBackendError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return updatedDoujin, nil
}

func (doujinService DoujinService) UpsertDoujin(melonbooksId int) (*models.DoujinWithMetadata, *models.ComiketBackendError) {
	existingDoujin, err := doujinService.GetDoujinByMelonbooksId(melonbooksId)
	if existingDoujin != nil {
		return doujinService.UpdateDoujin(melonbooksId)
	}

	var statusError models.ComiketBackendError
	if errors.As(err, &statusError) {
		if statusError.StatusCode == http.StatusNotFound {
			return doujinService.CreateDoujin(melonbooksId)
		}
	}

	return nil, err
}

func (doujinService DoujinService) DeleteDoujin(melonbooksId int) *models.ComiketBackendError {
	existingDoujin, err := doujinService.GetDoujinByMelonbooksId(melonbooksId)
	if existingDoujin != nil {
		err := doujinService.doujinRepository.DeleteDoujin(melonbooksId)
		if err != nil {
			return &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
		}
		return nil
	}

	return err
}
