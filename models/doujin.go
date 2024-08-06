package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Doujin struct {
	Id    primitive.ObjectID `bson:"_id" json:"_id"`
	Title string             `bson:"title" json:"title"`
}
