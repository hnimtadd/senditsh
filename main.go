package main

import (
	"log"

	"github.com/hnimtadd/senditsh/api"
	"github.com/hnimtadd/senditsh/config"
	server "github.com/hnimtadd/senditsh/server"
)

func main() {
	sshConfig, err := config.GetSSHConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	api, err := api.NewAPIHandlerImpl()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		sshServer, err := server.NewSSHServerImpl(api, sshConfig)
		if err != nil {
			log.Fatal(err)
		}
		if err := sshServer.Listen(); err != nil {
			log.Fatal(err)
		}
	}()

	httpConfig, err := config.GetHTTPConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	httpServer, err := server.NewHTTPServerImpl(api, httpConfig)
	if err != nil {
		log.Fatal(err)
	}
	if err := httpServer.Listen(); err != nil {
		log.Fatal(err)
	}
}
