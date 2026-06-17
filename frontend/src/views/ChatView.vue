<template>
  <div class="chat-view">
    <div class="chat-container">
      <div class="chat-header">
        <div class="header-left">
          <h2>AI 对话</h2>
          <span v-if="selectedModel" class="current-model">{{ selectedModel }}</span>
        </div>
        <div class="header-actions">
          <button class="new-chat-btn" @click="clearChat">新对话</button>
        </div>
      </div>

      <div v-if="showModelPicker" class="model-picker-overlay" @click.self="showModelPicker = false">
        <div class="model-picker">
          <h3>选择模型</h3>
          <div class="vendor-list">
            <button
              v-for="vendor in vendors"
              :key="vendor.name"
              :class="['vendor-btn', { active: activeVendor === vendor.name }]"
              @click="selectVendor(vendor.name)"
            >
              {{ vendor.label }}
              <span class="vendor-count">{{ vendor.count }}</span>
            </button>
          </div>
          <div v-if="vendorModels.length" class="model-list">
            <button
              v-for="m in vendorModels"
              :key="m"
              :class="['model-btn', { selected: selectedModel === m }]"
              @click="selectModel(m)"
            >
              {{ m }}
            </button>
          </div>
          <button class="close-picker" @click="showModelPicker = false">关闭</button>
        </div>
      </div>

      <div class="messages" ref="messagesRef">
        <div v-if="!messages.length" class="empty-chat">
          <p>{{ modelCategories ? '点击模型名称选择模型，开始对话' : '加载中...' }}</p>
        </div>
        <div
          v-for="(msg, index) in messages"
          :key="index"
          :class="['message', msg.role]"
        >
          <div class="message-role">{{ msg.role === 'user' ? '用户' : 'AI' }}</div>
          <div v-if="msg.role === 'user'" class="message-content">{{ msg.content }}</div>
          <div
            v-else
            ref="aiMessages"
            class="message-content markdown-body"
            v-html="renderMarkdown(msg.content)"
          ></div>
        </div>
        <div v-if="isStreaming" class="message ai streaming">
          <div class="message-role">AI</div>
          <div
            class="message-content markdown-body"
            v-html="renderMarkdown(streamingContent)"
          ></div>
          <span class="cursor">|</span>
        </div>
      </div>

      <form @submit.prevent="sendMessage" class="input-area">
        <input
          v-model="inputMessage"
          placeholder="输入消息..."
          :disabled="isStreaming"
          @keydown.enter.exact.prevent="sendMessage"
        />
        <button
          v-if="!isStreaming"
          type="button"
          class="model-select-btn"
          @click="showModelPicker = true"
        >
          {{ selectedModel || '选择模型' }}
        </button>
        <button v-if="isStreaming" type="button" class="stop-btn" @click="stopStreaming">
          停止生成
        </button>
        <button type="submit" :disabled="isStreaming || !inputMessage.trim()">
          发送
        </button>
      </form>

      <div v-if="error" class="error-message">
        {{ error }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted, computed } from 'vue'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'

interface Message {
  role: 'user' | 'assistant' | 'system'
  content: string
}

interface Vendor {
  name: string
  label: string
  count: number
  categories: string[]
}

const STORAGE_KEY = 'ai-chat-messages'

const messages = ref<Message[]>([])
const inputMessage = ref('')
const isStreaming = ref(false)
const streamingContent = ref('')
const error = ref('')
const selectedModel = ref('qwen3.7-plus')
const messagesRef = ref<HTMLElement>()
const modelCategories = ref<Record<string, string[]>>({})
const vendors = ref<Vendor[]>([])
const activeVendor = ref('')
const showModelPicker = ref(false)

let abortController: AbortController | null = null

marked.setOptions({
  breaks: true,
  gfm: true,
})

const renderer = new marked.Renderer()
renderer.code = function(token: { text: string; lang?: string; raw?: string }) {
  const text = typeof token === 'string' ? token : (token.text || token.raw || '')
  let lang = ''
  let code = text
  const langMatch = text.match(/^(\w+)\s*\n/)
  if (langMatch) {
    lang = langMatch[1]
    code = text.slice(langMatch[0].length)
  }
  const highlighted = lang && hljs.getLanguage(lang)
    ? hljs.highlight(code, { language: lang }).value
    : hljs.highlightAuto(code).value
  const langLabel = lang || 'code'
  return `<div class="code-block">
    <div class="code-header">
      <span class="code-lang">${langLabel}</span>
      <button class="copy-btn" onclick="navigator.clipboard.writeText(this.dataset.code);this.textContent='已复制';setTimeout(()=>this.textContent='复制',2000)" data-code="${escapeHtml(code)}">复制</button>
    </div>
    <pre><code class="hljs ${lang ? 'language-' + lang : ''}">${highlighted}</code></pre>
  </div>`
}

marked.use({ renderer })

function escapeHtml(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
}

function renderMarkdown(text: string): string {
  if (!text) return ''
  return marked.parse(text) as string
}

const vendorModels = computed(() => {
  if (!activeVendor.value || !modelCategories.value) return []
  const vendor = vendors.value.find(v => v.name === activeVendor.value)
  if (!vendor) return []
  const models: string[] = []
  for (const cat of vendor.categories) {
    if (modelCategories.value[cat]) {
      models.push(...modelCategories.value[cat])
    }
  }
  return models
})

async function loadModels() {
  try {
    const res = await fetch('/api/ai/models')
    const data = await res.json()
    modelCategories.value = data.categories
    buildVendors(data.categories)
  } catch (e) {
    console.error('Failed to load models:', e)
  }
}

function buildVendors(categories: Record<string, string[]>) {
  const vendorMap: Record<string, { label: string; count: number; categories: string[] }> = {
    'tongyi': { label: '通义千问', count: 0, categories: [] },
    'deepseek': { label: 'DeepSeek', count: 0, categories: [] },
    'kimi': { label: 'Kimi', count: 0, categories: [] },
    'glm': { label: 'GLM 智谱', count: 0, categories: [] },
    'minimax': { label: 'MiniMax', count: 0, categories: [] },
  }

  for (const [cat, models] of Object.entries(categories)) {
    if (cat.startsWith('通义千问')) {
      vendorMap.tongyi.count += models.length
      vendorMap.tongyi.categories.push(cat)
    } else if (cat === 'DeepSeek') {
      vendorMap.deepseek.count += models.length
      vendorMap.deepseek.categories.push(cat)
    } else if (cat.startsWith('Kimi')) {
      vendorMap.kimi.count += models.length
      vendorMap.kimi.categories.push(cat)
    } else if (cat.startsWith('GLM')) {
      vendorMap.glm.count += models.length
      vendorMap.glm.categories.push(cat)
    } else if (cat === 'MiniMax') {
      vendorMap.minimax.count += models.length
      vendorMap.minimax.categories.push(cat)
    }
  }

  vendors.value = Object.entries(vendorMap)
    .filter(([, v]) => v.count > 0)
    .map(([name, v]) => ({
      name,
      label: v.label,
      count: v.count,
      categories: v.categories,
    }))
}

function selectVendor(name: string) {
  activeVendor.value = name
}

function selectModel(model: string) {
  selectedModel.value = model
  showModelPicker.value = false
}

function clearChat() {
  messages.value = []
  error.value = ''
  saveToLocal()
}

async function sendMessage() {
  const content = inputMessage.value.trim()
  if (!content || isStreaming.value) return

  messages.value.push({ role: 'user', content })
  inputMessage.value = ''
  isStreaming.value = true
  streamingContent.value = ''
  error.value = ''
  saveToLocal()

  await nextTick()
  scrollToBottom()

  abortController = new AbortController()

  try {
    const response = await fetch('/api/ai/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        messages: messages.value,
        model: selectedModel.value,
        stream: true,
      }),
      signal: abortController.signal,
    })

    if (!response.ok) {
      const data = await response.json().catch(() => null)
      throw new Error(data?.details || data?.error || `HTTP ${response.status}`)
    }

    const reader = response.body?.getReader()
    if (!reader) throw new Error('No reader available')

    const decoder = new TextDecoder()
    let fullContent = ''
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const line of lines) {
        const trimmed = line.trim()
        if (!trimmed || !trimmed.startsWith('data: ')) continue
        const data = trimmed.slice(6)
        if (data === '[DONE]') continue

        try {
          const chunk = JSON.parse(data)
          const delta = chunk.choices?.[0]?.delta?.content
          if (delta) {
            fullContent += delta
            streamingContent.value = fullContent
            await nextTick()
            scrollToBottom()
          }
        } catch {
          // skip malformed chunks
        }
      }
    }

    if (fullContent) {
      messages.value.push({ role: 'assistant', content: fullContent })
      saveToLocal()
    }
  } catch (e: unknown) {
    if (e instanceof DOMException && e.name === 'AbortError') {
      if (streamingContent.value) {
        messages.value.push({ role: 'assistant', content: streamingContent.value })
        saveToLocal()
      }
    } else {
      error.value = e instanceof Error ? e.message : '请求失败'
    }
  } finally {
    isStreaming.value = false
    streamingContent.value = ''
    abortController = null
    await nextTick()
    scrollToBottom()
  }
}

function stopStreaming() {
  if (abortController) {
    abortController.abort()
  }
}

function scrollToBottom() {
  if (messagesRef.value) {
    messagesRef.value.scrollTop = messagesRef.value.scrollHeight
  }
}

function saveToLocal() {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(messages.value))
  } catch {
    // localStorage full
  }
}

function loadFromLocal() {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      messages.value = JSON.parse(saved)
    }
  } catch {
    // ignore
  }
}

onMounted(async () => {
  loadFromLocal()
  await loadModels()
  await nextTick()
  scrollToBottom()
})
</script>

<style scoped>
.chat-view {
  height: calc(100vh - 140px);
}

.chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  position: relative;
}

.chat-header {
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #eee;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.chat-header h2 {
  font-size: 1.25rem;
  color: #333;
}

.current-model {
  font-size: 0.8rem;
  background: #e8f0fe;
  color: #1a73e8;
  padding: 0.2rem 0.6rem;
  border-radius: 3px;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.new-chat-btn {
  padding: 0.4rem 1rem;
  background: #f0f0f0;
  color: #333;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background 0.2s;
}

.new-chat-btn:hover {
  background: #e0e0e0;
}

.model-picker-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.model-picker {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  width: 90%;
  max-width: 600px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
}

.model-picker h3 {
  margin-bottom: 1rem;
  color: #333;
}

.vendor-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.vendor-btn {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.vendor-btn:hover {
  border-color: #1a73e8;
  color: #1a73e8;
}

.vendor-btn.active {
  background: #1a73e8;
  color: white;
  border-color: #1a73e8;
}

.vendor-count {
  font-size: 0.75rem;
  opacity: 0.7;
}

.model-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.model-btn {
  padding: 0.35rem 0.75rem;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  background: #f9f9f9;
  cursor: pointer;
  font-size: 0.8rem;
  transition: all 0.2s;
}

.model-btn:hover {
  border-color: #1a73e8;
  color: #1a73e8;
}

.model-btn.selected {
  background: #1a73e8;
  color: white;
  border-color: #1a73e8;
}

.close-picker {
  display: block;
  width: 100%;
  padding: 0.5rem;
  background: #f0f0f0;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
}

.close-picker:hover {
  background: #e0e0e0;
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.empty-chat {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 1rem;
}

.message {
  max-width: 85%;
  padding: 1rem;
  border-radius: 8px;
  line-height: 1.6;
}

.message.user {
  align-self: flex-end;
  background: #1a73e8;
  color: white;
}

.message.ai {
  align-self: flex-start;
  background: #f5f5f5;
  color: #333;
}

.message.streaming {
  align-self: flex-start;
  background: #f5f5f5;
  color: #333;
  display: flex;
  gap: 0;
}

.message-role {
  font-size: 0.75rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
  opacity: 0.8;
}

.cursor {
  animation: blink 1s step-end infinite;
  color: #1a73e8;
  font-weight: bold;
}

@keyframes blink {
  50% { opacity: 0; }
}

.message-content {
  word-break: break-word;
}

.message-content :deep(p) {
  margin-bottom: 0.5rem;
}

.message-content :deep(p:last-child) {
  margin-bottom: 0;
}

.message-content :deep(ul), .message-content :deep(ol) {
  padding-left: 1.5rem;
  margin-bottom: 0.5rem;
}

.message-content :deep(h1), .message-content :deep(h2), .message-content :deep(h3),
.message-content :deep(h4), .message-content :deep(h5), .message-content :deep(h6) {
  margin: 0.75rem 0 0.5rem;
  font-weight: 600;
}

.message-content :deep(table) {
  border-collapse: collapse;
  margin: 0.5rem 0;
  width: 100%;
}

.message-content :deep(th), .message-content :deep(td) {
  border: 1px solid #ddd;
  padding: 0.5rem;
  text-align: left;
}

.message-content :deep(th) {
  background: #f5f5f5;
  font-weight: 600;
}

.message-content :deep(blockquote) {
  border-left: 3px solid #1a73e8;
  padding-left: 1rem;
  color: #666;
  margin: 0.5rem 0;
}

.message-content :deep(a) {
  color: #1a73e8;
  text-decoration: underline;
}

.message-content :deep(code:not(pre code)) {
  background: rgba(0, 0, 0, 0.06);
  padding: 0.15rem 0.35rem;
  border-radius: 3px;
  font-size: 0.875em;
}

.message-content :deep(.code-block) {
  margin: 0.5rem 0;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid #e0e0e0;
}

.message-content :deep(.code-header) {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.35rem 0.75rem;
  background: #f0f0f0;
  border-bottom: 1px solid #e0e0e0;
  font-size: 0.75rem;
}

.message-content :deep(.code-lang) {
  color: #666;
  font-weight: 500;
}

.message-content :deep(.copy-btn) {
  background: transparent;
  border: 1px solid #ccc;
  border-radius: 3px;
  padding: 0.15rem 0.5rem;
  cursor: pointer;
  font-size: 0.75rem;
  color: #666;
  transition: all 0.2s;
}

.message-content :deep(.copy-btn:hover) {
  background: #e0e0e0;
  border-color: #999;
}

.message-content :deep(pre) {
  margin: 0;
  padding: 0.75rem;
  overflow-x: auto;
  background: #1e1e1e;
}

.message-content :deep(pre code) {
  font-size: 0.85rem;
  line-height: 1.5;
}

.input-area {
  display: flex;
  gap: 0.5rem;
  padding: 1rem 1.5rem;
  border-top: 1px solid #eee;
}

.input-area input {
  flex: 1;
  padding: 0.75rem 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  outline: none;
  transition: border-color 0.2s;
}

.input-area input:focus {
  border-color: #1a73e8;
}

.input-area input:disabled {
  background: #f5f5f5;
}

.model-select-btn {
  padding: 0.75rem 1rem;
  background: #f0f0f0;
  color: #333;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  white-space: nowrap;
  transition: background 0.2s;
}

.model-select-btn:hover {
  background: #e0e0e0;
}

.input-area button[type="submit"] {
  padding: 0.75rem 1.5rem;
  background: #1a73e8;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

.input-area button[type="submit"]:hover:not(:disabled) {
  background: #1557b0;
}

.input-area button[type="submit"]:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.stop-btn {
  padding: 0.75rem 1rem;
  background: #d93025;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  white-space: nowrap;
}

.stop-btn:hover {
  background: #c5221f;
}

.error-message {
  padding: 0.75rem 1.5rem;
  background: #fee;
  color: #c00;
  border-top: 1px solid #fcc;
  font-size: 0.875rem;
}
</style>
