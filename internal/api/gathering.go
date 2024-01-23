package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *APIHandler) PostProcessingJob(ctx *gin.Context) {
	var req StartProcessingJob
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleError(ctx, errInvalidInput)
		return
	}
	resp, err := api.jobService.StartProcessingJob(ctx, req.UserID, req.Monster, req.JobDefID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	// resource created
	// set header to url with id
	ctx.Header("Location", fmt.Sprintf("/api/v1.0/jobs/mining/%d", resp))
	ctx.JSON(http.StatusCreated, resp)
}

func (api *APIHandler) PostGatheringJob(ctx *gin.Context) {
	var req StartGatheringJob
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleError(ctx, errInvalidInput)
		return
	}
	resp, err := api.jobService.StartGatheringJob(ctx, req.UserID, req.Monster, req.JobDefID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	// resource created
	// set header to url with id
	ctx.Header("Location", fmt.Sprintf("/api/v1.0/jobs/woodcutting/%d", resp))
	ctx.JSON(http.StatusCreated, resp)
}
