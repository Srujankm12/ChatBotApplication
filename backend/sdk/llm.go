package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"chatbot/internal/models"
)

const (
	defaultGeminiModel    = "gemini-1.5-flash"
	geminiProvider        = "google"
	geminiEndpointPattern = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent"
	retryDelay            = 400 * time.Millisecond
)

func getProvider() string {
	if p := os.Getenv("LLM_PROVIDER"); p != "" {
		return p
	}
	return "gemini"
}

func getModel() string {
	if m := os.Getenv("GEMINI_MODEL"); m != "" {
		return m
	}
	return defaultGeminiModel
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func buildPrompt(messages []models.Message) string {
	var sb strings.Builder
	for _, m := range messages {
		role := "User"
		if m.Role == "assistant" {
			role = "Assistant"
		}
		sb.WriteString(fmt.Sprintf("%s: %s\n", role, m.Content))
	}
	sb.WriteString("Assistant:")
	return sb.String()
}

func preview(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}

// estimateTokens approximates token count using the ~4 chars/token heuristic.
func estimateTokens(s string) int {
	return len(s) / 4
}

func sendLog(entry models.InferenceLog) {
	payload, err := json.Marshal(entry)
	if err != nil {
		return
	}

	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, backendURL+"/logs", bytes.NewReader(payload))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

// doHTTPPost sends a POST request and retries once after retryDelay on failure.
func doHTTPPost(ctx context.Context, url string, bodyBytes []byte) ([]byte, error) {
	attempt := func() ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return io.ReadAll(resp.Body)
	}

	data, err := attempt()
	if err != nil {
		if ctx.Err() != nil {
			return nil, err
		}
		time.Sleep(retryDelay)
		data, err = attempt()
	}
	return data, err
}

func callGemini(ctx context.Context, messages []models.Message, sessionID string) (string, error) {
	model := getModel()
	log.Println("Using model:", model)

	apiKey := os.Getenv("GEMINI_API_KEY")
	endpoint := fmt.Sprintf(geminiEndpointPattern, model)
	prompt := buildPrompt(messages)
	inputPreview := preview(prompt, 100)
	tokenUsage := estimateTokens(prompt)

	start := time.Now()

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		go sendLog(models.InferenceLog{
			Model: model, Provider: geminiProvider,
			LatencyMs: time.Since(start).Milliseconds(),
			InputPreview: inputPreview, OutputPreview: "", TokenUsage: tokenUsage,
			Timestamp: time.Now().UTC(), SessionID: sessionID, Status: "error",
		})
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s?key=%s", endpoint, apiKey)

	respBytes, err := doHTTPPost(ctx, url, bodyBytes)
	latencyMs := time.Since(start).Milliseconds()

	if err != nil {
		go sendLog(models.InferenceLog{
			Model: model, Provider: geminiProvider,
			LatencyMs: latencyMs, InputPreview: inputPreview, OutputPreview: "", TokenUsage: tokenUsage,
			Timestamp: time.Now().UTC(), SessionID: sessionID, Status: "error",
		})
		return "", fmt.Errorf("call gemini: %w", err)
	}

	var gemResp geminiResponse
	if err := json.Unmarshal(respBytes, &gemResp); err != nil {
		go sendLog(models.InferenceLog{
			Model: model, Provider: geminiProvider,
			LatencyMs: latencyMs, InputPreview: inputPreview, OutputPreview: "", TokenUsage: tokenUsage,
			Timestamp: time.Now().UTC(), SessionID: sessionID, Status: "error",
		})
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if gemResp.Error != nil {
		go sendLog(models.InferenceLog{
			Model: model, Provider: geminiProvider,
			LatencyMs: latencyMs, InputPreview: inputPreview, OutputPreview: "", TokenUsage: tokenUsage,
			Timestamp: time.Now().UTC(), SessionID: sessionID, Status: "error",
		})
		return "", fmt.Errorf("gemini API error: %s", gemResp.Error.Message)
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		go sendLog(models.InferenceLog{
			Model: model, Provider: geminiProvider,
			LatencyMs: latencyMs, InputPreview: inputPreview, OutputPreview: "", TokenUsage: tokenUsage,
			Timestamp: time.Now().UTC(), SessionID: sessionID, Status: "error",
		})
		return "", fmt.Errorf("empty response from gemini")
	}

	reply := gemResp.Candidates[0].Content.Parts[0].Text
	outputPreview := preview(reply, 100)

	go sendLog(models.InferenceLog{
		Model: model, Provider: geminiProvider,
		LatencyMs: latencyMs, InputPreview: inputPreview, OutputPreview: outputPreview, TokenUsage: tokenUsage,
		Timestamp: time.Now().UTC(), SessionID: sessionID, Status: "success",
	})

	return reply, nil
}

// CallLLM routes to the appropriate provider based on the LLM_PROVIDER env var.
func CallLLM(ctx context.Context, messages []models.Message, sessionID string) (string, error) {
	provider := getProvider()
	switch provider {
	case "gemini":
		return callGemini(ctx, messages, sessionID)
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}
