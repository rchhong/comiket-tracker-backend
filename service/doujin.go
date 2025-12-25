package service

import (
	"errors"
	"net/http"

	"github.com/rchhong/comiket-backend/dao"
	"github.com/rchhong/comiket-backend/models"
)

type DoujinService struct {
	doujinDAO                *dao.DoujinDAO
	melonbooksScraperService *MelonbooksScraperService
}

func NewDoujinService(doujinDAO *dao.DoujinDAO, melonbooksScraperService *MelonbooksScraperService) *DoujinService {
	return &DoujinService{
		doujinDAO:                doujinDAO,
		melonbooksScraperService: melonbooksScraperService,
	}
}

func (doujinService DoujinService) GetDoujinByMelonbooksId(melonbooksId int) (*models.DoujinWithMetadata, error) {
	return doujinService.doujinDAO.GetDoujinByMelonbooksId(melonbooksId)
}

func (doujinService DoujinService) CreateDoujin(melonbooksId int) (*models.DoujinWithMetadata, error) {
	scrapedData, err := doujinService.melonbooksScraperService.ScrapeMelonbooksProduct(melonbooksId)
	if err != nil {
		return nil, models.StatusError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return doujinService.doujinDAO.CreateDoujin(*scrapedData)
}

func (doujinService DoujinService) UpdateDoujin(melonbooksId int) (*models.DoujinWithMetadata, error) {
	scrapedData, err := doujinService.melonbooksScraperService.ScrapeMelonbooksProduct(melonbooksId)
	if err != nil {
		return nil, models.StatusError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return doujinService.doujinDAO.UpdateDoujin(melonbooksId, *scrapedData)
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
		return doujinService.doujinDAO.DeleteDoujin(melonbooksId)
	}

	return err
}
