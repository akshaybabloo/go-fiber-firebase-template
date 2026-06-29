package config

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

// Firebase holds the initialised Firebase clients. The underlying clients are
// safe for concurrent use and are meant to be created once at startup and
// reused for the lifetime of the process.
type Firebase struct {
	App       *firebase.App
	Auth      *auth.Client
	Firestore *firestore.Client

	zap *zap.Logger
}

// NewFirebase initialises the Firebase App together with its Auth and Firestore
// clients.
//
// Credentials are resolved in the following order:
//
//  1. the path in the FIREBASE_CREDENTIALS environment variable, if set;
//  2. Application Default Credentials — e.g. GOOGLE_APPLICATION_CREDENTIALS or
//     the workload identity / metadata server when running on Google Cloud.
//
// When debug mode is enabled the local Firestore emulator is used.
func NewFirebase(ctx context.Context, log *zap.Logger) (*Firebase, error) {
	if viper.GetBool("debug") {
		if err := os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080"); err != nil {
			return nil, fmt.Errorf("setting firestore emulator host: %w", err)
		}
	}

	var opts []option.ClientOption
	if path := os.Getenv("FIREBASE_CREDENTIALS"); path != "" {
		// Restrict to service-account credentials so an unexpected credential
		// type cannot be loaded from this path. Use ADC for other types.
		opts = append(opts, option.WithAuthCredentialsFile(option.ServiceAccount, path))
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, fmt.Errorf("initialising firebase app: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialising firebase auth client: %w", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialising firestore client: %w", err)
	}

	return &Firebase{
		App:       app,
		Auth:      authClient,
		Firestore: firestoreClient,
		zap:       log,
	}, nil
}

// Close releases the resources held by the Firebase clients. It should be
// called once during graceful shutdown.
func (f *Firebase) Close() error {
	if f.Firestore != nil {
		if err := f.Firestore.Close(); err != nil {
			return fmt.Errorf("closing firestore client: %w", err)
		}
	}
	return nil
}
