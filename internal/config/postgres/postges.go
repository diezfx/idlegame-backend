package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/diezfx/idlegame-backend/pkg/configloader"
	"github.com/diezfx/idlegame-backend/pkg/db"
)

const defaultNamespace = "postgres"

func LoadPostgresConfig(loader *configloader.Loader) (db.Config, error) {
	cfg := db.Config{}

	content, err := loader.LoadConfig(defaultNamespace)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("unmarshal postgres: %w", err)
	}

	username, err := loader.LoadSecret(defaultNamespace, "username")
	if err != nil {
		return cfg, err
	}
	password, err := loader.LoadSecret(defaultNamespace, "password")
	if err != nil {
		return cfg, err
	}
	cfg.Username = username
	cfg.Password = password
	return cfg, nil
}
