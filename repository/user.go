package repository

import (
	"context"
	"log"
	"time"

	"github.com/hnimtadd/senditsh/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo *repositoryImpl) CreateUser(user *data.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	user.Id = primitive.NewObjectID()
	_, err := repo.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return err
	}
	// repo.logger.DefaultLog(logger.Info, "Inserted User into repository")
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
	// repo.logger.DefaultLog(logger.Info, "Getted users from repository")
	return users, nil
}

func (repo *repositoryImpl) GetUserByUserName(userName string) (*data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{"userName": userName}
	res := repo.db.Collection("users").FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		log.Println("msg", "error while retriving user", "err", err)
		return nil, err
	}
	var user data.User
	if err := res.Decode(&user); err != nil {
		log.Println("msg", "error while decoding user", "err", err)
		return nil, err
	}
	return &user, nil

}

func (repo *repositoryImpl) GetUserBySSHKey(publicKey string) (*data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{"settings.sshKey": publicKey}
	res := repo.db.Collection("users").FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return nil, err
	}
	var user data.User
	if err := res.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *repositoryImpl) GetSettingOfUser(userName string) (*data.Settings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pipeline := []bson.M{{
		"$match": bson.M{
			"userName": userName,
		},
	},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$settings",
			},
		},
	}
	cur, err := repo.db.Collection("users").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(ctx)

	var settings []data.Settings
	if err := cur.All(ctx, &settings); err != nil {
		return nil, err
	}
	if len(settings) != 1 {
		log.Println("Debug")
	}

	setting := settings[0]
	return &setting, nil
}

func (repo *repositoryImpl) UpdateUserSetting(userName string, setting *data.Settings) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{
		"userName": userName,
	}
	update := bson.M{
		"$set": bson.M{
			"settings": setting,
		},
	}
	res := repo.db.Collection("users").FindOneAndUpdate(ctx, filter, update)
	if err := res.Err(); err != nil {
		return err
	}
	return nil
}

func (repo *repositoryImpl) InsertUserSSHKey(userName string, fullKey string, sshKey string, sshHash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{
		"userName": userName,
	}
	update := bson.M{
		"$set": bson.M{
			"settings.fullKey" : fullKey,
			"settings.sshKey":     sshKey,
			"settings.sshHash":    sshHash,
			"settings.modifiedAt": time.Now().Unix(),
		},
	}
	res := repo.db.Collection("users").FindOneAndUpdate(ctx, filter, update)
	if err := res.Err(); err != nil {
		return err
	}
	return nil
}

func (repo *repositoryImpl) InsertUserDomain(userName string, domain string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{
		"userName": userName,
	}
	update := bson.M{
		"$set": bson.M{
			"settings.subdomain":  domain,
			"settings.modifiedAt": time.Now().Unix(),
		},
	}
	res := repo.db.Collection("users").FindOneAndUpdate(ctx, filter, update)
	if err := res.Err(); err != nil {
		return err
	}
	return nil
}

func(repo *repositoryImpl) GetUserByDomain(domain string) (*data.User, error ) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second *5)
	defer cancel()
	filter := bson.M{
		"settings.subdomain":  domain,
	}
	res := repo.db.Collection("users").FindOne(ctx, filter)
	if err :=res.Err(); err != nil{
		return nil, err
	}
	user := &data.User{}
	if err := res.Decode(user); err != nil{
		return nil, err
	}
	return user, nil
}
func (repo *repositoryImpl) UpdateUserInformation(userName, fullName, email, location string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{
		"userName" : userName,
	}

	update := bson.M{
		"$set"  :bson.M{
			"fullName" : fullName,
			"email": email,
			"location": location,
		},
	}

	res := repo.db.Collection("users").FindOneAndUpdate(ctx, filter, update)
	if err :=res.Err(); err != nil{
		return err
	}
	return nil
}
