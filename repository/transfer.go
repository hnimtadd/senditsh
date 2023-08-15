package repository

import (
	"context"
	"log"
	"time"

	"github.com/hnimtadd/senditsh/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *repositoryImpl) InsertTransfer(transfer *data.Transfer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	log.Println("Insert transfer:", transfer)
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

func (repo *repositoryImpl) UpdateTransferStatus(transferId primitive.ObjectID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	log.Println("update transfer with id: ", transferId.String())
	filter := bson.M{
		"_id": transferId,
	}
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	cur := repo.db.Collection("transfers").FindOneAndUpdate(ctx, filter, update)
	if err := cur.Err(); err != nil {
		log.Println("Error while update", err)
		return err
	}
	return nil
}
func (repo *repositoryImpl) GetLastTransfer(userName string) (*data.Transfer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{
		"userName": userName,
	}
	res := repo.db.Collection("transfers").FindOne(ctx, filter, &options.FindOneOptions{Sort: bson.M{"createdAt": -1}})
	if err := res.Err(); err != nil {
		return nil, err
	}
	var transfer = new(data.Transfer)
	if err := res.Decode(transfer); err != nil {
		return nil, err
	}
	return transfer, nil
}
