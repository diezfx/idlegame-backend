package setup

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/diezfx/idlegame-backend/internal/api"
	"github.com/diezfx/idlegame-backend/internal/config"
	"github.com/diezfx/idlegame-backend/internal/service"
	"github.com/diezfx/idlegame-backend/internal/storage"
	"github.com/diezfx/idlegame-backend/pkg/logger"
	"github.com/diezfx/idlegame-backend/pkg/postgres"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupSplitService() (*http.Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	ctx := context.Background()
	if cfg.Environment == config.LocalEnv {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger.Info(ctx).String("config", fmt.Sprint(cfg)).Msg("Loaded config")

	psqlClient, err := postgres.New(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("create sqlite client: %w", err)
	}

	_, err = storage.New(ctx, psqlClient)
	if err != nil {
		return nil, fmt.Errorf("create storage client: %w", err)
	}

	projectService := service.New()

	router := api.InitAPI(&cfg, projectService)

	srv := &http.Server{
		Handler: router.Handler,
		Addr:    cfg.Addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv, nil
}