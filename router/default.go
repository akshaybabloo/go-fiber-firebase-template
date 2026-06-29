package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/akshaybabloo/go-fiber-template/pkg/factory"
)

func HomeRoute(f *factory.Factory) func(ctx *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer f.Zap.Sync()
		f.Zap.Info("Home page")

		return c.JSON(fiber.Map{"Hello": "World"})
	}
}
