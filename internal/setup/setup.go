package setup

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/diezfx/idlegame-backend/internal/api"
	"github.com/diezfx/idlegame-backend/internal/config"
	"github.com/diezfx/idlegame-backend/internal/service/inventory"
	"github.com/diezfx/idlegame-backend/internal/service/jobs"
	"github.com/diezfx/idlegame-backend/internal/service/monster"
	"github.com/diezfx/idlegame-backend/internal/storage"
	"github.com/diezfx/idlegame-backend/pkg/db"
	"github.com/diezfx/idlegame-backend/pkg/logger"
	"github.com/diezfx/idlegame-backend/pkg/masterdata"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupSplitService(ctx context.Context) (*http.Server, *jobs.Daemon, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("read config: %w", err)
	}
	if cfg.Environment == config.LocalEnv {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger.Info(ctx).String("config", fmt.Sprint(cfg)).Msg("Loaded config")

	psqlClient, err := db.New(cfg.DB)
	if err != nil {
		return nil, nil, fmt.Errorf("create sqlite client: %w", err)
	}

	storageClient, err := storage.New(ctx, psqlClient)
	if err != nil {
		return nil, nil, fmt.Errorf("create storage client: %w", err)
	}

	masterdataContainer, err := masterdata.New(cfg.Masterdata)
	if err != nil {
		return nil, nil, fmt.Errorf("create masterdata container: %w", err)
	}
	monsterService := monster.New(storageClient)
	jobService := jobs.New(storageClient, monsterService, storageClient, masterdataContainer)
	inventoryService := inventory.New(storageClient)

	jobDaemon := jobs.NewDaemon(jobService)

	router := api.InitAPI(&cfg, masterdataContainer, jobService, inventoryService, monsterService)

	srv := &http.Server{
		Handler: router.Handler,
		Addr:    cfg.Addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv, jobDaemon, nil
}
