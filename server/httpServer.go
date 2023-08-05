package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hnimtadd/senditsh/api"
	"github.com/hnimtadd/senditsh/config"
	log "github.com/hnimtadd/senditsh/logger"
)

var logger = log.GetLogger(log.Info, "SERVER")

type HTTPServer interface {
	Listen() error
}
type HTTPServerImpl struct {
	app    *fiber.App
	api    *api.ApiHandlerImpl
	config *config.HTTPConfig
}

func NewHTTPServerImpl(api *api.ApiHandlerImpl, config *config.HTTPConfig) (HTTPServer, error) {
	server := &HTTPServerImpl{
		api:    api,
		config: config,
	}
	if err := server.initConnection(); err != nil {
		return nil, err
	}
	return server, nil
}

func (server *HTTPServerImpl) Listen() error {
	logger.Info("Listening on address:", server.config.Port)
	if err := server.app.Listen(":" + server.config.Port); err != nil {
		return err
	}
	return nil

}

func (server *HTTPServerImpl) initConnection() error {
	config := fiber.Config{
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}
	server.app = fiber.New(config)
	server.initRoute()
	return nil
}

func (server *HTTPServerImpl) initRoute() {
	server.app.Get("/api/v1/transfer/:id", server.api.FileTransferHandler())
}
