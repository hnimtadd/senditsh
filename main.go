package main

import (
	"log"
	// "time"

	"github.com/hnimtadd/senditsh/api"
	"github.com/hnimtadd/senditsh/config"
	"github.com/hnimtadd/senditsh/repository"
	server "github.com/hnimtadd/senditsh/server"
	// "github.com/sujit-baniya/flash"
)

func main() {
	// flash.Default(flash.Config{
	// 	Name:     "defaultFlash",
	// 	Expires:  time.Now().Add(15 * time.Minute),
	// 	Secure:   false,
	// 	HTTPOnly: true,
	// 	SameSite: "Lax",
	// })
	repoConfig, err := config.GetMongoConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	repo, err := repository.NewRepositoryImpl(repoConfig)
	if err != nil {
		log.Fatal(err)
	}

	api, err := api.NewAPIHandlerImpl(repo)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		sshConfig, err := config.GetSSHConfig(".")
		if err != nil {
			log.Fatal(err)
		}
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

	oauthConfig, err := config.GetGithubConfig(".")
	if err != nil{
		log.Fatal(err)
	}
	httpServer, err := server.NewHTTPServerImpl(api, httpConfig, oauthConfig)
	if err != nil {
		log.Fatal(err)
	}
	if err := httpServer.Listen(); err != nil {
		log.Fatal(err)
	}
}
