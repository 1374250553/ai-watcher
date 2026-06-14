package services

import (
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"ai-watcher/internal/api"
	"ai-watcher/internal/config"
	"ai-watcher/internal/models"
)

type WechatFetcher struct {
	db     *api.Database
	config *config.Config
	client *http.Client
}

func NewWechatFetcher(db *api.Database, cfg *config.Config) *WechatFetcher {
	return &WechatFetcher{
		db:     db,
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (f *WechatFetcher) FetchAll() {
	if f.config.News.WechatCookie == "" {
		log.Println("WeChat cookie not set, skipping WeChat fetch")
		return
	}

	rssURL := f.config.News.RSSHubURL + "/wechat/mp/msgid"
	req, err := http.NewRequest("GET", rssURL, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("Cookie", f.config.News.WechatCookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := f.client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch WeChat RSS: %v", err)
		f.db.LogFetch("wechat", "error", err.Error(), 0)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Failed to parse WeChat RSS: %v", err)
		f.db.LogFetch("wechat", "error", err.Error(), 0)
		return
	}

	count := 0
	doc.Find("item").Each(func(i int, s *goquery.Selection) {
		title := s.Find("title").Text()
		link := s.Find("link").Text()
		description := s.Find("description").Text()

		if title == "" || link == "" {
			return
		}

		exists, err := f.db.NewsExists(link)
		if err != nil {
			log.Printf("Failed to check news existence: %v", err)
			return
		}

		if exists {
			return
		}

		summary := description
		if len(summary) > 200 {
			summary = summary[:200] + "..."
		}

		news := &models.News{
			Title:   title,
			URL:     link,
			Summary: summary,
			Source:  "微信公众号",
			Content: description,
		}

		if err := f.db.InsertNews(news); err != nil {
			log.Printf("Failed to insert news: %v", err)
			return
		}

		count++
		time.Sleep(100 * time.Millisecond)
	})

	status := "success"
	message := "Fetched successfully"
	if count == 0 {
		status = "warning"
		message = "No new articles found"
	}

	f.db.LogFetch("wechat", status, message, count)
	log.Printf("Fetched %d articles from WeChat", count)
}
