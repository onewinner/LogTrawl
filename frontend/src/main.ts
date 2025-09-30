import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import './styles/element-plus.css'
import './styles/dark-theme.css'
import { ElMessage } from 'element-plus'

// 创建应用实例
const app = createApp(App)

// 应用主题设置
const savedSettings = localStorage.getItem('logtrawl-settings')
if (savedSettings) {
  try {
    const settings = JSON.parse(savedSettings)
    if (settings.theme === 'dark') {
      document.documentElement.classList.add('dark')
      // 设置窗口主题为深色
      if (window.runtime) {
        window.runtime.WindowSetDarkTheme()
      }
    } else {
      // 设置窗口主题为浅色
      if (window.runtime) {
        window.runtime.WindowSetLightTheme()
      }
    }
  } catch (e) {
    console.error('解析保存的主题设置失败:', e)
    // 默认设置窗口主题为浅色
    if (window.runtime) {
      window.runtime.WindowSetLightTheme()
    }
  }
} else {
  // 默认设置窗口主题为浅色
  if (window.runtime) {
    window.runtime.WindowSetLightTheme()
  }
}

// 全局错误处理
app.config.errorHandler = (error, _instance, info) => {
  console.error('全局错误捕获:', error, info)

  // 如果是取消操作的错误，不显示错误消息
  if (error && typeof error === 'object' && 'message' in error &&
      typeof error.message === 'string' && error.message.includes('操作已取消')) {
    return
  }

  // 显示用户友好的错误消息
  ElMessage.error('应用发生错误，请检查控制台获取详细信息')
}

// 捕获未处理的Promise拒绝
window.addEventListener('unhandledrejection', (event) => {
  console.error('未处理的Promise拒绝:', event.reason)

  // 如果是取消操作的错误，不显示错误消息
  if (event.reason && typeof event.reason === 'object' && 'message' in event.reason &&
      typeof event.reason.message === 'string' && event.reason.message.includes('操作已取消')) {
    event.preventDefault()
    return
  }

  ElMessage.error('操作失败，请重试')
  event.preventDefault()
})

// 创建并使用Pinia状态管理
const pinia = createPinia()
app.use(pinia)

// 使用 Element Plus
app.use(ElementPlus)

// 挂载应用
app.mount('#app')
