package dao

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/rchhong/comiket-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DAO struct {
	client  *mongo.Client
	users   *mongo.Collection
	doujins *mongo.Collection
}

func NewDAO(uri string) DAO {
	options := options.Client()
	options.ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), options)

	if err != nil {
		panic(err)
	}

	return DAO{
		client:  client,
		users:   client.Database(os.Getenv("MONGO_DB_NAME")).Collection("users"),
		doujins: client.Database(os.Getenv("MONGO_DB_NAME")).Collection("doujins"),
	}
}

func (dao DAO) Close() {
	err := dao.client.Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}

}

func (dao DAO) RetrieveUserById(id primitive.ObjectID) (models.User, error) {
	var user models.User

	filter := bson.D{{Key: "_id", Value: id}}

	err := dao.users.FindOne(context.TODO(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return models.User{}, errors.New(fmt.Sprintf("Unable to find user with id %s", id))
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil

}

func (dao DAO) RetrieveDoujinById(id primitive.ObjectID) (models.Doujin, error) {
	var doujin models.Doujin

	filter := bson.D{{Key: "_id", Value: id}}

	err := dao.doujins.FindOne(context.TODO(), filter).Decode(&doujin)
	if err == mongo.ErrNoDocuments {
		return models.Doujin{}, errors.New(fmt.Sprintf("Unable to find user with id %s", id))
	}

	if err != nil {
		return models.Doujin{}, err
	}

	return doujin, nil

}
