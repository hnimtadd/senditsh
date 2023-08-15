package server

import (
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	loggmidleware "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/django/v3"
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
	oauthCofig *config.GithubConfig
	auth   *api.AuthenticationService
}

func NewHTTPServerImpl(api *api.ApiHandlerImpl, config *config.HTTPConfig, oauthConfig *config.GithubConfig) (HTTPServer, error) {
	server := &HTTPServerImpl{
		api:    api,
		config: config,
		oauthCofig: oauthConfig,
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
func (server *HTTPServerImpl) createEngine() *django.Engine {
	engine := django.New("views", ".html")
	engine.Reload(true)
	engine.AddFunc("css", func(name string) (res template.HTML) {
		filepath.Walk("static/assets", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == name {
				res = template.HTML("<link rel=\"stylesheet\" href=\"" + path + "\">")
			}
			return nil

		})
		return
	})
	engine.AddFunc("copyToClipboard", func(content string) {
		clipboard.WriteAll(content)
	})
	return engine

}

func (server *HTTPServerImpl) initConnection() error {
	// Initialize standard Go html template engine
	engine := server.createEngine()
	config := fiber.Config{
		ReadTimeout:           time.Second * 5,
		WriteTimeout:          time.Second * 5,
		ErrorHandler:          api.CustomHTTPErrorHandler(),
		PassLocalsToViews:     true,
		DisableStartupMessage: true,
		Views:                 engine,
	}

	server.app = fiber.New(config)

	j := api.NewJWTService(
		[]byte("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDwUPgPGvAdrasDvgQPkxFkb05HuLqVObkkRT9NO1/zf hhhh"),
		server.api,
	)
	auth := &api.AuthenticationService{
		JWT:              j,
		AuthorizationURL: "https://github.com/login/oauth/authorize",
		AccessTokenURL:   "https://github.com/login/oauth/access_token",
		ClientID: server.oauthCofig.ClientId,
		ClientSecret:server.oauthCofig.ClientSecret,
		RedirectURI:      "http://mysendit.sh/callback",
		Api:              server.api,
	}
	server.auth = auth
	server.initRoute()
	return nil
}

func (server *HTTPServerImpl) initRoute() {
	server.app.Static("/static", "./static")
	server.app.Use(favicon.New())
	server.app.Use(api.WithFlash)
	server.app.Use(server.auth.JWT.AuthMiddleware())
	server.app.Use(loggmidleware.New())

	server.app.Use("/usersubdomain/:u", server.api.GetUserDomainPageHandler())
	server.app.Get("/callback", server.auth.CallBackHandler())
	server.app.Get("/api/v1/transfer/:id", server.api.FileTransferHandler())
	server.app.Get("/api/v1/users", server.api.GetUsersHandler())
	server.app.Post("/api/v1/user", server.api.SignUpUserHandler())
	server.app.Get("/api/v1/get-transfers", server.api.GetTransfersHandler())

	server.app.Get("/", server.api.IndexPageHandler())
	server.app.Get("/signin", server.auth.LoginHandler())
	server.app.Get("/signout", server.auth.LogoutHandler())
	server.app.Get("/login", server.api.LoginPageHandler())
	// server.app.Get("/download/:userName", server.api.DownloadPageHandler())

	server.app.Get("/:u/download", server.api.DownloadPageHandler())
	server.app.Get("/user", server.api.MustAuthMiddleware(), server.api.UserPageHandler())
	server.app.Get("/user/transfers", server.api.MustAuthMiddleware(), server.api.GetUserTransferPagehandler())

	server.app.Post("/user/settings-edit", server.api.MustAuthMiddleware(), server.api.PostSettingsEditPageHandler())
	server.app.Get("/user/settings-edit", server.api.MustAuthMiddleware(), server.api.GetSettingsEditPageHandler())
	server.app.Get("/user/settings", server.api.MustAuthMiddleware(), server.api.GetSettingsPageHandler())

	server.app.Get("/user/information", server.api.MustAuthMiddleware(), server.api.GetUserInformationPageHandler())

	server.app.Get("/user/information-edit", server.api.MustAuthMiddleware(), server.api.GetUserInformationEditPageHandler())
	server.app.Post("/user/information-edit", server.api.MustAuthMiddleware(), server.api.PostUserInformationEditPageHandler())

	server.app.Get("/user/domain", server.api.MustAuthMiddleware(), server.api.GetUserDomainTrackingPageHandler())

	server.app.Put("/api/v1/register-domain", server.api.MustAuthMiddleware(), server.api.RegisterUserSubDomainSettingHandler())
	server.app.Put("/api/v1/register-sshKey", server.api.MustAuthMiddleware(), server.api.RegisterUserSSHSettingHandler())
	server.app.Get("/api/v1/get-transfer/", server.api.GetTransfersOfUserHandler())
	server.app.Get("/api/v1/get-setting/", server.api.MustAuthMiddleware(), server.api.GetUserSettingHandler())

	server.app.Use(server.api.NotFoundPageHandler())
}
