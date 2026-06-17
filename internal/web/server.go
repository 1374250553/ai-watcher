package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"ai-watcher/internal/api"
	"ai-watcher/internal/config"
	"ai-watcher/internal/services"
)

type Server struct {
	db        *api.Database
	config    *config.Config
	fetcher   *services.RSSFetcher
	templates *template.Template
	aiHandler *api.AIHandler
}

func NewServer(db *api.Database, cfg *config.Config) *Server {
	funcMap := template.FuncMap{
		"seq": func(start, end int) []int {
			var s []int
			for i := start; i <= end; i++ {
				s = append(s, i)
			}
			return s
		},
	}

	tmpl, err := template.New("*.html").Funcs(funcMap).ParseGlob("web/templates/*.html")
	if err != nil {
		log.Printf("Failed to parse templates: %v", err)
	}

	rssFetcher := services.NewRSSFetcher(db, cfg)

	return &Server{
		db:        db,
		config:    cfg,
		fetcher:   rssFetcher,
		templates: tmpl,
		aiHandler: api.NewAIHandler(cfg),
	}
}

func (s *Server) Start() {
	fs := http.FileServer(http.Dir("frontend/dist"))
	http.Handle("/assets/", fs)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/dist/favicon.ico")
	})

	http.HandleFunc("/api/fetch", s.handleFetch)
	http.HandleFunc("/api/clean", s.handleClean)
	http.HandleFunc("/api/news", s.handleAPINews)
	http.HandleFunc("/api/resources", s.handleAPIResourcesJSON)
	http.HandleFunc("/api/ai/models", s.aiHandler.ServeHTTP)
	http.HandleFunc("/api/ai/chat", s.aiHandler.ServeHTTP)
	http.HandleFunc("/api/stats", s.handleStats)
	http.HandleFunc("/api/health", s.handleHealth)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/dist/index.html")
	})

	addr := ":" + strconv.Itoa(s.config.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "AI 资讯监控",
	}
	s.templates.ExecuteTemplate(w, "index.html", data)
}

func (s *Server) handleNews(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize := 20
	source := r.URL.Query().Get("source")
	search := r.URL.Query().Get("search")

	news, total, err := s.db.GetNews(page, pageSize, source, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	data := map[string]interface{}{
		"News":       news,
		"Page":       page,
		"TotalPages": totalPages,
		"Total":      total,
		"Source":     source,
		"Search":     search,
	}
	s.templates.ExecuteTemplate(w, "news.html", data)
}

func (s *Server) handleAPIResources(w http.ResponseWriter, r *http.Request) {
	resources, err := s.db.GetAPIResources()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Resources": resources,
	}
	s.templates.ExecuteTemplate(w, "api_resources.html", data)
}

func (s *Server) handleFetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	go s.fetcher.FetchAll()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok", "message": "Fetch started"}`))
}

func (s *Server) handleClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := s.db.CleanOldNews(s.config.News.RetentionDays); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok", "message": "Clean completed"}`))
}

func (s *Server) handleAPINews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize := 20
	source := r.URL.Query().Get("source")
	search := r.URL.Query().Get("search")

	news, total, err := s.db.GetNews(page, pageSize, source, search)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	type NewsJSON struct {
		ID        int64  `json:"id"`
		Title     string `json:"title"`
		URL       string `json:"url"`
		Summary   string `json:"summary"`
		Source    string `json:"source"`
		CreatedAt string `json:"created_at"`
	}

	var items []NewsJSON
	for _, n := range news {
		items = append(items, NewsJSON{
			ID: n.ID, Title: n.Title, URL: n.URL, Summary: n.Summary,
			Source: n.Source, CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp := fmt.Sprintf(`{"news":%s,"total":%d,"page":%d,"totalPages":%d}`, marshalJSON(items), total, page, totalPages)
	w.Write([]byte(resp))
}

func (s *Server) handleAPIResourcesJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resources, err := s.db.GetAPIResources()
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	type ResJSON struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Provider    string `json:"provider"`
		Description string `json:"description"`
		Endpoint    string `json:"endpoint"`
		FreeQuota   string `json:"free_quota"`
		DocURL      string `json:"doc_url"`
		LastUpdated string `json:"last_updated"`
	}

	var items []ResJSON
	for _, r := range resources {
		items = append(items, ResJSON{
			ID: r.ID, Name: r.Name, Provider: r.Provider, Description: r.Description,
			Endpoint: r.Endpoint, FreeQuota: r.FreeQuota, DocURL: r.DocURL,
			LastUpdated: r.LastUpdated.Format("2006-01-02 15:04:05"),
		})
	}

	resp := fmt.Sprintf(`{"resources":%s}`, marshalJSON(items))
	w.Write([]byte(resp))
}

func marshalJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok","service":"ai-watcher"}`))
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats, err := s.db.GetStats()
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	latestNews, _ := s.db.GetLatestNews(3)
	var latestItems []map[string]interface{}
	for _, n := range latestNews {
		latestItems = append(latestItems, map[string]interface{}{
			"id":         n.ID,
			"title":      n.Title,
			"url":        n.URL,
			"summary":    n.Summary,
			"source":     n.Source,
			"created_at": n.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resources, _ := s.db.GetAPIResources()

	resp := map[string]interface{}{
		"model_count":    api.GetModelCount(),
		"news_count":     stats.NewsCount,
		"resource_count": len(resources),
		"last_fetch":     stats.LastFetch,
		"source_counts":  stats.SourceCounts,
		"latest_news":    latestItems,
	}

	json.NewEncoder(w).Encode(resp)
}
