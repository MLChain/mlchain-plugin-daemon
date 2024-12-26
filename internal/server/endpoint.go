package server

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mlchain/mlchain-plugin-daemon/internal/db"
	"github.com/mlchain/mlchain-plugin-daemon/internal/service"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/plugin_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/exception"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/models"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/log"
)

// DifyPlugin supports register and use endpoint to improve the plugin's functionality
// you can use it to do some magics, looking forward to your imagination, Ciallo～(∠·ω< )⌒
// - Yeuoly

// EndpointHandler is a function type that can be used to handle endpoint requests
type EndpointHandler func(ctx *gin.Context, hookId string, path string)

func (app *App) Endpoint() func(c *gin.Context) {
	return func(c *gin.Context) {
		hookId := c.Param("hook_id")
		path := c.Param("path")

		if app.endpointHandler != nil {
			app.endpointHandler(c, hookId, path)
		} else {
			app.EndpointHandler(c, hookId, path)
		}
	}
}

func (app *App) EndpointHandler(ctx *gin.Context, hookId string, path string) {
	endpoint, err := db.GetOne[models.Endpoint](
		db.Equal("hook_id", hookId),
	)
	if err == db.ErrDatabaseNotFound {
		ctx.JSON(404, exception.BadRequestError(errors.New("endpoint not found")).ToResponse())
		return
	}

	if err != nil {
		log.Error("get endpoint error %v", err)
		ctx.JSON(500, exception.InternalServerError(errors.New("internal server error")).ToResponse())
		return
	}

	// get plugin installation
	pluginInstallation, err := db.GetOne[models.PluginInstallation](
		db.Equal("plugin_id", endpoint.PluginID),
		db.Equal("tenant_id", endpoint.TenantID),
	)
	if err != nil {
		ctx.JSON(404, exception.BadRequestError(errors.New("plugin installation not found")).ToResponse())
		return
	}

	pluginUniqueIdentifier, err := plugin_entities.NewPluginUniqueIdentifier(
		pluginInstallation.PluginUniqueIdentifier,
	)
	if err != nil {
		ctx.JSON(400, exception.PluginUniqueIdentifierError(
			errors.New("invalid plugin unique identifier"),
		).ToResponse())
		return
	}

	// check if plugin exists in current node
	if ok, originalError := app.cluster.IsPluginOnCurrentNode(pluginUniqueIdentifier); !ok {
		app.redirectPluginInvokeByPluginIdentifier(ctx, pluginUniqueIdentifier, originalError)
	} else {
		service.Endpoint(ctx, &endpoint, &pluginInstallation, path)
	}
}
