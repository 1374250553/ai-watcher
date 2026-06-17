<template>
  <div class="api-view">
    <div class="api-header">
      <h2>免费 API 资源清单</h2>
      <p class="subtitle">当前可用的免费大模型 API 资源</p>
    </div>

    <div class="resource-grid">
      <div v-for="r in resources" :key="r.id" class="resource-card">
        <div class="resource-header">
          <div>
            <h3>{{ r.name }}</h3>
            <span class="provider">{{ r.provider }}</span>
          </div>
          <span :class="['status-dot', healthStatus[r.id] || 'unknown']" :title="healthStatus[r.id]"></span>
        </div>
        <p class="description">{{ r.description }}</p>
        <div class="quota">
          <span class="label">免费额度：</span>
          <span class="value">{{ r.free_quota }}</span>
        </div>
        <div class="endpoint">
          <span class="label">端点：</span>
          <code>{{ r.endpoint }}</code>
        </div>

        <div class="resource-actions">
          <a v-if="r.doc_url" :href="r.doc_url" target="_blank" class="doc-link">查看文档</a>
          <button
            v-if="r.model"
            class="trial-btn"
            @click="goChat(r.model)"
          >
            一键试用
          </button>
          <button
            class="code-toggle-btn"
            @click="toggleCode(r.id)"
          >
            {{ codeExpanded[r.id] ? '收起示例' : '代码示例' }}
          </button>
        </div>

        <div v-if="codeExpanded[r.id]" class="code-examples">
          <div class="code-block">
            <div class="code-header">
              <span>curl</span>
              <button class="copy-btn" @click="copyCode(r.id, 'curl')">复制</button>
            </div>
            <pre><code>{{ curlExample(r) }}</code></pre>
          </div>
          <div class="code-block">
            <div class="code-header">
              <span>Python</span>
              <button class="copy-btn" @click="copyCode(r.id, 'python')">复制</button>
            </div>
            <pre><code>{{ pythonExample(r) }}</code></pre>
          </div>
        </div>

        <div class="updated">上次更新：{{ formatTime(r.last_updated) }}</div>
      </div>
      <div v-if="!resources.length && !loading" class="empty">暂无 API 资源</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

interface APIResource {
  id: number
  name: string
  provider: string
  description: string
  endpoint: string
  free_quota: string
  doc_url: string
  model: string
  last_updated: string
}

const resources = ref<APIResource[]>([])
const loading = ref(true)
const healthStatus = ref<Record<number, string>>({})
const codeExpanded = ref<Record<number, boolean>>({})

async function loadResources() {
  loading.value = true
  try {
    const res = await fetch('/api/resources')
    const data = await res.json()
    resources.value = data.resources || []
    checkAllHealth()
  } catch (e) {
    console.error('Failed to load API resources:', e)
  } finally {
    loading.value = false
  }
}

async function checkAllHealth() {
  for (const r of resources.value) {
    try {
      const res = await fetch(`/api/resources/health?id=${r.id}`)
      const data = await res.json()
      healthStatus.value[r.id] = data.status || 'unknown'
    } catch {
      healthStatus.value[r.id] = 'unknown'
    }
  }
}

function goChat(model: string) {
  router.push(`/chat?model=${model}`)
}

function toggleCode(id: number) {
  codeExpanded.value[id] = !codeExpanded.value[id]
}

function curlExample(r: APIResource): string {
  const model = r.model || 'MODEL_NAME'
  const authHeader = r.provider === '阿里云' ? 'Authorization: Bearer $DASHSCOPE_API_KEY' : 'Authorization: Bearer $API_KEY'
  return `curl -X POST "${r.endpoint}" \\
  -H "Content-Type: application/json" \\
  -H "${authHeader}" \\
  -d '{
    "model": "${model}",
    "messages": [{"role":"user","content":"你好"}]
  }'`
}

function pythonExample(r: APIResource): string {
  const model = r.model || 'MODEL_NAME'
  return `import requests

response = requests.post(
    "${r.endpoint}",
    headers={
        "Content-Type": "application/json",
        "Authorization": "Bearer $API_KEY"
    },
    json={
        "model": "${model}",
        "messages": [{"role": "user", "content": "你好"}]
    }
)
print(response.json())`
}

function copyCode(id: number, lang: string) {
  const r = resources.value.find(x => x.id === id)
  if (!r) return
  const code = lang === 'curl' ? curlExample(r) : pythonExample(r)
  navigator.clipboard.writeText(code)
}

function formatTime(t: string) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

onMounted(loadResources)
</script>

<style scoped>
.api-header {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  margin-bottom: 1rem;
}

.api-header h2 { color: #333; }
.subtitle { color: #999; margin-top: 0.25rem; }

.resource-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  gap: 1rem;
}

.resource-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.resource-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.resource-header h3 { color: #1a73e8; font-size: 1.1rem; }

.provider {
  display: inline-block;
  background: #e8f0fe;
  color: #1a73e8;
  padding: 0.125rem 0.5rem;
  border-radius: 3px;
  font-size: 0.7rem;
  margin-top: 0.25rem;
}

.status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 0.35rem;
}

.status-dot.online { background: #34a853; }
.status-dot.offline { background: #d93025; }
.status-dot.unknown { background: #fbbc04; }

.description { color: #666; margin-bottom: 0.75rem; line-height: 1.5; }

.quota, .endpoint {
  font-size: 0.875rem;
  margin-bottom: 0.4rem;
}

.label { color: #999; }

.quota .value { color: #0d904f; font-weight: 500; }

.endpoint code {
  background: #f5f5f5;
  padding: 0.125rem 0.375rem;
  border-radius: 3px;
  font-size: 0.78rem;
  word-break: break-all;
}

.resource-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.75rem;
  flex-wrap: wrap;
}

.doc-link {
  display: inline-block;
  padding: 0.35rem 0.75rem;
  color: #1a73e8;
  text-decoration: none;
  font-size: 0.85rem;
  border: 1px solid #1a73e8;
  border-radius: 4px;
  transition: background 0.2s;
}

.doc-link:hover { background: #e8f0fe; }

.trial-btn {
  padding: 0.35rem 0.75rem;
  background: #34a853;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: background 0.2s;
}

.trial-btn:hover { background: #2d9249; }

.code-toggle-btn {
  padding: 0.35rem 0.75rem;
  background: #f0f0f0;
  color: #666;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.2s;
}

.code-toggle-btn:hover {
  background: #e0e0e0;
  color: #333;
}

.code-examples {
  margin-top: 0.75rem;
}

.code-block {
  margin-bottom: 0.75rem;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid #e0e0e0;
}

.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.35rem 0.75rem;
  background: #f5f5f5;
  border-bottom: 1px solid #e0e0e0;
  font-size: 0.75rem;
  font-weight: 500;
  color: #666;
}

.copy-btn {
  background: transparent;
  border: 1px solid #ccc;
  border-radius: 3px;
  padding: 0.15rem 0.5rem;
  cursor: pointer;
  font-size: 0.7rem;
  color: #666;
}

.copy-btn:hover { background: #e0e0e0; }

.code-block pre {
  margin: 0;
  padding: 0.75rem;
  overflow-x: auto;
  background: #1e1e1e;
  color: #e0e0e0;
  font-size: 0.8rem;
  line-height: 1.5;
}

.code-block code {
  font-family: 'Menlo', 'Monaco', 'Consolas', monospace;
}

.updated {
  margin-top: 0.75rem;
  font-size: 0.75rem;
  color: #bbb;
}

.empty {
  grid-column: 1 / -1;
  text-align: center;
  color: #999;
  padding: 3rem;
}

@media (max-width: 768px) {
  .resource-grid {
    grid-template-columns: 1fr;
  }
}
</style>
