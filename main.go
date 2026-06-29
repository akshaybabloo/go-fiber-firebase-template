package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/akshaybabloo/go-fiber-template/pkg/factory"
	"github.com/akshaybabloo/go-fiber-template/pkg/middleware/auth"
	"github.com/akshaybabloo/go-fiber-template/pkg/problem"
	"github.com/akshaybabloo/go-fiber-template/router"
)

func init() {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config.yaml file not found: %w", err))
		}
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log, err := newLogger()
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	// Cancel the context on SIGINT/SIGTERM so the server shuts down gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	f, err := factory.New(ctx, log)
	if err != nil {
		log.Fatal("initialising dependencies", zap.Error(err))
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Error("closing dependencies", zap.Error(err))
		}
	}()

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler(f),
	})

	if viper.GetBool("debug") {
		app.Use(logger.New())
	} else {
		app.Use(logger.New(logger.Config{Format: `{"pid":${pid}, "timestamp":"${time}", "status":${status}, "latency":"${latency}", "method":"${method}", "path":"${path}"}` + "\n"}))
	}
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	// NOTE: cors.New() allows all origins by default. Restrict this for production.
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(requestid.New())
	app.Use(auth.Auth(f))

	// TODO: Add routes here
	app.Get("/", router.HomeRoute(f))

	// handles 404s
	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(
			f.Problems.PageNotFoundProblem(c.OriginalURL()),
			problem.MIMEApplicationProblemJSON,
		)
	})

	listenConfig := fiber.ListenConfig{
		GracefulContext:       ctx,
		DisableStartupMessage: !viper.GetBool("debug"),
	}

	if err := app.Listen(":"+port, listenConfig); err != nil {
		log.Fatal("server stopped", zap.Error(err))
	}
}

// newLogger builds a zap logger appropriate for the configured mode.
func newLogger() (*zap.Logger, error) {
	if viper.GetBool("debug") {
		log, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		log.Info("Running in debug mode")
		return log, nil
	}

	log, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	log.Info("Running in production mode")
	return log, nil
}

// errorHandler returns a Fiber error handler that renders errors as RFC 9457
// problem details, preserving the HTTP status from *fiber.Error.
func errorHandler(f *factory.Factory) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		title := "Unknown error"
		detail := "An error occurred internally and we are looking into it"

		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
			title = "Request error"
			detail = e.Message
		} else {
			// Unexpected error: keep the client-facing response generic but
			// preserve the underlying failure in the logs.
			f.Zap.Error("unhandled error", zap.String("path", c.OriginalURL()), zap.Error(err))
		}

		return c.Status(code).JSON(
			f.Problems.Problem(detail, c.OriginalURL(), code, title, ""),
			problem.MIMEApplicationProblemJSON,
		)
	}
}
