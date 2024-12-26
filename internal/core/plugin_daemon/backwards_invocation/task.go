package backwards_invocation

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/mlchain/mlchain-plugin-daemon/internal/core/mlchain_invocation"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/persistence"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/plugin_daemon/access_types"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/session_manager"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/plugin_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/log"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/parser"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/routine"
)

// returns error only if payload is not correct
func InvokeDify(
	declaration *plugin_entities.PluginDeclaration,
	invoke_from access_types.PluginAccessType,
	session *session_manager.Session,
	writer BackwardsInvocationWriter,
	data []byte,
) error {
	// unmarshal invoke data
	request, err := parser.UnmarshalJsonBytes2Map(data)
	if err != nil {
		return fmt.Errorf("unmarshal invoke request failed: %s", err.Error())
	}

	if request == nil {
		return fmt.Errorf("invoke request is empty")
	}

	// prepare invocation arguments
	requestHandle, err := prepareDifyInvocationArguments(
		session,
		writer,
		request,
	)
	if err != nil {
		return err
	}

	if invoke_from == access_types.PLUGIN_ACCESS_TYPE_MODEL {
		requestHandle.WriteError(fmt.Errorf("you can not invoke mlchain from %s", invoke_from))
		requestHandle.EndResponse()
		return nil
	}

	// check permission
	if err := checkPermission(declaration, requestHandle); err != nil {
		requestHandle.WriteError(err)
		requestHandle.EndResponse()
		return nil
	}

	// dispatch invocation task
	routine.Submit(map[string]string{
		"module":   "plugin_daemon",
		"function": "InvokeDify",
	}, func() {
		dispatchDifyInvocationTask(requestHandle)
		defer requestHandle.EndResponse()
	})

	return nil
}

var (
	permissionMapping = map[mlchain_invocation.InvokeType]map[string]any{
		mlchain_invocation.INVOKE_TYPE_TOOL: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return declaration.Resource.Permission.AllowInvokeTool()
			},
			"error": "permission denied, you need to enable tool access in plugin manifest",
		},
		mlchain_invocation.INVOKE_TYPE_LLM: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return declaration.Resource.Permission.AllowInvokeLLM()
			},
			"error": "permission denied, you need to enable llm access in plugin manifest",
		},
		mlchain_invocation.INVOKE_TYPE_TEXT_EMBEDDING: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return declaration.Resource.Permission.AllowInvokeTextEmbedding()
			},
			"error": "permission denied, you need to enable text-embedding access in plugin manifest",
		},
		mlchain_invocation.INVOKE_TYPE_RERANK: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return declaration.Resource.Permission.AllowInvokeRerank()
			},
			"error": "permission denied, you need to enable rerank access in plugin manifest",
		},
		mlchain_invocation.INVOKE_TYPE_TTS: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return declaration.Resource.Permission.AllowInvokeTTS()
			},
			"error": "permission denied, you need to enable tts access in plugin manifest",
		},
		mlchain_invocation.INVOKE_TYPE_SPEECH2TEXT: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return declaration.Resource.Permission.AllowInvokeSpeech2Text()
			},
			"error": "permission denied, you need to enable speech2text access in plugin manifest",
		},
		mlchain_invocation.INVOKE_TYPE_MODERATION: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return declaration.Resource.Permission.AllowInvokeModeration()
			},
			"error": "permission denied, you need to enable moderation access in plugin manifest",
		},
		mlchain_invocation.INVOKE_TYPE_NODE_PARAMETER_EXTRACTOR: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return declaration.Resource.Permission.AllowInvokeNode()
			},
			"error": "permission denied, you need to enable node access in plugin manifest",
		},
		mlchain_invocation.INVOKE_TYPE_NODE_QUESTION_CLASSIFIER: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return declaration.Resource.Permission.AllowInvokeNode()
			},
			"error": "permission denied, you need to enable node access in plugin manifest",
		},
		mlchain_invocation.INVOKE_TYPE_APP: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return declaration.Resource.Permission.AllowInvokeApp()
			},
			"error": "permission denied, you need to enable app access in plugin manifest",
		},
		mlchain_invocation.INVOKE_TYPE_STORAGE: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return declaration.Resource.Permission.AllowInvokeStorage()
			},
			"error": "permission denied, you need to enable storage access in plugin manifest",
		},
		mlchain_invocation.INVOKE_TYPE_SYSTEM_SUMMARY: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return declaration.Resource.Permission.AllowInvokeLLM()
			},
			"error": "permission denied, you need to enable llm access in plugin manifest",
		},
		mlchain_invocation.INVOKE_TYPE_UPLOAD_FILE: {
			"func": func(declaration *plugin_entities.PluginDeclaration) bool {
				return true
			},
			"error": "permission denied, you need to enable storage access in plugin manifest",
		},
	}
)

func checkPermission(runtime *plugin_entities.PluginDeclaration, requestHandle *BackwardsInvocation) error {
	permission, ok := permissionMapping[requestHandle.Type()]
	if !ok {
		return fmt.Errorf("unsupported invoke type: %s", requestHandle.Type())
	}

	permissionFunc, ok := permission["func"].(func(runtime *plugin_entities.PluginDeclaration) bool)
	if !ok {
		return fmt.Errorf("permission function not found: %s", requestHandle.Type())
	}

	if !permissionFunc(runtime) {
		return fmt.Errorf(permission["error"].(string))
	}

	return nil
}

func prepareDifyInvocationArguments(
	session *session_manager.Session,
	writer BackwardsInvocationWriter,
	request map[string]any,
) (*BackwardsInvocation, error) {
	typ, ok := request["type"].(string)
	if !ok {
		return nil, fmt.Errorf("invoke request missing type: %s", request)
	}

	// get request id
	backwardsRequestId, ok := request["backwards_request_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invoke request missing request_id: %s", request)
	}

	// get request
	detailedRequest, ok := request["request"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invoke request missing request: %s", request)
	}

	return NewBackwardsInvocation(
		BackwardsInvocationType(typ),
		backwardsRequestId,
		session,
		writer,
		detailedRequest,
	), nil
}

var (
	dispatchMapping = map[mlchain_invocation.InvokeType]func(handle *BackwardsInvocation){
		mlchain_invocation.INVOKE_TYPE_TOOL: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationToolTask)
		},
		mlchain_invocation.INVOKE_TYPE_LLM: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationLLMTask)
		},
		mlchain_invocation.INVOKE_TYPE_TEXT_EMBEDDING: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationTextEmbeddingTask)
		},
		mlchain_invocation.INVOKE_TYPE_RERANK: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationRerankTask)
		},
		mlchain_invocation.INVOKE_TYPE_TTS: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationTTSTask)
		},
		mlchain_invocation.INVOKE_TYPE_SPEECH2TEXT: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationSpeech2TextTask)
		},
		mlchain_invocation.INVOKE_TYPE_MODERATION: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationModerationTask)
		},
		mlchain_invocation.INVOKE_TYPE_APP: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationAppTask)
		},
		mlchain_invocation.INVOKE_TYPE_NODE_PARAMETER_EXTRACTOR: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationParameterExtractor)
		},
		mlchain_invocation.INVOKE_TYPE_NODE_QUESTION_CLASSIFIER: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationQuestionClassifier)
		},
		mlchain_invocation.INVOKE_TYPE_STORAGE: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationStorageTask)
		},
		mlchain_invocation.INVOKE_TYPE_SYSTEM_SUMMARY: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationSystemSummaryTask)
		},
		mlchain_invocation.INVOKE_TYPE_UPLOAD_FILE: func(handle *BackwardsInvocation) {
			genericDispatchTask(handle, executeDifyInvocationUploadFileTask)
		},
	}
)

func genericDispatchTask[T any](
	handle *BackwardsInvocation,
	dispatch func(
		handle *BackwardsInvocation,
		request *T,
	),
) {
	r, err := parser.MapToStruct[T](handle.RequestData())
	if err != nil {
		handle.WriteError(fmt.Errorf("unmarshal backwards invoke request failed: %s", err.Error()))
		return
	}
	dispatch(handle, r)
}

func dispatchDifyInvocationTask(handle *BackwardsInvocation) {
	requestData := handle.RequestData()
	tenantId, err := handle.TenantID()
	if err != nil {
		handle.WriteError(fmt.Errorf("get tenant id failed: %s", err.Error()))
		return
	}
	requestData["tenant_id"] = tenantId
	userId, err := handle.UserID()
	if err != nil {
		handle.WriteError(fmt.Errorf("get user id failed: %s", err.Error()))
		return
	}
	requestData["user_id"] = userId
	typ := handle.Type()
	requestData["type"] = typ

	for t, v := range dispatchMapping {
		if t == handle.Type() {
			v(handle)
			return
		}
	}

	handle.WriteError(fmt.Errorf("unsupported invoke type: %s", handle.Type()))
}

func executeDifyInvocationToolTask(
	handle *BackwardsInvocation,
	request *mlchain_invocation.InvokeToolRequest,
) {
	response, err := handle.backwardsInvocation.InvokeTool(request)
	if err != nil {
		handle.WriteError(fmt.Errorf("invoke tool failed: %s", err.Error()))
		return
	}

	for response.Next() {
		value, err := response.Read()
		if err != nil {
			handle.WriteError(fmt.Errorf("read tool response failed: %s", err.Error()))
			return
		}

		handle.WriteResponse("stream", value)
	}
}

func executeDifyInvocationLLMTask(
	handle *BackwardsInvocation,
	request *mlchain_invocation.InvokeLLMRequest,
) {
	response, err := handle.backwardsInvocation.InvokeLLM(request)
	if err != nil {
		handle.WriteError(fmt.Errorf("invoke llm model failed: %s", err.Error()))
		return
	}

	for response.Next() {
		value, err := response.Read()
		if err != nil {
			handle.WriteError(fmt.Errorf("read llm model failed: %s", err.Error()))
			return
		}

		handle.WriteResponse("stream", value)
	}
}

func executeDifyInvocationTextEmbeddingTask(
	handle *BackwardsInvocation,
	request *mlchain_invocation.InvokeTextEmbeddingRequest,
) {
	response, err := handle.backwardsInvocation.InvokeTextEmbedding(request)
	if err != nil {
		handle.WriteError(fmt.Errorf("invoke text-embedding model failed: %s", err.Error()))
		return
	}

	handle.WriteResponse("struct", response)
}

func executeDifyInvocationRerankTask(
	handle *BackwardsInvocation,
	request *mlchain_invocation.InvokeRerankRequest,
) {
	response, err := handle.backwardsInvocation.InvokeRerank(request)
	if err != nil {
		handle.WriteError(fmt.Errorf("invoke rerank model failed: %s", err.Error()))
		return
	}

	handle.WriteResponse("struct", response)
}

func executeDifyInvocationTTSTask(
	handle *BackwardsInvocation,
	request *mlchain_invocation.InvokeTTSRequest,
) {
	response, err := handle.backwardsInvocation.InvokeTTS(request)
	if err != nil {
		handle.WriteError(fmt.Errorf("invoke tts model failed: %s", err.Error()))
		return
	}

	for response.Next() {
		value, err := response.Read()
		if err != nil {
			handle.WriteError(fmt.Errorf("read tts model failed: %s", err.Error()))
			return
		}

		handle.WriteResponse("stream", value)
	}
}

func executeDifyInvocationSpeech2TextTask(
	handle *BackwardsInvocation,
	request *mlchain_invocation.InvokeSpeech2TextRequest,
) {
	response, err := handle.backwardsInvocation.InvokeSpeech2Text(request)
	if err != nil {
		handle.WriteError(fmt.Errorf("invoke speech2text model failed: %s", err.Error()))
		return
	}

	handle.WriteResponse("struct", response)
}

func executeDifyInvocationModerationTask(
	handle *BackwardsInvocation,
	request *mlchain_invocation.InvokeModerationRequest,
) {
	response, err := handle.backwardsInvocation.InvokeModeration(request)
	if err != nil {
		handle.WriteError(fmt.Errorf("invoke moderation model failed: %s", err.Error()))
		return
	}

	handle.WriteResponse("struct", response)
}

func executeDifyInvocationAppTask(
	handle *BackwardsInvocation,
	request *mlchain_invocation.InvokeAppRequest,
) {
	response, err := handle.backwardsInvocation.InvokeApp(request)
	if err != nil {
		handle.WriteError(fmt.Errorf("invoke app failed: %s", err.Error()))
		return
	}

	userId, err := handle.UserID()
	if err != nil {
		handle.WriteError(fmt.Errorf("get user id failed: %s", err.Error()))
		return
	}

	request.User = userId

	for response.Next() {
		value, err := response.Read()
		if err != nil {
			handle.WriteError(fmt.Errorf("read app failed: %s", err.Error()))
			return
		}

		handle.WriteResponse("stream", value)
	}
}

func executeDifyInvocationParameterExtractor(
	handle *BackwardsInvocation,
	request *mlchain_invocation.InvokeParameterExtractorRequest,
) {
	response, err := handle.backwardsInvocation.InvokeParameterExtractor(request)
	if err != nil {
		handle.WriteError(fmt.Errorf("invoke parameter extractor failed: %s", err.Error()))
		return
	}

	handle.WriteResponse("struct", response)
}

func executeDifyInvocationQuestionClassifier(
	handle *BackwardsInvocation,
	request *mlchain_invocation.InvokeQuestionClassifierRequest,
) {
	response, err := handle.backwardsInvocation.InvokeQuestionClassifier(request)
	if err != nil {
		handle.WriteError(fmt.Errorf("invoke question classifier failed: %s", err.Error()))
		return
	}

	handle.WriteResponse("struct", response)
}

func executeDifyInvocationStorageTask(
	handle *BackwardsInvocation,
	request *mlchain_invocation.InvokeStorageRequest,
) {
	if handle.session == nil {
		handle.WriteError(fmt.Errorf("session not found"))
		return
	}

	persistence := persistence.GetPersistence()
	if persistence == nil {
		handle.WriteError(fmt.Errorf("persistence not found"))
		return
	}

	tenantId, err := handle.TenantID()
	if err != nil {
		handle.WriteError(fmt.Errorf("get tenant id failed: %s", err.Error()))
		return
	}

	pluginId := handle.session.PluginUniqueIdentifier

	if request.Opt == mlchain_invocation.STORAGE_OPT_GET {
		data, err := persistence.Load(tenantId, pluginId.PluginID(), request.Key)
		if err != nil {
			log.Error("load data failed: %s", err.Error())
			handle.WriteError(errors.New("load data failed, please check if the key is correct or you have not set it"))
			return
		}

		handle.WriteResponse("struct", map[string]any{
			"data": hex.EncodeToString(data),
		})
	} else if request.Opt == mlchain_invocation.STORAGE_OPT_SET {
		data, err := hex.DecodeString(request.Value)
		if err != nil {
			handle.WriteError(fmt.Errorf("decode data failed: %s", err.Error()))
			return
		}

		session := handle.session
		if session == nil {
			handle.WriteError(fmt.Errorf("session not found"))
			return
		}

		declaration := session.Declaration
		if declaration == nil {
			handle.WriteError(fmt.Errorf("declaration not found"))
			return
		}

		resource := declaration.Resource.Permission
		if resource == nil {
			handle.WriteError(fmt.Errorf("resource not found"))
			return
		}

		maxStorageSize := int64(-1)

		storage := resource.Storage
		if storage != nil {
			maxStorageSize = int64(storage.Size)
		}

		if err := persistence.Save(tenantId, pluginId.PluginID(), maxStorageSize, request.Key, data); err != nil {
			handle.WriteError(fmt.Errorf("save data failed: %s", err.Error()))
			return
		}

		handle.WriteResponse("struct", map[string]any{
			"data": "ok",
		})
	} else if request.Opt == mlchain_invocation.STORAGE_OPT_DEL {
		if err := persistence.Delete(tenantId, pluginId.PluginID(), request.Key); err != nil {
			handle.WriteError(fmt.Errorf("delete data failed: %s", err.Error()))
			return
		}

		handle.WriteResponse("struct", map[string]any{
			"data": "ok",
		})
	}
}

func executeDifyInvocationSystemSummaryTask(
	handle *BackwardsInvocation,
	request *mlchain_invocation.InvokeSummaryRequest,
) {
	response, err := handle.backwardsInvocation.InvokeSummary(request)
	if err != nil {
		handle.WriteError(fmt.Errorf("invoke summary failed: %s", err.Error()))
		return
	}

	handle.WriteResponse("struct", response)
}

func executeDifyInvocationUploadFileTask(
	handle *BackwardsInvocation,
	request *mlchain_invocation.UploadFileRequest,
) {
	response, err := handle.backwardsInvocation.UploadFile(request)
	if err != nil {
		handle.WriteError(fmt.Errorf("upload file failed: %s", err.Error()))
		return
	}

	handle.WriteResponse("struct", response)
}
