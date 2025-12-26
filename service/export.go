package service

import (
	"maps"
	"net/http"
	"slices"
	"strconv"

	"github.com/rchhong/comiket-backend/models"
	"github.com/rchhong/comiket-backend/repositories"
)

type ExportService struct {
	exportRepository repositories.ExportRepository
}

func NewExportService(exportRepository repositories.ExportRepository) *ExportService {
	return &ExportService{
		exportRepository: exportRepository,
	}
}

func (exportService ExportService) GenerateExport() ([][]string, *models.ComiketBackendError) {
	rawRows, err := exportService.exportRepository.GetRawExportData()
	if err != nil {
		return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	headers := []string{
		"melonbooks_id",
		"url",
		"title",
		"price_in_yen",
		"price_in_usd",
	}
	// Add a column for each user
	// TODO: Maybe use separate database calls for this + make this less spaghetti
	// but for now I just want this to work, we can make it better later
	// "writing software is an iterative process"
	users := make(map[string]bool)
	for _, row := range rawRows {
		users[row.DiscordName] = true
	}
	sortedUsers := slices.Sorted(maps.Keys(users))

	headers = slices.Concat(headers, sortedUsers)
	outputRows := make(map[int]map[string]string)
	for _, row := range rawRows {
		existingOutputRow, exist := outputRows[row.MelonbooksId]
		if !exist {
			toAdd := map[string]string{
				"melonbooks_id": strconv.Itoa(row.MelonbooksId),
				"url":           row.Url,
				"title":         row.Title,
				"price_in_yen":  strconv.Itoa(row.PriceInYen),
				"price_in_usd":  strconv.FormatFloat(row.PriceInUsd, 'f', -1, 64),
			}
			for _, user := range sortedUsers {
				toAdd[user] = ""
			}
			toAdd[row.DiscordName] = "X"
			outputRows[row.MelonbooksId] = toAdd
		} else {
			existingOutputRow[row.DiscordName] = "X"
		}
	}

	var ret [][]string
	ret = append(ret, headers)
	for row := range maps.Values(outputRows) {
		var temp []string
		for _, header := range headers {
			temp = append(temp, row[header])
		}
		ret = append(ret, temp)
	}

	return ret, nil
}
