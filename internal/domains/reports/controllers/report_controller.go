package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"football-api/internal/domains/reports/services"
	"football-api/internal/shared/response"
)

type Controller struct{ service services.Service }

func New(service services.Service) *Controller { return &Controller{service: service} }

func (ctl *Controller) List(c *gin.Context) {
	data, err := ctl.service.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed generate report", "REPORT_FAILED")
		return
	}
	c.Header("Cache-Control", "public, max-age=30, stale-while-revalidate=60")
	response.Success(c, "report fetched", data)
}

func (ctl *Controller) Revalidate(c *gin.Context) {
	ctl.service.Revalidate()
	response.Success(c, "report cache revalidated", gin.H{"ok": true})
}
