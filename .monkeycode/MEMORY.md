# 用户指令记忆

本文件记录了用户的指令、偏好和教导，用于在未来的交互中提供参考。

## 格式

### 用户指令条目
用户指令条目应遵循以下格式：

[用户指令摘要]
- Date: [YYYY-MM-DD]
- Context: [提及的场景或时间]
- Instructions:
  - [用户教导或指示的内容，逐行描述]

### 项目知识条目
Agent 在任务执行过程中发现的条目应遵循以下格式：

[项目知识摘要]
- Date: [YYYY-MM-DD]
- Context: Agent 在执行 [具体任务描述] 时发现
- Category: [运维部署|构建方法|测试方法|排错调试|工作流协作|环境配置]
- Instructions:
  - [具体的知识点，逐行描述]

## 去重策略
- 添加新条目前，检查是否存在相似或相同的指令
- 若发现重复，跳过新条目或与已有条目合并
- 合并时，更新上下文或日期信息
- 这有助于避免冗余条目，保持记忆文件整洁

## 条目

### 启动 MariaDB 服务
- Date: 2026-06-14
- Context: 项目需要数据库时
- Category: 运维部署
- Instructions:
  - 当前环境需要手动启动 MariaDB 服务：`service mariadb start`
  - 数据库初始化命令：`mysql -u root < migrations/001_init.sql`

### Go 编译命令
- Date: 2026-06-14
- Context: 构建后端程序时
- Category: 构建方法
- Instructions:
  - 编译命令：`go build -o ai-watcher ./cmd/main.go`
  - 编译输出文件：`/workspace/ai-watcher`

### 前端 pnpm workspace 配置
- Date: 2026-06-14
- Context: 前端依赖安装时发现 workspace 配置缺失
- Category: 环境配置
- Instructions:
  - 如果 `pnpm install` 后 node_modules 仍缺失，检查根目录是否存在 `pnpm-workspace.yaml`
  - 缺失时需创建，内容：
    ```yaml
    packages:
      - 'frontend'
      - 'backend'
    ```

### Go 包循环依赖问题
- Date: 2026-06-14
- Context: Agent 在添加 AI 对话功能时发现
- Category: 排错调试
- Instructions:
  - Go 项目中 `internal/api` 和 `internal/services` 不能相互引用，会导致 import cycle not allowed 编译错误
  - 解决方法：将相关逻辑合并到一个包中，或引入第三个中间包解耦

### 前后端代理配置
- Date: 2026-06-14
- Context: Vite 配置反向代理时
- Category: 环境配置
- Instructions:
  - 前端 Vite 代理配置在 `frontend/vite.config.ts` 中
  - `/api` 路径代理到 `http://localhost:8080`（Go 后端）
  - 所有 API 请求都经过前端代理转发到 Go 后端

### 项目架构现状
- Date: 2026-06-14
- Context: Agent 部署过程中发现
- Category: 运维部署
- Instructions:
  - 项目已从 Express.js 后端完全迁移到 Go 后端
  - `backend/` 目录下的 Express.js 代码已废弃，所有 API 由 Go 提供
  - 启动顺序：1) MariaDB -> 2) Go 后端 (8080) -> 3) Vite 前端 (5173)
  - Go 后端配置文件：`config.json`
