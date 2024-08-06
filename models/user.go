package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoujinReservation struct {
	Doujin_Id      primitive.ObjectID `bson:"doujin_id" json:"doujin_id"`
	Datetime_added time.Time          `bson:"datetime_added" json:"datetime_added"`
}

type User struct {
	Id           primitive.ObjectID  `bson:"_id" json:"_id"`
	Discord_Id   int                 `bson:"discord_id" json:"discord_id"`
	Reservations []DoujinReservation `bson:"reservations" json:"reservations"`
}

func (user User) String() string {
	return fmt.Sprintf("{Id: %s, Discord_Id: %d}", user.Id, user.Discord_Id)
}
