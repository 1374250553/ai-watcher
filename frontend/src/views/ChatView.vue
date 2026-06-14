<template>
  <div class="chat-view">
    <div class="chat-container">
      <div class="chat-header">
        <h2>对话测试</h2>
        <div class="model-selector">
          <select v-model="selectedModel" @change="onModelChange">
            <optgroup v-for="(models, category) in modelCategories" :key="category" :label="category">
              <option v-for="model in models" :key="model" :value="model">{{ model }}</option>
            </optgroup>
          </select>
          <span class="model-count">{{ totalModels }} 个模型</span>
        </div>
      </div>

      <div class="messages" ref="messagesRef">
        <div
          v-for="(msg, index) in messages"
          :key="index"
          :class="['message', msg.role]"
        >
          <div class="message-role">{{ msg.role === 'user' ? '用户' : 'AI' }} · {{ selectedModel }}</div>
          <div class="message-content">{{ msg.content }}</div>
        </div>
        <div v-if="isLoading" class="message ai loading">
          <div class="message-role">AI · {{ selectedModel }}</div>
          <div class="message-content">思考中...</div>
        </div>
      </div>

      <form @submit.prevent="sendMessage" class="input-area">
        <input
          v-model="inputMessage"
          placeholder="输入消息测试 AI 响应..."
          :disabled="isLoading"
          @keydown.enter.exact.prevent="sendMessage"
        />
        <button type="submit" :disabled="isLoading || !inputMessage.trim()">
          {{ isLoading ? '发送中...' : '发送' }}
        </button>
      </form>

      <div v-if="error" class="error-message">
        {{ error }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'

interface Message {
  role: 'user' | 'assistant' | 'system'
  content: string
}

const messages = ref<Message[]>([])
const inputMessage = ref('')
const isLoading = ref(false)
const error = ref('')
const selectedModel = ref('qwen3.7-plus')
const messagesRef = ref<HTMLElement>()
const modelCategories = ref<Record<string, string[]>>({})
const totalModels = ref(0)

onMounted(async () => {
  try {
    const res = await fetch('/api/ai/models')
    const data = await res.json()
    modelCategories.value = data.categories
    totalModels.value = data.total
  } catch (e) {
    console.error('Failed to load models:', e)
  }
})

function onModelChange() {
  messages.value = []
  error.value = ''
}

async function sendMessage() {
  const content = inputMessage.value.trim()
  if (!content || isLoading.value) return

  messages.value.push({ role: 'user', content })
  inputMessage.value = ''
  isLoading.value = true
  error.value = ''

  try {
    const response = await fetch('/api/ai/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        messages: messages.value,
        model: selectedModel.value
      })
    })

    if (!response.ok) {
      const data = await response.json().catch(() => null)
      throw new Error(data?.details || data?.error || `HTTP ${response.status}`)
    }

    const data = await response.json()
    
    if (data.choices?.[0]?.message?.content) {
      messages.value.push({
        role: 'assistant',
        content: data.choices[0].message.content
      })
    } else {
      throw new Error('响应格式异常')
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : '请求失败'
    console.error('Chat error:', e)
  } finally {
    isLoading.value = false
    await nextTick()
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight
    }
  }
}
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
}

.chat-header {
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #eee;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.chat-header h2 {
  font-size: 1.25rem;
  color: #333;
  white-space: nowrap;
}

.model-selector {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.model-selector select {
  max-width: 280px;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.875rem;
  background: white;
}

.model-count {
  font-size: 0.75rem;
  color: #999;
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.message {
  max-width: 80%;
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
  background: #f0f0f0;
  color: #333;
}

.message-role {
  font-size: 0.75rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
  opacity: 0.8;
}

.message-content {
  white-space: pre-wrap;
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

.input-area button {
  padding: 0.75rem 1.5rem;
  background: #1a73e8;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

.input-area button:hover:not(:disabled) {
  background: #1557b0;
}

.input-area button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-message {
  padding: 0.75rem 1.5rem;
  background: #fee;
  color: #c00;
  border-top: 1px solid #fcc;
  font-size: 0.875rem;
}
</style>
