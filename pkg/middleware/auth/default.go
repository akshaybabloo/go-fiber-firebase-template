package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/akshaybabloo/go-fiber-template/model"
	"github.com/akshaybabloo/go-fiber-template/pkg/factory"
)

// Auth middleware verifies your access tokens before proceeding to the request
func Auth(f *factory.Factory) fiber.Handler {

	return func(c fiber.Ctx) error {

		firebaseConfig := f.FirebaseConfig()
		defer f.Zap.Sync()

		authToken := c.Get(fiber.HeaderAuthorization)
		if authToken == "" {
			return c.Status(http.StatusForbidden).JSON(model.ProblemDetails{
				Type:     "https://tools.ietf.org/html/rfc7231#section-6.5.3",
				Status:   http.StatusForbidden,
				Detail:   "Authorization header is missing",
				Instance: c.OriginalURL(),
				Title:    "Header missing",
			})
		}

		splitToken := strings.Split(authToken, "Bearer")
		if len(splitToken) != 2 {
			return c.Status(http.StatusForbidden).JSON(model.ProblemDetails{
				Type:     "https://tools.ietf.org/html/rfc7231#section-6.5.3",
				Status:   http.StatusForbidden,
				Detail:   "Malformed token, make sure you follow - 'Bearer <token>' in 'Authorization' header",
				Instance: c.OriginalURL(),
				Title:    "Invalid token",
			})
		}

		reqToken := strings.TrimSpace(splitToken[1])

		ctx := context.Background()
		client, err := firebaseConfig.FirebaseAuth()
		if err != nil {
			f.Zap.Sugar().Errorf("error getting auth app: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(model.ProblemDetails{
				Type:     "https://tools.ietf.org/html/rfc7231#section-6.6.1",
				Status:   http.StatusInternalServerError,
				Detail:   "Internal server error",
				Instance: c.OriginalURL(),
				Title:    "Server error",
			})
		}

		token, err := client.VerifyIDToken(ctx, reqToken)
		if err != nil {
			f.Zap.Sugar().Errorf("error verifying ID token: %v", err)
			return c.Status(http.StatusUnauthorized).JSON(model.ProblemDetails{
				Type:     "https://tools.ietf.org/html/rfc7235#section-3.1",
				Status:   http.StatusUnauthorized,
				Detail:   "Invalid token provided",
				Instance: c.OriginalURL(),
				Title:    "Unauthorised request",
			})
		}

		c.Locals("userToken", token)

		return c.Next()
	}
}
