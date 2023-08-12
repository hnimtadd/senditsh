package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hnimtadd/senditsh/utils"
	"github.com/sujit-baniya/flash"
)

func (api *ApiHandlerImpl) IndexPageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if err := ctx.Render("index", fiber.Map{}); err != nil {
			logger.Error("Error", err)
		}
		return nil
	}
}

func (api *ApiHandlerImpl) DownloadPageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger.Info("pass into download")

		flashData := flash.Get(ctx)

		user := utils.GetLocalWithType[User](ctx, "user")
		if user == nil {
			logger.Error("cannot parse user")
		} else {
			logger.Info("userInfo", user)
		}
		flashData["user"] = user
		return ctx.Render("download", flashData)
	}
}

func (api *ApiHandlerImpl) GetSettingsPageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		data := fiber.Map{}
		user := utils.GetLocalWithType[User](ctx, "user")
		if user == nil {
			data["settingError"] = "Cannot get user from context"
			return flash.WithData(ctx, data).Redirect("/user")
		}
		logger.Info("user", user)
		settings, err := api.GetSettingOfUser(user.Username)
		if err != nil {
			data["settingError"] = err
			return flash.WithData(ctx, data).Redirect("/user")
		}
		logger.Info("settings", settings)

		return ctx.Render("user/settings", fiber.Map{
			"userSetting": settings,
		})
	}
}

func (api *ApiHandlerImpl) GetUserTransferPagehandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		flashData := fiber.Map{}
		user := utils.GetLocalWithType[User](ctx, "user")
		if user == nil {
			flashData["transfersError"] = "Cann't get user"
			return flash.WithData(ctx, flashData).Redirect("/user/transfer")
		}
		flashData["tab"] = "transfer"

		transfers, err := api.GetTransfersOfUser(user.Username)
		if err != nil {
			flashData["error"] = err
			return flash.WithData(ctx, flashData).Redirect("/user")
		}
		flashData["transfers"] = transfers

		return ctx.Render("user/transfers", flashData)
	}
}
func (api *ApiHandlerImpl) UserPageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		data := flash.Get(ctx)
		logger.Info("data", data)
		cookieValue := ctx.Cookies("defaultFlash")
		logger.Info("cookieValue", cookieValue)
		return ctx.Render("user/index", fiber.Map{})
	}
}

func (api *ApiHandlerImpl) GetUserInformationPageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func (api *ApiHandlerImpl) NotFoundPageHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("error/404", fiber.Map{})
	}
}

func (api *ApiHandlerImpl) GetUserDomainTrackingPageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		context := fiber.Map{}
		user := utils.GetLocalWithType[User](ctx, "user")
		if user == nil {
			context["domainError"] = "Cann't get user"
			return ctx.Render("user/domain", context)
		}
		logger.Info("user", &user)
		context["user"] = *user
		return ctx.Render("user/domain", context)
	}
}
