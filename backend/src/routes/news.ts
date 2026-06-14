import { Router } from 'express'

const router = Router()

const MOCK_NEWS = [
  {
    id: 1,
    title: '阿里通义千问 Qwen3.7 发布：性能全面超越前代',
    url: 'https://example.com/qwen3.7-release',
    summary: '阿里云发布通义千问 Qwen3.7 系列模型，在推理、代码生成等任务上表现显著提升。',
    source: '机器之心',
    created_at: '2026-06-14T08:00:00Z',
  },
  {
    id: 2,
    title: 'DeepSeek V4 模型正式开源，支持超长上下文窗口',
    url: 'https://example.com/deepseek-v4',
    summary: 'DeepSeek 宣布 V4 模型正式开源，支持 256K 上下文窗口和多轮对话优化。',
    source: '量子位',
    created_at: '2026-06-13T12:00:00Z',
  },
  {
    id: 3,
    title: '智谱 AI GLM-5 发布：多模态能力大幅增强',
    url: 'https://example.com/glm5-release',
    summary: '智谱 AI 发布新一代 GLM-5 模型，在图像理解和生成方面表现突出。',
    source: '机器之心',
    created_at: '2026-06-12T15:30:00Z',
  },
  {
    id: 4,
    title: 'OpenAI 发布 GPT-5 技术报告：模型架构全解析',
    url: 'https://example.com/gpt5-report',
    summary: 'OpenAI 发布 GPT-5 技术报告，详细阐述了新模型的架构设计和训练策略。',
    source: '知乎',
    created_at: '2026-06-11T09:00:00Z',
  },
  {
    id: 5,
    title: 'Anthropic Claude 4 正式商用：更强的推理和编程能力',
    url: 'https://example.com/claude4',
    summary: 'Anthropic 宣布 Claude 4 正式商用，在代码生成和数学推理上表现优异。',
    source: 'V2EX',
    created_at: '2026-06-10T11:45:00Z',
  },
  {
    id: 6,
    title: 'Kimi K2.6 升级：支持实时联网搜索和文件解析',
    url: 'https://example.com/kimi-k26',
    summary: '月之暗面 Kimi 升级至 K2.6 版本，新增实时联网搜索和复杂文件解析能力。',
    source: '微信公众号',
    created_at: '2026-06-09T16:20:00Z',
  },
  {
    id: 7,
    title: 'MiniMax M2.5 发布：角色扮演和语音交互能力升级',
    url: 'https://example.com/minimax-m25',
    summary: 'MiniMax 发布 M2.5 模型，在角色扮演、语音交互和多轮对话上表现突出。',
    source: '机器之心',
    created_at: '2026-06-08T10:00:00Z',
  },
  {
    id: 8,
    title: '国内 AI 大模型市场竞争格局分析',
    url: 'https://example.com/ai-market-analysis',
    summary: '2026 年国内 AI 大模型市场呈现百花齐放态势，多家厂商竞相发布新产品。',
    source: '量子位',
    created_at: '2026-06-07T14:30:00Z',
  },
]

router.get('/news', (req, res) => {
  const page = parseInt(req.query.page as string) || 1
  const pageSize = 10
  const source = req.query.source as string
  const search = req.query.search as string

  let filtered = MOCK_NEWS

  if (source) {
    filtered = filtered.filter((item) => item.source === source)
  }

  if (search) {
    const keyword = search.toLowerCase()
    filtered = filtered.filter(
      (item) =>
        item.title.toLowerCase().includes(keyword) ||
        item.summary.toLowerCase().includes(keyword)
    )
  }

  const total = filtered.length
  const totalPages = Math.max(1, Math.ceil(total / pageSize))
  const start = (page - 1) * pageSize
  const news = filtered.slice(start, start + pageSize)

  res.json({
    news,
    total,
    totalPages,
    page,
  })
})

export { router as newsRouter }
