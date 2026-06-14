package main

import (
	"log"
	"os"
	"time"

	"ai-watcher/internal/api"
	"ai-watcher/internal/config"
	"ai-watcher/internal/models"
	"ai-watcher/internal/services"
	"ai-watcher/internal/web"
)

func main() {
	configPath := "config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("Failed to load config, using defaults: %v", err)
		cfg = config.DefaultConfig()
	}

	db, err := api.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	loadSampleData(db)
	loadDefaultResources(db)

	rssFetcher := services.NewRSSFetcher(db, cfg)
	wechatFetcher := services.NewWechatFetcher(db, cfg)
	communityFetcher := services.NewCommunityFetcher(db, cfg)
	summaryService := services.NewSummaryService(db, cfg)

	go func() {
		for {
			log.Println("Starting news fetch...")
			rssFetcher.FetchAll()
			wechatFetcher.FetchAll()
			communityFetcher.FetchAll()

			time.Sleep(time.Duration(cfg.News.FetchInterval) * time.Minute)
		}
	}()

	go func() {
		for {
			time.Sleep(1 * time.Hour)
			log.Println("Starting summarization...")
			summaryService.SummarizeAll()
		}
	}()

	go func() {
		for {
			time.Sleep(24 * time.Hour)
			log.Println("Cleaning old news...")
			db.CleanOldNews(cfg.News.RetentionDays)
		}
	}()

	server := web.NewServer(db, cfg)
	server.Start()
}

func loadSampleData(db *api.Database) {
	samples := []models.News{
		{Title: "OpenAI 发布 GPT-5：推理能力大幅提升", URL: "https://example.com/gpt5", Summary: "OpenAI 正式发布 GPT-5，在数学推理和代码生成方面有重大突破。", Source: "机器之心"},
		{Title: "DeepSeek-V3 开源：性能超越 Llama 3", URL: "https://example.com/deepseek-v3", Summary: "DeepSeek 开源 V3 模型，多项基准测试成绩超过 Llama 3。", Source: "量子位"},
		{Title: "百度文心一言 4.5 发布：多模态能力升级", URL: "https://example.com/ernie45", Summary: "百度发布文心一言 4.5 版本，图像理解和生成能力显著提升。", Source: "机器之心"},
		{Title: "阿里通义千问 Qwen3：支持 128K 上下文", URL: "https://example.com/qwen3", Summary: "阿里发布 Qwen3 系列模型，支持 128K 上下文窗口。", Source: "量子位"},
		{Title: "智谱 GLM-4-Plus 发布：推理速度提升 3 倍", URL: "https://example.com/glm4plus", Summary: "智谱 AI 发布 GLM-4-Plus，推理速度提升 3 倍，成本降低 50%。", Source: "知乎"},
		{Title: "字节豆包大模型 API 限时免费开放", URL: "https://example.com/doubao-free", Summary: "字节跳动宣布豆包大模型 API 限时免费开放，吸引开发者。", Source: "V2EX"},
		{Title: "Meta Llama 4 预训练数据集规模达 15T tokens", URL: "https://example.com/llama4", Summary: "Meta 透露 Llama 4 预训练数据集规模达到 15T tokens。", Source: "机器之心"},
		{Title: "国产 AI 芯片算力突破 1000 PFLOPS", URL: "https://example.com/chip", Summary: "国产 AI 芯片单集群算力突破 1000 PFLOPS，打破国外垄断。", Source: "量子位"},
		{Title: "零一万物 Yi-Lightning 模型评测成绩亮眼", URL: "https://example.com/yi-lightning", Summary: "零一万物发布 Yi-Lightning 模型，在多项评测中表现优异。", Source: "微信公众号"},
		{Title: "AI 编程助手 Cursor 用户突破 1000 万", URL: "https://example.com/cursor", Summary: "AI 编程工具 Cursor 宣布用户数突破 1000 万。", Source: "知乎"},
	}

	for i := range samples {
		db.InsertNews(&samples[i])
	}
	log.Printf("Loaded %d sample news articles", len(samples))
}

func loadDefaultResources(db *api.Database) {
	var count int
	row := db.Query("SELECT COUNT(*) FROM api_resources")
	row.Scan(&count)
	if count > 0 {
		return
	}

	resources := []models.APIResource{
		{Name: "通义千问 Qwen-Turbo", Provider: "阿里云", Description: "阿里云通义千问系列轻量级模型，适合摘要生成", Endpoint: "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation", FreeQuota: "100万tokens/月免费", DocURL: "https://help.aliyun.com/zh/dashscope/developer-reference/api-details", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
		{Name: "DeepSeek V3", Provider: "DeepSeek", Description: "DeepSeek 大语言模型，支持多种任务", Endpoint: "https://api.deepseek.com/v1/chat/completions", FreeQuota: "500万tokens/月免费", DocURL: "https://platform.deepseek.com/api-docs", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
		{Name: "智谱 GLM-4-Flash", Provider: "智谱AI", Description: "智谱AI GLM-4 系列轻量模型", Endpoint: "https://open.bigmodel.cn/api/paas/v4/chat/completions", FreeQuota: "100万tokens/月免费", DocURL: "https://bigmodel.cn/dev/api", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
		{Name: "百度文心一言 ERNIE-Speed", Provider: "百度智能云", Description: "百度文心一言轻量版", Endpoint: "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie-speed-128k", FreeQuota: "每天免费10000次", DocURL: "https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Hlr74s9x2", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
		{Name: "讯飞星火 Spark Lite", Provider: "科大讯飞", Description: "科大讯飞星火认知大模型轻量版", Endpoint: "https://spark-api-open.xf-yun.com/v1/chat/completions", FreeQuota: "每天免费5000次", DocURL: "https://www.xfyun.cn/doc/spark/HTTP%E8%B0%83%E7%94%A8%E6%96%87%E6%A1%A3.html", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
		{Name: "字节豆包 Doubao-lite", Provider: "字节跳动", Description: "字节跳动豆包系列轻量模型", Endpoint: "https://ark.cn-beijing.volces.com/api/v3/chat/completions", FreeQuota: "50万tokens/月免费", DocURL: "https://www.volcengine.com/docs/82379/1298454", IsActive: true, LastUpdated: time.Now(), CreatedAt: time.Now()},
	}

	for i := range resources {
		db.InsertAPIResource(&resources[i])
	}
	log.Printf("Loaded %d default API resources", len(resources))
}
