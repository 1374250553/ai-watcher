package api

import (
	"sort"
	"strings"
	"sync"
	"time"

	"ai-watcher/internal/config"
	"ai-watcher/internal/models"
)

type Database struct {
	mu       sync.RWMutex
	news     map[string]*models.News
	resources []models.APIResource
	logs     []models.FetchLog
	nextID   int64
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	db := &Database{
		news:     make(map[string]*models.News),
		resources: defaultAPIResources(),
		nextID:   1,
	}
	return db, nil
}

func (db *Database) Close() error { return nil }

func (db *Database) InsertNews(news *models.News) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if existing, ok := db.news[news.URL]; ok {
		existing.Summary = news.Summary
		existing.Content = news.Content
		existing.UpdatedAt = time.Now()
		return nil
	}

	news.ID = db.nextID
	db.nextID++
	news.CreatedAt = time.Now()
	news.UpdatedAt = time.Now()
	db.news[news.URL] = news
	return nil
}

func (db *Database) NewsExists(url string) (bool, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	_, ok := db.news[url]
	return ok, nil
}

func (db *Database) GetNews(page, pageSize int, source, search string) ([]models.News, int, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var filtered []*models.News
	for _, n := range db.news {
		if source != "" && n.Source != source {
			continue
		}
		if search != "" {
			lowerSearch := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(n.Title), lowerSearch) && !strings.Contains(strings.ToLower(n.Summary), lowerSearch) {
				continue
			}
		}
		filtered = append(filtered, n)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := len(filtered)
	start := (page - 1) * pageSize
	if start >= total {
		return nil, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	var result []models.News
	for _, n := range filtered[start:end] {
		result = append(result, *n)
	}
	return result, total, nil
}

func (db *Database) GetAPIResources() ([]models.APIResource, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	result := make([]models.APIResource, len(db.resources))
	copy(result, db.resources)
	return result, nil
}

func (db *Database) InsertAPIResource(resource *models.APIResource) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	resource.ID = db.nextID
	db.nextID++
	resource.CreatedAt = time.Now()
	resource.LastUpdated = time.Now()
	db.resources = append(db.resources, *resource)
	return nil
}

func (db *Database) LogFetch(source, status, message string, count int) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.logs = append(db.logs, models.FetchLog{
		ID:        db.nextID,
		Source:    source,
		Status:    status,
		Message:   message,
		Count:     count,
		CreatedAt: time.Now(),
	})
	db.nextID++
	return nil
}

func (db *Database) CleanOldNews(days int) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	cutoff := time.Now().AddDate(0, 0, -days)
	for url, n := range db.news {
		if n.CreatedAt.Before(cutoff) {
			delete(db.news, url)
		}
	}
	return nil
}

func defaultAPIResources() []models.APIResource {
	return []models.APIResource{
		{ID: 1, Name: "Qwen-Turbo", Provider: "阿里云", Description: "通义千问轻量级模型，适合摘要生成", Endpoint: "https://dashscope.aliyuncs.com/api/v1", FreeQuota: "100万tokens/月免费", DocURL: "https://help.aliyun.com/zh/dashscope", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
		{ID: 2, Name: "DeepSeek V3", Provider: "DeepSeek", Description: "DeepSeek 大语言模型", Endpoint: "https://api.deepseek.com/v1", FreeQuota: "500万tokens/月免费", DocURL: "https://platform.deepseek.com", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
		{ID: 3, Name: "GLM-4-Flash", Provider: "智谱AI", Description: "GLM-4 系列轻量模型", Endpoint: "https://open.bigmodel.cn/api/paas/v4", FreeQuota: "100万tokens/月免费", DocURL: "https://bigmodel.cn", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
		{ID: 4, Name: "ERNIE-Speed-128K", Provider: "百度智能云", Description: "文心一言轻量版", Endpoint: "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie-speed-128k", FreeQuota: "每天免费10000次", DocURL: "https://cloud.baidu.com/doc/WENXINWORKSHOP", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
		{ID: 5, Name: "Spark Lite", Provider: "科大讯飞", Description: "星火认知大模型轻量版", Endpoint: "https://spark-api-open.xf-yun.com/v1/chat/completions", FreeQuota: "每天免费5000次", DocURL: "https://www.xfyun.cn/doc/spark/", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
		{ID: 6, Name: "Doubao-lite", Provider: "字节跳动", Description: "豆包系列轻量模型", Endpoint: "https://ark.cn-beijing.volces.com/api/v3", FreeQuota: "50万tokens/月免费", DocURL: "https://www.volcengine.com/docs/82379", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
		{ID: 7, Name: "Yi-Lightning", Provider: "零一万物", Description: "零一万物轻量模型", Endpoint: "https://api.lingyiwanwu.com/v1", FreeQuota: "50万tokens/月免费", DocURL: "https://platform.lingyiwanwu.com", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
		{ID: 8, Name: "MiniMax-abab6.5s", Provider: "MiniMax", Description: "MiniMax 轻量模型", Endpoint: "https://api.minimax.chat/v1", FreeQuota: "每天免费1000次", DocURL: "https://platform.minimaxi.com", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
	}
}
