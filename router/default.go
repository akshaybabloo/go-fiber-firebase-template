package router

import (
	"github.com/gofiber/fiber/v3"

	"github.com/akshaybabloo/go-fiber-template/pkg/factory"
)

func HomeRoute(f *factory.Factory) fiber.Handler {
	return func(c fiber.Ctx) error {
		f.Zap.Info("Home page")

		return c.JSON(fiber.Map{"Hello": "World"})
	}
}
