package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mlchain/mlchain-plugin-daemon/internal/service"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/requests"
)

func GetRemoteDebuggingKey(c *gin.Context) {
	BindRequest(
		c, func(request requests.RequestGetRemoteDebuggingKey) {
			c.JSON(200, service.GetRemoteDebuggingKey(request.TenantID))
		},
	)
}
