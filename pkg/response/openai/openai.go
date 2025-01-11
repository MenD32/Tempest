package openai

import (
	"bufio"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MenD32/Tempest/pkg/response"
	"k8s.io/klog/v2"
)

const (
	END_TOKEN    = "data: [DONE]"
	TOKEN_PREFIX = "data: "
)

type OpenAIResponse struct {
	Sent   time.Time `json:"sent"`
	Tokens []Token   `json:"tokens"`
}

type Token struct {
	Timestamp         time.Time `json:"timestamp"`
	ID                string    `json:"id"`
	Choices           []Choice  `json:"choices"`
	Created           int64     `json:"created"`
	Model             string    `json:"model"`
	ServiceTier       string    `json:"service_tier,omitempty"`
	SystemFingerprint string    `json:"system_fingerprint"`
	Object            string    `json:"object"`
	Usage             Usage     `json:"usage,omitempty"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason,omitempty"`
	Index        int    `json:"index"`
}

type Delta struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role"`
}

func (m OpenAIResponse) Metrics() (*response.Metrics, error) {

	body, err := m.Body()
	if err != nil {
		return nil, err
	}

	usage := m.GetUsage()

	ttft := m.Tokens[0].Timestamp.Sub(m.Sent)
	e2e := m.Tokens[len(m.Tokens)-1].Timestamp.Sub(m.Sent)

	metrics := map[string]interface{}{
		"input_tokens":  usage.PromptTokens,
		"output_tokens": usage.CompletionTokens,
		"ttft_ms":       ttft.Abs().Milliseconds(),
		"e2e_ms":        e2e.Abs().Milliseconds(),
	}

	return &response.Metrics{
		Sent:    m.Sent,
		Body:    body,
		Metrics: metrics,
	}, nil
}

func (m OpenAIResponse) GetUsage() Usage {
	return m.Tokens[len(m.Tokens)-1].Usage
}

func NewToken(chunk []byte) (*Token, error) {
	var NewToken Token
	err := json.Unmarshal(chunk, &NewToken)
	if err != nil {
		return nil, err
	}
	return &NewToken, nil
}

func OpenAIResponseBuilder(resp *http.Response, sent time.Time) (response.Response, error) {
	var tokens = []Token{}
	var tokenTimestamp time.Time

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		raw := scanner.Text()
		klog.V(1).Infof("got raw: '%s'\n", raw)
		tokenTimestamp = time.Now()
		if len(raw) == 0 || raw == END_TOKEN {
			continue
		}
		if raw[:len(TOKEN_PREFIX)] == TOKEN_PREFIX {
			raw = raw[len(TOKEN_PREFIX):]
		}
		token, err := NewToken([]byte(raw))
		if err != nil {
			return nil, err
		}
		token.Timestamp = tokenTimestamp
		tokens = append(tokens, *token)
	}

	return OpenAIResponse{Tokens: tokens, Sent: sent}, nil
}

func GetMilliseconds(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

func (m OpenAIResponse) Body() ([]byte, error) {
	body, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (m OpenAIResponse) Verify() error {
	if len(m.Tokens) == 0 {
		return ErrNoTokens
	}
	for _, token := range m.Tokens {
		if len(token.Choices) == 0 && token.Usage == (Usage{}) {
			klog.Errorf("%+v\n", token)
			return ErrInvalidToken
		}
	}
	return nil
}
