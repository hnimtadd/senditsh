package repository

import (
	"context"
	"time"

	"github.com/hnimtadd/senditsh/config"
	"github.com/hnimtadd/senditsh/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Repository interface {
}

type repositoryImpl struct {
	db     *mongo.Database
	config *config.MongoConfig
	logger *logger.Logger
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
	logg := &logger.Logger{}
	logg.SetLogLevel(logger.Info).SetLogScope("repository").Create()
	repo.logger = logg

	// TODO: init connection to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().ApplyURI(repo.config.Source).SetAuth(
		options.Credential{
			AuthSource: repo.config.AuthSource,
			Username:   repo.config.Username,
			Password:   repo.config.Password,
		},
	).SetTLSConfig(nil).SetTimeout(5 * time.Second)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	repo.db = client.Database(repo.config.Database)
	repo.logger.DefaultLog(logger.Info, "Initialized database")
	return nil
}
