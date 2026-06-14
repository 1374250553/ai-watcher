import { Router } from 'express'

const router = Router()

const DASHSCOPE_API_KEY = process.env.DASHSCOPE_API_KEY
const DASHSCOPE_BASE_URL = 'https://dashscope.aliyuncs.com/compatible-mode/v1'

if (!DASHSCOPE_API_KEY) {
  console.warn('Warning: DASHSCOPE_API_KEY environment variable is not set')
}

const TEXT_CHAT_MODELS: Record<string, string[]> = {
  '通义千问 Qwen3.7': ['qwen3.7-max', 'qwen3.7-max-preview', 'qwen3.7-plus'],
  '通义千问 Qwen3.6': ['qwen3.6-max-preview', 'qwen3.6-plus', 'qwen3.6-flash'],
  '通义千问 Qwen3.5': ['qwen3.5-plus', 'qwen3.5-flash', 'qwen3.5-397b-a17b', 'qwen3.5-122b-a10b', 'qwen3.5-35b-a3b', 'qwen3.5-27b'],
  '通义千问 Qwen3': ['qwen3-max', 'qwen3-max-preview', 'qwen3-plus', 'qwen3-235b-a22b', 'qwen3-32b', 'qwen3-30b-a3b', 'qwen3-14b', 'qwen3-8b'],
  '通义千问经典': ['qwen-plus', 'qwen-turbo', 'qwen-max', 'qwen-flash', 'qwen-long', 'qwq-plus'],
  '通义千问代码': ['qwen3-coder-next', 'qwen3-coder-plus', 'qwen3-coder-flash'],
  '通义千问思考': ['qwen3-next-80b-a3b-thinking', 'qwen3-235b-a22b-thinking-2507', 'qwen3-30b-a3b-thinking-2507'],
  'DeepSeek': ['deepseek-v4-flash', 'deepseek-v4-pro', 'deepseek-v3.2', 'deepseek-v3.1', 'deepseek-v3', 'deepseek-r1', 'deepseek-r1-distill-qwen-32b', 'deepseek-r1-distill-qwen-14b', 'deepseek-r1-distill-llama-70b'],
  'Kimi 月之暗面': ['kimi-k2.6', 'kimi-k2.5', 'kimi-k2-thinking'],
  'GLM 智谱': ['glm-5', 'glm-5.1', 'glm-4.7'],
  'MiniMax': ['MiniMax-M2.5', 'MiniMax-M2.1'],
}

router.get('/models', (_req, res) => {
  res.json({
    categories: TEXT_CHAT_MODELS,
    total: Object.values(TEXT_CHAT_MODELS).flat().length
  })
})

router.post('/chat', async (req, res) => {
  try {
    if (!DASHSCOPE_API_KEY) {
      return res.status(500).json({
        error: 'API key not configured',
        details: 'Please set DASHSCOPE_API_KEY environment variable'
      })
    }

    const { messages, model = 'qwen-turbo' } = req.body

    if (!messages || !Array.isArray(messages)) {
      return res.status(400).json({
        error: 'Invalid request',
        details: 'messages field is required and must be an array'
      })
    }

    const response = await fetch(`${DASHSCOPE_BASE_URL}/chat/completions`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${DASHSCOPE_API_KEY}`
      },
      body: JSON.stringify({
        model,
        messages
      })
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => null)
      return res.status(response.status).json({
        error: `DashScope API error: ${response.statusText}`,
        details: errorData
      })
    }

    const data = await response.json()
    res.json(data)
  } catch (error) {
    console.error('Chat endpoint error:', error)
    res.status(500).json({
      error: 'Internal server error',
      details: error instanceof Error ? error.message : String(error)
    })
  }
})

export { router as aiRouter }
