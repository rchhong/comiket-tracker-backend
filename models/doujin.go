package models

import (
	"fmt"
	"time"
)

type Doujin struct {
	Melonbooks_id     int64      `json:"melonbooks_id"`
	Title             string     `json:"title"`
	Price_In_Yen      int        `json:"price_in_yen"`
	Price_In_Usd      float64    `json:"price_in_usd"`
	Is_R18            bool       `json:"is_r18"`
	Image_Preview_Url string     `json:"image_preview_url"`
	URL               string     `json:"url"`
	Circle_Name       string     `json:"circle_name"`
	Author_Names      [10]string `json:"author_names"`
	Genres            [10]string `json:"genres"`
	Events            [10]string `json:"events"`
	Created_At        time.Time  `json:"created_at"`
	Updated_At        time.Time  `json:"updated_at"`
}

func (doujin Doujin) String() string {
	return fmt.Sprintf("{Melonbooks_id: %d, Title: %s}\n", doujin.Melonbooks_id, doujin.Title)
}
