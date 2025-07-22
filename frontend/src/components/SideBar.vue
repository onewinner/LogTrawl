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
        <div class="action-item" @click="openFromUrl">
          <el-icon class="action-icon"><Link /></el-icon>
          <span class="action-text">从URL导入</span>
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
          :key="file.id"
          class="recent-file-item"
          @click="openRecentFile(file)"
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
import { ElMessage } from 'element-plus'
import {
  Document,
  FolderOpened,
  Folder,
  DocumentCopy,
  Link,
  Collection,
  Close
} from '@element-plus/icons-vue'
import { OpenFileDialog, OpenDirectoryDialog, GetFileInfo, GetFilesInDirectory, GetRecentFiles, AddRecentFile } from 'wailsjs/go/main/App'
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

// 加载最近文件
const loadRecentFiles = async () => {
  try {
    const recentFiles = await GetRecentFiles()
    if (recentFiles && recentFiles.length > 0) {
      // 清空现有的最近文件列表
      appStore.recentFiles.length = 0
      // 添加新的最近文件
      recentFiles.forEach(file => {
        const logFile: LogFile = {
          id: file.id,
          name: file.name,
          path: file.path,
          size: file.size,
          lastModified: new Date(file.lastModified),
          isOpen: false
        }
        appStore.recentFiles.push(logFile)
      })
    }
  } catch (error) {
    console.error('加载最近文件失败:', error)
  }
}

// 组件挂载时加载最近文件
onMounted(() => {
  loadRecentFiles()
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

        // 添加到最近文件列表
        await AddRecentFile(filePath)
        // 重新加载最近文件列表
        await loadRecentFiles()

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

const openFromClipboard = () => {
  ElMessage.info('剪贴板导入功能开发中...')
}

const openFromUrl = () => {
  ElMessage.info('URL导入功能开发中...')
}

const openProject = () => {
  ElMessage.info('项目功能开发中...')
}

const openRecentFile = async (file: LogFile) => {
  try {
    const fileInfo = await GetFileInfo(file.path)
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

      // 更新最近文件列表
      await AddRecentFile(file.path)
      await loadRecentFiles()
    }
  } catch (error) {
    console.error('打开最近文件失败:', error)
    ElMessage.error('文件可能已被移动或删除')
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
  z-index: 10;
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

.action-icon {
  color: #3b82f6;
  font-size: 16px;
}

.action-text {
  font-size: 13px;
  color: #374151;
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
