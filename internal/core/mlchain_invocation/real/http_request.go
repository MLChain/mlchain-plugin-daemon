package real

import (
	"fmt"
	"reflect"

	"github.com/mlchain/mlchain-plugin-daemon/internal/core/mlchain_invocation"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/model_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/tool_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/validators"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/http_requests"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/routine"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/stream"
)

// Send a request to mlchain inner api and validate the response
func Request[T any](i *RealBackwardsInvocation, method string, path string, options ...http_requests.HttpOptions) (*T, error) {
	options = append(options,
		http_requests.HttpHeader(map[string]string{
			"X-Inner-Api-Key": i.mlchainInnerApiKey,
		}),
		http_requests.HttpWriteTimeout(5000),
		http_requests.HttpReadTimeout(240000),
	)

	req, err := http_requests.RequestAndParse[BaseBackwardsInvocationResponse[T]](i.client, i.mlchainPath(path), method, options...)
	if err != nil {
		return nil, err
	}

	if req.Error != "" {
		return nil, fmt.Errorf("request failed: %s", req.Error)
	}

	if req.Data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	// check if req.Data is a map[string]any
	if reflect.TypeOf(*req.Data).Kind() == reflect.Map {
		return req.Data, nil
	}

	if err := validators.GlobalEntitiesValidator.Struct(req.Data); err != nil {
		return nil, fmt.Errorf("validate request failed: %s", err.Error())
	}

	return req.Data, nil
}

func StreamResponse[T any](i *RealBackwardsInvocation, method string, path string, options ...http_requests.HttpOptions) (
	*stream.Stream[T], error,
) {
	options = append(
		options, http_requests.HttpHeader(map[string]string{
			"X-Inner-Api-Key": i.mlchainInnerApiKey,
		}),
		http_requests.HttpWriteTimeout(5000),
		http_requests.HttpReadTimeout(240000),
	)

	response, err := http_requests.RequestAndParseStream[BaseBackwardsInvocationResponse[T]](i.client, i.mlchainPath(path), method, options...)
	if err != nil {
		return nil, err
	}

	newResponse := stream.NewStream[T](1024)
	newResponse.OnClose(func() {
		response.Close()
	})
	routine.Submit(map[string]string{
		"module":   "mlchain_invocation",
		"function": "StreamResponse",
	}, func() {
		defer newResponse.Close()
		for response.Next() {
			t, err := response.Read()
			if err != nil {
				newResponse.WriteError(err)
				break
			}

			if t.Error != "" {
				newResponse.WriteError(fmt.Errorf("request failed: %s", t.Error))
				break
			}

			if t.Data == nil {
				newResponse.WriteError(fmt.Errorf("data is nil"))
				break
			}

			// check if t.Data is a map[string]any, skip validation if it is
			if reflect.TypeOf(*t.Data).Kind() != reflect.Map {
				if err := validators.GlobalEntitiesValidator.Struct(t.Data); err != nil {
					newResponse.WriteError(fmt.Errorf("validate request failed: %s", err.Error()))
					break
				}
			}

			newResponse.Write(*t.Data)
		}
	})

	return newResponse, nil
}

func (i *RealBackwardsInvocation) InvokeLLM(payload *mlchain_invocation.InvokeLLMRequest) (*stream.Stream[model_entities.LLMResultChunk], error) {
	return StreamResponse[model_entities.LLMResultChunk](i, "POST", "invoke/llm", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeTextEmbedding(payload *mlchain_invocation.InvokeTextEmbeddingRequest) (*model_entities.TextEmbeddingResult, error) {
	return Request[model_entities.TextEmbeddingResult](i, "POST", "invoke/text-embedding", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeRerank(payload *mlchain_invocation.InvokeRerankRequest) (*model_entities.RerankResult, error) {
	return Request[model_entities.RerankResult](i, "POST", "invoke/rerank", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeTTS(payload *mlchain_invocation.InvokeTTSRequest) (*stream.Stream[model_entities.TTSResult], error) {
	return StreamResponse[model_entities.TTSResult](i, "POST", "invoke/tts", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeSpeech2Text(payload *mlchain_invocation.InvokeSpeech2TextRequest) (*model_entities.Speech2TextResult, error) {
	return Request[model_entities.Speech2TextResult](i, "POST", "invoke/speech2text", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeModeration(payload *mlchain_invocation.InvokeModerationRequest) (*model_entities.ModerationResult, error) {
	return Request[model_entities.ModerationResult](i, "POST", "invoke/moderation", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeTool(payload *mlchain_invocation.InvokeToolRequest) (*stream.Stream[tool_entities.ToolResponseChunk], error) {
	return StreamResponse[tool_entities.ToolResponseChunk](i, "POST", "invoke/tool", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeApp(payload *mlchain_invocation.InvokeAppRequest) (*stream.Stream[map[string]any], error) {
	return StreamResponse[map[string]any](i, "POST", "invoke/app", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeParameterExtractor(payload *mlchain_invocation.InvokeParameterExtractorRequest) (*mlchain_invocation.InvokeNodeResponse, error) {
	return Request[mlchain_invocation.InvokeNodeResponse](i, "POST", "invoke/parameter-extractor", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeQuestionClassifier(payload *mlchain_invocation.InvokeQuestionClassifierRequest) (*mlchain_invocation.InvokeNodeResponse, error) {
	return Request[mlchain_invocation.InvokeNodeResponse](i, "POST", "invoke/question-classifier", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeEncrypt(payload *mlchain_invocation.InvokeEncryptRequest) (map[string]any, error) {
	if !payload.EncryptRequired(payload.Data) {
		return payload.Data, nil
	}

	type resp struct {
		Data map[string]any `json:"data,omitempty"`
	}

	data, err := Request[resp](i, "POST", "invoke/encrypt", http_requests.HttpPayloadJson(payload))
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}

func (i *RealBackwardsInvocation) InvokeSummary(payload *mlchain_invocation.InvokeSummaryRequest) (*mlchain_invocation.InvokeSummaryResponse, error) {
	return Request[mlchain_invocation.InvokeSummaryResponse](i, "POST", "invoke/summary", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) UploadFile(payload *mlchain_invocation.UploadFileRequest) (*mlchain_invocation.UploadFileResponse, error) {
	return Request[mlchain_invocation.UploadFileResponse](i, "POST", "upload/file/request", http_requests.HttpPayloadJson(payload))
}
