<template>
  <div class="toolbar">
    <!-- 左侧按钮组 -->
    <div class="toolbar-left">
      <!-- 侧边栏切换 -->
      <el-button 
        :icon="Menu" 
        size="small" 
        text
        @click="appStore.toggleSidebar()"
        title="切换侧边栏"
      />
      
      <!-- 文件操作 -->
      <el-button 
        :icon="FolderOpened" 
        size="small"
        @click="openFile"
        title="打开文件"
      >
        文件
      </el-button>
      
      <el-button 
        :icon="Refresh" 
        size="small"
        @click="refreshFile"
        :disabled="!appStore.currentFile"
        title="刷新文件"
      />
      
      <!-- 分隔符 -->
      <el-divider direction="vertical" />
      
      <!-- 编码选择 -->
      <el-select 
        v-model="encoding" 
        size="small" 
        style="width: 100px"
        title="文件编码"
      >
        <el-option label="UTF-8" value="utf-8" />
        <el-option label="GBK" value="gbk" />
        <el-option label="ASCII" value="ascii" />
      </el-select>
    </div>
    
    <!-- 中间搜索区域 -->
    <div class="toolbar-center">
      <div class="search-container">
        <!-- 搜索输入框 -->
        <el-input
          v-model="searchQuery"
          placeholder="搜索日志内容..."
          size="small"
          clearable
          @keyup.enter="performSearch"
          @clear="clearSearch"
          class="search-input"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        
        <!-- 搜索选项 -->
        <div class="search-options">
          <el-tooltip content="正则表达式" placement="bottom">
            <el-button 
              :type="searchOptions.isRegex ? 'primary' : ''"
              size="small"
              text
              @click="toggleSearchOption('isRegex')"
              class="search-option-btn"
            >
              .*
            </el-button>
          </el-tooltip>
          
          <el-tooltip content="区分大小写" placement="bottom">
            <el-button 
              :type="searchOptions.caseSensitive ? 'primary' : ''"
              size="small"
              text
              @click="toggleSearchOption('caseSensitive')"
              class="search-option-btn"
            >
              Aa
            </el-button>
          </el-tooltip>
          
          <el-tooltip content="全词匹配" placement="bottom">
            <el-button 
              :type="searchOptions.wholeWord ? 'primary' : ''"
              size="small"
              text
              @click="toggleSearchOption('wholeWord')"
              class="search-option-btn"
            >
              Ab
            </el-button>
          </el-tooltip>
        </div>
        
        <!-- 搜索导航 -->
        <div class="search-navigation" v-if="searchResults.total > 0">
          <span class="search-count">{{ searchResults.current }}/{{ searchResults.total }}</span>
          <el-button 
            :icon="ArrowUp" 
            size="small" 
            text
            @click="previousMatch"
            :disabled="searchResults.current <= 1"
          />
          <el-button 
            :icon="ArrowDown" 
            size="small" 
            text
            @click="nextMatch"
            :disabled="searchResults.current >= searchResults.total"
          />
        </div>
      </div>
    </div>
    
    <!-- 右侧按钮组 -->
    <div class="toolbar-right">
      <!-- 过滤器 -->
      <el-button 
        :icon="Filter" 
        size="small"
        @click="showFilterDialog = true"
        title="过滤器"
      />
      
      <!-- 设置 -->
      <el-button 
        :icon="Setting" 
        size="small"
        @click="showSettingsDialog = true"
        title="设置"
      />
      
      <!-- 视图选项 -->
      <el-dropdown @command="handleViewCommand">
        <el-button :icon="View" size="small" title="视图选项" />
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item 
              :icon="appStore.showLineNumbers ? Check : ''"
              command="toggleLineNumbers"
            >
              显示行号
            </el-dropdown-item>
            <el-dropdown-item 
              :icon="appStore.syntaxHighlighting ? Check : ''"
              command="toggleSyntaxHighlighting"
            >
              语法高亮
            </el-dropdown-item>
            <el-dropdown-item 
              :icon="appStore.wordWrap ? Check : ''"
              command="toggleWordWrap"
            >
              自动换行
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      
      <!-- 导出 -->
      <el-button 
        :icon="Download" 
        size="small"
        @click="exportData"
        :disabled="!appStore.currentFile"
        title="导出"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAppStore } from '@/stores/app'
import { ElMessage } from 'element-plus'
import {
  Menu,
  FolderOpened,
  Refresh,
  Search,
  ArrowUp,
  ArrowDown,
  Filter,
  Setting,
  View,
  Download,
  Check
} from '@element-plus/icons-vue'
import { OpenFileDialog, GetFileInfo, SearchInFile, SaveFileDialog, ExportLogLines } from 'wailsjs/go/main/App'

const appStore = useAppStore()

// 搜索相关
const searchQuery = ref('')
const encoding = ref('utf-8')
const showFilterDialog = ref(false)
const showSettingsDialog = ref(false)

// 搜索结果
const searchResults = ref({
  total: 0,
  current: 0
})

// 计算属性
const searchOptions = computed(() => appStore.searchOptions)

// 方法
const openFile = async () => {
  try {
    const filePath = await OpenFileDialog()
    if (filePath) {
      const fileInfo = await GetFileInfo(filePath)
      if (fileInfo) {
        const logFile = {
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

const refreshFile = async () => {
  if (!appStore.currentFile) return

  try {
    const fileInfo = await GetFileInfo(appStore.currentFile.path)
    if (fileInfo) {
      // 更新文件信息
      const updatedFile = {
        ...appStore.currentFile,
        size: fileInfo.size,
        lastModified: new Date(fileInfo.lastModified)
      }
      appStore.addLogFile(updatedFile)
      ElMessage.success('文件已刷新')
    }
  } catch (error) {
    console.error('刷新文件失败:', error)
    ElMessage.error('刷新文件失败')
  }
}

const performSearch = async () => {
  if (!searchQuery.value.trim() || !appStore.currentFile) return

  try {
    appStore.updateSearchOptions({ query: searchQuery.value })
    const results = await SearchInFile(
      appStore.currentFile.path,
      searchQuery.value,
      searchOptions.value.caseSensitive
    )

    if (results && results.length > 0) {
      // 存储搜索结果到 store
      appStore.setSearchResults(results)
      searchResults.value = {
        total: results.length,
        current: 1
      }
      // 跳转到第一个匹配项
      jumpToSearchResult(0)
      ElMessage({
        message: `找到 ${results.length} 个匹配项`,
        type: 'success',
        duration: 2000,
        showClose: true,
        offset: 20,
        customClass: 'message-bottom-right'
      })
    } else {
      appStore.setSearchResults([])
      searchResults.value = { total: 0, current: 0 }
      ElMessage({
        message: '未找到匹配项',
        type: 'info',
        duration: 2000,
        showClose: true,
        offset: 20,
        customClass: 'message-bottom-right'
      })
    }
  } catch (error) {
    console.error('搜索失败:', error)
    ElMessage.error('搜索失败')
  }
}

const clearSearch = () => {
  searchQuery.value = ''
  appStore.updateSearchOptions({ query: '' })
  searchResults.value = { total: 0, current: 0 }
}

const toggleSearchOption = (option: keyof typeof searchOptions.value) => {
  const currentValue = searchOptions.value[option]
  appStore.updateSearchOptions({ [option]: !currentValue })
}

const jumpToSearchResult = (index: number) => {
  if (appStore.searchResults.length > 0 && index >= 0 && index < appStore.searchResults.length) {
    appStore.setCurrentSearchIndex(index)
    // 触发 LogViewer 跳转到对应行
    const lineNumber = appStore.searchResults[index].lineNumber
    // 通过事件总线或者直接调用 LogViewer 的方法来跳转
    window.dispatchEvent(new CustomEvent('jumpToLine', { detail: { lineNumber, isSearchResult: true } }))
  }
}

const previousMatch = () => {
  if (searchResults.value.current > 1) {
    searchResults.value.current--
    jumpToSearchResult(searchResults.value.current - 1)
  }
}

const nextMatch = () => {
  if (searchResults.value.current < searchResults.value.total) {
    searchResults.value.current++
    jumpToSearchResult(searchResults.value.current - 1)
  }
}

const handleViewCommand = (command: string) => {
  switch (command) {
    case 'toggleLineNumbers':
      appStore.showLineNumbers = !appStore.showLineNumbers
      break
    case 'toggleSyntaxHighlighting':
      appStore.syntaxHighlighting = !appStore.syntaxHighlighting
      break
    case 'toggleWordWrap':
      appStore.wordWrap = !appStore.wordWrap
      break
  }
}

const exportData = async () => {
  if (!appStore.currentFile || !appStore.logContent.length) {
    ElMessage.warning('没有可导出的数据')
    return
  }

  try {
    const savePath = await SaveFileDialog()
    if (savePath) {
      await ExportLogLines(appStore.logContent, savePath)
      ElMessage.success('导出成功')
    }
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败')
  }
}
</script>

<style scoped>
.toolbar {
  height: 48px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  background-color: #ffffff;
  border-bottom: 1px solid #e5e7eb;
  gap: 8px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-center {
  flex: 1;
  display: flex;
  justify-content: center;
  max-width: 600px;
  margin: 0 auto;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-container {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.search-input {
  flex: 1;
  max-width: 300px;
}

.search-options {
  display: flex;
  gap: 2px;
}

.search-option-btn {
  width: 32px;
  height: 32px;
  padding: 0;
  font-size: 12px;
  font-weight: bold;
}

.search-navigation {
  display: flex;
  align-items: center;
  gap: 4px;
}

.search-count {
  font-size: 12px;
  color: #6b7280;
  min-width: 40px;
  text-align: center;
}

/* 消息提示位置调整 - 右侧中间 */
:deep(.message-bottom-right) {
  position: fixed !important;
  top: 50% !important;
  bottom: auto !important;
  right: 20px !important;
  left: auto !important;
  transform: translateY(-50%) !important;
  z-index: 9999 !important;
}
</style>
