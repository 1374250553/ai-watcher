package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"ai-watcher/internal/config"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type AIHandler struct {
	config  *config.Config
	baseURL string
	apiKey  string
}

var textChatModels = map[string][]string{
	"通义千问 Qwen3.7": {"qwen3.7-max", "qwen3.7-max-preview", "qwen3.7-plus"},
	"通义千问 Qwen3.6": {"qwen3.6-max-preview", "qwen3.6-plus", "qwen3.6-flash"},
	"通义千问 Qwen3.5": {"qwen3.5-plus", "qwen3.5-flash", "qwen3.5-397b-a17b", "qwen3.5-122b-a10b", "qwen3.5-35b-a3b", "qwen3.5-27b"},
	"通义千问 Qwen3":   {"qwen3-max", "qwen3-max-preview", "qwen3-plus", "qwen3-235b-a22b", "qwen3-32b", "qwen3-30b-a3b", "qwen3-14b", "qwen3-8b"},
	"通义千问经典":     {"qwen-plus", "qwen-turbo", "qwen-max", "qwen-flash", "qwen-long", "qwq-plus"},
	"通义千问代码":     {"qwen3-coder-next", "qwen3-coder-plus", "qwen3-coder-flash"},
	"通义千问思考":     {"qwen3-next-80b-a3b-thinking", "qwen3-235b-a22b-thinking-2507", "qwen3-30b-a3b-thinking-2507"},
	"DeepSeek":         {"deepseek-v4-flash", "deepseek-v4-pro", "deepseek-v3.2", "deepseek-v3.1", "deepseek-v3", "deepseek-r1", "deepseek-r1-distill-qwen-32b", "deepseek-r1-distill-qwen-14b", "deepseek-r1-distill-llama-70b"},
	"Kimi 月之暗面":     {"kimi-k2.6", "kimi-k2.5", "kimi-k2-thinking"},
	"GLM 智谱":         {"glm-5", "glm-5.1", "glm-4.7"},
	"MiniMax":          {"MiniMax-M2.5", "MiniMax-M2.1"},
}

func NewAIHandler(cfg *config.Config) *AIHandler {
	return &AIHandler{
		config:  cfg,
		baseURL: cfg.AI.BaseURL,
		apiKey:  cfg.AI.APIKey,
	}
}

func (h *AIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/ai/models":
		h.handleModels(w, r)
	case "/api/ai/chat":
		h.handleChat(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *AIHandler) handleModels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	total := 0
	for _, models := range textChatModels {
		total += len(models)
	}

	resp := map[string]interface{}{
		"categories": textChatModels,
		"total":      total,
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *AIHandler) handleChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	if h.apiKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "DASHSCOPE_API_KEY not configured",
			"details": "Please set ai.api_key in config.json",
		})
		return
	}

	var req struct {
		Model    string        `json:"model"`
		Messages []ChatMessage `json:"messages"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Messages == nil || len(req.Messages) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "messages field is required"})
		return
	}

	if req.Model == "" {
		req.Model = "qwen-turbo"
	}

	reqBody := ChatRequest{
		Model:    req.Model,
		Messages: req.Messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal error"})
		return
	}

	apiReq, err := http.NewRequest("POST", h.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal error"})
		return
	}

	apiReq.Header.Set("Content-Type", "application/json")
	apiReq.Header.Set("Authorization", "Bearer "+h.apiKey)

	client := &http.Client{}
	apiResp, err := client.Do(apiReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Request failed: %v", err),
		})
		return
	}
	defer apiResp.Body.Close()

	if apiResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(apiResp.Body)
		w.WriteHeader(apiResp.StatusCode)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   fmt.Sprintf("DashScope API error: %s", apiResp.Status),
			"details": string(body),
		})
		return
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(apiResp.Body).Decode(&chatResp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to decode response"})
		return
	}

	json.NewEncoder(w).Encode(chatResp)
}
