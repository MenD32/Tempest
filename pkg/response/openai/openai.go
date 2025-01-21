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

	MIN_TOKEN_COUNT = 2 // 1 for at least 1 token, 1 for usage
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
	tokens := m.GetTokens()

	ttft := tokens[0].Timestamp.Sub(m.Sent)
	e2e := tokens[len(tokens)-1].Timestamp.Sub(m.Sent)
	itl := time.Duration(0)
	if len(tokens) > 1 {
		itl = (e2e - ttft) / time.Duration(len(tokens)-1)
	}

	metrics := map[string]interface{}{
		"input_tokens":  usage.PromptTokens,
		"output_tokens": usage.CompletionTokens,
		"ttft_ms":       getDurationMilliseconds(ttft.Abs()),
		"e2e_ms":        getDurationMilliseconds(e2e.Abs()),
		"itl_ms":        getDurationMilliseconds(itl.Abs()),
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

func (m OpenAIResponse) GetTokens() []Token {
	return m.Tokens[1 : len(m.Tokens)-1] // start token is empty, usage token is unecessary
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
	klog.Info("Creating OpenAI response")
	var tokens = []Token{}
	var tokenTimestamp time.Time

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		raw := scanner.Text()
		klog.V(9).Infof("got raw: '%s'\n", raw)
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
	body, err := json.Marshal(m.Tokens)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (m OpenAIResponse) Verify() error {
	if len(m.Tokens) < MIN_TOKEN_COUNT {
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

func getDurationMilliseconds(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}
