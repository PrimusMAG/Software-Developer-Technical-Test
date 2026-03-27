package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"football-api/internal/domains/users/dtos"
	"football-api/internal/domains/users/services"
	"football-api/internal/shared/pagination"
	"football-api/internal/shared/response"
	sharedvalidator "football-api/internal/shared/validator"
)

type Controller struct {
	service services.Service
}

func New(service services.Service) *Controller {
	return &Controller{service: service}
}

func (ctl *Controller) Create(c *gin.Context) {
	var req dtos.CreateUserRequest
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
		response.ErrorFrom(c, err, http.StatusBadRequest, "CREATE_USER_FAILED", "failed to create user")
		return
	}
	response.Created(c, "user created", data)
}

func (ctl *Controller) List(c *gin.Context) {
	p := pagination.Parse(c)
	search := c.Query("search")
	data, total, err := ctl.service.List(p.Offset, p.Limit, search)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed fetch users", "FETCH_USER_FAILED")
		return
	}
	response.SuccessWithMeta(c, "users fetched", data, response.Meta{
		Page: p.Page, Limit: p.Limit, Total: total, TotalPages: pagination.TotalPages(total, p.Limit),
	})
}
