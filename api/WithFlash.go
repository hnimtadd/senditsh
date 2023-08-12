package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"
)

func WithFlash(ctx *fiber.Ctx) error {
	log.Println("WithFlash")
	values := flash.Get(ctx)
	if len(values) == 0 {
		return ctx.Next()
	}
	ctx.Locals("flash", values)
	log.Printf("flash=%v\n", values)
	return ctx.Next()
}
