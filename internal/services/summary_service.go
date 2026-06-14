package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"ai-watcher/internal/api"
	"ai-watcher/internal/config"
)

type SummaryService struct {
	db     *api.Database
	config *config.Config
	client *http.Client
}

func NewSummaryService(db *api.Database, cfg *config.Config) *SummaryService {
	return &SummaryService{
		db:     db,
		config: cfg,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

type QwenRequest struct {
	Model    string        `json:"model"`
	Messages []QwenMessage `json:"messages"`
}

type QwenMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type QwenResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (s *SummaryService) GenerateSummary(title, content string) (string, error) {
	if s.config.Summary.APIKey == "" {
		return "", fmt.Errorf("API key not configured")
	}

	prompt := fmt.Sprintf("请用一句话（50字以内）总结以下文章的核心内容：\n\n标题：%s\n\n内容：%s", title, content)

	reqBody := QwenRequest{
		Model: s.config.Summary.Model,
		Messages: []QwenMessage{
			{Role: "system", Content: "你是一个专业的AI资讯摘要助手，请用简洁的语言总结文章核心内容。"},
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.Summary.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var qwenResp QwenResponse
	if err := json.Unmarshal(body, &qwenResp); err != nil {
		return "", err
	}

	if len(qwenResp.Choices) == 0 {
		return "", fmt.Errorf("no response from Qwen")
	}

	summary := qwenResp.Choices[0].Message.Content
	if len(summary) > s.config.Summary.MaxLen {
		summary = summary[:s.config.Summary.MaxLen]
	}

	return summary, nil
}

func (s *SummaryService) SummarizeAll() {
	news, _, err := s.db.GetNews(1, 100, "", "")
	if err != nil {
		log.Printf("Failed to get news for summarization: %v", err)
		return
	}

	for _, n := range news {
		if n.Summary != "" && len(n.Summary) <= s.config.Summary.MaxLen {
			continue
		}

		summary, err := s.GenerateSummary(n.Title, n.Content)
		if err != nil {
			log.Printf("Failed to generate summary for news %d: %v", n.ID, err)
			continue
		}

		n.Summary = summary
		if err := s.db.InsertNews(&n); err != nil {
			log.Printf("Failed to update news summary: %v", err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
