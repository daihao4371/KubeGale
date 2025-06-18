<template>
  <div class="terminal-container">
    <div class="header">
      <el-button v-if="isConnected" type="danger" @click="handleClose">关闭连接</el-button>
    </div>
    <div v-if="isConnected" ref="terminalRef" class="terminal"></div>
    <div v-else class="reconnect-message">
      <el-result
        icon="warning"
        title="连接已断开"
        sub-title="终端连接已关闭,请重新连接"
      >
        <template #extra>
          <el-button type="primary" @click="handleReconnect">重新连接</el-button>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { ElMessage } from 'element-plus'
import 'xterm/css/xterm.css'

// 常量定义
const INACTIVITY_TIMEOUT = 60000 // 不活动超时时间(毫秒)
const TERMINAL_CONFIG = {
  cols: 120,
  rows: 30,
  convertEol: true,
  scrollback: 1000,
  disableStdin: false,
  cursorStyle: 'block' as const,
  cursorBlink: true,
  fontFamily: 'Menlo, Monaco, Consolas, monospace',
  fontSize: 14,
  theme: {
    foreground: '#ffffff',
    background: '#000000',
    cursor: '#ffffff',
    black: '#000000',
    red: '#cd3131',
    green: '#0dbc79',
    yellow: '#e5e510',
    blue: '#2472c8',
    magenta: '#bc3fbc',
    cyan: '#11a8cd',
    white: '#e5e5e5',
    brightBlack: '#666666',
    brightRed: '#f14c4c',
    brightGreen: '#23d18b',
    brightYellow: '#f5f543',
    brightBlue: '#3b8eea',
    brightMagenta: '#d670d6',
    brightCyan: '#29b8db',
    brightWhite: '#e5e5e5'
  }
}

defineOptions({
  name: 'TerminalView'
})

const route = useRoute()
let terminal: Terminal | null = null
let fitAddon: FitAddon | null = null
let ws: WebSocket | null = null
const isConnected = ref(true)
let inactivityTimer: ReturnType<typeof setTimeout> | null = null
const terminalRef = ref<HTMLElement | null>(null)

// 重置不活动计时器
const resetInactivityTimer = () => {
  if (inactivityTimer) {
    clearTimeout(inactivityTimer)
  }
  inactivityTimer = setTimeout(() => {
    if (ws) {
      ElMessage.warning('由于长时间未操作,连接已自动关闭')
      ws.close(1000, '用户不活动超时')
    }
  }, INACTIVITY_TIMEOUT)
}

// 手动关闭连接
const handleClose = () => {
  if (ws) {
    ws.close(1000, '用户手动关闭连接')
    ElMessage.success('已关闭终端连接')
  }
}

// 重新连接
const handleReconnect = () => {
  const hostId = route.query.id
  if (hostId) {
    isConnected.value = true
    window.location.reload()
  }
}

// 处理终端输入
const handleTerminalInput = (data: string) => {
  if (!ws || ws.readyState !== WebSocket.OPEN) return;
  resetInactivityTimer();
  ws.send(data); // 只发送，不写入终端，由后端回显
}

onMounted(() => {
  const hostId = route.query.id
  if (!hostId) {
    ElMessage.error('缺少主机ID')
    return
  }

  // 获取 token
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.error('未登录或token失效')
    return
  }

  // 初始化终端
  terminal = new Terminal(TERMINAL_CONFIG)
  fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)

  if (terminalRef.value) {
    terminal.open(terminalRef.value)
    fitAddon.fit()
  }

  // 拼接 WebSocket 地址
  const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${wsProtocol}//127.0.0.1:8080/cmdb/hosts/terminal?id=${hostId}`
  ws = new WebSocket(wsUrl)

  ws.onopen = () => {
    terminal?.write('\r\n\x1B[1;32m正在连接...\x1B[0m\r\n')
    // 发送认证信息
    ws?.send(JSON.stringify({
      type: 'auth',
      token: token
    }))
    resetInactivityTimer()
  }

  ws.onmessage = (event) => {
    try {
      // 尝试解析消息是否为JSON
      const data = JSON.parse(event.data)
      if (data.type === 'auth_result') {
        if (data.success) {
          terminal?.write('\r\n\x1B[1;32m认证成功\x1B[0m\r\n')
        } else {
          terminal?.write('\r\n\x1B[1;31m认证失败: ' + data.message + '\x1B[0m\r\n')
          ws?.close()
        }
      } else {
        // 如果不是认证结果，则直接显示消息
        terminal?.write(event.data)
      }
    } catch {
      // 如果不是JSON，则直接显示消息
      terminal?.write(event.data)
    }
  }

  ws.onerror = (error) => {
    console.error('WebSocket错误:', error)
    ElMessage.error('终端连接错误')
    isConnected.value = false
  }

  ws.onclose = () => {
    isConnected.value = false
    if (inactivityTimer) {
      clearTimeout(inactivityTimer)
    }
  }

  // 处理终端输入
  terminal.onData(handleTerminalInput)

  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    fitAddon?.fit()
  })
})

onUnmounted(() => {
  if (ws) {
    ws.close()
  }
  if (terminal) {
    terminal.dispose()
  }
  if (inactivityTimer) {
    clearTimeout(inactivityTimer)
  }
  window.removeEventListener('resize', () => {
    fitAddon?.fit()
  })
})
</script>

<style scoped>
.terminal-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  padding: 0;
  margin: 0;
  background: #000;
  border-radius: 0;
  box-shadow: none;
  z-index: 10;
}

.header {
  position: absolute;
  top: var(--spacing-lg); /* Use theme variable */
  right: var(--spacing-xl); /* Use theme variable */
  margin: 0;
  z-index: 20;
}

.terminal {
  width: 100vw;
  height: 100vh;
  border-radius: 0;
  overflow: hidden;
}

.reconnect-message {
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 0;
}
</style>
