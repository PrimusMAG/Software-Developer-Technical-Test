package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"football-api/internal/domains/matches/dtos"
	"football-api/internal/domains/matches/services"
	"football-api/internal/shared/pagination"
	"football-api/internal/shared/response"
	sharedvalidator "football-api/internal/shared/validator"
)

type Controller struct{ service services.Service }

func New(service services.Service) *Controller { return &Controller{service: service} }

func (ctl *Controller) Create(c *gin.Context) {
	var req dtos.CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", "INVALID_BODY")
		return
	}
	if err := sharedvalidator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}
	data, err := ctl.service.Create(req)
	if err != nil {
		response.ErrorFrom(c, err, http.StatusBadRequest, "CREATE_MATCH_FAILED", "failed to create match")
		return
	}
	response.Created(c, "match created", data)
}

func (ctl *Controller) List(c *gin.Context) {
	p := pagination.Parse(c)
	status := c.Query("status")
	homeID, _ := strconv.ParseUint(c.Query("home_team_id"), 10, 64)
	awayID, _ := strconv.ParseUint(c.Query("away_team_id"), 10, 64)

	var dateFrom, dateTo *time.Time
	if v := c.Query("date_from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			dateFrom = &t
		}
	}
	if v := c.Query("date_to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			dateTo = &t
		}
	}

	data, total, err := ctl.service.List(p.Offset, p.Limit, status, uint(homeID), uint(awayID), dateFrom, dateTo)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed fetch matches", "FETCH_MATCH_FAILED")
		return
	}
	c.Header("Cache-Control", "public, max-age=60, stale-while-revalidate=120")
	response.SuccessWithMeta(c, "matches fetched", data, response.Meta{
		Page: p.Page, Limit: p.Limit, Total: total, TotalPages: pagination.TotalPages(total, p.Limit),
	})
}

func (ctl *Controller) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	data, err := ctl.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "match not found", "MATCH_NOT_FOUND")
		return
	}
	response.Success(c, "match fetched", data)
}

func (ctl *Controller) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := ctl.service.Delete(uint(id)); err != nil {
		response.ErrorFrom(c, err, http.StatusBadRequest, "DELETE_MATCH_FAILED", "failed to delete match")
		return
	}
	response.Success(c, "match deleted", gin.H{"id": id})
}

func (ctl *Controller) SubmitResult(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.SubmitResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", "INVALID_BODY")
		return
	}
	if err := sharedvalidator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}
	data, err := ctl.service.SubmitResult(uint(id), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), "SUBMIT_RESULT_FAILED")
		return
	}
	response.Success(c, "result submitted", data)
}

func (ctl *Controller) RollbackResult(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	data, err := ctl.service.RollbackResult(uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), "ROLLBACK_FAILED")
		return
	}
	response.Success(c, "result rolled back", data)
}
