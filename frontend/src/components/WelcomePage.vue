<template>
  <div class="welcome-page">
    <div class="welcome-content">
      <!-- 主标题区域 -->
      <div class="header-section">
        <div class="app-logo">
          <LogIcon :size="48" color="#f59e0b" />
        </div>
        <h1 class="app-title">LogTrawl</h1>
        <p class="app-subtitle">专业的日志分析工具</p>
      </div>
      
      <!-- 快速操作区域 -->
      <div class="quick-actions">
        <div class="action-grid">
          <div class="action-card" @click="openFile">
            <div class="action-icon">
              <el-icon><FolderOpened /></el-icon>
            </div>
            <div class="action-content">
              <h3 class="action-title">打开文件</h3>
              <p class="action-desc">选择本地日志文件进行分析</p>
            </div>
          </div>
          
          <div class="action-card" @click="openFolder">
            <div class="action-icon">
              <el-icon><Folder /></el-icon>
            </div>
            <div class="action-content">
              <h3 class="action-title">打开文件夹</h3>
              <p class="action-desc">批量导入文件夹中的日志文件</p>
            </div>
          </div>
          
          <div class="action-card" @click="openFromUrl">
            <div class="action-icon">
              <el-icon><Link /></el-icon>
            </div>
            <div class="action-content">
              <h3 class="action-title">从URL导入</h3>
              <p class="action-desc">从远程URL获取日志数据</p>
            </div>
          </div>
          
          <div class="action-card" @click="openProject">
            <div class="action-icon">
              <el-icon><Collection /></el-icon>
            </div>
            <div class="action-content">
              <h3 class="action-title">打开项目</h3>
              <p class="action-desc">加载已保存的分析项目</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 功能特性展示 -->
      <div class="features-section">
        <h2 class="features-title">主要功能</h2>
        <div class="features-grid">
          <div class="feature-item">
            <el-icon class="feature-icon"><Search /></el-icon>
            <h4>高级搜索</h4>
            <p>支持正则表达式、大小写敏感、全词匹配等多种搜索模式</p>
          </div>
          
          <div class="feature-item">
            <el-icon class="feature-icon"><Filter /></el-icon>
            <h4>实时过滤</h4>
            <p>创建自定义过滤规则，快速筛选关键日志信息</p>
          </div>
          
          <div class="feature-item">
            <el-icon class="feature-icon"><View /></el-icon>
            <h4>语法高亮</h4>
            <p>智能识别日志格式，提供清晰的语法高亮显示</p>
          </div>
          
          <div class="feature-item">
            <el-icon class="feature-icon"><Clock /></el-icon>
            <h4>时间线追踪</h4>
            <p>按时间顺序追踪重要日志条目，快速定位问题</p>
          </div>
          
          <div class="feature-item">
            <el-icon class="feature-icon"><Files /></el-icon>
            <h4>多文件支持</h4>
            <p>同时打开和分析多个日志文件，提高工作效率</p>
          </div>
          
          <div class="feature-item">
            <el-icon class="feature-icon"><Download /></el-icon>
            <h4>导出功能</h4>
            <p>导出过滤结果和分析报告，便于分享和存档</p>
          </div>
        </div>
      </div>
      
      <!-- 最近文件快速访问 -->
      <div class="recent-section" v-if="appStore.recentFiles.length > 0">
        <h2 class="recent-title">最近使用</h2>
        <div class="recent-files-grid">
          <div 
            v-for="file in appStore.recentFiles.slice(0, 6)" 
            :key="file.id"
            class="recent-file-card"
            @click="openRecentFile(file)"
          >
            <div class="file-icon">
              <LogIcon :size="20" color="#6b7280" />
            </div>
            <div class="file-details">
              <div class="file-name" :title="file.name">{{ file.name }}</div>
              <div class="file-path" :title="file.path">{{ file.path }}</div>
              <div class="file-meta">
                <span class="file-size">{{ formatFileSize(file.size) }}</span>
                <span class="file-date">{{ formatDate(file.lastModified) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAppStore, type LogFile } from '@/stores/app'
import {
  FolderOpened,
  Folder,
  Link,
  Collection,
  Search,
  Filter,
  View,
  Clock,
  Files,
  Download
} from '@element-plus/icons-vue'
import LogIcon from './LogIcon.vue'

const appStore = useAppStore()

// 方法
const openFile = () => {
  console.log('打开文件')
  // TODO: 实现文件选择逻辑
}

const openFolder = () => {
  console.log('打开文件夹')
  // TODO: 实现文件夹选择逻辑
}

const openFromUrl = () => {
  console.log('从URL导入')
  // TODO: 实现URL导入逻辑
}

const openProject = () => {
  console.log('打开项目')
  // TODO: 实现项目打开逻辑
}

const openRecentFile = (file: LogFile) => {
  appStore.addLogFile(file)
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatDate = (date: Date): string => {
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}
</script>

<style scoped>
.welcome-page {
  height: 100%;
  overflow-y: auto;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
}

.welcome-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 24px;
}

.header-section {
  text-align: center;
  margin-bottom: 48px;
}

.app-logo {
  margin-bottom: 16px;
}

.logo-icon {
  font-size: 64px;
  color: #f59e0b;
}

.app-title {
  font-size: 48px;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.app-subtitle {
  font-size: 18px;
  color: #6b7280;
  margin: 0;
}

.quick-actions {
  margin-bottom: 48px;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
}

.action-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.action-icon {
  width: 48px;
  height: 48px;
  background: #dbeafe;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #3b82f6;
  font-size: 24px;
}

.action-content {
  flex: 1;
}

.action-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 4px 0;
}

.action-desc {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

.features-section {
  margin-bottom: 48px;
}

.features-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  text-align: center;
  margin: 0 0 32px 0;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 24px;
}

.feature-item {
  text-align: center;
  padding: 24px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.feature-icon {
  font-size: 32px;
  color: #3b82f6;
  margin-bottom: 12px;
}

.feature-item h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.feature-item p {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
}

.recent-section {
  margin-bottom: 24px;
}

.recent-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 24px 0;
}

.recent-files-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
}

.recent-file-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.2s ease;
}

.recent-file-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.file-icon {
  width: 40px;
  height: 40px;
  background: #f3f4f6;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
  font-size: 20px;
}

.file-details {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.file-path {
  font-size: 12px;
  color: #6b7280;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.file-meta {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: #9ca3af;
}
</style>
