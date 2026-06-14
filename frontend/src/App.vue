<template>
  <div id="app-root">
    <header class="app-header">
      <h1>AI 接入平台</h1>
      <nav>
        <router-link to="/">首页</router-link>
        <router-link to="/news">AI 资讯</router-link>
        <router-link to="/resources">API 资源</router-link>
        <router-link to="/chat">对话测试</router-link>
      </nav>
    </header>
    <main class="app-main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: #f5f5f5;
  color: #333;
}

.app-header {
  background: #1a73e8;
  color: white;
  padding: 1rem 2rem;
  display: flex;
  align-items: center;
  gap: 2rem;
}

.app-header h1 {
  font-size: 1.5rem;
}

.app-header nav {
  display: flex;
  gap: 1rem;
}

.app-header a {
  color: rgba(255, 255, 255, 0.9);
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background 0.2s;
}

.app-header a:hover,
.app-header a.router-link-active {
  background: rgba(255, 255, 255, 0.2);
}

.app-main {
  max-width: 1200px;
  margin: 2rem auto;
  padding: 0 1rem;
}
</style>
