package models

import (
	"fmt"
	"time"
)

type Reservation struct {
	Reservation_Id int8      `json:"reservation_id"`
	Discord_Id     int64     `json:"discord_id"`
	Melonbooks_Id  int64     `json:"melonbooks_id"`
	Created_At     time.Time `json:"created_at"`
}

func (reservation Reservation) String() string {
	return fmt.Sprintf("{Reservation_Id: %d Discord_Id: %d Melonbooks_Id: %d}\n", reservation.Reservation_Id, reservation.Discord_Id, reservation.Melonbooks_Id)
}
