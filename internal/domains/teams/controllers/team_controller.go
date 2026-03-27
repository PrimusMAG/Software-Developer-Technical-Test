package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"football-api/internal/domains/teams/dtos"
	"football-api/internal/domains/teams/services"
	"football-api/internal/shared/pagination"
	"football-api/internal/shared/response"
	sharedvalidator "football-api/internal/shared/validator"
)

type Controller struct{ service services.Service }

func New(service services.Service) *Controller { return &Controller{service: service} }

func (ctl *Controller) Create(c *gin.Context) {
	var req dtos.TeamRequest
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
		response.ErrorFrom(c, err, http.StatusBadRequest, "CREATE_TEAM_FAILED", "failed to create team")
		return
	}
	response.Created(c, "team created", data)
}

func (ctl *Controller) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.TeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", "INVALID_BODY")
		return
	}
	if err := sharedvalidator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}
	data, err := ctl.service.Update(uint(id), req)
	if err != nil {
		response.ErrorFrom(c, err, http.StatusBadRequest, "UPDATE_TEAM_FAILED", "failed to update team")
		return
	}
	response.Success(c, "team updated", data)
}

func (ctl *Controller) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := ctl.service.Delete(uint(id)); err != nil {
		response.ErrorFrom(c, err, http.StatusBadRequest, "DELETE_TEAM_FAILED", "failed to delete team")
		return
	}
	response.Success(c, "team deleted", gin.H{"id": id})
}

func (ctl *Controller) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	data, err := ctl.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "team not found", "TEAM_NOT_FOUND")
		return
	}
	response.Success(c, "team fetched", data)
}

func (ctl *Controller) List(c *gin.Context) {
	p := pagination.Parse(c)
	search := c.Query("search")
	city := c.Query("city")
	data, total, err := ctl.service.List(p.Offset, p.Limit, search, city)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed fetch teams", "FETCH_TEAM_FAILED")
		return
	}
	response.SuccessWithMeta(c, "teams fetched", data, response.Meta{
		Page: p.Page, Limit: p.Limit, Total: total, TotalPages: pagination.TotalPages(total, p.Limit),
	})
}
