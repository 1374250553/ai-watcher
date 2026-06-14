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
        <button class="refresh-btn" @click="loadNews(page)">刷新</button>
      </div>
    </div>

    <div class="news-list">
      <div v-for="item in newsList" :key="item.id" class="news-item">
        <a :href="item.url" target="_blank" class="news-title">{{ item.title }}</a>
        <div class="news-meta">
          <span class="news-source">{{ item.source }}</span>
          <span class="news-time">{{ formatTime(item.created_at) }}</span>
        </div>
        <div v-if="item.summary" class="news-summary">{{ item.summary }}</div>
      </div>
      <div v-if="!newsList.length" class="empty">暂无资讯</div>
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
const sources = ref<string[]>(['机器之心', '量子位', '微信公众号', '知乎', 'V2EX'])

async function loadNews(p: number) {
  page.value = p
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
  }
}

function formatTime(t: string) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

onMounted(() => loadNews(1))
</script>

<style scoped>
.news-header {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  margin-bottom: 1rem;
}

.news-header h2 { margin-bottom: 1rem; color: #333; }
.news-header .count { font-size: 0.875rem; color: #999; font-weight: normal; }

.controls {
  display: flex;
  gap: 0.75rem;
  align-items: center;
  flex-wrap: wrap;
}

.controls select, .controls input {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.875rem;
}

.search-box {
  display: flex;
  gap: 0.5rem;
  flex: 1;
  min-width: 200px;
}

.search-box input { flex: 1; }

.controls button, .refresh-btn {
  padding: 0.5rem 1rem;
  background: #1a73e8;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
  white-space: nowrap;
}

.controls button:hover { background: #1557b0; }

.news-list {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.news-item {
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #f0f0f0;
}

.news-item:last-child { border-bottom: none; }

.news-title {
  font-size: 1rem;
  color: #1a73e8;
  text-decoration: none;
  font-weight: 500;
  display: block;
  margin-bottom: 0.25rem;
}

.news-title:hover { text-decoration: underline; }

.news-meta {
  font-size: 0.75rem;
  color: #999;
  display: flex;
  gap: 1rem;
}

.news-source {
  background: #e8f0fe;
  color: #1a73e8;
  padding: 0.125rem 0.5rem;
  border-radius: 3px;
}

.news-summary {
  margin-top: 0.5rem;
  color: #666;
  font-size: 0.875rem;
  line-height: 1.5;
}

.empty {
  text-align: center;
  color: #999;
  padding: 3rem;
}

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
}

.pagination button:disabled { opacity: 0.5; cursor: not-allowed; }
.pagination button:hover:not(:disabled) { background: #f5f5f5; }
</style>
