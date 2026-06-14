#!/bin/bash
set -e

echo "=== AI Watcher 启动脚本 ==="

echo "[1/4] 初始化数据库..."
sudo mysql <<'SQL'
CREATE DATABASE IF NOT EXISTS ai_watcher CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE ai_watcher;

CREATE TABLE IF NOT EXISTS news (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    url VARCHAR(2000) NOT NULL,
    summary VARCHAR(200) DEFAULT '',
    source VARCHAR(100) NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY idx_url (url(191)),
    INDEX idx_source (source),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS api_resources (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    description TEXT,
    endpoint VARCHAR(500),
    free_quota VARCHAR(200),
    doc_url VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_provider (provider),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS fetch_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    source VARCHAR(100) NOT NULL,
    status ENUM('success', 'error', 'warning') NOT NULL,
    message TEXT,
    count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_source (source),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO api_resources (name, provider, description, endpoint, free_quota, doc_url) VALUES
('通义千问 Qwen-Turbo', '阿里云', '阿里云通义千问系列轻量级模型，适合摘要生成', 'https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation', '100万tokens/月免费', 'https://help.aliyun.com/zh/dashscope/developer-reference/api-details'),
('DeepSeek V3', 'DeepSeek', 'DeepSeek 大语言模型，支持多种任务', 'https://api.deepseek.com/v1/chat/completions', '500万tokens/月免费', 'https://platform.deepseek.com/api-docs'),
('智谱 GLM-4-Flash', '智谱AI', '智谱AI GLM-4 系列轻量模型', 'https://open.bigmodel.cn/api/paas/v4/chat/completions', '100万tokens/月免费', 'https://bigmodel.cn/dev/api'),
('百度文心一言 ERNIE-Speed', '百度智能云', '百度文心一言轻量版', 'https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie-speed-128k', '每天免费10000次', 'https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Hlr74s9x2'),
('讯飞星火 Spark Lite', '科大讯飞', '科大讯飞星火认知大模型轻量版', 'https://spark-api-open.xf-yun.com/v1/chat/completions', '每天免费5000次', 'https://www.xfyun.cn/doc/spark/HTTP%E8%B0%83%E7%94%A8%E6%96%87%E6%A1%A3.html'),
('字节豆包 Doubao-lite', '字节跳动', '字节跳动豆包系列轻量模型', 'https://ark.cn-beijing.volces.com/api/v3/chat/completions', '50万tokens/月免费', 'https://www.volcengine.com/docs/82379/1298454');
SQL
echo "   数据库初始化完成 ✓"

echo "[2/4] 配置文件已就绪"
cd /projects/mimo/ai-watcher

echo "[3/4] 启动 AI 资讯监控服务..."
echo "   访问地址: http://localhost:8080"
echo ""
./ai-watcher
