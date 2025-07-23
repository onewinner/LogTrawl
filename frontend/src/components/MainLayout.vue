<template>
  <div
    class="main-layout wails-drop-target"
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
        <AnalysisPage
          v-else-if="appStore.currentView === 'analysis'"
          :file-path="analysisFilePath"
          :file-name="analysisFileName"
        />
      </div>
    </div>

    <!-- 时间线面板 -->
    <TimelinePanel />

    <!-- 全局加载组件 -->
    <GlobalLoading />

    <!-- 文件分片对话框 -->
    <FileSplitterDialog
      v-model="showFileSplitter"
      :file-path="splitterFilePath"
      :file-name="splitterFileName"
      :file-size="splitterFileSize"
      @split-complete="handleSplitComplete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import { useAppStore } from '@/stores/app'
import ToolBar from './ToolBar.vue'
import SideBar from './SideBar.vue'
import WelcomePage from './WelcomePage.vue'
import LogViewer from './LogViewer.vue'
import TimelinePanel from './TimelinePanel.vue'
import GlobalLoading from './GlobalLoading.vue'
import AnalysisPage from './AnalysisPage.vue'
import FileSplitterDialog from './FileSplitterDialog.vue'

const appStore = useAppStore()

// 拖拽状态
const isDragOver = ref(false)
const isWailsEnvironment = ref(false)

// 分析页面状态
const analysisFilePath = ref('')
const analysisFileName = ref('')

// 分片对话框状态
const showFileSplitter = ref(false)
const splitterFilePath = ref('')
const splitterFileName = ref('')
const splitterFileSize = ref(0)

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

// 处理文件拖拽 - 在 Wails 环境下由 OnFileDrop 处理，在浏览器环境下使用前端处理
const handleDrop = async (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false

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

  // 检查文件大小（限制为1GB）
  const maxSize = 1024 * 1024 * 1024 // 1GB
  if (file.size > maxSize) {
    ElMessage({
      message: '文件过大，超过1GB的文件需要使用分片功能处理',
      type: 'warning',
      duration: 3000,
      showClose: true,
      offset: 20,
      customClass: 'message-bottom-right'
    })

    // 触发分片处理流程
    showFileSplitter.value = true
    splitterFilePath.value = file.name // 在浏览器环境中使用文件名
    splitterFileName.value = file.name
    splitterFileSize.value = file.size
    return
  }

  // 在 Wails 环境下，OnFileDrop 会自动处理，这里只需要阻止浏览器默认行为
}





// 处理 Wails 文件拖拽 - 优化版本
const handleWailsFileDrop = async (filePath: string) => {
  console.log('🚀 开始处理 Wails 文件拖拽:', filePath)

  // 提取路径信息
  const pathInfo = getPathInfo(filePath)
  console.log('📁 路径信息:', pathInfo)

  // 检查文件类型
  const fullFileName = pathInfo.Name + '.' + pathInfo.Ext
  console.log('🔍 检查文件类型:', fullFileName)

  if (!isSupportedFile(fullFileName)) {
    console.warn('❌ 不支持的文件类型:', fullFileName)
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

  console.log('✅ 文件类型检查通过，开始处理文件')

  try {
    // 如果有正在进行的加载，先中断
    if (appStore.isGlobalLoading) {
      console.log('🛑 中断之前的文件加载操作')
      appStore.setGlobalLoading(false)
      // 给一点时间让之前的操作清理
      await new Promise(resolve => setTimeout(resolve, 100))
    }

    // 启动加载状态
    appStore.setGlobalLoading(true, '正在检查文件信息...', 10)

    // 先获取文件信息以判断文件大小
    const { GetFileInfo } = await import('wailsjs/go/main/App')
    const fileInfo = await GetFileInfo(filePath)
    const fileSizeMB = fileInfo.size / (1024 * 1024)
    const fileSizeGB = fileInfo.size / (1024 * 1024 * 1024)

    console.log('📊 文件信息:', {
      name: fileInfo.name,
      size: `${fileSizeMB.toFixed(2)} MB`,
      lines: fileInfo.lines || '未知'
    })

    // 检查文件大小是否超过1GB限制
    if (fileSizeGB > 1) {
      appStore.setGlobalLoading(false) // 关闭加载状态

      ElMessage({
        message: `文件大小 ${fileSizeGB.toFixed(2)} GB 超过1GB限制，请使用分片功能处理`,
        type: 'warning',
        duration: 5000,
        showClose: true,
        offset: 20,
        customClass: 'message-bottom-right'
      })

      // 触发分片处理流程
      showFileSplitter.value = true
      splitterFilePath.value = filePath
      splitterFileName.value = fileInfo.name
      splitterFileSize.value = fileInfo.size
      return // 重要：直接返回，不继续执行后续的文件打开逻辑
    }

    // 根据文件大小显示不同的加载信息
    if (fileSizeMB > 200) {
      appStore.setGlobalLoading(true, `正在处理大文件 (${fileSizeMB.toFixed(1)} MB)...`, 20)
    } else if (fileSizeMB > 10) {
      appStore.setGlobalLoading(true, `正在加载文件 (${fileSizeMB.toFixed(1)} MB)...`, 20)
    } else {
      appStore.setGlobalLoading(true, '正在加载文件...', 30)
    }

    // 使用 setTimeout 让 UI 有时间更新
    await new Promise(resolve => setTimeout(resolve, 100))

    // 使用 AppStore 的 openFile 方法
    appStore.updateLoadingProgress(50, '正在读取文件内容...')
    await appStore.openFile(filePath)

    appStore.updateLoadingProgress(90, '正在初始化界面...')

    // 给界面一些时间来渲染
    await new Promise(resolve => setTimeout(resolve, 200))

    console.log('✅ 文件拖拽处理完成')

    // 关闭加载状态
    appStore.setGlobalLoading(false)

    ElMessage({
      message: `文件 ${pathInfo.Name}.${pathInfo.Ext} 已成功打开`,
      type: 'success',
      duration: 3000,
      showClose: true,
      offset: 20,
      customClass: 'message-bottom-right'
    })
  } catch (error) {
    console.error('❌ 文件拖拽处理失败:', error)

    // 关闭加载状态
    appStore.setGlobalLoading(false)

    ElMessage({
      message: `打开文件失败: ${error}`,
      type: 'error',
      duration: 5000,
      showClose: true,
      offset: 20,
      customClass: 'message-bottom-right'
    })
  }
}

// 处理分片完成
const handleSplitComplete = (result) => {
  console.log('分片完成:', result)
  ElMessage.success(`文件已分片为 ${result.totalFiles} 个文件`)

  // 可以在这里提供选择分片文件进行分析的选项
  // 或者自动打开分片文件所在目录
}

// 前端路径信息提取
const getPathInfo = (filePath: string) => {
  const parts = filePath.split(/[\/\\]/)
  const fileName = parts[parts.length - 1]
  const extParts = fileName.split('.')
  const ext = extParts.length > 1 ? extParts[extParts.length - 1] : ''

  return {
    Name: fileName.replace(/\.[^/.]+$/, ""), // 移除扩展名
    Ext: ext.toLowerCase(),
    Path: filePath
  }
}

// Wails 文件拖拽处理

// 生命周期钩子
// 打开分析页面
const openAnalysisPage = (filePath: string, fileName: string) => {
  console.log('📊 打开分析页面:', { filePath, fileName })
  analysisFilePath.value = filePath
  analysisFileName.value = fileName
  appStore.currentView = 'analysis'
}

// 关闭分析页面
const closeAnalysisPage = () => {
  console.log('📊 关闭分析页面')
  appStore.currentView = 'log-viewer'
  analysisFilePath.value = ''
  analysisFileName.value = ''
}

onMounted(async () => {
  console.log('MainLayout 组件已挂载，尝试初始化 Wails 拖拽功能')

  // 监听分析页面事件
  window.addEventListener('openAnalysisPage', (event: any) => {
    const { filePath, fileName } = event.detail
    openAnalysisPage(filePath, fileName)
  })

  window.addEventListener('closeAnalysisPage', () => {
    closeAnalysisPage()
  })

  try {
    // 导入 Wails 运行时
    const runtime = await import('wailsjs/runtime/runtime')
    console.log('Wails 运行时导入成功:', runtime)

    const { OnFileDrop } = runtime

    // 设置 Wails 环境标志
    isWailsEnvironment.value = true
    console.log('Wails 环境已设置为 true')

    // 使用 Wails 提供的 OnFileDrop API
    console.log('正在设置 OnFileDrop 监听器...')

    // 尝试不同的 OnFileDrop 调用方式
    try {
      OnFileDrop((x: number, y: number, paths: string[]) => {
        console.log('🎯 Wails OnFileDrop 触发!', { x, y, paths })

        // 处理每个拖放的文件路径
        paths.forEach(async (filePath) => {
          console.log('处理文件路径:', filePath)
          await handleWailsFileDrop(filePath)
        })
      }, true) // true 表示启用拖放目标检测

      console.log('OnFileDrop 监听器设置成功 (带目标检测)')
    } catch (dropError) {
      console.warn('带目标检测的 OnFileDrop 设置失败，尝试不带目标检测:', dropError)

      // 尝试不带目标检测参数
      OnFileDrop((x: number, y: number, paths: string[]) => {
        console.log('🎯 Wails OnFileDrop 触发! (无目标检测)', { x, y, paths })

        // 处理每个拖放的文件路径
        paths.forEach(async (filePath) => {
          console.log('处理文件路径:', filePath)
          await handleWailsFileDrop(filePath)
        })
      }, false)

      console.log('OnFileDrop 监听器设置成功 (无目标检测)')
    }

    console.log('OnFileDrop 监听器设置完成')

  } catch (error) {
    console.error('Wails 运行时初始化失败:', error)
    // Wails 运行时不可用，保持默认的 false 值
    isWailsEnvironment.value = false
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

/* Wails 拖拽目标标记 */
.wails-drop-target {
  --wails-drop-target: drop;
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
  transition: margin-right 0.3s ease;
}

/* 当时间线面板显示时，为主内容区域添加右边距 */
.main-layout:has(.timeline-panel-visible) .content-area {
  margin-right: 350px; /* 调整为时间线面板的实际宽度 */
}

/* 当有时间线条目但面板关闭时，为侧边栏留出空间 */
.main-layout:has(.timeline-sidebar:not(.timeline-sidebar-hidden)) .content-area {
  margin-right: 40px;
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
