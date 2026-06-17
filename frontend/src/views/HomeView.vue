<template>
  <div class="home-view">
    <div class="hero">
      <h2>欢迎使用 AI 接入平台</h2>
      <p>提供统一的 AI 能力接入与管理服务</p>
      <router-link to="/chat" class="cta-btn">立即体验 AI 对话</router-link>
    </div>

    <div class="stats-dashboard">
      <div class="stat-card">
        <div class="stat-icon models-icon">M</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.model_count }}</span>
          <span class="stat-label">可用模型</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon news-icon">N</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.news_count }}</span>
          <span class="stat-label">资讯总数</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon time-icon">T</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.last_fetch ? '已采集' : '暂无' }}</span>
          <span class="stat-label">{{ stats.last_fetch || '最近采集' }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon api-icon">A</div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.resource_count }}</span>
          <span class="stat-label">API 资源</span>
        </div>
      </div>
    </div>

    <div v-if="latestNews.length" class="latest-news">
      <h3>最新资讯</h3>
      <div class="news-grid">
        <div v-for="item in latestNews" :key="item.id" class="news-card" @click="goNews">
          <h4>{{ item.title }}</h4>
          <div class="news-card-meta">
            <span class="news-source-tag">{{ item.source }}</span>
            <span class="news-time">{{ formatTime(item.created_at) }}</span>
          </div>
          <p class="news-excerpt">{{ truncate(item.summary, 60) }}</p>
        </div>
      </div>
    </div>

    <div class="features">
      <div class="feature-card">
        <h3>通义千问全系</h3>
        <p>qwen3.7 / qwen3.6 / qwen3.5 / qwen3 等 20+ 个版本，包含 max、plus、flash 不同规格</p>
      </div>
      <div class="feature-card">
        <h3>DeepSeek 全系</h3>
        <p>deepseek-v4-flash、v4-pro、v3.2、v3.1、v3，以及 deepseek-r1 思考模型</p>
      </div>
      <div class="feature-card">
        <h3>多厂商模型</h3>
        <p>Kimi (月之暗面)、GLM (智谱)、MiniMax 等主流国产大模型全覆盖</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

interface NewsItem {
  id: number
  title: string
  url: string
  summary: string
  source: string
  created_at: string
}

interface Stats {
  model_count: number
  news_count: number
  resource_count: number
  last_fetch: string
  source_counts: Record<string, number>
  latest_news: NewsItem[]
}

const stats = ref<Stats>({
  model_count: 0,
  news_count: 0,
  resource_count: 0,
  last_fetch: '',
  source_counts: {},
  latest_news: [],
})

const latestNews = ref<NewsItem[]>([])

async function loadStats() {
  try {
    const [statsRes, modelsRes] = await Promise.all([
      fetch('/api/stats'),
      fetch('/api/ai/models'),
    ])
    const statsData = await statsRes.json()
    const modelsData = await modelsRes.json()

    stats.value = { ...stats.value, ...statsData }
    stats.value.model_count = modelsData.total || stats.value.model_count
    latestNews.value = statsData.latest_news || []
  } catch (e) {
    console.error('Failed to load stats:', e)
  }
}

function goNews() {
  router.push('/news')
}

function formatTime(t: string) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

function truncate(s: string, maxLen: number) {
  if (!s) return '暂无摘要'
  return s.length > maxLen ? s.slice(0, maxLen) + '...' : s
}

onMounted(loadStats)
</script>

<style scoped>
.home-view {
  text-align: center;
}

.hero {
  background: linear-gradient(135deg, #1a73e8 0%, #4a90d9 100%);
  color: white;
  padding: 3rem 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(26, 115, 232, 0.3);
  margin-bottom: 2rem;
}

.hero h2 {
  font-size: 2rem;
  margin-bottom: 1rem;
}

.hero p {
  font-size: 1.125rem;
  opacity: 0.9;
  margin-bottom: 1.5rem;
}

.cta-btn {
  display: inline-block;
  background: white;
  color: #1a73e8;
  padding: 0.85rem 2.5rem;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  font-size: 1.1rem;
  transition: transform 0.2s, box-shadow 0.2s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.cta-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.25);
}

.stats-dashboard {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 1rem;
  text-align: left;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 1.25rem;
  color: white;
  flex-shrink: 0;
}

.models-icon { background: #1a73e8; }
.news-icon { background: #34a853; }
.time-icon { background: #fb8c00; }
.api-icon { background: #7c4dff; }

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: #333;
}

.stat-label {
  font-size: 0.8rem;
  color: #999;
  margin-top: 0.15rem;
}

.latest-news {
  margin-bottom: 2rem;
}

.latest-news h3 {
  text-align: left;
  margin-bottom: 1rem;
  color: #333;
  font-size: 1.2rem;
}

.news-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
}

.news-card {
  background: white;
  padding: 1.25rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: left;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.news-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.news-card h4 {
  font-size: 1rem;
  color: #333;
  margin-bottom: 0.5rem;
  line-height: 1.4;
}

.news-card-meta {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 0.5rem;
  font-size: 0.75rem;
  color: #999;
}

.news-source-tag {
  background: #e8f0fe;
  color: #1a73e8;
  padding: 0.1rem 0.4rem;
  border-radius: 3px;
}

.news-excerpt {
  font-size: 0.85rem;
  color: #666;
  line-height: 1.5;
}

.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
}

.feature-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: left;
}

.feature-card h3 {
  color: #1a73e8;
  margin-bottom: 0.5rem;
}

.feature-card p {
  color: #666;
  line-height: 1.6;
}

@media (max-width: 768px) {
  .hero {
    padding: 2rem 1rem;
  }
  .hero h2 {
    font-size: 1.4rem;
  }
  .hero p {
    font-size: 0.95rem;
  }
  .cta-btn {
    font-size: 1rem;
    padding: 0.7rem 1.5rem;
  }
  .stats-dashboard {
    grid-template-columns: repeat(2, 1fr);
    gap: 0.75rem;
  }
  .stat-card {
    padding: 1rem;
    gap: 0.5rem;
  }
  .stat-icon {
    width: 36px;
    height: 36px;
    font-size: 1rem;
    border-radius: 8px;
  }
  .stat-value {
    font-size: 1.2rem;
  }
  .news-grid {
    grid-template-columns: 1fr;
  }
  .features {
    grid-template-columns: 1fr;
  }
}
</style>
