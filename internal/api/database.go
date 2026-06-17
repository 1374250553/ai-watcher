package api

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"ai-watcher/internal/config"
	"ai-watcher/internal/models"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	log.Printf("Connected to database %s at %s:%d", cfg.DBName, cfg.Host, cfg.Port)
	return &Database{db: db}, nil
}

func (db *Database) Close() error {
	return db.db.Close()
}

func (db *Database) Query(query string, args ...interface{}) *sql.Row {
	return db.db.QueryRow(query, args...)
}

func (db *Database) Exec(query string, args ...interface{}) error {
	_, err := db.db.Exec(query, args...)
	return err
}

func (db *Database) InsertNews(news *models.News) error {
	_, err := db.db.Exec(
		`INSERT IGNORE INTO news (title, url, summary, source, content, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		news.Title, news.URL, news.Summary, news.Source, news.Content,
		time.Now(), time.Now(),
	)
	return err
}

func (db *Database) NewsExists(url string) (bool, error) {
	var count int
	err := db.db.QueryRow("SELECT COUNT(*) FROM news WHERE url = ?", url).Scan(&count)
	return count > 0, err
}

func (db *Database) GetNews(page, pageSize int, source, search string) ([]models.News, int, error) {
	conditions := []string{}
	args := []interface{}{}

	if source != "" {
		conditions = append(conditions, "source = ?")
		args = append(args, source)
	}
	if search != "" {
		conditions = append(conditions, "(title LIKE ? OR summary LIKE ?)")
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	where := ""
	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM news" + where
	if err := db.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := `SELECT id, title, url, summary, source, content, created_at, updated_at
		FROM news` + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	rows, err := db.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []models.News
	for rows.Next() {
		var n models.News
		if err := rows.Scan(&n.ID, &n.Title, &n.URL, &n.Summary, &n.Source, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, n)
	}

	return result, total, rows.Err()
}

func (db *Database) GetAPIResources() ([]models.APIResource, error) {
	rows, err := db.db.Query(`SELECT id, name, provider, description, endpoint, free_quota, doc_url,
		COALESCE(model, ''), is_active, last_updated, created_at FROM api_resources ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.APIResource
	for rows.Next() {
		var r models.APIResource
		if err := rows.Scan(&r.ID, &r.Name, &r.Provider, &r.Description, &r.Endpoint,
			&r.FreeQuota, &r.DocURL, &r.Model, &r.IsActive, &r.LastUpdated, &r.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, rows.Err()
}

func (db *Database) InsertAPIResource(resource *models.APIResource) error {
	_, err := db.db.Exec(
		`INSERT INTO api_resources (name, provider, description, endpoint, free_quota, doc_url, model, is_active, last_updated, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		resource.Name, resource.Provider, resource.Description, resource.Endpoint,
		resource.FreeQuota, resource.DocURL, resource.Model, resource.IsActive, resource.LastUpdated, resource.CreatedAt,
	)
	return err
}

func (db *Database) LogFetch(source, status, message string, count int) error {
	_, err := db.db.Exec(
		`INSERT INTO fetch_logs (source, status, message, count, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		source, status, message, count, time.Now(),
	)
	return err
}

func (db *Database) CleanOldNews(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	_, err := db.db.Exec("DELETE FROM news WHERE created_at < ?", cutoff)
	return err
}

type Stats struct {
	NewsCount    int            `json:"news_count"`
	LastFetch    string         `json:"last_fetch"`
	SourceCounts map[string]int `json:"source_counts"`
}

func (db *Database) GetStats() (*Stats, error) {
	stats := &Stats{SourceCounts: make(map[string]int)}

	db.db.QueryRow("SELECT COUNT(*) FROM news").Scan(&stats.NewsCount)

	var lastFetch sql.NullTime
	db.db.QueryRow("SELECT MAX(created_at) FROM fetch_logs WHERE status = 'success'").Scan(&lastFetch)
	if lastFetch.Valid {
		stats.LastFetch = lastFetch.Time.Format("2006-01-02 15:04:05")
	}

	rows, err := db.db.Query("SELECT source, COUNT(*) as cnt FROM news GROUP BY source ORDER BY cnt DESC")
	if err != nil {
		return stats, nil
	}
	defer rows.Close()
	for rows.Next() {
		var source string
		var count int
		if err := rows.Scan(&source, &count); err != nil {
			continue
		}
		stats.SourceCounts[source] = count
	}

	return stats, nil
}

func (db *Database) GetLatestNews(limit int) ([]models.News, error) {
	rows, err := db.db.Query(
		"SELECT id, title, url, summary, source, content, created_at, updated_at FROM news ORDER BY created_at DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.News
	for rows.Next() {
		var n models.News
		if err := rows.Scan(&n.ID, &n.Title, &n.URL, &n.Summary, &n.Source, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, rows.Err()
}
