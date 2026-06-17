<template>
  <div id="app-root">
    <header class="app-header">
      <h1>AI 接入平台</h1>
      <button class="hamburger" @click="menuOpen = !menuOpen" :aria-label="menuOpen ? '关闭菜单' : '打开菜单'">
        <span></span><span></span><span></span>
      </button>
      <nav :class="{ open: menuOpen }">
        <router-link to="/" @click="menuOpen = false">首页</router-link>
        <router-link to="/news" @click="menuOpen = false">AI 资讯</router-link>
        <router-link to="/resources" @click="menuOpen = false">API 资源</router-link>
        <router-link to="/chat" @click="menuOpen = false">AI 对话</router-link>
      </nav>
    </header>
    <main class="app-main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const menuOpen = ref(false)
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
  padding: 0.75rem 1.5rem;
  display: flex;
  align-items: center;
  gap: 2rem;
  position: sticky;
  top: 0;
  z-index: 100;
}

.app-header h1 {
  font-size: 1.25rem;
  white-space: nowrap;
}

.hamburger {
  display: none;
  flex-direction: column;
  gap: 4px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
}

.hamburger span {
  display: block;
  width: 22px;
  height: 2px;
  background: white;
  border-radius: 1px;
  transition: transform 0.2s;
}

.app-header nav {
  display: flex;
  gap: 0.5rem;
}

.app-header a {
  color: rgba(255, 255, 255, 0.9);
  text-decoration: none;
  padding: 0.4rem 0.75rem;
  border-radius: 4px;
  transition: background 0.2s;
  font-size: 0.9rem;
}

.app-header a:hover,
.app-header a.router-link-active {
  background: rgba(255, 255, 255, 0.2);
}

.app-main {
  max-width: 1200px;
  margin: 1.5rem auto;
  padding: 0 1rem;
}

@media (max-width: 768px) {
  .app-header {
    justify-content: space-between;
  }

  .hamburger {
    display: flex;
  }

  .app-header nav {
    display: none;
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: #1557b0;
    flex-direction: column;
    padding: 0.5rem;
    gap: 0.25rem;
    box-shadow: 0 4px 12px rgba(0,0,0,0.2);
  }

  .app-header nav.open {
    display: flex;
  }

  .app-header a {
    padding: 0.6rem 1rem;
    border-radius: 4px;
  }

  .app-main {
    margin: 1rem auto;
    padding: 0 0.75rem;
  }
}
</style>
