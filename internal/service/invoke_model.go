package service

import (
	"github.com/gin-gonic/gin"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/plugin_daemon"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/plugin_daemon/access_types"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/session_manager"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/model_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/plugin_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/requests"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/exception"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/stream"
)

func InvokeLLM(
	r *plugin_entities.InvokePluginRequest[requests.RequestInvokeLLM],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	// create session
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_MODEL,
		access_types.PLUGIN_ACCESS_ACTION_INVOKE_LLM,
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
		func() (*stream.Stream[model_entities.LLMResultChunk], error) {
			return plugin_daemon.InvokeLLM(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func InvokeTextEmbedding(
	r *plugin_entities.InvokePluginRequest[requests.RequestInvokeTextEmbedding],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	// create session
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_MODEL,
		access_types.PLUGIN_ACCESS_ACTION_INVOKE_TEXT_EMBEDDING,
		ctx.GetString("cluster_id"))
	if err != nil {
		ctx.JSON(500, exception.InternalServerError(err).ToResponse())
		return
	}
	defer session.Close(session_manager.CloseSessionPayload{
		IgnoreCache: false,
	})

	baseSSEService(
		func() (*stream.Stream[model_entities.TextEmbeddingResult], error) {
			return plugin_daemon.InvokeTextEmbedding(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func InvokeRerank(
	r *plugin_entities.InvokePluginRequest[requests.RequestInvokeRerank],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	// create session
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_MODEL,
		access_types.PLUGIN_ACCESS_ACTION_INVOKE_RERANK,
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
		func() (*stream.Stream[model_entities.RerankResult], error) {
			return plugin_daemon.InvokeRerank(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func InvokeTTS(
	r *plugin_entities.InvokePluginRequest[requests.RequestInvokeTTS],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	// create session
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_MODEL,
		access_types.PLUGIN_ACCESS_ACTION_INVOKE_TTS,
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
		func() (*stream.Stream[model_entities.TTSResult], error) {
			return plugin_daemon.InvokeTTS(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func InvokeSpeech2Text(
	r *plugin_entities.InvokePluginRequest[requests.RequestInvokeSpeech2Text],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	// create session
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_MODEL,
		access_types.PLUGIN_ACCESS_ACTION_INVOKE_SPEECH2TEXT,
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
		func() (*stream.Stream[model_entities.Speech2TextResult], error) {
			return plugin_daemon.InvokeSpeech2Text(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func InvokeModeration(
	r *plugin_entities.InvokePluginRequest[requests.RequestInvokeModeration],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	// create session
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_MODEL,
		access_types.PLUGIN_ACCESS_ACTION_INVOKE_MODERATION,
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
		func() (*stream.Stream[model_entities.ModerationResult], error) {
			return plugin_daemon.InvokeModeration(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func ValidateProviderCredentials(
	r *plugin_entities.InvokePluginRequest[requests.RequestValidateProviderCredentials],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	// create session
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_MODEL,
		access_types.PLUGIN_ACCESS_ACTION_VALIDATE_PROVIDER_CREDENTIALS,
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
		func() (*stream.Stream[model_entities.ValidateCredentialsResult], error) {
			return plugin_daemon.ValidateProviderCredentials(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func ValidateModelCredentials(
	r *plugin_entities.InvokePluginRequest[requests.RequestValidateModelCredentials],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	// create session
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_MODEL,
		access_types.PLUGIN_ACCESS_ACTION_VALIDATE_MODEL_CREDENTIALS,
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
		func() (*stream.Stream[model_entities.ValidateCredentialsResult], error) {
			return plugin_daemon.ValidateModelCredentials(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func GetTTSModelVoices(
	r *plugin_entities.InvokePluginRequest[requests.RequestGetTTSModelVoices],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_MODEL,
		access_types.PLUGIN_ACCESS_ACTION_GET_TTS_MODEL_VOICES,
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
		func() (*stream.Stream[model_entities.GetTTSVoicesResponse], error) {
			return plugin_daemon.GetTTSModelVoices(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func GetTextEmbeddingNumTokens(
	r *plugin_entities.InvokePluginRequest[requests.RequestGetTextEmbeddingNumTokens],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_MODEL,
		access_types.PLUGIN_ACCESS_ACTION_GET_TEXT_EMBEDDING_NUM_TOKENS,
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
		func() (*stream.Stream[model_entities.GetTextEmbeddingNumTokensResponse], error) {
			return plugin_daemon.GetTextEmbeddingNumTokens(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func GetAIModelSchema(
	r *plugin_entities.InvokePluginRequest[requests.RequestGetAIModelSchema],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_MODEL,
		access_types.PLUGIN_ACCESS_ACTION_GET_AI_MODEL_SCHEMAS,
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
		func() (*stream.Stream[model_entities.GetModelSchemasResponse], error) {
			return plugin_daemon.GetAIModelSchema(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}

func GetLLMNumTokens(
	r *plugin_entities.InvokePluginRequest[requests.RequestGetLLMNumTokens],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	session, err := createSession(
		r,
		access_types.PLUGIN_ACCESS_TYPE_MODEL,
		access_types.PLUGIN_ACCESS_ACTION_GET_LLM_NUM_TOKENS,
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
		func() (*stream.Stream[model_entities.LLMGetNumTokensResponse], error) {
			return plugin_daemon.GetLLMNumTokens(session, &r.Data)
		},
		ctx,
		max_timeout_seconds,
	)
}
