package repository

import (
	"context"
	"time"

	"github.com/hnimtadd/senditsh/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo *repositoryImpl) InsertTransfer(transfer *data.Transfer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if _, err := repo.db.Collection("transfers").InsertOne(ctx, transfer); err != nil {
		return err
	}
	return nil
}

func (repo *repositoryImpl) GetTransfers() ([]data.Transfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	cur, err := repo.db.Collection("transfers").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var transfers = []data.Transfer{}
	for cur.Next(ctx) {
		var transfer data.Transfer
		if err := cur.Decode(&transfer); err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return transfers, nil
}

func (repo *repositoryImpl) GetTransfersOfUser(userName string) ([]data.Transfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.D{primitive.E{
		Key: "userName", Value: userName,
	}}
	cur, err := repo.db.Collection("transfers").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var transfers = []data.Transfer{}
	for cur.Next(ctx) {
		var transfer data.Transfer
		if err := cur.Decode(&transfer); err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return transfers, nil
}
