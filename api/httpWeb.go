package api

import (
	"fmt"

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
		context := fiber.Map{}
		user := utils.GetLocalWithType[User](ctx, "user")
		if user == nil {
			context["informationError"] = "Can't get user"
			return ctx.Render("user/information", context)
		}
		logger.Info("user", &user)
		context["user"] = *user
		return ctx.Render("user/information", context)
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

func (api *ApiHandlerImpl) GetSettingsEditPageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		context := fiber.Map{}
		user := utils.GetLocalWithType[User](ctx, "user")
		if user == nil {
			context["settingError"] = "Cann't get user"
			return ctx.Render("user/settingsEdit", context)
		}
		logger.Info("user", &user)
		context["user"] = *user
		return ctx.Render("user/settingsEdit", context)
	}
}

func (api *ApiHandlerImpl) PostSettingsEditPageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		errContext := fiber.Map{}

		user := utils.GetLocalWithType[User](ctx, "user")
		if user == nil {
			errContext["settingError"] = "Cann't get user information"
			logger.Error("err", "Cann't get user information")
			return ctx.Render("user/settingsEdit", errContext)
		}

		domain := ctx.FormValue("domain")
		if domain != "" {
			if err := api.RegisterUserDomainSetting(user.Username, domain); err != nil {
				errContext["settingError"] = err.Error()
				logger.Error("err", err.Error())
				return ctx.Render("user/settingsEdit", errContext)
			}
		}

		publickey := ctx.FormValue("publicKey")
		if publickey != "" {
			if err := api.RegisterUserSSHKeySetting(user.Username, publickey); err != nil {
				errContext["settingError"] = err.Error()
				logger.Error("err", err.Error())
				return ctx.Render("user/settingsEdit", errContext)
			}
		}

		return ctx.Redirect("/user/settings")
	}
}
func (api *ApiHandlerImpl) LoginPageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Render("login/index", fiber.Map{})
	}
}
func (api *ApiHandlerImpl) GetUserInformationEditPageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Render("user/informationEdit", fiber.Map{})
	}
}

func (api *ApiHandlerImpl) PostUserInformationEditPageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		errContext := fiber.Map{}
		user := utils.GetLocalWithType[User](ctx, "user")
		if user == nil {
			errContext["informationError"] = "Can't get user"
			return flash.WithData(ctx, errContext).Redirect("/user/information-edit")
		}
		fullName := ctx.FormValue("fullName")
		email := ctx.FormValue("email")
		location := ctx.FormValue("location")
		logger.Info("msg", fmt.Sprintf("Update user with information: Username: %v, email: %v, location: %v", fullName, email, location))
		if err := api.UpdateUserInformation(user.Username, fullName, email, location); err != nil {
			errContext["informationError"] = err.Error()
			return flash.WithData(ctx, errContext).Redirect("/user/information-edit")
		}
		return ctx.Redirect("/user/information")
	}
}

func (api *ApiHandlerImpl) GetUserDomainPageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		errContext := fiber.Map{}
		userName := ctx.Params("u")
		if userName == "" {
			errContext["userDomainError"] = "UserName must valid"
			return ctx.Render("domain/notfound", errContext)
		}

		user, err := api.GetUserByDomain(userName)
		if err != nil {
			errContext["userDomainError"] = "User not exits"
			return ctx.Render("domain/notfound", errContext)
		}

		sharing, err := api.GetUserSharingLink(user.Username)
		if err != nil {
			errContext["userDomainError"] = err.Error()
		}
		return ctx.Render("domain/domain", fiber.Map{
			"domainUser":    user,
			"domainSharing": sharing,
		})
	}
}
