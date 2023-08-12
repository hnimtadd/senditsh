package utils

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func GetContextVariableWithType[T any](ctx context.Context, v string, def T) T {
	if p := ctx.Value(v); p != nil {
		return p.(T)
	}
	return def
}

func GetLocalWithType[T any](ctx *fiber.Ctx, key string) *T {
	val, ok := ctx.Locals(key).(*T)
	if !ok {
		return nil
	}
	return val
}
