package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"football-api/internal/domains/players/dtos"
	"football-api/internal/domains/players/services"
	"football-api/internal/shared/pagination"
	"football-api/internal/shared/response"
	sharedvalidator "football-api/internal/shared/validator"
)

type Controller struct{ service services.Service }

func New(service services.Service) *Controller { return &Controller{service: service} }

func (ctl *Controller) Create(c *gin.Context) {
	var req dtos.PlayerRequest
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
		response.ErrorFrom(c, err, http.StatusBadRequest, "CREATE_PLAYER_FAILED", "failed to create player")
		return
	}
	response.Created(c, "player created", data)
}

func (ctl *Controller) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dtos.PlayerRequest
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
		response.ErrorFrom(c, err, http.StatusBadRequest, "UPDATE_PLAYER_FAILED", "failed to update player")
		return
	}
	response.Success(c, "player updated", data)
}

func (ctl *Controller) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := ctl.service.Delete(uint(id)); err != nil {
		response.ErrorFrom(c, err, http.StatusBadRequest, "DELETE_PLAYER_FAILED", "failed to delete player")
		return
	}
	response.Success(c, "player deleted", gin.H{"id": id})
}

func (ctl *Controller) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	data, err := ctl.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "player not found", "PLAYER_NOT_FOUND")
		return
	}
	response.Success(c, "player fetched", data)
}

func (ctl *Controller) List(c *gin.Context) {
	p := pagination.Parse(c)
	search := c.Query("search")
	position := c.Query("position")
	teamIDRaw := c.Query("team_id")
	var teamID uint64
	if teamIDRaw != "" {
		teamID, _ = strconv.ParseUint(teamIDRaw, 10, 64)
	}
	data, total, err := ctl.service.List(p.Offset, p.Limit, search, position, uint(teamID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed fetch players", "FETCH_PLAYER_FAILED")
		return
	}
	response.SuccessWithMeta(c, "players fetched", data, response.Meta{
		Page: p.Page, Limit: p.Limit, Total: total, TotalPages: pagination.TotalPages(total, p.Limit),
	})
}
