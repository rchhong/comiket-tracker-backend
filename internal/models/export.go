package models

type ExportRow struct {
	MelonbooksId int     `json:"melonbooks_id"`
	DiscordId    int64   `json:"discord_id"`
	Url          string  `json:"url"`
	Title        string  `json:"title"`
	PriceInYen   int     `json:"price_in_yen"`
	PriceInUsd   float64 `json:"price_in_usd"`
	DiscordName  string  `json:"discord_name"`
}
