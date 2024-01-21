package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *APIHandler) GetMiningJob(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		handleError(ctx, errInvalidInput)
		return
	}
	resp, err := api.jobService.GetMiningJob(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (api *APIHandler) GetHarvestingJob(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		handleError(ctx, errInvalidInput)
		return
	}
	resp, err := api.jobService.GetWoodcuttingJob(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (api *APIHandler) PostHarvestingJob(ctx *gin.Context) {
	var req StartHarvestingJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleError(ctx, errInvalidInput)
		return
	}
	resp, err := api.jobService.StartHarvestingJob(ctx, req.UserID, req.Monster, req.CropType)
	if err != nil {
		handleError(ctx, err)
		return
	}
	// resource created
	// set header to url with id
	ctx.Header("Location", fmt.Sprintf("/api/v1.0/jobs/mining/%d", resp))
	ctx.JSON(http.StatusCreated, resp)
}

// postminingjob
func (api *APIHandler) PostMiningJob(ctx *gin.Context) {
	var req StartMiningJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleError(ctx, errInvalidInput)
		return
	}
	resp, err := api.jobService.StartMiningJob(ctx, req.UserID, req.Monster, req.OreType)
	if err != nil {
		handleError(ctx, err)
		return
	}
	// resource created
	// set header to url with id
	ctx.Header("Location", fmt.Sprintf("/api/v1.0/jobs/mining/%d", resp))
	ctx.JSON(http.StatusCreated, resp)
}

func (api *APIHandler) GetWoodcuttingJob(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		handleError(ctx, errInvalidInput)
		return
	}
	resp, err := api.jobService.GetWoodcuttingJob(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// postwoodcuttingjob
func (api *APIHandler) PostWoodcuttingJob(ctx *gin.Context) {
	var req StartWoodCuttingJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleError(ctx, errInvalidInput)
		return
	}
	resp, err := api.jobService.StartWoodCuttingJob(ctx, req.UserID, req.Monster, req.TreeType)
	if err != nil {
		handleError(ctx, err)
		return
	}
	// resource created
	// set header to url with id
	ctx.Header("Location", fmt.Sprintf("/api/v1.0/jobs/woodcutting/%d", resp))
	ctx.JSON(http.StatusCreated, resp)
}
