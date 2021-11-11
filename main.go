package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/akshaybabloo/go-fiber-template/pkg/factory"
	"github.com/akshaybabloo/go-fiber-template/pkg/middleware/auth"
	"github.com/akshaybabloo/go-fiber-template/router"
)

func init() {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config.yaml file not found: %w \n", err))
		} else {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	var err error
	var log *zap.Logger
	if viper.GetBool("debug") {
		log, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		log.Sugar().Debug("Running in debug mode")
	} else {
		log, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
		log.Sugar().Info("Running in production mode")
	}

	f := factory.New(log)
	problems := f.Problems()

	config := fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			// Check if it's an fiber.Error type
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code

				// Return HTTP response
				return c.Status(code).JSON(problems.InternalServerErrorProblem("Internal server error", err.Error(), c.OriginalURL()))
			}

			return c.Status(fiber.StatusInternalServerError).JSON(problems.InternalServerErrorProblem("Unknown error", "An error occurred internally and we are looking into it", c.OriginalURL()))
		},
	}

	if !viper.GetBool("debug") {
		config.DisableStartupMessage = true
	}

	app := fiber.New(config)
	if !viper.GetBool("debug") {
		app.Use(logger.New(logger.Config{Format: `{"pid":${pid}, "timestamp":"${time}", "status":${status}, "latency":"${latency}", "method":"${method}", "path":"${path}"}` + "\n"}))
	} else {
		app.Use(logger.New())
	}
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(requestid.New())

	if !viper.GetBool("debug") {
		app.Use(auth.Auth(f))
	}

	// TODO: Add routers here
	app.Get("/", router.Home(f))

	// handles 404s
	app.Use(func(c *fiber.Ctx) error {
		problems := f.Problems()
		return c.Status(fiber.StatusNotFound).JSON(problems.PageNotFoundProblem(c.OriginalURL()))
	})

	log.Sugar().Fatal(app.Listen(port))

}
