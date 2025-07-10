<template>
  <div class="cloudtty-container">
    <div class="terminal-header">
      <h3>CloudTTY 终端</h3>
      <div class="header-actions">
        <el-button size="small" @click="handleConnect" :disabled="isConnected">连接</el-button>
        <el-button size="small" @click="handleDisconnect" :disabled="!isConnected">断开</el-button>
        <el-button size="small" @click="handleClear">清屏</el-button>
      </div>
    </div>
    
    <div class="terminal-content">
      <div class="terminal-window" ref="terminalRef">
        <div v-if="!isConnected" class="connection-prompt">
          <el-empty description="点击连接按钮开始终端会话">
            <el-button type="primary" @click="handleConnect">连接终端</el-button>
          </el-empty>
        </div>
        <div v-else class="terminal-output">
          <div v-for="(line, index) in terminalOutput" :key="index" class="terminal-line">
            {{ line }}
          </div>
          <div class="terminal-input-line">
            <span class="prompt">$ </span>
            <input 
              ref="inputRef"
              v-model="currentInput"
              @keydown.enter="handleCommand"
              @keydown="handleKeydown"
              class="terminal-input"
              placeholder="输入命令..."
            />
          </div>
        </div>
      </div>
    </div>
    
    <div class="terminal-footer">
      <span class="status">{{ isConnected ? '已连接' : '未连接' }}</span>
      <el-button @click="handleClose">关闭</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

interface Props {
  clusterId?: string
  namespace?: string
  podName?: string
  containerName?: string
}

interface Emits {
  (e: 'close'): void
}

const props = withDefaults(defineProps<Props>(), {
  clusterId: '',
  namespace: '',
  podName: '',
  containerName: ''
})

const emit = defineEmits<Emits>()

const isConnected = ref(false)
const terminalOutput = ref<string[]>([])
const currentInput = ref('')
const terminalRef = ref<HTMLElement>()
const inputRef = ref<HTMLInputElement>()

const handleConnect = () => {
  isConnected.value = true
  terminalOutput.value.push('连接到终端...')
  terminalOutput.value.push('欢迎使用 CloudTTY 终端')
  terminalOutput.value.push('')
  
  // 聚焦到输入框
  setTimeout(() => {
    inputRef.value?.focus()
  }, 100)
}

const handleDisconnect = () => {
  isConnected.value = false
  terminalOutput.value.push('终端连接已断开')
}

const handleClear = () => {
  terminalOutput.value = []
}

const handleCommand = () => {
  if (!currentInput.value.trim()) return
  
  const command = currentInput.value
  terminalOutput.value.push(`$ ${command}`)
  
  // 模拟命令执行
  setTimeout(() => {
    terminalOutput.value.push(`执行命令: ${command}`)
    terminalOutput.value.push('命令执行完成')
    terminalOutput.value.push('')
  }, 100)
  
  currentInput.value = ''
  
  // 滚动到底部
  setTimeout(() => {
    if (terminalRef.value) {
      terminalRef.value.scrollTop = terminalRef.value.scrollHeight
    }
  }, 50)
}

const handleKeydown = (event: KeyboardEvent) => {
  // 处理特殊按键
  if (event.key === 'Tab') {
    event.preventDefault()
    currentInput.value += '  '
  }
}

const handleClose = () => {
  if (isConnected.value) {
    handleDisconnect()
  }
  emit('close')
}

onMounted(() => {
  // 初始化终端
})

onUnmounted(() => {
  // 清理资源
  if (isConnected.value) {
    handleDisconnect()
  }
})
</script>

<style scoped>
.cloudtty-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  background-color: var(--el-bg-color-page);
  border-bottom: 1px solid var(--el-border-color-light);
}

.terminal-header h3 {
  margin: 0;
  color: var(--el-text-color-primary);
}

.header-actions {
  display: flex;
  gap: 10px;
}

.terminal-content {
  flex: 1;
  overflow: hidden;
}

.terminal-window {
  height: 100%;
  background-color: #1e1e1e;
  color: #ffffff;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  padding: 10px;
  overflow-y: auto;
}

.connection-prompt {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.terminal-output {
  min-height: 100%;
}

.terminal-line {
  white-space: pre-wrap;
  word-break: break-all;
}

.terminal-input-line {
  display: flex;
  align-items: center;
}

.prompt {
  color: #00ff00;
  margin-right: 5px;
}

.terminal-input {
  flex: 1;
  background: transparent;
  border: none;
  color: #ffffff;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  outline: none;
}

.terminal-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  background-color: var(--el-bg-color-page);
  border-top: 1px solid var(--el-border-color-light);
}

.status {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}
</style> 