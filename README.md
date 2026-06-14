# AI 接入平台

提供统一的 AI 能力接入与管理服务。

## 技术栈

- **前端**：Vue 3 + TypeScript + Vite
- **后端**：Express.js + TypeScript
- **数据库**：PostgreSQL + Prisma（待启用）
- **包管理**：pnpm

## AI 服务

已接入阿里云百炼平台，覆盖 5 大厂商 40+ 款模型，每款模型赠送 200 万 Tokens（输入输出各 100 万），有效期 90 天。

**支持的模型家族：**

| 厂商 | 模型 |
|------|------|
| 通义千问 | qwen3.7 / qwen3.6 / qwen3.5 / qwen3 全系列，plus/max/flash 规格 |
| DeepSeek | v4-flash、v4-pro、v3.2、v3.1、v3、r1 |
| Kimi | k2.6、k2.5、k2-thinking |
| GLM 智谱 | glm-5、glm-5.1、glm-4.7 |
| MiniMax | M2.5、M2.1 |

## 快速开始

```bash
# 安装依赖
pnpm install

# 配置环境变量（复制并编辑）
cp backend/.env.example backend/.env

# 启动开发服务（前端 + 后端）
pnpm dev
```

## 项目结构

```
├── backend/          # 后端服务
│   ├── src/
│   │   ├── index.ts      # 服务入口
│   │   └── routes/
│   │       └── ai.ts     # AI 相关路由
│   └── .env              # 环境变量
├── frontend/         # 前端应用
│   └── src/
│       ├── views/
│       │   ├── HomeView.vue
│       │   └── ChatView.vue
│       └── router/
│           └── index.ts
└── package.json      # workspace 配置
```

## API 接口

### POST /api/ai/chat

发送对话消息到 AI 模型。

请求体：
```json
{
  "messages": [
    { "role": "user", "content": "你好" }
  ],
  "model": "qwen-turbo"
}
```

### GET /api/health

健康检查接口。
