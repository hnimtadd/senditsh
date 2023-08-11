package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func CustomHTTPErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if he, ok := err.(*fiber.Error); ok {
			code = he.Code
		}
		// fmt.Println(err.Error())
		if e := ctx.Status(code).SendFile(fmt.Sprintf("./views/%d.html", code)); e != nil {
			// fmt.Println(e.Error())
			ctx.Render("./views/error.html", fiber.Map{
				"Error": err.Error(),
			})
			return nil
		}
		return nil
	}
}
