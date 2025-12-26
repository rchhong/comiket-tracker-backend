package dto

type MelonbooksData struct {
	MelonbooksId    int      `json:"melonbooks_id"`
	Title           string   `json:"title"`
	PriceInYen      int      `json:"price_in_yen"`
	IsR18           bool     `json:"is_r18"`
	ImagePreviewURL string   `json:"image_preview_url"`
	URL             string   `json:"url"`
	Circle          string   `json:"circle"`
	Authors         []string `json:"authors"`
	Genres          []string `json:"genres"`
	Events          []string `json:"events"`
}
