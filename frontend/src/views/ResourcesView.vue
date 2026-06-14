<template>
  <div class="api-view">
    <div class="api-header">
      <h2>免费 API 资源清单</h2>
      <p class="subtitle">当前可用的免费大模型 API 资源</p>
    </div>

    <div class="resource-list">
      <div v-for="r in resources" :key="r.id" class="resource-card">
        <div class="resource-header">
          <h3>{{ r.name }}</h3>
          <span class="provider">{{ r.provider }}</span>
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
        <a v-if="r.doc_url" :href="r.doc_url" target="_blank" class="doc-link">查看文档</a>
        <div class="updated">上次更新：{{ formatTime(r.last_updated) }}</div>
      </div>
      <div v-if="!resources.length" class="empty">暂无 API 资源</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface APIResource {
  id: number
  name: string
  provider: string
  description: string
  endpoint: string
  free_quota: string
  doc_url: string
  last_updated: string
}

const resources = ref<APIResource[]>([])

async function loadResources() {
  try {
    const res = await fetch('/api/resources')
    const data = await res.json()
    resources.value = data.resources || []
  } catch (e) {
    console.error('Failed to load API resources:', e)
  }
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

.resource-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
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
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.resource-header h3 { color: #1a73e8; font-size: 1.1rem; }

.provider {
  background: #e8f0fe;
  color: #1a73e8;
  padding: 0.125rem 0.5rem;
  border-radius: 3px;
  font-size: 0.75rem;
}

.description { color: #666; margin-bottom: 0.75rem; line-height: 1.5; }

.quota, .endpoint {
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
}

.label { color: #999; }

.quota .value {
  color: #0d904f;
  font-weight: 500;
}

.endpoint code {
  background: #f5f5f5;
  padding: 0.125rem 0.375rem;
  border-radius: 3px;
  font-size: 0.8rem;
  word-break: break-all;
}

.doc-link {
  display: inline-block;
  margin-top: 0.5rem;
  color: #1a73e8;
  text-decoration: none;
  font-size: 0.875rem;
}

.doc-link:hover { text-decoration: underline; }

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
</style>
