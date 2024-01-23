package config

import (
	"os"

	masterdatacfg "github.com/diezfx/idlegame-backend/internal/config/masterdata"
	postgrescfg "github.com/diezfx/idlegame-backend/internal/config/postgres"
	supabasecfg "github.com/diezfx/idlegame-backend/internal/config/supabase"
	"github.com/diezfx/idlegame-backend/pkg/auth"
	"github.com/diezfx/idlegame-backend/pkg/configloader"
	"github.com/diezfx/idlegame-backend/pkg/db"
	"github.com/diezfx/idlegame-backend/pkg/masterdata"
)

type Environment string

const (
	LocalEnv       Environment = "local"
	DevelopmentEnv Environment = "dev"
)

type Config struct {
	Addr        string
	Environment Environment
	LogLevel    string
	Auth        auth.Config
	DB          db.Config
	Masterdata  masterdata.Config
}

func Load() (Config, error) {
	env := os.Getenv("ENVIRONMENT")

	// read from stuff json
	loader := configloader.NewFileLoader("/etc/config", "/etc/secrets")

	mstrdata := masterdatacfg.LoadConfig()

	// read postgres secrets

	if env == string(DevelopmentEnv) {
		pgDB, err := postgrescfg.LoadPostgresConfig(loader)
		if err != nil {
			return Config{}, err
		}

		authCfg, err := supabasecfg.LoadSupabaseConfig(loader)
		if err != nil {
			return Config{}, err
		}

		return Config{
			Addr:        ":8080",
			Environment: DevelopmentEnv,
			LogLevel:    "debug",
			DB:          pgDB,
			Auth:        authCfg,
			Masterdata:  mstrdata,
		}, nil
	}

	return Config{
		Addr:        "localhost:5002",
		Environment: LocalEnv,
		LogLevel:    "debug",
		DB: db.Config{
			Port: 5432, Host: "localhost", Database: "postgres",
			Username: "postgres", Password: "postgres",
			MigrationsDir: "db/migrations",
		},
		Masterdata: mstrdata,
	}, nil
}

func (cfg *Config) IsLocal() bool {
	return cfg.Environment == LocalEnv
}
