<template>
  <div class="home-view">
    <div class="hero">
      <h2>欢迎使用 AI 接入平台</h2>
      <p>提供统一的 AI 能力接入与管理服务</p>
      <div class="hero-actions">
        <router-link to="/news" class="primary-btn">AI 资讯</router-link>
        <router-link to="/resources" class="primary-btn">API 资源</router-link>
        <router-link to="/chat" class="primary-btn">对话测试</router-link>
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
</script>

<style scoped>
.home-view {
  text-align: center;
}

.hero {
  background: white;
  padding: 3rem 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
}

.hero h2 {
  font-size: 2rem;
  margin-bottom: 1rem;
  color: #1a73e8;
}

.hero p {
  font-size: 1.125rem;
  color: #666;
  margin-bottom: 1.5rem;
}

.hero-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.primary-btn {
  display: inline-block;
  background: #1a73e8;
  color: white;
  padding: 0.75rem 2rem;
  border-radius: 4px;
  text-decoration: none;
  font-weight: 500;
  transition: background 0.2s;
}

.primary-btn:hover {
  background: #1557b0;
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
</style>
