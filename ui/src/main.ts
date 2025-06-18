import './assets/main.css'
import './assets/utilities.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { createPinia } from 'pinia'
import permission from '@/directives/permission'

// 创建应用实例
const app = createApp(App)

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 使用Element Plus
app.use(ElementPlus)

// 使用路由
app.use(router)

// 使用 pinia
app.use(createPinia())

// 注册权限指令
app.directive('permission', permission)

// 挂载应用
app.mount('#app')
