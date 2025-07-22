import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import './styles/element-plus.css'

// 创建应用实例
const app = createApp(App)

// 创建并使用Pinia状态管理
const pinia = createPinia()
app.use(pinia)

// 使用 Element Plus
app.use(ElementPlus)

// 挂载应用
app.mount('#app')
