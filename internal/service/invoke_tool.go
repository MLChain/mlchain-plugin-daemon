package service

import (
	"github.com/gin-gonic/gin"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/plugin_daemon"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/plugin_daemon/access_types"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/session_manager"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/plugin_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/requests"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/tool_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/exception"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/stream"
)

func InvokeTool(
	r *plugin_entities.InvokePluginRequest[requests.RequestInvokeTool],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	// create session
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_TOOL,
		access_types.PLUGIN_ACCESS_ACTION_INVOKE_TOOL,
		ctx.GetString("cluster_id"),
	)
	if err != nil {
		ctx.JSON(500, exception.InternalServerError(err).ToResponse())
		return
	}
	defer session.Close(session_manager.CloseSessionPayload{
		IgnoreCache: false,
	})

	baseSSEService(
		func() (*stream.Stream[tool_entities.ToolResponseChunk], error) {
			return plugin_daemon.InvokeTool(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func ValidateToolCredentials(
	r *plugin_entities.InvokePluginRequest[requests.RequestValidateToolCredentials],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	// create session
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_TOOL,
		access_types.PLUGIN_ACCESS_ACTION_VALIDATE_TOOL_CREDENTIALS,
		ctx.GetString("cluster_id"),
	)
	if err != nil {
		ctx.JSON(500, exception.InternalServerError(err).ToResponse())
		return
	}
	defer session.Close(session_manager.CloseSessionPayload{
		IgnoreCache: false,
	})

	baseSSEService(
		func() (*stream.Stream[tool_entities.ValidateCredentialsResult], error) {
			return plugin_daemon.ValidateToolCredentials(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func GetToolRuntimeParameters(
	r *plugin_entities.InvokePluginRequest[requests.RequestGetToolRuntimeParameters],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	// create session
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_TOOL,
		access_types.PLUGIN_ACCESS_ACTION_GET_TOOL_RUNTIME_PARAMETERS,
		ctx.GetString("cluster_id"),
	)
	if err != nil {
		ctx.JSON(500, exception.InternalServerError(err).ToResponse())
		return
	}
	defer session.Close(session_manager.CloseSessionPayload{
		IgnoreCache: false,
	})

	baseSSEService(
		func() (*stream.Stream[tool_entities.GetToolRuntimeParametersResponse], error) {
			return plugin_daemon.GetToolRuntimeParameters(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}
