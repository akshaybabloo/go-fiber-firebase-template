package config

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

// type FirebaseConfig struct {
// 	Type                    string `json:"type"`
// 	ProjectID               string `json:"project_id"`
// 	PrivateKeyID            string `json:"private_key_id"`
// 	PrivateKey              string `json:"private_key"`
// 	ClientEmail             string `json:"client_email"`
// 	ClientID                string `json:"client_id"`
// 	AuthURL                 string `json:"auth_url"`
// 	TokenURL                string `json:"token_url"`
// 	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
// 	ClientX509CertURL       string `json:"client_x509_cert_url"`
// }

type FirebaseConfig interface {
	FirebaseApp() (*firebase.App, error)
	FirebaseAuth() (*auth.Client, error)
	FirebaseFirestore() (*firestore.Client, error)
}

func NewFirebaseConfig(zap *zap.Logger) FirebaseConfig {
	return &firebaseStubs{
		zap: zap,
	}
}

type firebaseStubs struct {
	FirebaseConfig

	zap *zap.Logger
}

func (f *firebaseStubs) FirebaseApp() (*firebase.App, error) {

	if viper.GetBool("debug") {
		err := os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")
		// err = os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099")
		if err != nil {
			return nil, err
		}
	}

	ctx := context.Background()
	opt := option.WithCredentialsJSON(FireBaseConfig)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (f *firebaseStubs) FirebaseAuth() (*auth.Client, error) {
	ctx := context.Background()
	app, err := f.FirebaseApp()
	if err != nil {
		return nil, err
	}
	// Get an auth client from the firebase.App
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (f *firebaseStubs) FirebaseFirestore() (*firestore.Client, error) {
	ctx := context.Background()
	app, err := f.FirebaseApp()
	if err != nil {
		return nil, err
	}
	// Get an auth client from the firebase.App
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// func FirebaseStorage() (*storage.Client, error) {
//	ctx := context.Background()
//	app, err := FirebaseApp()
//	if err != nil {
//		return nil, err
//	}
//	// Get an auth client from the firebase.App
//	client, err := app.Storage(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	return client, nil
// }
//
// func FirebaseConfiguration() (FirebaseConfig, error) {
//	var config FirebaseConfig
//	jsonFile, err := os.Open("firebase_config.json")
//	if err != nil {
//		return config, err
//	}
//	all, err := ioutil.ReadAll(jsonFile)
//	if err != nil {
//		return config, err
//	}
//	err = json.Unmarshal(all, &config)
//	if err != nil {
//		return config, err
//	}
//	return config, nil
// }
