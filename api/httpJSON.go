package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hnimtadd/senditsh/data"
	"github.com/hnimtadd/senditsh/utils"
)

func (api *ApiHandlerImpl) GetUsersHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		users, err := api.GetUsers()
		if err != nil {
			logger.Error("err", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		if err := json.NewEncoder(ctx.Response().BodyWriter()).Encode(users); err != nil {
			logger.Error("err", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		ctx.Response().Header.Set("Content-Type", "application/json")
		return ctx.SendStatus(fiber.StatusOK)
	}
}

type signUpUserReq struct {
	Email     string `json:"email,omiempty"`
	FullName  string `json:"fullName,omiempty"`
	Username  string `json:"userName,omiempty"`
	Location  string `json:"location,omiempty"`
	PublicKey string `json:"publicKey,omiempty"`
}

func (api *ApiHandlerImpl) SignUpUserHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			req     signUpUserReq
			user    data.User
			sshHash string
			sshKey  string
			err     error
		)
		if err := json.Unmarshal(ctx.Request().Body(), &req); err != nil {
			logger.Error("err", err)
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if req.PublicKey != "" {
			sshHash, sshKey, err = utils.ParsePublicKey(req.PublicKey)
			if err != nil {
				logger.Error("err", err)
				json.NewEncoder(ctx.Response().BodyWriter()).Encode(err)
				return ctx.SendStatus(fiber.StatusBadRequest)
			}
		}

		user = data.User{
			FullName: req.FullName,
			Username: req.Username,
			Location: req.Location,
			Email:    req.Email,
			Settings: data.Settings{
				SSHKey:  sshKey,
				SSHHash: sshHash,
			}}

		if err = api.CreateUser(&user); err != nil {
			logger.Error("err", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (api *ApiHandlerImpl) FileTransferHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id", "")
		if id == "" {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		writer := ctx.Response().BodyWriter()
		// doneCh := make(chan struct{})
		// tunnel := Tunnel{
		// 	Writer: writer,
		// 	DoneCh: doneCh,
		// }
		// api.tunnels[id] <- tunnel
		tunnel, err := api.GetTunnelWithID(id)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		sCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := tunnel.PipeWriter(writer); err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if err := api.WaitForCopyDone(sCtx, id); err != nil {
			logger.Error("err", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		ctx.Response().Header.Add("Content-Type", "application/zip")
		ctx.Response().Header.Add("Content-Disposition", "attachment; filename=sendit.zip")
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (api *ApiHandlerImpl) GetTransfersHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		transfers, err := api.GetTransfers()
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		ctx.Response().Header.Set("Content-Type", "application/json")
		if err := json.NewEncoder(ctx.Response().BodyWriter()).Encode(&transfers); err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (api *ApiHandlerImpl) GetTransfersOfUserHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userName := ctx.Params("userName")
		if userName == "" {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		transfers, err := api.GetTransfersOfUser(userName)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		ctx.Response().Header.Set("Content-Type", "application/json")
		if err := json.NewEncoder(ctx.Response().BodyWriter()).Encode(&transfers); err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}
func (api *ApiHandlerImpl) GetUserSettingHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		userName := ctx.Params("userName")
		if userName == "" {
			logger.Info("msg", "userNull")
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		setting, err := api.GetSettingOfUser(userName)
		if err != nil {
			logger.Error("msg", "Error why retrive user setting", "err", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		ctx.Response().Header.Set("Content-Type", "application/json")
		if err := json.NewEncoder(ctx.Response().BodyWriter()).Encode(&setting); err != nil {

			logger.Error("msg", "Error why decode user setting to http response", "err", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}

type addSSHSettingHTTPReq struct {
	Key string `json:"sshKey"`
}

func (api *ApiHandlerImpl) RegisterUserSSHSettingHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userName := ctx.Params("userName")
		if userName == "" {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		var req addSSHSettingHTTPReq
		if err := json.Unmarshal(ctx.Body(), &req); err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if err := api.RegisterUserSSHKeySetting(userName, req.Key); err != nil {
			logger.Error("msg", fmt.Sprintf("Error occur while register ssh key for user: %v", userName), "err", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}

type addDomainSettingHTTPReq struct {
	Domain string `json:"domain"`
}

func (api *ApiHandlerImpl) RegisterUserSubDomainSettingHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userName := ctx.Params("userName")
		if userName == "" {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		var req addDomainSettingHTTPReq
		if err := json.Unmarshal(ctx.Body(), &req); err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		if err := api.RegisterUserDomainSetting(userName, req.Domain); err != nil {
			logger.Error("msg", fmt.Sprintf("Error occur while register ssh key for user: %v", userName), "err", err)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
