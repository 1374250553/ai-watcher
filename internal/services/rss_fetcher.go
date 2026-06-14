package services

import (
	"log"
	"time"

	"github.com/mmcdole/gofeed"
	"ai-watcher/internal/api"
	"ai-watcher/internal/config"
	"ai-watcher/internal/models"
)

type RSSFetcher struct {
	db       *api.Database
	config   *config.Config
	parser   *gofeed.Parser
}

func NewRSSFetcher(db *api.Database, cfg *config.Config) *RSSFetcher {
	return &RSSFetcher{
		db:     db,
		config: cfg,
		parser: gofeed.NewParser(),
	}
}

func (f *RSSFetcher) FetchAll() {
	for _, source := range f.config.News.RSSSources {
		go f.FetchSource(source)
	}
}

func (f *RSSFetcher) FetchSource(sourceURL string) {
	feed, err := f.parser.ParseURL(sourceURL)
	if err != nil {
		log.Printf("Failed to fetch RSS from %s: %v", sourceURL, err)
		f.db.LogFetch(sourceURL, "error", err.Error(), 0)
		return
	}

	count := 0
	for _, item := range feed.Items {
		exists, err := f.db.NewsExists(item.Link)
		if err != nil {
			log.Printf("Failed to check news existence: %v", err)
			continue
		}

		if exists {
			continue
		}

		summary := item.Description
		if len(summary) > 200 {
			summary = summary[:200] + "..."
		}

		news := &models.News{
			Title:   item.Title,
			URL:     item.Link,
			Summary: summary,
			Source:  feed.Title,
			Content: item.Description,
		}

		if err := f.db.InsertNews(news); err != nil {
			log.Printf("Failed to insert news: %v", err)
			continue
		}

		count++
		time.Sleep(100 * time.Millisecond)
	}

	status := "success"
	message := "Fetched successfully"
	if count == 0 {
		status = "warning"
		message = "No new articles found"
	}

	f.db.LogFetch(sourceURL, status, message, count)
	log.Printf("Fetched %d articles from %s", count, feed.Title)
}
