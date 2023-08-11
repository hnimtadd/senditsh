package repository

import (
	"context"
	"log"
	"time"

	"github.com/hnimtadd/senditsh/config"
	"github.com/hnimtadd/senditsh/data"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Repository interface {
	InsertTransfer(transfer *data.Transfer) error
	GetTransfersOfUser(id string) ([]data.Transfer, error)
	GetTransfers() ([]data.Transfer, error)

	CreateUser(user *data.User) error
	GetUsers() ([]data.User, error)
	GetUserByUserName(userName string) (*data.User, error)
	GetUserByPublicKey(publicKey string) (*data.User, error)
	GetSettingOfUser(userName string) (*data.Settings, error)
	UpdateUserSetting(userName string, setting *data.Settings) error
	InsertUserSSHKey(userName string, sshKey string, sshHash string) error
	InsertUserDomain(userName string, domain string) error
}

type repositoryImpl struct {
	db     *mongo.Database
	config *config.MongoConfig
}

func NewRepositoryImpl(config *config.MongoConfig) (Repository, error) {
	repo := &repositoryImpl{
		config: config,
	}
	if err := repo.InitRepo(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (repo *repositoryImpl) InitRepo() error {
	// TODO: init connection to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().ApplyURI(repo.config.Source)
	// .SetAuth(
	// 	options.Credential{
	// 		AuthSource:    repo.config.AuthSource,
	// 		Username:      repo.config.Username,
	// 		Password:      repo.config.Password,
	// 		AuthMechanism: repo.config.AuthMechanism,
	// 	},
	// ).SetTLSConfig(nil).SetTimeout(5 * time.Second)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	repo.db = client.Database(repo.config.Database)
	log.Printf("Connected to mongodb with config: %v\n", repo.config)
	return nil
}
