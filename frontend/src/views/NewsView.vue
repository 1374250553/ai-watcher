<template>
  <div class="news-view">
    <div class="news-header">
      <h2>AI 资讯列表 <span class="count">(共 {{ total }} 条)</span></h2>
      <div class="controls">
        <select v-model="selectedSource" @change="loadNews(1)">
          <option value="">全部来源</option>
          <option v-for="s in sources" :key="s" :value="s">{{ s }}</option>
        </select>
        <div class="search-box">
          <input v-model="searchText" placeholder="搜索标题或摘要..." @keydown.enter="loadNews(1)" />
          <button @click="loadNews(1)">搜索</button>
        </div>
        <button class="fetch-btn" @click="triggerFetch" :disabled="fetching">
          {{ fetching ? '采集中...' : '立即采集' }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="skeleton-list">
      <div v-for="i in 3" :key="i" class="skeleton-card">
        <div class="skeleton-line skeleton-title"></div>
        <div class="skeleton-line skeleton-meta"></div>
        <div class="skeleton-line skeleton-summary"></div>
        <div class="skeleton-line skeleton-summary short"></div>
      </div>
    </div>

    <div v-else-if="!newsList.length" class="empty-state">
      <div class="empty-icon">-</div>
      <p class="empty-title">暂无资讯</p>
      <p class="empty-desc">还没有采集到任何 AI 资讯，点击下方按钮开始采集</p>
      <button class="empty-fetch-btn" @click="triggerFetch" :disabled="fetching">
        {{ fetching ? '采集中...' : '立即采集最新资讯' }}
      </button>
    </div>

    <div v-else class="news-grid">
      <div v-for="item in newsList" :key="item.id" class="news-card">
        <a :href="item.url" target="_blank" class="card-title">{{ item.title }}</a>
        <div class="card-meta">
          <span :class="['source-tag', sourceClass(item.source)]">{{ item.source }}</span>
          <span class="card-time">{{ relativeTime(item.created_at) }}</span>
        </div>
        <p class="card-summary">{{ truncate(item.summary, 120) }}</p>
        <div class="card-footer">
          <span class="footer-source">{{ item.source }}</span>
          <span class="footer-time">{{ formatTime(item.created_at) }}</span>
          <a :href="item.url" target="_blank" class="read-more">阅读原文</a>
        </div>
      </div>
    </div>

    <div v-if="totalPages > 1" class="pagination">
      <button :disabled="page <= 1" @click="loadNews(page - 1)">上一页</button>
      <span>{{ page }} / {{ totalPages }}</span>
      <button :disabled="page >= totalPages" @click="loadNews(page + 1)">下一页</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface NewsItem {
  id: number
  title: string
  url: string
  summary: string
  source: string
  created_at: string
}

const newsList = ref<NewsItem[]>([])
const total = ref(0)
const page = ref(1)
const totalPages = ref(1)
const selectedSource = ref('')
const searchText = ref('')
const fetching = ref(false)
const loading = ref(true)
const sources = ref<string[]>(['机器之心', '量子位', '微信公众号', '知乎', 'V2EX'])

const sourceColors: Record<string, string> = {
  '量子位': 'source-blue',
  '机器之心': 'source-green',
  '微信公众号': 'source-orange',
  '知乎': 'source-purple',
  'V2EX': 'source-teal',
}

function sourceClass(source: string) {
  return sourceColors[source] || 'source-default'
}

async function loadNews(p: number) {
  page.value = p
  loading.value = true
  const params = new URLSearchParams({ page: String(p) })
  if (selectedSource.value) params.set('source', selectedSource.value)
  if (searchText.value) params.set('search', searchText.value)

  try {
    const res = await fetch(`/api/news?${params}`)
    const data = await res.json()
    newsList.value = data.news || []
    total.value = data.total || 0
    totalPages.value = data.totalPages || 1
  } catch (e) {
    console.error('Failed to load news:', e)
  } finally {
    loading.value = false
  }
}

async function triggerFetch() {
  fetching.value = true
  try {
    await fetch('/api/fetch', { method: 'POST' })
    setTimeout(() => loadNews(1), 2000)
  } catch (e) {
    console.error('Failed to trigger fetch:', e)
  } finally {
    fetching.value = false
  }
}

function relativeTime(t: string): string {
  if (!t) return ''
  const now = Date.now()
  const then = new Date(t).getTime()
  const diff = Math.floor((now - then) / 1000)
  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)} 分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)} 小时前`
  if (diff < 172800) return '昨天'
  return formatTime(t)
}

function formatTime(t: string) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

function truncate(s: string, maxLen: number) {
  if (!s) return ''
  return s.length > maxLen ? s.slice(0, maxLen) + '...' : s
}

onMounted(() => loadNews(1))
</script>

<style scoped>
.news-header {
  background: white;
  padding: 1.25rem 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  margin-bottom: 1rem;
}

.news-header h2 { margin-bottom: 0.75rem; color: #333; }
.news-header .count { font-size: 0.875rem; color: #999; font-weight: normal; }

.controls {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  flex-wrap: wrap;
}

.controls select, .controls input {
  padding: 0.45rem 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.85rem;
}

.search-box {
  display: flex;
  gap: 0.4rem;
  flex: 1;
  min-width: 180px;
}

.search-box input { flex: 1; }

.search-box button {
  padding: 0.45rem 0.75rem;
  background: #1a73e8;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  white-space: nowrap;
}

.search-box button:hover { background: #1557b0; }

.fetch-btn {
  padding: 0.45rem 0.75rem;
  background: #34a853;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  white-space: nowrap;
  margin-left: auto;
}

.fetch-btn:hover:not(:disabled) { background: #2d9249; }
.fetch-btn:disabled { opacity: 0.6; cursor: not-allowed; }

.news-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1rem;
}

.news-card {
  background: white;
  padding: 1.25rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transition: transform 0.15s, box-shadow 0.15s;
  display: flex;
  flex-direction: column;
}

.news-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
}

.card-title {
  font-size: 1.05rem;
  font-weight: 600;
  color: #1a73e8;
  text-decoration: none;
  line-height: 1.4;
  margin-bottom: 0.5rem;
}

.card-title:hover { text-decoration: underline; }

.card-meta {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.source-tag {
  font-size: 0.7rem;
  padding: 0.15rem 0.45rem;
  border-radius: 3px;
  font-weight: 500;
}

.source-blue { background: #e8f0fe; color: #1a73e8; }
.source-green { background: #e6f4ea; color: #137333; }
.source-orange { background: #fef7e0; color: #e37400; }
.source-purple { background: #f3e8fd; color: #7c3aed; }
.source-teal { background: #e0f2f1; color: #00796b; }
.source-default { background: #f0f0f0; color: #666; }

.card-time {
  font-size: 0.75rem;
  color: #999;
}

.card-summary {
  font-size: 0.875rem;
  color: #666;
  line-height: 1.5;
  flex: 1;
  margin-bottom: 0.75rem;
}

.card-footer {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px solid #f0f0f0;
  font-size: 0.75rem;
  color: #999;
}

.footer-source { color: #888; }
.footer-time { flex: 1; }

.read-more {
  color: #1a73e8;
  text-decoration: none;
  font-weight: 500;
}

.read-more:hover { text-decoration: underline; }

.skeleton-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1rem;
}

.skeleton-card {
  background: white;
  padding: 1.25rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}

.skeleton-line {
  height: 14px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
  margin-bottom: 0.6rem;
}

.skeleton-title { width: 70%; height: 18px; }
.skeleton-meta { width: 35%; height: 12px; }
.skeleton-summary { width: 100%; height: 14px; }
.skeleton-summary.short { width: 60%; }

@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 1rem;
  background: #f0f0f0;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  color: #999;
}

.empty-title {
  font-size: 1.25rem;
  color: #333;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.empty-desc {
  color: #999;
  margin-bottom: 1.5rem;
}

.empty-fetch-btn {
  padding: 0.75rem 2rem;
  background: #34a853;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

.empty-fetch-btn:hover:not(:disabled) { background: #2d9249; }
.empty-fetch-btn:disabled { opacity: 0.6; cursor: not-allowed; }

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 1rem;
}

.pagination button {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
}

.pagination button:disabled { opacity: 0.5; cursor: not-allowed; }
.pagination button:hover:not(:disabled) { background: #f5f5f5; }

@media (max-width: 768px) {
  .news-grid, .skeleton-list {
    grid-template-columns: 1fr;
  }
  .controls {
    flex-direction: column;
  }
  .fetch-btn {
    margin-left: 0;
    width: 100%;
  }
  .search-box {
    width: 100%;
  }
}
</style>
