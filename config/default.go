package config

import (
	_ "embed"
)

//go:embed firebase_config.json
var FireBaseConfig []byte
