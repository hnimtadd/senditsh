package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (api *ApiHandlerImpl) FileTransferHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id", "")
		if id == "" {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		writer := ctx.Response().BodyWriter()
		doneCh := make(chan struct{})
		tunnel := Tunnel{
			Writer: writer,
			DoneCh: doneCh,
		}
		api.tunnels[id] <- tunnel
		timeout := time.NewTicker(time.Second * 5)
		select {
		case <-doneCh:
			logger.Info("msg", "done")
			ctx.Response().Header.Add("Content-Type", "application/zip")
			ctx.Response().Header.Add("Content-Disposition", "attachment; filename=sendit.zip")
			// ctx.Response().Header.Add("filename", "sendit.zip")
			return ctx.SendStatus(fiber.StatusOK)
		case <-timeout.C:
			logger.Info("msg", "not done")
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
	}
}
