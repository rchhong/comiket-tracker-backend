package service

import (
	"errors"
	"net/http"

	"github.com/rchhong/comiket-backend/models"
	"github.com/rchhong/comiket-backend/repositories"
)

type DoujinService struct {
	doujinRepository         repositories.DoujinRepository
	melonbooksScraperService *MelonbooksScraperService
}

func NewDoujinService(doujinRepository repositories.DoujinRepository, melonbooksScraperService *MelonbooksScraperService) *DoujinService {
	return &DoujinService{
		doujinRepository:         doujinRepository,
		melonbooksScraperService: melonbooksScraperService,
	}
}

func (doujinService DoujinService) GetDoujinByMelonbooksId(melonbooksId int) (*models.DoujinWithMetadata, error) {
	return doujinService.doujinRepository.GetDoujinByMelonbooksId(melonbooksId)
}

func (doujinService DoujinService) CreateDoujin(melonbooksId int) (*models.DoujinWithMetadata, error) {
	scrapedData, err := doujinService.melonbooksScraperService.ScrapeMelonbooksProduct(melonbooksId)
	if err != nil {
		return nil, models.StatusError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return doujinService.doujinRepository.CreateDoujin(*scrapedData)
}

func (doujinService DoujinService) UpdateDoujin(melonbooksId int) (*models.DoujinWithMetadata, error) {
	scrapedData, err := doujinService.melonbooksScraperService.ScrapeMelonbooksProduct(melonbooksId)
	if err != nil {
		return nil, models.StatusError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return doujinService.doujinRepository.UpdateDoujin(melonbooksId, *scrapedData)
}

func (doujinService DoujinService) UpsertDoujin(melonbooksId int) (*models.DoujinWithMetadata, error) {
	_, err := doujinService.GetDoujinByMelonbooksId(melonbooksId)
	if err == nil {
		return doujinService.UpdateDoujin(melonbooksId)
	}

	var statusError models.StatusError
	if errors.As(err, &statusError) {
		if statusError.StatusCode == http.StatusNotFound {
			return doujinService.CreateDoujin(melonbooksId)
		}
	}
	return nil, err
}

func (doujinService DoujinService) DeleteDoujin(melonbooksId int) error {
	_, err := doujinService.GetDoujinByMelonbooksId(melonbooksId)
	if err == nil {
		return doujinService.doujinRepository.DeleteDoujin(melonbooksId)
	}

	return err
}
