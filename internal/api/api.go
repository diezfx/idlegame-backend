package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/diezfx/idlegame-backend/internal/config"
	"github.com/diezfx/idlegame-backend/internal/service"
	"github.com/diezfx/idlegame-backend/pkg/auth"
	"github.com/diezfx/idlegame-backend/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	projectService ProjectService
}

func newAPIHandler(projectService ProjectService) *APIHandler {
	return &APIHandler{projectService: projectService}
}

func InitAPI(cfg *config.Config, projectService ProjectService) *http.Server {
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
	_ = newAPIHandler(projectService)
	// r.GET("projects/:id", apiHandler.getProjectByIDHandler)

	return &http.Server{
		Handler: mr,
		Addr:    "localhost:5002",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
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
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusNotFound,
			Reason:    "not found",
		})
	default:
		logger.Error(ctx, err).Msg("unexpected error occurred")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Reason:    "unexpected",
		})
	}
}
