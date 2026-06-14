package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	News     NewsConfig     `json:"news"`
	Summary  SummaryConfig  `json:"summary"`
	AI       AIConfig       `json:"ai"`
}

type ServerConfig struct {
	Port int `json:"port"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type NewsConfig struct {
	RSSSources      []string `json:"rss_sources"`
	WechatCookie    string   `json:"wechat_cookie"`
	RSSHubURL       string   `json:"rsshub_url"`
	RetentionDays   int      `json:"retention_days"`
	FetchInterval   int      `json:"fetch_interval"`
}

type SummaryConfig struct {
	APIKey  string `json:"api_key"`
	Model   string `json:"model"`
	MaxLen  int    `json:"max_len"`
}

type AIConfig struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "",
			DBName:   "ai_watcher",
		},
		News: NewsConfig{
			RSSSources: []string{
				"https://www.jiqizhixin.com/rss",
				"https://www.qbitai.com/feed",
			},
			WechatCookie:  "",
			RSSHubURL:     "http://localhost:1200",
			RetentionDays: 30,
			FetchInterval: 30,
		},
		Summary: SummaryConfig{
			APIKey: "",
			Model:  "qwen-turbo",
			MaxLen: 50,
		},
		AI: AIConfig{
			APIKey:  "",
			BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		},
	}
}
