<template>
  <div
    class="sidebar"
    :class="{ collapsed: appStore.sidebarCollapsed }"
    :style="!appStore.sidebarCollapsed ? { width: sidebarWidth + 'px' } : {}"
  >
    <!-- 标题区域 -->
    <div class="sidebar-header">
      <div class="app-title">
        <el-icon class="app-icon"><Document /></el-icon>
        <span class="title-text">LogTrawl</span>
      </div>
    </div>
    
    <!-- 启动区域 -->
    <div class="section">
      <div class="section-title">启动</div>
      <div class="action-list">
        <div class="action-item" @click="openFile">
          <el-icon class="action-icon"><FolderOpened /></el-icon>
          <span class="action-text">打开文件</span>
        </div>
        <div class="action-item" @click="openFolder">
          <el-icon class="action-icon"><Folder /></el-icon>
          <span class="action-text">打开文件夹</span>
        </div>
        <div class="action-item" @click="openFromClipboard">
          <el-icon class="action-icon"><DocumentCopy /></el-icon>
          <span class="action-text">从剪贴板导入</span>
        </div>
        <div class="action-item disabled">
          <el-icon class="action-icon"><Link /></el-icon>
          <span class="action-text">从URL导入 (暂未开放)</span>
        </div>
        <div class="action-item" @click="openProject">
          <el-icon class="action-icon"><Collection /></el-icon>
          <span class="action-text">打开项目</span>
        </div>
      </div>
    </div>
    
    <!-- 最近打开区域 -->
    <div class="section">
      <div class="section-title">最近打开</div>
      <div class="recent-files">
        <div
          v-for="file in appStore.recentFiles"
          :key="file.path"
          class="recent-file-item"
          @click="openRecentFile(file.path)"
        >
          <el-icon class="file-icon">
            <Document />
          </el-icon>
          <div class="file-info">
            <div class="file-name" :title="file.name">{{ file.name }}</div>
            <div class="file-path" :title="file.path">{{ file.path }}</div>
          </div>
        </div>
        
        <!-- 空状态 -->
        <div v-if="appStore.recentFiles.length === 0" class="empty-state">
          <el-icon class="empty-icon"><Document /></el-icon>
          <span class="empty-text">暂无最近文件</span>
        </div>
      </div>
    </div>
    
    <!-- 已打开文件标签页 -->
    <div class="section" v-if="appStore.hasOpenFiles">
      <div class="section-title">已打开</div>
      <div class="open-files">
        <div 
          v-for="file in appStore.openFiles" 
          :key="file.id"
          class="open-file-item"
          :class="{ active: file.id === appStore.currentFileId }"
          @click="appStore.setCurrentFile(file.id)"
        >
          <el-icon class="file-icon">
            <Document />
          </el-icon>
          <div class="file-info">
            <div class="file-name" :title="file.name">{{ file.name }}</div>
          </div>
          <el-button 
            :icon="Close" 
            size="small" 
            text
            class="close-btn"
            @click.stop="appStore.closeFile(file.id)"
          />
        </div>
      </div>
    </div>

    <!-- 拖拽调整手柄 -->
    <div
      class="resize-handle"
      @mousedown="startResize"
      v-if="!appStore.sidebarCollapsed"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { useAppStore, type LogFile } from '@/stores/app'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Document,
  FolderOpened,
  Folder,
  DocumentCopy,
  Link,
  Collection,
  Close
} from '@element-plus/icons-vue'
import { OpenFileDialog, OpenDirectoryDialog, GetFileInfo, GetFilesInDirectory } from 'wailsjs/go/main/App'
import { onMounted, ref, onUnmounted } from 'vue'

const appStore = useAppStore()

// 侧边栏宽度调整
const sidebarWidth = ref(200) // 默认宽度
const isResizing = ref(false)

const startResize = (e: MouseEvent) => {
  isResizing.value = true
  document.addEventListener('mousemove', handleResize)
  document.addEventListener('mouseup', stopResize)
  e.preventDefault()
}

const handleResize = (e: MouseEvent) => {
  if (!isResizing.value) return

  const newWidth = e.clientX
  // 限制最小和最大宽度
  if (newWidth >= 200 && newWidth <= 500) {
    sidebarWidth.value = newWidth
  }
}

const stopResize = () => {
  isResizing.value = false
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
}

// 清理事件监听器
onUnmounted(() => {
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
})

// 组件挂载时加载最近文件
onMounted(() => {
  appStore.loadRecentFiles()
})

// 方法
const openFile = async () => {
  try {
    const filePath = await OpenFileDialog()
    if (filePath) {
      const fileInfo = await GetFileInfo(filePath)
      if (fileInfo) {
        const logFile: LogFile = {
          id: fileInfo.id,
          name: fileInfo.name,
          path: fileInfo.path,
          size: fileInfo.size,
          lastModified: new Date(fileInfo.lastModified),
          isOpen: true
        }
        appStore.addLogFile(logFile)

        // 自动隐藏侧边栏
        if (!appStore.sidebarCollapsed) {
          appStore.toggleSidebar()
        }

        ElMessage.success(`已打开文件: ${fileInfo.name}`)
      }
    }
  } catch (error) {
    console.error('打开文件失败:', error)
    ElMessage.error('打开文件失败')
  }
}

const openFolder = async () => {
  try {
    const dirPath = await OpenDirectoryDialog()
    if (dirPath) {
      const files = await GetFilesInDirectory(dirPath)
      if (files && files.length > 0) {
        files.forEach(file => {
          const logFile: LogFile = {
            id: file.id,
            name: file.name,
            path: file.path,
            size: file.size,
            lastModified: new Date(file.lastModified),
            isOpen: false
          }
          appStore.addLogFile(logFile)
        })
        ElMessage.success(`已发现 ${files.length} 个日志文件`)
      } else {
        ElMessage.warning('该文件夹中没有找到日志文件')
      }
    }
  } catch (error) {
    console.error('打开文件夹失败:', error)
    ElMessage.error('打开文件夹失败')
  }
}

const openFromClipboard = async () => {
  try {
    // 读取剪切板内容
    const clipboardText = await navigator.clipboard.readText()

    if (!clipboardText || clipboardText.trim() === '') {
      ElMessage.warning('剪切板为空或无法读取')
      return
    }

    // 弹出对话框让用户输入文件名
    const { value: fileName } = await ElMessageBox.prompt(
      '请输入文件名（不需要扩展名）',
      '从剪切板导入',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPattern: /^[^<>:"/\\|?*]+$/,
        inputErrorMessage: '文件名不能包含特殊字符',
        inputValue: `clipboard_${new Date().getTime()}`
      }
    )

    if (!fileName) {
      return
    }

    // 创建临时文件
    const logContent = clipboardText.trim()
    const tempFileName = `${fileName}.log`

    // 创建一个虚拟的文件对象
    const logFile: LogFile = {
      id: `clipboard_${Date.now()}`,
      name: tempFileName,
      path: `temp://${tempFileName}`,
      size: new Blob([logContent]).size,
      lastModified: new Date(),
      isOpen: true,
      content: logContent // 添加内容字段用于临时文件
    }

    // 添加到应用状态
    appStore.addLogFile(logFile)

    // 自动隐藏侧边栏
    if (!appStore.sidebarCollapsed) {
      appStore.toggleSidebar()
    }

    ElMessage.success(`已从剪切板导入: ${tempFileName}`)

  } catch (error) {
    if (error === 'cancel') {
      return // 用户取消
    }
    console.error('从剪切板导入失败:', error)
    ElMessage.error('从剪切板导入失败: ' + (error as Error).message)
  }
}

// const openFromUrl = () => {
//   ElMessage.info('URL导入功能开发中...')
// }

const openProject = async () => {
  try {
    // 创建文件输入元素
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = '.ltproj'
    input.style.display = 'none'

    // 监听文件选择
    input.onchange = async (event) => {
      const file = (event.target as HTMLInputElement).files?.[0]
      if (!file) return

      try {
        const text = await file.text()
        const projectData = JSON.parse(text)

        // 验证项目文件格式
        if (!validateProjectData(projectData)) {
          ElMessage.error('无效的项目文件格式')
          return
        }

        // 加载项目
        await loadProject(projectData)
        ElMessage.success(`项目 "${file.name}" 加载成功`)

      } catch (error) {
        console.error('加载项目失败:', error)
        ElMessage.error('加载项目失败: ' + (error as Error).message)
      } finally {
        document.body.removeChild(input)
      }
    }

    // 触发文件选择
    document.body.appendChild(input)
    input.click()

  } catch (error) {
    console.error('打开项目失败:', error)
    ElMessage.error('打开项目失败')
  }
}

// 验证项目数据格式
const validateProjectData = (data: any): boolean => {
  try {
    return (
      data &&
      typeof data === 'object' &&
      data.version &&
      data.createdAt &&
      data.metadata &&
      Array.isArray(data.openFiles) &&
      data.fileStates &&
      data.globalSettings
    )
  } catch {
    return false
  }
}

// 加载项目
const loadProject = async (projectData: any) => {
  try {
    console.log('📂 开始加载项目:', projectData.metadata)

    // 1. 清空当前状态
    appStore.clearAllFiles()

    // 2. 恢复全局设置
    if (projectData.globalSettings) {
      const settings = projectData.globalSettings
      appStore.showLineNumbers = settings.showLineNumbers ?? true
      appStore.wordWrap = settings.wordWrap ?? true
      appStore.syntaxHighlight = settings.syntaxHighlight ?? true

      // 通知工具栏恢复其他设置（高亮词现在是文件独立的）
      window.dispatchEvent(new CustomEvent('restoreProjectSettings', {
        detail: {
          currentWindow: settings.currentWindow || 'main'
        }
      }))
    }

    // 3. 打开文件
    let activeFilePath = ''
    console.log('📁 准备打开文件:', projectData.openFiles.length, '个文件')

    for (const fileInfo of projectData.openFiles) {
      try {
        console.log('📄 正在打开文件:', fileInfo.path)
        const result = await appStore.openFile(fileInfo.path, fileInfo.encoding)
        console.log('✅ 文件打开成功:', result?.name)

        if (fileInfo.isActive) {
          activeFilePath = fileInfo.path
          console.log('🎯 设置活动文件:', activeFilePath)
        }
      } catch (error) {
        console.warn(`❌ 无法打开文件 ${fileInfo.path}:`, error)
        ElMessage.warning(`文件 ${fileInfo.name} 无法打开，可能已被移动或删除`)
      }
    }

    // 4. 等待所有文件加载完成后再恢复状态
    await new Promise(resolve => setTimeout(resolve, 200))

    // 5. 切换到活动文件
    if (activeFilePath && appStore.openFiles.find(f => f.path === activeFilePath)) {
      const activeFile = appStore.openFiles.find(f => f.path === activeFilePath)
      if (activeFile) {
        console.log('🎯 切换到活动文件:', activeFile.name)
        appStore.setCurrentFile(activeFile.id)

        // 等待文件内容加载
        await new Promise(resolve => setTimeout(resolve, 300))
      }
    } else if (appStore.openFiles.length > 0) {
      // 如果没有指定活动文件，选择第一个文件
      console.log('📄 选择第一个文件作为活动文件:', appStore.openFiles[0].name)
      appStore.setCurrentFile(appStore.openFiles[0].id)
      activeFilePath = appStore.openFiles[0].path

      // 等待文件内容加载
      await new Promise(resolve => setTimeout(resolve, 300))
    }

    // 6. 恢复文件状态
    if (projectData.fileStates && activeFilePath) {
      console.log('📁 开始恢复文件状态...')
      window.dispatchEvent(new CustomEvent('restoreProjectFileStates', {
        detail: {
          fileStates: projectData.fileStates,
          activeFilePath: activeFilePath
        }
      }))

      // 等待状态恢复完成
      await new Promise(resolve => setTimeout(resolve, 500))
    }

    // 7. 强制刷新界面
    console.log('🔄 强制刷新界面状态...')
    window.dispatchEvent(new CustomEvent('forceRefreshUI'))

    console.log('✅ 项目加载完成')

  } catch (error) {
    console.error('加载项目过程中出错:', error)
    throw error
  }
}

const openRecentFile = async (filePath: string) => {
  try {
    const fileInfo = await GetFileInfo(filePath)
    if (fileInfo) {
      const logFile: LogFile = {
        id: fileInfo.id,
        name: fileInfo.name,
        path: fileInfo.path,
        size: fileInfo.size,
        lastModified: new Date(fileInfo.lastModified),
        isOpen: true
      }
      appStore.addLogFile(logFile)

      // 自动隐藏侧边栏
      if (!appStore.sidebarCollapsed) {
        appStore.toggleSidebar()
      }
    } else {
      ElMessage.error('文件不存在或无法访问')
      // 从最近文件列表中移除无效文件
      appStore.removeFromRecentFiles(filePath)
    }
  } catch (error) {
    console.error('打开最近文件失败:', error)
    ElMessage.error('文件可能已被移动或删除')
    // 从最近文件列表中移除无效文件
    appStore.removeFromRecentFiles(filePath)
  }
}
</script>

<style scoped>
.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #f8f9fa;
  border-right: 1px solid #e5e7eb;
  transition: width 0.3s ease;
  overflow: hidden;
  position: relative;
  min-width: 100px;
  max-width: 500px;
  z-index: 500; /* 与时间线侧边栏相同层级 */
}

.sidebar.collapsed {
  width: 0 !important;
  min-width: 0 !important;
  max-width: 0 !important;
  overflow: hidden;
  border-right: none;
}

.resize-handle {
  position: absolute;
  top: 0;
  right: 0;
  width: 4px;
  height: 100%;
  background-color: transparent;
  cursor: col-resize;
  z-index: 510; /* 确保在侧边栏之上 */
}

.resize-handle:hover {
  background-color: #3b82f6;
}

.resize-handle:active {
  background-color: #2563eb;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.app-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-icon {
  color: #f59e0b;
  font-size: 20px;
}

.title-text {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.section {
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 12px;
}

.action-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.action-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.action-item:hover {
  background-color: #e5e7eb;
}

.action-item.disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.action-item.disabled:hover {
  background-color: transparent;
}

.action-icon {
  color: #3b82f6;
  font-size: 16px;
}

.action-item.disabled .action-icon {
  color: #9ca3af;
}

.action-text {
  font-size: 13px;
  color: #374151;
}

.action-item.disabled .action-text {
  color: #9ca3af;
}

.recent-files,
.open-files {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 300px;
  overflow-y: auto;
}

.recent-file-item,
.open-file-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.recent-file-item:hover,
.open-file-item:hover {
  background-color: #e5e7eb;
}

.open-file-item.active {
  background-color: #dbeafe;
  border: 1px solid #3b82f6;
}

.file-icon {
  color: #6b7280;
  font-size: 16px;
  flex-shrink: 0;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-size: 13px;
  color: #1f2937;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-path {
  font-size: 11px;
  color: #6b7280;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 2px;
}

.close-btn {
  opacity: 0;
  transition: opacity 0.2s;
  padding: 4px;
  width: 24px;
  height: 24px;
}

.open-file-item:hover .close-btn {
  opacity: 1;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px 12px;
  color: #9ca3af;
}

.empty-icon {
  font-size: 32px;
}

.empty-text {
  font-size: 12px;
}
</style>
