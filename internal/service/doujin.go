package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/rchhong/comiket-backend/internal/models"
	"github.com/rchhong/comiket-backend/internal/repositories"
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

func (doujinService DoujinService) GetDoujinByMelonbooksId(ctx context.Context, melonbooksId int) (*models.DoujinWithMetadata, *models.ComiketBackendError) {
	doujin, err := doujinService.doujinRepository.GetDoujinByMelonbooksId(ctx, melonbooksId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusNotFound}
		} else {
			return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
		}
	}
	return doujin, nil
}

func (doujinService DoujinService) CreateDoujin(ctx context.Context, melonbooksId int) (*models.DoujinWithMetadata, *models.ComiketBackendError) {
	scrapedData, err := doujinService.melonbooksScraperService.ScrapeMelonbooksProduct(melonbooksId)
	if err != nil {
		return nil, &models.ComiketBackendError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	doujin, err := doujinService.doujinRepository.CreateDoujin(ctx, *scrapedData)
	if err != nil {
		return nil, &models.ComiketBackendError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return doujin, nil
}

func (doujinService DoujinService) UpdateDoujin(ctx context.Context, melonbooksId int) (*models.DoujinWithMetadata, *models.ComiketBackendError) {
	scrapedData, err := doujinService.melonbooksScraperService.ScrapeMelonbooksProduct(melonbooksId)
	if err != nil {
		return nil, &models.ComiketBackendError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	updatedDoujin, err := doujinService.doujinRepository.UpdateDoujin(ctx, melonbooksId, *scrapedData)
	if err != nil {
		return nil, &models.ComiketBackendError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return updatedDoujin, nil
}

func (doujinService DoujinService) UpsertDoujin(ctx context.Context, melonbooksId int) (*models.DoujinWithMetadata, *models.ComiketBackendError) {
	existingDoujin, err := doujinService.GetDoujinByMelonbooksId(ctx, melonbooksId)
	if existingDoujin != nil {
		return doujinService.UpdateDoujin(ctx, melonbooksId)
	}

	var statusError models.ComiketBackendError
	if errors.As(err, &statusError) {
		if statusError.StatusCode == http.StatusNotFound {
			return doujinService.CreateDoujin(ctx, melonbooksId)
		}
	}

	return nil, err
}

func (doujinService DoujinService) DeleteDoujin(ctx context.Context, melonbooksId int) *models.ComiketBackendError {
	existingDoujin, err := doujinService.GetDoujinByMelonbooksId(ctx, melonbooksId)
	if existingDoujin != nil {
		err := doujinService.doujinRepository.DeleteDoujin(ctx, melonbooksId)
		if err != nil {
			return &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
		}
		return nil
	}

	return err
}
