<template>
  <div
    class="main-layout"
    @dragover.prevent="handleDragOver"
    @dragleave.prevent="handleDragLeave"
    @drop.prevent="handleDrop"
    :class="{ 'drag-over': isDragOver }"
  >
    <!-- 拖拽覆盖层 -->
    <div v-if="isDragOver" class="drag-overlay">
      <div class="drag-message">
        <el-icon class="drag-icon"><Upload /></el-icon>
        <p class="drag-text">拖拽文件到此处打开</p>
        <p class="drag-hint" v-if="isWailsEnvironment">
          支持 .log, .txt 等文本文件 (原生拖拽)
        </p>
        <p class="drag-hint" v-else>
          支持 .log, .txt 等文本文件 (浏览器模式)
        </p>
      </div>
    </div>

    <!-- 顶部工具栏 -->
    <ToolBar />

    <!-- 主内容区域 -->
    <div class="main-content">
      <!-- 左侧边栏 -->
      <SideBar />

      <!-- 内容区域 -->
      <div class="content-area">
        <WelcomePage v-if="appStore.currentView === 'welcome'" />
        <LogViewer v-else-if="appStore.currentView === 'log-viewer'" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import { useAppStore } from '@/stores/app'
import ToolBar from './ToolBar.vue'
import SideBar from './SideBar.vue'
import WelcomePage from './WelcomePage.vue'
import LogViewer from './LogViewer.vue'

const appStore = useAppStore()

// 拖拽状态
const isDragOver = ref(false)
const isWailsEnvironment = ref(false)
let hasReceivedBackendResponse = false

// 支持的文件类型
const supportedFileTypes = ['.log', '.txt', '.text', '.out', '.err', '.trace']

// 检查文件类型是否支持
const isSupportedFile = (fileName: string): boolean => {
  const extension = fileName.toLowerCase().substring(fileName.lastIndexOf('.'))
  return supportedFileTypes.includes(extension) || !fileName.includes('.')
}

// 拖拽进入
const handleDragOver = (event: DragEvent) => {
  event.preventDefault()
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'copy'
  }
  isDragOver.value = true
}

// 拖拽离开
const handleDragLeave = (event: DragEvent) => {
  event.preventDefault()
  // 只有当离开整个拖拽区域时才隐藏覆盖层
  if (!event.relatedTarget || !(event.currentTarget as Element).contains(event.relatedTarget as Node)) {
    isDragOver.value = false
  }
}

// 处理文件拖拽 - 在 Wails 环境下等待后端处理，在浏览器环境下使用前端处理
const handleDrop = async (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false

  // 重置后端响应标志
  hasReceivedBackendResponse = false

  if (!event.dataTransfer?.files.length) {
    return
  }

  const file = event.dataTransfer.files[0]

  // 检查文件类型
  if (!isSupportedFile(file.name)) {
    ElMessage({
      message: `不支持的文件类型。支持的格式：${supportedFileTypes.join(', ')}`,
      type: 'warning',
      duration: 3000,
      showClose: true,
      offset: 20,
      customClass: 'message-bottom-right'
    })
    return
  }

  // 检查文件大小（限制为100MB）
  const maxSize = 100 * 1024 * 1024 // 100MB
  if (file.size > maxSize) {
    ElMessage({
      message: '文件过大，请选择小于100MB的文件',
      type: 'warning',
      duration: 3000,
      showClose: true,
      offset: 20,
      customClass: 'message-bottom-right'
    })
    return
  }

  if (isWailsEnvironment.value) {
    // 在 Wails 环境下，后端的原生拖拽处理会自动触发
    // 这里我们只需要阻止浏览器的默认行为，等待后端事件
    console.log('等待后端处理文件拖拽:', file.name)

    // 设置一个超时，如果后端没有响应，则回退到前端处理
    setTimeout(() => {
      if (!hasReceivedBackendResponse) {
        console.log('后端处理超时，回退到前端处理')
        handleFileWithFrontend(file)
      }
    }, 2000) // 2秒超时
  } else {
    // 在浏览器环境中，直接使用前端处理
    await handleFileWithFrontend(file)
  }
}

// 前端处理文件的备用方法
const handleFileWithFrontend = async (file: File) => {
  const fileContent = await readFileContent(file)

  // 创建文件对象
  const logFile = {
    id: Date.now().toString(),
    name: file.name,
    path: file.name,
    size: file.size,
    lastModified: new Date(file.lastModified),
    isOpen: true
  }

  // 添加到store
  appStore.addLogFile(logFile)
  appStore.setLogContent(fileContent.split('\n'))
  appStore.setCurrentFile(logFile.id)

  ElMessage({
    message: `成功打开文件：${file.name}`,
    type: 'success',
    duration: 2000,
    showClose: true,
    offset: 20,
    customClass: 'message-bottom-right'
  })
}

// 读取文件内容
const readFileContent = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      resolve(e.target?.result as string || '')
    }
    reader.onerror = () => {
      reject(new Error('文件读取失败'))
    }
    reader.readAsText(file, 'utf-8')
  })
}

// 后端事件监听
let fileDropSuccessListener: any = null
let fileDropErrorListener: any = null

// 处理后端文件拖拽成功事件
const handleFileDropSuccess = (data: any) => {
  hasReceivedBackendResponse = true
  ElMessage({
    message: `成功打开文件：${data.file?.name || '未知文件'}`,
    type: 'success',
    duration: 2000,
    showClose: true,
    offset: 20,
    customClass: 'message-bottom-right'
  })
}

// 处理后端文件拖拽错误事件
const handleFileDropError = (data: any) => {
  hasReceivedBackendResponse = true
  ElMessage({
    message: `打开文件失败：${data.error || '未知错误'}`,
    type: 'error',
    duration: 3000,
    showClose: true,
    offset: 20,
    customClass: 'message-bottom-right'
  })
}

// 生命周期钩子
onMounted(async () => {
  try {
    // 导入 Wails 事件系统
    const { EventsOn } = await import('wailsjs/runtime/runtime')

    // 设置 Wails 环境标志
    isWailsEnvironment.value = true

    // 监听后端文件拖拽事件
    fileDropSuccessListener = EventsOn('file-drop-success', handleFileDropSuccess)
    fileDropErrorListener = EventsOn('file-drop-error', handleFileDropError)
  } catch (error) {
    // Wails 事件系统不可用，保持默认的 false 值
    isWailsEnvironment.value = false
  }
})

onUnmounted(() => {
  // 清理事件监听器
  if (fileDropSuccessListener) {
    fileDropSuccessListener()
  }
  if (fileDropErrorListener) {
    fileDropErrorListener()
  }
})
</script>

<style scoped>
.main-layout {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #ffffff;
  margin: 0;
  padding: 0;
  position: fixed;
  top: 0;
  left: 0;
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}



.content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: #ffffff;
}

/* 拖拽相关样式 */
.main-layout.drag-over {
  position: relative;
}

.drag-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(59, 130, 246, 0.1);
  backdrop-filter: blur(2px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  border: 3px dashed #3b82f6;
  animation: pulse 2s infinite;
}

.drag-message {
  text-align: center;
  padding: 40px;
  background-color: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  border: 2px solid #3b82f6;
}

.drag-icon {
  font-size: 48px;
  color: #3b82f6;
  margin-bottom: 16px;
}

.drag-text {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.drag-hint {
  font-size: 16px;
  color: #6b7280;
  margin: 0;
}

@keyframes pulse {
  0%, 100% {
    border-color: #3b82f6;
    background-color: rgba(59, 130, 246, 0.1);
  }
  50% {
    border-color: #1d4ed8;
    background-color: rgba(59, 130, 246, 0.2);
  }
}
</style>
