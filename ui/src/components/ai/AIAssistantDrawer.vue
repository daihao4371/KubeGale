<template>
    <el-drawer
      v-model="visible"
      direction="rtl"
      size="40%"
      :modal="true"
      :close-on-press-escape="true"
      :close-on-click-modal="true"
      :with-header="false"
      custom-class="ai-drawer"
    >
      <div class="ai-drawer-header">
        <span class="ai-avatar">A</span>
        <span class="ai-title">KubeGale AI助手</span>
        <span class="ai-beta">Beta</span>
        <el-icon class="ai-close" @click="visible = false"><Close /></el-icon>
      </div>
      <el-scrollbar class="ai-drawer-body" ref="scrollbarRef">
        <div v-for="msg in messages" :key="msg.id" :class="['ai-msg', msg.role]">
          <span v-if="msg.role==='ai'" class="ai-avatar">A</span>
          <div class="ai-msg-content">
            <div class="ai-msg-bubble" v-html="msg.html"></div>
            <div class="ai-msg-time">{{ msg.time }}</div>
          </div>
          <span v-if="msg.role==='user'" class="ai-avatar user">U</span>
        </div>
      </el-scrollbar>
      <div class="ai-drawer-footer">
        <el-input
          v-model="input"
          type="textarea"
          :rows="2"
          placeholder="Shift+Enter 换行，Enter 发送"
          @keydown.enter.prevent="onEnter"
          @keydown.shift.enter.stop
          class="ai-input"
        />
        <el-button type="primary" class="ai-send-btn" @click="send">发送</el-button>
      </div>
    </el-drawer>
  </template>
  
  <script setup lang="ts">
  import { ref, watch, nextTick } from 'vue'
  import { Close } from '@element-plus/icons-vue'
  import markdownit from 'markdown-it'
  
  const md = markdownit({ breaks: true })
  const props = defineProps<{ modelValue: boolean }>()
  const emit = defineEmits(['update:modelValue'])
  const visible = ref(props.modelValue)
  watch(() => props.modelValue, v => visible.value = v)
  watch(visible, v => emit('update:modelValue', v))
  
  const input = ref('')
  const scrollbarRef = ref()
  function now() {
    const d = new Date()
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }
  
  const messages = ref([
    {
      id: 1,
      role: 'ai',
      content: '你好，我是 KubeGale AI 助手！Kubernetes 相关问题都可以问我哟~',
      html: md.render('你好，我是 KubeGale AI 助手！Kubernetes 相关问题都可以问我哟~'),
      time: now()
    }
  ])
  
  function send() {
    if (!input.value.trim()) return
    const content = input.value
    messages.value.push({ id: Date.now(), role: 'user', content, html: md.render(content), time: now() })
    input.value = ''
    nextTick(() => scrollbarRef.value?.setScrollTop(99999))
    setTimeout(() => {
      const reply = '收到：' + content + '\n\n- 这是一个列表\n- 支持 **Markdown** 渲染\n- 代码块：\n```shell\necho hello\n```'
      messages.value.push({ id: Date.now() + 1, role: 'ai', content: reply, html: md.render(reply), time: now() })
      nextTick(() => scrollbarRef.value?.setScrollTop(99999))
    }, 800)
  }
  
  function onEnter(e: KeyboardEvent) {
    if (e.shiftKey) return
    send()
  }
  </script>
  
  <style scoped>
  .ai-drawer-header {
    display: flex;
    align-items: center;
    padding: 16px;
    border-bottom: 1px solid #eee;
    font-weight: bold;
    font-size: 18px;
    background: #fff;
  }
  
  .ai-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: #409eff;
    color: #fff;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    margin: 0 16px;
    flex-shrink: 0;
    box-shadow: 0 2px 8px rgba(64,158,255,0.08);
  }
  
  .ai-avatar.user {
    background: #e6f7ff;
    color: #409eff;
  }
  
  .ai-title {
    margin-right: 8px;
  }
  
  .ai-beta {
    font-size: 12px;
    color: #67c23a;
    background: #f0f9eb;
    border-radius: 4px;
    padding: 2px 6px;
    margin-right: auto;
  }
  
  .ai-close {
    cursor: pointer;
    margin-left: 8px;
  }
  
  .ai-drawer-body {
    background: #f7f8fa;
    padding: 32px 24px 24px 24px;
    min-height: 300px;
    max-height: calc(100vh - 180px);
    box-sizing: border-box;
  }
  
  .ai-msg {
    display: flex;
    align-items: center;
    margin-bottom: 24px;
  }
  
  .ai-msg.ai {
    flex-direction: row;
    justify-content: flex-start;
  }
  
  .ai-msg.user {
    flex-direction: row-reverse;
    justify-content: flex-end;
  }
  
  .ai-msg-content {
    max-width: 70vw;
    min-width: 120px;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    position: relative;
  }
  
  .ai-msg.user .ai-msg-content {
    align-items: flex-end;
  }
  
  .ai-msg-bubble {
    background: #fff;
    border-radius: 20px;
    padding: 18px 24px;
    font-size: 16px;
    box-shadow: 0 4px 24px rgba(0,0,0,0.10);
    word-break: break-word;
    line-height: 1.7;
    min-height: 40px;
    min-width: 80px;
    margin-bottom: 4px;
  }
  
  .ai-msg.user .ai-msg-bubble {
    background: #e6f7ff;
  }
  
  .ai-msg-time {
    font-size: 12px;
    color: #c0c4cc;
    position: absolute;
    right: 18px;
    bottom: 6px;
    background: transparent;
    padding-left: 8px;
  }
  
  .ai-drawer-footer {
    background: #f7f8fa;
    border-top: 1px solid #eee;
    padding: 16px;
    display: flex;
    align-items: center; /* 改为居中对齐 */
    gap: 12px;
    box-sizing: border-box;
    width: 100%;
    position: relative;
  }
  
  .ai-input {
    flex: 1;
    min-width: 0;
    max-width: calc(100% - 90px); /* 限制输入框宽度 */
    border-radius: 8px; /* 可选：增加输入框圆角 */
  }
  
  .ai-send-btn {
    height: 40px;
    padding: 0 18px;
    font-size: 15px;
    flex-shrink: 0;
    white-space: nowrap;
    border-radius: 8px; /* 与输入框圆角保持一致 */
    min-width: 70px; /* 确保按钮有足够宽度 */
  }
  
  @media (max-width: 600px) {
    .ai-drawer-footer {
      flex-direction: column;
      align-items: stretch;
      padding: 16px 16px 24px; /* 增加底部内边距 */
    }
  
    .ai-input {
      max-width: 100%; /* 在移动端下占满宽度 */
    }
  
    .ai-send-btn {
      width: 100%;
      margin-top: 8px;
    }
  
    .ai-drawer-body {
      padding: 16px 4px 16px 4px;
    }
    .ai-msg-content {
      max-width: 95vw;
      margin: 0 4px;
    }
  }
  </style>
  