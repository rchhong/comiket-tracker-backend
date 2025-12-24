package service

import (
	"maps"
	"net/http"
	"slices"
	"strconv"

	"github.com/rchhong/comiket-backend/dao"
	"github.com/rchhong/comiket-backend/models"
)

type AdminService struct {
	reservationDao *dao.ReservationDAO
}

func NewAdminService(reservationDao *dao.ReservationDAO) *AdminService {
	return &AdminService{
		reservationDao: reservationDao,
	}
}

func (adminService AdminService) GenerateExport() ([][]string, error) {
	rawRows, err := adminService.reservationDao.GetRawExportData()
	if err != nil {
		return nil, models.StatusError{Err: err, StatusCode: http.StatusInternalServerError}
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

	for _, header := range headers {
		println(header)
	}

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
