package repository

import (
	"context"
	"time"

	"github.com/hnimtadd/senditsh/data"
	"github.com/hnimtadd/senditsh/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo *repositoryImpl) CreateUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := repo.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return err
	}
	repo.logger.DefaultLog(logger.Info, "Inserted User into repository")
	return nil
}
func (repo *repositoryImpl) GetUsers() ([]data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	cur, err := repo.db.Collection("users").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var users []data.User
	for cur.Next(ctx) {
		var user data.User
		if err := cur.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	repo.logger.DefaultLog(logger.Info, "Getted users from repository")
	return users, nil
}

func (repo *repositoryImpl) GetUserById(id string) (*data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.D{primitive.E{Key: "id", Value: id}}
	res := repo.db.Collection("users").FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return nil, err
	}
	var user *data.User
	if err := res.Decode(user); err != nil {
		return nil, err
	}
	return user, nil

}
