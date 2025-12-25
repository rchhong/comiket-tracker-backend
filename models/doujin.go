package models

import (
	"fmt"
	"time"
)

type Doujin struct {
	MelonbooksId    int      `json:"melonbooks_id"`
	Title           string   `json:"title"`
	PriceInYen      int      `json:"price_in_yen"`
	PriceInUsd      float64  `json:"price_in_usd"`
	IsR18           bool     `json:"is_r18"`
	ImagePreviewURL string   `json:"image_preview_url"`
	URL             string   `json:"url"`
	Circle          string   `json:"circle"`
	Authors         []string `json:"authors"`
	Genres          []string `json:"genres"`
	Events          []string `json:"events"`
}

type DoujinWithMetadata struct {
	Doujin
	CreatedAt  time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

func (doujin Doujin) String() string {
	return fmt.Sprintf(`{Melonbooks_id: %d, Title: %s}\n`, doujin.MelonbooksId, doujin.Title)
}
