import { Router } from 'express'

const router = Router()

const MOCK_RESOURCES = [
  {
    id: 1,
    name: '通义千问 Qwen-Turbo',
    provider: '阿里云',
    description: '阿里云通义千问系列轻量级模型，适合摘要生成',
    endpoint: 'https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation',
    free_quota: '100万tokens/月免费',
    doc_url: 'https://help.aliyun.com/zh/dashscope/developer-reference/api-details',
    last_updated: '2026-06-14T00:00:00Z',
  },
  {
    id: 2,
    name: 'DeepSeek V3',
    provider: 'DeepSeek',
    description: 'DeepSeek 大语言模型，支持多种任务',
    endpoint: 'https://api.deepseek.com/v1/chat/completions',
    free_quota: '500万tokens/月免费',
    doc_url: 'https://platform.deepseek.com/api-docs',
    last_updated: '2026-06-13T00:00:00Z',
  },
  {
    id: 3,
    name: '智谱 GLM-4-Flash',
    provider: '智谱AI',
    description: '智谱AI GLM-4 系列轻量模型',
    endpoint: 'https://open.bigmodel.cn/api/paas/v4/chat/completions',
    free_quota: '100万tokens/月免费',
    doc_url: 'https://bigmodel.cn/dev/api',
    last_updated: '2026-06-12T00:00:00Z',
  },
  {
    id: 4,
    name: '百度文心一言 ERNIE-Speed',
    provider: '百度智能云',
    description: '百度文心一言轻量版',
    endpoint: 'https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie-speed-128k',
    free_quota: '每天免费10000次',
    doc_url: 'https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Hlr74s9x2',
    last_updated: '2026-06-11T00:00:00Z',
  },
  {
    id: 5,
    name: '讯飞星火 Spark Lite',
    provider: '科大讯飞',
    description: '科大讯飞星火认知大模型轻量版',
    endpoint: 'https://spark-api-open.xf-yun.com/v1/chat/completions',
    free_quota: '每天免费5000次',
    doc_url: 'https://www.xfyun.cn/doc/spark/HTTP%E8%B0%83%E7%94%A8%E6%96%87%E6%A1%A3.html',
    last_updated: '2026-06-10T00:00:00Z',
  },
  {
    id: 6,
    name: '字节豆包 Doubao-lite',
    provider: '字节跳动',
    description: '字节跳动豆包系列轻量模型',
    endpoint: 'https://ark.cn-beijing.volces.com/api/v3/chat/completions',
    free_quota: '50万tokens/月免费',
    doc_url: 'https://www.volcengine.com/docs/82379/1298454',
    last_updated: '2026-06-09T00:00:00Z',
  },
]

router.get('/resources', (_req, res) => {
  res.json({
    resources: MOCK_RESOURCES,
  })
})

export { router as resourcesRouter }
