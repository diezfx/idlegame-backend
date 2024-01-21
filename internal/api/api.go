package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/diezfx/idlegame-backend/internal/config"
	"github.com/diezfx/idlegame-backend/internal/service"
	"github.com/diezfx/idlegame-backend/pkg/auth"
	"github.com/diezfx/idlegame-backend/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	jobService       JobService
	inventoryService InventoryService
	monsterService   MonsterService
}

func newAPIHandler(jobService JobService, inventoryService InventoryService, monsterService MonsterService) *APIHandler {
	return &APIHandler{
		jobService:       jobService,
		inventoryService: inventoryService,
		monsterService:   monsterService}
}

func InitAPI(cfg *config.Config, jobService JobService, inventoryService InventoryService, monsterService MonsterService) *http.Server {
	mr := gin.New()
	mr.Use(gin.Recovery())
	mr.Use(logger.HTTPLoggingMiddleware())
	mr.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "OPTION"},
		AllowHeaders:     []string{"Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))
	r := mr.Group("/api/v1.0/")
	if !cfg.IsLocal() {
		r.Use(auth.AuthMiddleware(cfg.Auth))
	}

	// jobHandlers
	api := newAPIHandler(jobService, inventoryService, monsterService)

	r.GET("/monsters/:id", api.GetMonster)
	r.GET("/jobs/:id", api.GetJob)
	woodcuttingRouter := r.Group("/jobs/woodcutting")
	woodcuttingRouter.POST("/", api.PostWoodcuttingJob)
	woodcuttingRouter.GET("/:id", api.GetWoodcuttingJob)

	miningRouter := r.Group("/jobs/mining")
	miningRouter.POST("/", api.PostMiningJob)
	miningRouter.GET("/:id", api.GetMiningJob)

	harvestingRouter := r.Group("/jobs/harvesting")
	harvestingRouter.POST("/", api.PostHarvestingJob)
	harvestingRouter.GET("/:id", api.GetHarvestingJob)

	smeltingRouter := r.Group("/jobs/smelting")
	smeltingRouter.POST("/", api.PostSmeltingJob)
	smeltingRouter.GET("/:id", api.GetHarvestingJob)

	r.DELETE("/jobs/:id", api.DeleteJob)
	r.GET("/inventory/:userID", api.GetInventory)
	return &http.Server{
		Handler: mr,
		Addr:    "localhost:5002",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func (api *APIHandler) GetJob(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		handleError(ctx, errInvalidInput)
		return
	}

	resp, err := api.jobService.GetJob(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (api *APIHandler) DeleteJob(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		handleError(ctx, errInvalidInput)
		return
	}
	err = api.jobService.StopJob(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (api *APIHandler) GetInventory(ctx *gin.Context) {
	userIDStr := ctx.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if userIDStr == "" || err != nil {
		handleError(ctx, errInvalidInput)
		return
	}
	resp, err := api.inventoryService.GetInventory(ctx, userID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (api *APIHandler) GetMonster(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		handleError(ctx, errInvalidInput)
		return
	}
	resp, err := api.monsterService.GetMonsterByID(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toMonster(resp))
}

func handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, errInvalidInput):
		logger.Info(ctx).Err(err).Msg("request failed with invalid input")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			Reason:    "invalid input",
		})
	case errors.Is(err, service.ErrProjectNotFound):
		logger.Info(ctx).Err(err).Msg("not found")
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			ErrorCode: http.StatusNotFound,
			Reason:    "not found",
		})
	case errors.Is(err, service.ErrLevelRequirementNotMet):
		logger.Info(ctx).Err(err).Msg("level requirement not met")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			Reason:    "level requirement not met",
		})
	case errors.Is(err, service.ErrJobTypeNotFound):
		logger.Info(ctx).Err(err).Msg("job type not found")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			Reason:    "job type not found",
		})
	case errors.Is(err, service.ErrNotEnoughItems):
		logger.Info(ctx).Err(err).Msg("not enough items for job")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusNotFound,
			Reason:    "not enough items to start job",
		})
	default:
		logger.Error(ctx, err).Msg("unexpected error occurred")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Reason:    "unexpected",
		})
	}
}
