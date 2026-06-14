# AI 接入平台

提供统一的 AI 能力接入、资讯监控与对话管理服务。

## 技术栈

- **前端**：Vue 3 + TypeScript + Vite
- **后端**：Go 1.21+
- **数据库**：MariaDB / MySQL
- **包管理**：pnpm（前端）

## 功能模块

### 1. AI 资讯监控

定时从 RSS 源、微信公众号、知乎、V2EX 等渠道采集 AI 行业新闻，支持搜索、按来源过滤、分页浏览。

### 2. AI 对话

接入阿里云百炼平台，支持通义千问、DeepSeek、Kimi、智谱 GLM、MiniMax 等 11 个分类 49+ 款模型。

### 3. 免费 API 资源

收录各厂商免费大模型 API 资源信息，包括额度、端点和文档链接。

## 支持的模型家族

| 厂商 | 模型 |
|------|------|
| 通义千问 | qwen3.7 / qwen3.6 / qwen3.5 / qwen3 全系列，plus/max/flash 规格 |
| 通义经典 | qwen-plus、qwen-turbo、qwen-max、qwen-flash 等 |
| 通义代码 | qwen3-coder-next、qwen3-coder-plus、qwen3-coder-flash |
| DeepSeek | v4-flash、v4-pro、v3.2、v3.1、v3、r1 系列 |
| Kimi 月之暗面 | k2.6、k2.5、k2-thinking |
| GLM 智谱 | glm-5、glm-5.1、glm-4.7 |
| MiniMax | M2.5、M2.1 |

## 快速开始

### 1. 安装依赖

```bash
# 前端依赖
pnpm install
```

### 2. 配置数据库

确保 MariaDB / MySQL 已运行，执行初始化脚本：

```bash
mysql -u root < migrations/001_init.sql
```

### 3. 编译 Go 后端

```bash
go build -o ai-watcher ./cmd/main.go
```

### 4. 配置

复制并编辑配置文件：

```bash
cp config.example.json config.json
```

编辑 `config.json`，重点配置以下字段：

```json
{
  "database": {
    "host": "localhost",
    "port": 3306,
    "user": "root",
    "password": "",
    "dbname": "ai_watcher"
  },
  "ai": {
    "api_key": "你的 DashScope API Key",
    "base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1"
  }
}
```

### 5. 启动服务

```bash
# 启动 Go 后端（8080 端口，含新闻采集服务）
./ai-watcher config.json

# 启动前端（5173 端口）
cd frontend && pnpm dev
```

访问 [http://localhost:5173](http://localhost:5173) 查看站点。

## 项目结构

```
├── cmd/
│   └── main.go            # Go 程序入口
├── internal/
│   ├── api/                # 数据库操作与 API 处理器
│   │   ├── database.go     # 数据库连接与查询
│   │   └── ai_handler.go   # AI 对话接口
│   ├── config/             # 配置加载
│   │   └── config.go
│   ├── models/             # 数据模型
│   ├── services/           # 业务服务
│   │   ├── rss_fetcher.go      # RSS 采集
│   │   ├── wechat_fetcher.go   # 微信公众号采集
│   │   ├── community_fetcher.go # 社区采集（知乎、V2EX）
│   │   └── summary_service.go  # AI 摘要生成
│   └── web/
│       └── server.go        # Web 服务器与路由
├── frontend/                # Vue 3 前端
│   └── src/
│       ├── views/
│       │   ├── HomeView.vue       # 首页
│       │   ├── NewsView.vue       # 资讯列表
│       │   ├── ResourcesView.vue  # API 资源清单
│       │   └── ChatView.vue       # AI 对话
│       └── router/
│           └── index.ts
├── migrations/              # 数据库迁移
│   └── 001_init.sql         # 建表脚本
├── web/templates/            # Go HTML 模板
├── config.example.json       # 配置示例
├── go.mod
└── package.json              # 前端 workspace 配置
```

## API 接口

### 健康检查

`GET /api/health`

响应：
```json
{ "status": "ok", "service": "ai-watcher" }
```

### 模型列表

`GET /api/ai/models`

返回所有支持的模型分类及模型名称。

### AI 对话

`POST /api/ai/chat`

请求体：
```json
{
  "messages": [
    { "role": "user", "content": "你好" }
  ],
  "model": "qwen-turbo"
}
```

响应：
```json
{
  "choices": [{
    "message": { "role": "assistant", "content": "你好！有什么可以帮助你的？" },
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 15,
    "total_tokens": 25
  }
}
```

### 资讯列表

`GET /api/news?page=1&source=&search=`

参数：
- `page` - 页码，默认 1
- `source` - 按来源过滤（机器之心、量子位等）
- `search` - 按标题/摘要搜索关键词

响应：
```json
{
  "news": [...],
  "total": 100,
  "page": 1,
  "totalPages": 5
}
```

### API 资源清单

`GET /api/resources`

返回所有免费的 AI 大模型 API 资源信息。

## 配置说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `server.port` | 后端服务端口 | 8080 |
| `database.host` | 数据库地址 | localhost |
| `database.port` | 数据库端口 | 3306 |
| `news.fetch_interval` | 新闻采集间隔（分钟） | 30 |
| `news.retention_days` | 新闻保留天数 | 30 |
| `news.wechat_cookie` | 微信公众号 Cookie（可选） | 空 |
| `ai.api_key` | DashScope API Key（必填） | 空 |
| `ai.base_url` | DashScope API 地址 | dashscope.aliyuncs.com |
