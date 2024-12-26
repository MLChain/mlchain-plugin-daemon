package tester

import (
	"time"

	"github.com/mlchain/mlchain-plugin-daemon/internal/core/mlchain_invocation"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/model_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/tool_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/routine"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/stream"
)

type MockedDifyInvocation struct{}

func NewMockedDifyInvocation() mlchain_invocation.BackwardsInvocation {
	return &MockedDifyInvocation{}
}

func (m *MockedDifyInvocation) InvokeLLM(payload *mlchain_invocation.InvokeLLMRequest) (*stream.Stream[model_entities.LLMResultChunk], error) {
	stream := stream.NewStream[model_entities.LLMResultChunk](5)
	routine.Submit(nil, func() {
		stream.Write(model_entities.LLMResultChunk{
			Model:             model_entities.LLMModel(payload.Model),
			PromptMessages:    payload.PromptMessages,
			SystemFingerprint: "test",
			Delta: model_entities.LLMResultChunkDelta{
				Index: &[]int{1}[0],
				Message: model_entities.PromptMessage{
					Role:    model_entities.PROMPT_MESSAGE_ROLE_ASSISTANT,
					Content: "hello",
					Name:    "test",
				},
			},
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(model_entities.LLMResultChunk{
			Model:             model_entities.LLMModel(payload.Model),
			PromptMessages:    payload.PromptMessages,
			SystemFingerprint: "test",
			Delta: model_entities.LLMResultChunkDelta{
				Index: &[]int{1}[0],
				Message: model_entities.PromptMessage{
					Role:    model_entities.PROMPT_MESSAGE_ROLE_ASSISTANT,
					Content: " world",
					Name:    "test",
				},
			},
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(model_entities.LLMResultChunk{
			Model:             model_entities.LLMModel(payload.Model),
			PromptMessages:    payload.PromptMessages,
			SystemFingerprint: "test",
			Delta: model_entities.LLMResultChunkDelta{
				Index: &[]int{2}[0],
				Message: model_entities.PromptMessage{
					Role:    model_entities.PROMPT_MESSAGE_ROLE_ASSISTANT,
					Content: " world",
					Name:    "test",
				},
			},
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(model_entities.LLMResultChunk{
			Model:             model_entities.LLMModel(payload.Model),
			PromptMessages:    payload.PromptMessages,
			SystemFingerprint: "test",
			Delta: model_entities.LLMResultChunkDelta{
				Index: &[]int{3}[0],
				Message: model_entities.PromptMessage{
					Role:    model_entities.PROMPT_MESSAGE_ROLE_ASSISTANT,
					Content: " !",
					Name:    "test",
				},
			},
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(model_entities.LLMResultChunk{
			Model:             model_entities.LLMModel(payload.Model),
			PromptMessages:    payload.PromptMessages,
			SystemFingerprint: "test",
			Delta: model_entities.LLMResultChunkDelta{
				Index: &[]int{3}[0],
				Usage: &model_entities.LLMUsage{
					PromptTokens:     &[]int{100}[0],
					CompletionTokens: &[]int{100}[0],
					TotalTokens:      &[]int{200}[0],
					Latency:          &[]float64{0.4}[0],
					Currency:         &[]string{"USD"}[0],
				},
			},
		})
		stream.Close()
	})
	return stream, nil
}

func (m *MockedDifyInvocation) InvokeTextEmbedding(payload *mlchain_invocation.InvokeTextEmbeddingRequest) (*model_entities.TextEmbeddingResult, error) {
	result := model_entities.TextEmbeddingResult{
		Model: payload.Model,
		Usage: model_entities.EmbeddingUsage{
			Tokens:      &[]int{100}[0],
			TotalTokens: &[]int{100}[0],
			Latency:     &[]float64{0.1}[0],
			Currency:    &[]string{"USD"}[0],
		},
	}

	for range payload.Texts {
		result.Embeddings = append(result.Embeddings, []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0})
	}

	return &result, nil
}

func (m *MockedDifyInvocation) InvokeRerank(payload *mlchain_invocation.InvokeRerankRequest) (*model_entities.RerankResult, error) {
	result := model_entities.RerankResult{
		Model: payload.Model,
	}
	for i, doc := range payload.Docs {
		result.Docs = append(result.Docs, model_entities.RerankDocument{
			Index: &[]int{i}[0],
			Score: &[]float64{0.1}[0],
			Text:  &doc,
		})
	}
	return &result, nil
}

func (m *MockedDifyInvocation) InvokeTTS(payload *mlchain_invocation.InvokeTTSRequest) (*stream.Stream[model_entities.TTSResult], error) {
	stream := stream.NewStream[model_entities.TTSResult](5)
	routine.Submit(nil, func() {
		for i := 0; i < 10; i++ {
			stream.Write(model_entities.TTSResult{
				Result: "a1b2c3d4",
			})
			time.Sleep(100 * time.Millisecond)
		}
		stream.Close()
	})
	return stream, nil
}

func (m *MockedDifyInvocation) InvokeSpeech2Text(payload *mlchain_invocation.InvokeSpeech2TextRequest) (*model_entities.Speech2TextResult, error) {
	result := model_entities.Speech2TextResult{
		Result: "hello world",
	}
	return &result, nil
}

func (m *MockedDifyInvocation) InvokeModeration(payload *mlchain_invocation.InvokeModerationRequest) (*model_entities.ModerationResult, error) {
	result := model_entities.ModerationResult{
		Result: true,
	}
	return &result, nil
}

func (m *MockedDifyInvocation) InvokeTool(payload *mlchain_invocation.InvokeToolRequest) (*stream.Stream[tool_entities.ToolResponseChunk], error) {
	stream := stream.NewStream[tool_entities.ToolResponseChunk](5)
	routine.Submit(nil, func() {
		for i := 0; i < 10; i++ {
			stream.Write(tool_entities.ToolResponseChunk{
				Type: tool_entities.ToolResponseChunkTypeText,
				Message: map[string]any{
					"text": "hello world",
				},
			})
			time.Sleep(100 * time.Millisecond)
		}
		stream.Close()
	})

	return stream, nil
}

func (m *MockedDifyInvocation) InvokeApp(payload *mlchain_invocation.InvokeAppRequest) (*stream.Stream[map[string]any], error) {
	stream := stream.NewStream[map[string]any](5)
	routine.Submit(nil, func() {
		stream.Write(map[string]any{
			"event":           "agent_message",
			"message_id":      "5ad4cb98-f0c7-4085-b384-88c403be6290",
			"conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
			"answer":          "なんで",
			"created_at":      time.Now().Unix(),
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(map[string]any{
			"event":           "agent_message",
			"message_id":      "5ad4cb98-f0c7-4085-b384-88c403be6290",
			"conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
			"answer":          "春日影",
			"created_at":      time.Now().Unix(),
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(map[string]any{
			"event":           "agent_message",
			"message_id":      "5ad4cb98-f0c7-4085-b384-88c403be6290",
			"conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
			"answer":          "やったの",
			"created_at":      time.Now().Unix(),
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(map[string]any{
			"event":           "message_end",
			"id":              "5e52ce04-874b-4d27-9045-b3bc80def685",
			"conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
			"created_at":      time.Now().Unix(),
			"metadata": map[string]any{
				"retriever_resources": []map[string]any{
					{
						"position":      1,
						"dataset_id":    "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
						"dataset_name":  "あなた",
						"document_id":   "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
						"document_name": "ご自分のことばかりですのね",
						"score":         0.98457545,
						"content":       "CRYCHICは壊れてしまいましたわ",
					},
				},
				"usage": map[string]any{
					"prompt_tokens":         1033,
					"prompt_unit_price":     "0.001",
					"prompt_price_unit":     "0.001",
					"prompt_price":          "0.0010330",
					"completion_tokens":     135,
					"completion_unit_price": "0.002",
					"completion_price_unit": "0.001",
					"completion_price":      "0.0002700",
					"total_tokens":          1168,
					"total_price":           "0.0013030",
					"currency":              "USD",
					"latency":               1.381760165997548,
				},
			},
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(map[string]any{
			"event":           "message_file",
			"id":              "5e52ce04-874b-4d27-9045-b3bc80def685",
			"conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
			"belongs_to":      "assistant",
			"url":             "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			"created_at":      time.Now().Unix(),
		})
		time.Sleep(100 * time.Millisecond)
		stream.Close()
	})
	return stream, nil
}

func (m *MockedDifyInvocation) InvokeEncrypt(payload *mlchain_invocation.InvokeEncryptRequest) (map[string]any, error) {
	return payload.Data, nil
}

func (m *MockedDifyInvocation) InvokeParameterExtractor(payload *mlchain_invocation.InvokeParameterExtractorRequest) (*mlchain_invocation.InvokeNodeResponse, error) {
	resp := &mlchain_invocation.InvokeNodeResponse{
		ProcessData: map[string]any{},
		Outputs:     map[string]any{},
		Inputs: map[string]any{
			"query": payload.Query,
		},
	}

	for _, parameter := range payload.Parameters {
		typ := parameter.Type
		if typ == "string" {
			resp.Outputs[parameter.Name] = "Never gonna give you up ~"
		} else if typ == "number" {
			resp.Outputs[parameter.Name] = 1234567890
		} else if typ == "bool" {
			resp.Outputs[parameter.Name] = true
		} else if typ == "select" {
			options := parameter.Options
			if len(options) == 0 {
				resp.Outputs[parameter.Name] = "Never gonna let you down ~"
			} else {
				resp.Outputs[parameter.Name] = options[0]
			}
		} else if typ == "array[string]" {
			resp.Outputs[parameter.Name] = []string{
				"Never gonna run around and desert you ~",
				"Never gonna make you cry ~",
				"Never gonna say goodbye ~",
				"Never gonna tell a lie and hurt you ~",
			}
		} else if typ == "array[number]" {
			resp.Outputs[parameter.Name] = []int{114, 514, 1919, 810}
		} else if typ == "array[bool]" {
			resp.Outputs[parameter.Name] = []bool{true, false, true, false, true, false, true, false, true, false}
		} else if typ == "array[object]" {
			resp.Outputs[parameter.Name] = []map[string]any{
				{
					"name": "お願い",
					"age":  55555,
				},
				{
					"name": "何でもするがら",
					"age":  99999,
				},
			}
		}
	}

	return resp, nil
}

func (m *MockedDifyInvocation) InvokeQuestionClassifier(payload *mlchain_invocation.InvokeQuestionClassifierRequest) (*mlchain_invocation.InvokeNodeResponse, error) {
	return &mlchain_invocation.InvokeNodeResponse{
		ProcessData: map[string]any{},
		Outputs: map[string]any{
			"class_name": payload.Classes[0].Name,
		},
		Inputs: map[string]any{},
	}, nil
}

func (m *MockedDifyInvocation) InvokeSummary(payload *mlchain_invocation.InvokeSummaryRequest) (*mlchain_invocation.InvokeSummaryResponse, error) {
	return &mlchain_invocation.InvokeSummaryResponse{
		Summary: payload.Text,
	}, nil
}

func (m *MockedDifyInvocation) UploadFile(payload *mlchain_invocation.UploadFileRequest) (*mlchain_invocation.UploadFileResponse, error) {
	return &mlchain_invocation.UploadFileResponse{
		URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	}, nil
}