package factory

import (
	"go.uber.org/zap"

	"github.com/akshaybabloo/go-fiber-template/config"
	"github.com/akshaybabloo/go-fiber-template/pkg/problem"
)

type Factory struct {
	FirebaseConfig func() config.FirebaseConfig
	Problems       func() problem.Problems
	Zap            *zap.Logger
}

func New(zap *zap.Logger) *Factory {
	return &Factory{
		FirebaseConfig: func() config.FirebaseConfig {
			return config.NewFirebaseConfig(zap)
		},
		Problems: problem.NewProblems,
		Zap:      zap,
	}
}
