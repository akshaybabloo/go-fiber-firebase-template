package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	fbauth "firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/akshaybabloo/go-fiber-template/pkg/factory"
	"github.com/akshaybabloo/go-fiber-template/pkg/problem"
)

// newTestFactory builds a Factory with an injectable token verifier so the auth
// path can be exercised without a real Firebase client.
func newTestFactory(verify func(ctx context.Context, idToken string) (*fbauth.Token, error)) *factory.Factory {
	return &factory.Factory{
		Zap:           zap.NewNop(),
		Problems:      problem.NewProblems(),
		VerifyIDToken: verify,
	}
}

// errVerify is a verifier that should never be reached.
func errVerify(context.Context, string) (*fbauth.Token, error) {
	return nil, errors.New("verifier should not be called")
}

func TestAuthRejectsMissingHeader(t *testing.T) {
	app := fiber.New()
	app.Use(Auth(newTestFactory(errVerify)))
	app.Get("/", func(c fiber.Ctx) error { return c.SendString("ok") })

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
	if got := resp.Header.Get("WWW-Authenticate"); got != "Bearer" {
		t.Errorf("WWW-Authenticate = %q, want %q", got, "Bearer")
	}
	if got := resp.Header.Get("Content-Type"); got != problem.MIMEApplicationProblemJSON {
		t.Errorf("Content-Type = %q, want %q", got, problem.MIMEApplicationProblemJSON)
	}
}

func TestAuthRejectsMalformedHeader(t *testing.T) {
	app := fiber.New()
	app.Use(Auth(newTestFactory(errVerify)))
	app.Get("/", func(c fiber.Ctx) error { return c.SendString("ok") })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Token abc123")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("status = %d, want %d (body: %s)", resp.StatusCode, http.StatusUnauthorized, body)
	}
}

func TestAuthRejectsBlankBearerToken(t *testing.T) {
	app := fiber.New()
	app.Use(Auth(newTestFactory(errVerify))) // verifier must not be called
	app.Get("/", func(c fiber.Ctx) error { return c.SendString("ok") })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer    ")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestAuthRejectsInvalidToken(t *testing.T) {
	verify := func(context.Context, string) (*fbauth.Token, error) {
		return nil, errors.New("invalid token")
	}
	app := fiber.New()
	app.Use(Auth(newTestFactory(verify)))
	app.Get("/", func(c fiber.Ctx) error { return c.SendString("ok") })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer sometoken")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestAuthAcceptsValidToken(t *testing.T) {
	verify := func(_ context.Context, idToken string) (*fbauth.Token, error) {
		if idToken != "good-token" {
			return nil, errors.New("unexpected token")
		}
		return &fbauth.Token{UID: "user-123"}, nil
	}
	app := fiber.New()
	app.Use(Auth(newTestFactory(verify)))
	app.Get("/", func(c fiber.Ctx) error {
		tok, _ := c.Locals("userToken").(*fbauth.Token)
		if tok == nil || tok.UID != "user-123" {
			return fiber.NewError(http.StatusInternalServerError, "token not stored")
		}
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer good-token")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("status = %d, want %d (body: %s)", resp.StatusCode, http.StatusOK, body)
	}
}
