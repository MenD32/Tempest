package client

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"os"
	"time"

	"math/rand"
)

const TOKEN_PREFIX = "data: "
const END_TOKEN = "data: [DONE]"

type Token struct {
	Timestamp time.Time
	Token     string
}

func (t Token) Trim() (string, error) {
	if len(t.Token) < len(TOKEN_PREFIX) {
		return "", fmt.Errorf("Token is malformed")
	}
	return t.Token[len(TOKEN_PREFIX):], nil
}

func (t Token) Unmarshal() (map[string]interface{}, error) {
	var data map[string]interface{}
	trim, err := t.Trim()
	if err != nil {
		return data, err
	}
	err = json.Unmarshal([]byte(trim), &data)
	return data, err
}

type ParsedResponse struct {
	Timestamp time.Time
	Tokens    []Token
}

type Request struct {
	Timedelta time.Duration `json:"timedelta"`
}

type Response struct {
	Timestamp time.Time     `json:"timestamp"`
	Response  http.Response `json:"response"`
}

type Output struct {
	Timestamp time.Time `json:"timestamp"`
	Metrics   metrics   `json:"metrics"`
}

func (pr ParsedResponse) ParseMetrics() *metrics {
	var ttft time.Duration = pr.Tokens[0].Timestamp.Sub(pr.Timestamp)
	var e2e time.Duration = pr.Tokens[len(pr.Tokens)-1].Timestamp.Sub(pr.Timestamp)
	var itl time.Duration = (e2e - ttft) / time.Duration(len(pr.Tokens)-1)

	var usage_token Token = pr.Tokens[len(pr.Tokens)-1]
	usage, err := usage_token.Unmarshal()
	if err != nil {
		log.Printf("failed to unmarshal usage token: %v", err)
		return nil
	}

	usageDict := usage["usage"].(map[string]interface{})
	usageJson, err := json.Marshal(usageDict)
	if err != nil {
		log.Printf("failed to marshal usage to JSON: %v", err)
		return nil
	}
	fmt.Printf("Usage: %s\n", usageJson)

	input_tokens := int(usageDict["completion_tokens"].(float64))
	output_tokens := int(usageDict["prompt_tokens"].(float64))

	return &metrics{
		InputTokens:  input_tokens,
		OutputTokens: output_tokens,
		TTFT:         time.Duration(ttft.Milliseconds()),
		E2E:          time.Duration(e2e.Milliseconds()),
		ITL:          time.Duration(itl.Milliseconds()),
	}
}

func (r *Request) Send(resp chan<- *ParsedResponse, wg *sync.WaitGroup) {

	defer wg.Done()

	body, err := buildRequestbody()
	if err != nil {
		log.Printf("failed to build request body: %v", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/v1/chat/completions", strings.NewReader(body))
	if err != nil {
		log.Printf("failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	send_timestamp := time.Now()
	response, err := client.Do(req)
	if err != nil {
		log.Printf("failed to send request: %v", err)
		return
	}

	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	TokenList := []Token{}
	for scanner.Scan() {
		token := scanner.Text()
		fmt.Printf("%s\n", token)
		if len(token) == 0 || token == END_TOKEN {
			continue
		}
		TokenList = append(TokenList, Token{Timestamp: time.Now(), Token: token})
	}

	pr := ParsedResponse{Timestamp: send_timestamp, Tokens: TokenList}
	resp <- &pr

}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func getPrompt() []string {
	prompts := readCsvFile("/Users/mend/Downloads/prompts.csv")
	if len(prompts) == 0 {
		return []string{}
	}
	randomIndex := rand.Intn(len(prompts))
	return prompts[randomIndex]
}

func buildRequestbody() (string, error) {

	prompt := getPrompt()
	if len(prompt) != 2 {
		return "", fmt.Errorf("prompt is Malformed")
	}

	body := make(map[string]interface{})
	body["model"] = "Qwen/Qwen2-7B-Instruct"
	body["messages"] = []map[string]string{
		{
			"role":    "user",
			"content": prompt[1],
		},
	}
	body["stream"] = true
	body["stream_options"] = map[string]interface{}{
		"include_usage": true,
	}

	body["max_completion_tokens"] = 100

	jsonData, err := json.Marshal(body)
	fmt.Printf("Request Body: %s\n", jsonData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body to JSON: %v", err)
	}

	return string(jsonData), nil
}

type metrics struct {
	InputTokens  int           `json:"input_tokens"`
	OutputTokens int           `json:"output_tokens"`
	TTFT         time.Duration `json:"ttft"`
	E2E          time.Duration `json:"e2e"`
	ITL          time.Duration `json:"itl"`
}
