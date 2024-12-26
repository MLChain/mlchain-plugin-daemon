package mlchain_invocation

import (
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/model_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/tool_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/stream"
)

type BackwardsInvocation interface {
	// InvokeLLM
	InvokeLLM(payload *InvokeLLMRequest) (*stream.Stream[model_entities.LLMResultChunk], error)
	// InvokeTextEmbedding
	InvokeTextEmbedding(payload *InvokeTextEmbeddingRequest) (*model_entities.TextEmbeddingResult, error)
	// InvokeRerank
	InvokeRerank(payload *InvokeRerankRequest) (*model_entities.RerankResult, error)
	// InvokeTTS
	InvokeTTS(payload *InvokeTTSRequest) (*stream.Stream[model_entities.TTSResult], error)
	// InvokeSpeech2Text
	InvokeSpeech2Text(payload *InvokeSpeech2TextRequest) (*model_entities.Speech2TextResult, error)
	// InvokeModeration
	InvokeModeration(payload *InvokeModerationRequest) (*model_entities.ModerationResult, error)
	// InvokeTool
	InvokeTool(payload *InvokeToolRequest) (*stream.Stream[tool_entities.ToolResponseChunk], error)
	// InvokeApp
	InvokeApp(payload *InvokeAppRequest) (*stream.Stream[map[string]any], error)
	// InvokeParameterExtractor
	InvokeParameterExtractor(payload *InvokeParameterExtractorRequest) (*InvokeNodeResponse, error)
	// InvokeQuestionClassifier
	InvokeQuestionClassifier(payload *InvokeQuestionClassifierRequest) (*InvokeNodeResponse, error)
	// InvokeEncrypt
	InvokeEncrypt(payload *InvokeEncryptRequest) (map[string]any, error)
	// InvokeSummary
	InvokeSummary(payload *InvokeSummaryRequest) (*InvokeSummaryResponse, error)
	// UploadFile
	UploadFile(payload *UploadFileRequest) (*UploadFileResponse, error)
}
