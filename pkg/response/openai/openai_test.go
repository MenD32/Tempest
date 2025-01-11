package openai_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/MenD32/Tempest/pkg/response/openai"
	"github.com/stretchr/testify/assert"
)

func TestOpenAIResponseBuilder(t *testing.T) {
	tests := []struct {
		name       string
		httpBody   string
		sent       time.Time
		wantTokens int
		wantErr    bool
	}{
		{
			name:       "first assistant token", // idk why but first token is always empty, see openai docs https://platform.openai.com/docs/api-reference/chat/create, with streaming enabled
			httpBody:   `data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1694268190,"model":"gpt-4o-mini", "system_fingerprint": "fp_44709d6fcb", "choices":[{"index":0,"delta":{"role":"assistant","content":""},"logprobs":null,"finish_reason":null}]}`,
			sent:       time.Now(),
			wantTokens: 1,
			wantErr:    false,
		},
		{
			name:       "valid token",
			httpBody:   `data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1694268190,"model":"gpt-4o-mini", "system_fingerprint": "fp_44709d6fcb", "choices":[{"index":0,"delta":{"content":"Hello"},"logprobs":null,"finish_reason":null}]}`,
			sent:       time.Now(),
			wantTokens: 1,
			wantErr:    false,
		},
		{
			name:       "finish reason: stop token",
			httpBody:   `data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1694268190,"model":"gpt-4o-mini", "system_fingerprint": "fp_44709d6fcb", "choices":[{"index":0,"delta":{},"logprobs":null,"finish_reason":"stop"}]}`,
			sent:       time.Now(),
			wantTokens: 1,
			wantErr:    false,
		},
		{
			name:       "stream end token",
			httpBody:   `data: [DONE]`,
			sent:       time.Now(),
			wantTokens: 0,
			wantErr:    false,
		},
		{
			name:       "empty response",
			httpBody:   "",
			sent:       time.Now(),
			wantTokens: 0,
			wantErr:    false,
		},
		{
			name: "Multiple tokens",
			httpBody: `
data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1694268190,"model":"gpt-4o-mini", "system_fingerprint": "fp_44709d6fcb", "choices":[{"index":0,"delta":{"role":"assistant","content":""},"logprobs":null,"finish_reason":null}]}
data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1694268190,"model":"gpt-4o-mini", "system_fingerprint": "fp_44709d6fcb", "choices":[{"index":0,"delta":{"content":"Hello"},"logprobs":null,"finish_reason":null}]}
data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1728933352,"model":"gpt-4o-mini", "system_fingerprint": "fp_44709d6fcb", "choices":[],"usage":{"prompt_tokens":19,"completion_tokens":10,"total_tokens":29}}
data: [DONE]
`,
			sent:       time.Now(),
			wantTokens: 3,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				Body:       io.NopCloser(bytes.NewBufferString(tt.httpBody)),
				StatusCode: http.StatusOK,
			}
			got, err := openai.OpenAIResponseBuilder(resp, tt.sent)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenAIResponseBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				openAIResp, ok := got.(openai.OpenAIResponse)
				assert.True(t, ok)
				assert.Equal(t, tt.wantTokens, len(openAIResp.Tokens))
			}
		})
	}
}

func TestOpenAIResponse_Verify(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []openai.Token
		wantErr bool
	}{
		{
			name: "valid tokens",
			tokens: []openai.Token{
				{
					ID: "token1",
					Choices: []openai.Choice{
						{Delta: openai.Delta{Content: "Hello"}},
					},
				},
				{
					ID: "token2",
					Usage: openai.Usage{
						CompletionTokens: 1,
						PromptTokens:     1,
						TotalTokens:      2,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no tokens",
			tokens: []openai.Token{
				{
					ID: "token1",
					Choices: []openai.Choice{
						{Delta: openai.Delta{Content: "Hello"}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid token",
			tokens: []openai.Token{
				{
					ID: "token1",
					Choices: []openai.Choice{
						{Delta: openai.Delta{Content: "Hello"}},
					},
				},
				{
					ID: "token2",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := openai.OpenAIResponse{
				Tokens: tt.tokens,
			}
			err := resp.Verify()
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
