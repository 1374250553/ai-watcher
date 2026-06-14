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

type CommunityFetcher struct {
	db     *api.Database
	config *config.Config
	client *http.Client
}

func NewCommunityFetcher(db *api.Database, cfg *config.Config) *CommunityFetcher {
	return &CommunityFetcher{
		db:     db,
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (f *CommunityFetcher) FetchAll() {
	go f.FetchZhihu()
	go f.FetchV2EX()
}

func (f *CommunityFetcher) FetchZhihu() {
	url := "https://www.zhihu.com/topic/19552832/hot"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create Zhihu request: %v", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := f.client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch Zhihu: %v", err)
		f.db.LogFetch("zhihu", "error", err.Error(), 0)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Failed to parse Zhihu: %v", err)
		f.db.LogFetch("zhihu", "error", err.Error(), 0)
		return
	}

	count := 0
	doc.Find(".ContentItem").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".ContentItem-title a").Text()
		link, _ := s.Find(".ContentItem-title a").Attr("href")
		description := s.Find(".RichContent-inner").Text()

		if title == "" || link == "" {
			return
		}

		if len(link) > 0 && link[0] == '/' {
			link = "https://www.zhihu.com" + link
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
			Source:  "知乎",
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

	f.db.LogFetch("zhihu", status, message, count)
	log.Printf("Fetched %d articles from Zhihu", count)
}

func (f *CommunityFetcher) FetchV2EX() {
	url := "https://www.v2ex.com/?tab=creative"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create V2EX request: %v", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := f.client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch V2EX: %v", err)
		f.db.LogFetch("v2ex", "error", err.Error(), 0)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Failed to parse V2EX: %v", err)
		f.db.LogFetch("v2ex", "error", err.Error(), 0)
		return
	}

	count := 0
	doc.Find(".item_title a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		link, _ := s.Attr("href")

		if title == "" || link == "" {
			return
		}

		if len(link) > 0 && link[0] == '/' {
			link = "https://www.v2ex.com" + link
		}

		exists, err := f.db.NewsExists(link)
		if err != nil {
			log.Printf("Failed to check news existence: %v", err)
			return
		}

		if exists {
			return
		}

		news := &models.News{
			Title:   title,
			URL:     link,
			Summary: "",
			Source:  "V2EX",
			Content: "",
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

	f.db.LogFetch("v2ex", status, message, count)
	log.Printf("Fetched %d articles from V2EX", count)
}
