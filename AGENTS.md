# Agents

## 项目概述

AI 接入平台 —— 提供统一的 AI 能力接入、资讯监控与对话管理服务。

## 技术栈

- **前端**：Vue 3 + TypeScript + Vite
- **后端**：Go 1.21+
- **数据库**：MariaDB / MySQL
- **包管理**：pnpm（前端）

## 项目架构

```
├── cmd/main.go              # Go 程序入口
├── internal/
│   ├── api/                  # 数据库操作与 API 处理器
│   │   ├── database.go       # 数据库连接与查询
│   │   └── ai_handler.go     # AI 对话接口
│   ├── config/               # 配置加载
│   ├── models/               # 数据模型
│   ├── services/             # 业务服务（采集、摘要）
│   └── web/                  # Web 服务器与路由
├── frontend/                 # Vue 3 前端
├── migrations/               # 数据库迁移脚本
├── web/templates/            # Go HTML 模板
├── config.json               # 运行时配置
├── config.example.json       # 配置示例
└── package.json              # 前端 workspace 配置
```

## 端口配置

| 服务 | 端口 | 说明 |
|------|------|------|
| Vite 前端 | 5173 | 开发服务器 |
| Go 后端 | 8080 | Web 服务 + API + 新闻采集 |
| MariaDB | 3306 | 数据库 |

## 快速启动

```bash
# 1. 启动数据库
service mariadb start

# 2. 初始化数据库（首次运行）
mysql -u root < migrations/001_init.sql

# 3. 编译并启动 Go 后端
go build -o ai-watcher ./cmd/main.go
./ai-watcher config.json

# 4. 启动前端
cd frontend && pnpm dev
```

## 开发规范

1. 后端使用 Go，注意避免 `internal/api` 与 `internal/services` 之间的循环依赖
2. 前端使用 Vue 3 Composition API + TypeScript
3. 修改代码后确保 `go build` 编译通过
4. 配置文件使用 `config.json`，敏感信息（API Key）不要提交

## 工作流

- 功能开发完成后先部署本地预览确认效果
- 需求相关自动调用 feature-design skill 生成文档
- 涉及项目文档生成时调用 project-wiki skill
