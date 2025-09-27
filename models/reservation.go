package models

import (
	"fmt"
	"time"
)

type Reservation struct {
	ReservationId int8  `json:"reservation_id"`
	DiscordId     int64 `json:"discord_id"`
	MelonbooksId  int64 `json:"melonbooks_id"`
}

type ReservationWithMetadata struct {
	Reservation
	CreatedAt time.Time `json:"created_at"`
}

func (reservation Reservation) String() string {
	return fmt.Sprintf("{Reservation_Id: %d DiscordId: %d Melonbooks_Id: %d}\n", reservation.ReservationId, reservation.DiscordId, reservation.MelonbooksId)
}
