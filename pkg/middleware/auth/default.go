package auth

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/akshaybabloo/go-fiber-template/pkg/factory"
	"github.com/akshaybabloo/go-fiber-template/pkg/problem"
)

const bearerPrefix = "Bearer "

// Auth middleware verifies the Firebase ID token in the Authorization header
// before proceeding to the request. The verified token is stored in the
// request locals under "userToken".
func Auth(f *factory.Factory) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get(fiber.HeaderAuthorization)
		if authHeader == "" {
			return unauthorized(c, f, "Header missing", "Authorization header is missing")
		}

		if len(authHeader) <= len(bearerPrefix) || !strings.EqualFold(authHeader[:len(bearerPrefix)], bearerPrefix) {
			return unauthorized(c, f, "Invalid token",
				"Malformed token, make sure you follow - 'Bearer <token>' in 'Authorization' header")
		}

		reqToken := strings.TrimSpace(authHeader[len(bearerPrefix):])
		if reqToken == "" {
			return unauthorized(c, f, "Invalid token",
				"Malformed token, make sure you follow - 'Bearer <token>' in 'Authorization' header")
		}

		token, err := f.VerifyIDToken(c.Context(), reqToken)
		if err != nil {
			f.Zap.Error("verifying ID token", zap.Error(err))
			return unauthorized(c, f, "Unauthorised request", "Invalid token provided")
		}

		c.Locals("userToken", token)

		return c.Next()
	}
}

// unauthorized writes an RFC 9457 problem response with a 401 status and the
// WWW-Authenticate header, as required for missing or invalid credentials.
func unauthorized(c fiber.Ctx, f *factory.Factory, title, detail string) error {
	c.Set(fiber.HeaderWWWAuthenticate, "Bearer")
	return c.Status(fiber.StatusUnauthorized).JSON(
		f.Problems.UnauthorizedProblem(title, detail, c.OriginalURL()),
		problem.MIMEApplicationProblemJSON,
	)
}
