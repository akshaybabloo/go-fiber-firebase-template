package factory

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"go.uber.org/zap"

	"github.com/akshaybabloo/go-fiber-template/config"
	"github.com/akshaybabloo/go-fiber-template/pkg/problem"
)

// Factory holds the application's shared, long-lived dependencies. It is created
// once at startup and reused for the lifetime of the process.
type Factory struct {
	Zap      *zap.Logger
	Firebase *config.Firebase
	Problems problem.Problems

	// VerifyIDToken verifies a Firebase ID token. It is a seam so the auth path
	// can be tested with a fake verifier instead of a real Firebase client.
	VerifyIDToken func(ctx context.Context, idToken string) (*auth.Token, error)
}

// New initialises the application dependencies. The provided context is used for
// the one-off client initialisation only.
func New(ctx context.Context, log *zap.Logger) (*Factory, error) {
	firebaseClients, err := config.NewFirebase(ctx, log)
	if err != nil {
		return nil, err
	}

	return &Factory{
		Zap:           log,
		Firebase:      firebaseClients,
		Problems:      problem.NewProblems(),
		VerifyIDToken: firebaseClients.Auth.VerifyIDToken,
	}, nil
}

// Close releases resources held by the factory's dependencies and should be
// called during graceful shutdown.
func (f *Factory) Close() error {
	return f.Firebase.Close()
}
