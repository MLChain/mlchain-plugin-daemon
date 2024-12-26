package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mlchain/mlchain-plugin-daemon/internal/service"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/app"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/plugin_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/requests"
)

func InvokeAgentStrategy(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestInvokeAgentStrategy]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.InvokeAgentStrategy(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func ListAgentStrategies(c *gin.Context) {
	BindRequest(c, func(request struct {
		TenantID string `uri:"tenant_id" validate:"required"`
		Page     int    `form:"page" validate:"required,min=1"`
		PageSize int    `form:"page_size" validate:"required,min=1,max=256"`
	}) {
		c.JSON(http.StatusOK, service.ListAgentStrategies(request.TenantID, request.Page, request.PageSize))
	})
}

func GetAgentStrategy(c *gin.Context) {
	BindRequest(c, func(request struct {
		TenantID string `uri:"tenant_id" validate:"required"`
		PluginID string `form:"plugin_id" validate:"required"`
		Provider string `form:"provider" validate:"required"`
	}) {
		c.JSON(http.StatusOK, service.GetAgentStrategy(request.TenantID, request.PluginID, request.Provider))
	})
}