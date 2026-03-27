package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"football-api/internal/domains/auth/dtos"
	"football-api/internal/domains/auth/services"
	"football-api/internal/shared/response"
	sharedvalidator "football-api/internal/shared/validator"
)

type Controller struct {
	service services.Service
}

func New(service services.Service) *Controller {
	return &Controller{service: service}
}

func (ctl *Controller) Login(c *gin.Context) {
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", "INVALID_BODY")
		return
	}
	if err := sharedvalidator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	data, err := ctl.service.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), "INVALID_CREDENTIAL")
		return
	}
	response.Success(c, "login success", data)
}

func (ctl *Controller) Refresh(c *gin.Context) {
	var req dtos.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", "INVALID_BODY")
		return
	}
	if err := sharedvalidator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}
	data, err := ctl.service.Refresh(req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), "INVALID_REFRESH_TOKEN")
		return
	}
	response.Success(c, "refresh success", data)
}
