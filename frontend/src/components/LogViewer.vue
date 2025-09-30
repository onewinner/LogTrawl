<template>
  <div class="log-viewer">
    <!-- 文件标签页 -->
    <div class="file-tabs" v-if="appStore.openFiles.length > 0">
      <div 
        v-for="file in appStore.openFiles" 
        :key="file.id"
        class="file-tab"
        :class="{ active: file.id === appStore.currentFileId }"
        @click="appStore.setCurrentFile(file.id)"
      >
        <span class="tab-name">{{ file.name }}</span>
        <el-button
          :icon="Close"
          size="small"
          text
          class="tab-close"
          @click.stop="appStore.closeFile(file.id)"
        />
      </div>
    </div>
    
    <!-- 日志内容区域 -->
    <div class="log-content">
      <!-- 状态栏 -->
      <div class="status-bar">
        <div class="status-left">
          <span class="file-info">
            {{ appStore.currentFile?.name }} 
            <span class="file-size">({{ formatFileSize(appStore.currentFile?.size || 0) }})</span>
          </span>
          <span class="line-count">
            <template v-if="isLargeFile">
              共 {{ totalLines }} 行，已加载 {{ loadedChunks.size }} 块
              <span class="performance-info" v-if="performanceStats.chunksLoaded > 0">
                (平均加载: {{ performanceStats.averageLoadTime.toFixed(1) }}ms/块)
              </span>
            </template>
            <template v-else>
              共 {{ logLines.length }} 行，{{ logEntryCount }} 条日志
            </template>
          </span>
          <span class="encoding">UTF-8</span>
        </div>
        <div class="status-right">
          <span class="system-info" v-if="systemInfo">
            应用CPU: {{ systemInfo.cpuUsage.toFixed(1) }}% |
            应用内存: {{ formatMemoryUsage(systemInfo.memoryUsage) }}% ({{ formatBytes(systemInfo.appMemoryUsage) }})
          </span>
          <span class="system-info-debug" v-else>
            系统信息加载中...
          </span>
          <span class="cursor-position">行 {{ currentLine }}, 列 {{ currentColumn }}</span>
        </div>
      </div>



      <!-- 主日志区域容器 -->
      <div class="main-log-container">
        <!-- 日志显示区域 -->
        <div
          class="log-display"
          ref="logDisplayRef"
          tabindex="0"
          @keydown="handleKeyDown"
          @contextmenu="handleContextMenu"
          @focus="handleMainWindowFocus"
          @click="handleMainWindowFocus"
        >
        <div class="log-lines">
          <!-- 行号列 -->
          <div
            class="line-numbers"
            v-if="appStore.showLineNumbers"
            :class="{
              'focused': focusedWindow === 'main',
              'selected': selectedWindow === 'main'
            }"
          >
            <div
              v-for="(_, index) in visibleLines"
              :key="index + startLine"
              class="line-number-container"
              :class="{
                highlighted: highlightedLines.includes(index + startLine)
              }"
            >
              <div class="line-number-content">
                <span class="line-number-text">{{ visibleLineNumbers[index] }}</span>
              </div>
            </div>
          </div>

          <!-- 内容列 -->
          <div class="log-content-lines">
            <div
              v-for="(line, index) in visibleLines"
              :key="index + startLine"
              class="log-line"
              :class="{
                highlighted: highlightedLines.includes(index + startLine),
                'word-wrap': appStore.wordWrap
              }"
              @click="selectLine(index + startLine)"
            >
              <span
                v-if="appStore.syntaxHighlight"
                v-html="highlightSyntax(line, startLine + index)"
              ></span>
              <span v-else>{{ line }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 共用滚动条 - 根据聚焦窗口控制不同内容 -->
      <div
        class="scroll-info"
        v-if="shouldShowScrollBar"
        @mousedown="startUnifiedScrollDrag"
        ref="scrollBarRef"
      >
        <div
          class="scroll-thumb"
          :style="unifiedScrollThumbStyle"
          @mousedown.stop="startUnifiedThumbDrag"
        ></div>
      </div>
    </div>
    </div> <!-- 结束 main-log-container -->

    <!-- 过滤窗口区域 - 挤压式布局 -->
    <div
      class="filter-windows"
      v-if="filterWindows.length > 0"
      :style="{
        height: filterWindowsHeight + 'px'
      }"
    >
      <!-- 拖拽分隔条 -->
      <div
        class="filter-splitter"
        @mousedown="startFilterSplitterDrag"
      >
        <div class="filter-splitter-line"></div>
      </div>

      <!-- 统一的过滤窗口容器 -->
      <div class="filter-container">
        <!-- 过滤窗口标题栏 - 包含所有标签页 -->
        <div class="filter-window-header">
          <div class="filter-window-tabs">
            <div
              v-for="filterWindow in filterWindows"
              :key="filterWindow.id"
              class="filter-window-tab"
              :class="{ active: activeFilterWindow === filterWindow.id }"
              @click="activeFilterWindow = filterWindow.id"
            >
              <span class="filter-tab-name">{{ filterWindow.name }}</span>
              <span class="filter-tab-count">({{ filterWindow.filteredLines.length }})</span>
              <el-button
                :icon="Close"
                size="small"
                text
                class="filter-tab-close"
                @click.stop="closeFilterWindow(filterWindow.id)"
              />
            </div>
          </div>
          <div class="filter-window-controls">
            <span class="filter-result-count">
              {{ activeFilterWindowData?.filteredLines.length || 0 }} 条结果
            </span>
          </div>
        </div>

        <!-- 过滤窗口内容 - 完全复用主窗口的显示逻辑 -->
        <div class="filter-window-content" v-if="activeFilterWindowData">
          <div
            class="filter-log-display"
            ref="filterLogDisplayRef"
            tabindex="0"
            @keydown="handleFilterKeyDown"
            @focus="handleFilterWindowFocus"
            @click="handleFilterWindowFocus"
          >
            <div class="log-lines">
              <!-- 行号列 -->
              <div
                class="line-numbers"
                v-if="appStore.showLineNumbers"
                :class="{
                  'focused': focusedWindow === 'filter',
                  'selected': selectedWindow === 'filter'
                }"
              >
                <div
                  v-for="(_, index) in visibleFilterLines"
                  :key="index + filterStartLine"
                  class="line-number-container"
                  :class="{
                    highlighted: highlightedFilterLines.includes(index + filterStartLine)
                  }"
                >
                  <div class="line-number-content">
                    <span class="line-number-text">{{ visibleFilterLineNumbers[index] }}</span>
                  </div>
                </div>
              </div>

              <!-- 内容列 -->
              <div class="log-content-lines">
                <div
                  v-for="(line, index) in visibleFilterLines"
                  :key="index + filterStartLine"
                  class="log-line"
                  :class="{
                    highlighted: highlightedFilterLines.includes(index + filterStartLine),
                    'word-wrap': appStore.wordWrap
                  }"
                  @click="selectFilterLine(index + filterStartLine)"
                >
                  <span
                    v-if="appStore.syntaxHighlight"
                    v-html="highlightFilterSyntax(line, activeFilterWindowData.originalLineNumbers[index + filterStartLine])"
                  ></span>
                  <span v-else>{{ line }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 过滤窗口不再有独立滚动条，使用共用滚动条 -->

        </div>
      </div>
    </div>

    <!-- 右键菜单 -->
    <div
      v-if="contextMenu.visible"
      class="context-menu"
      :style="{
        left: contextMenu.x + 'px',
        top: contextMenu.y + 'px'
      }"
      @click.stop
    >
      <div class="context-menu-item" @click="copySelectedText" v-if="contextMenu.hasSelection">
        <el-icon><DocumentCopy /></el-icon>
        <span>复制</span>
      </div>

      <div class="context-menu-item" @click="addHighlight" v-if="contextMenu.hasSelection">
        <el-icon><Edit /></el-icon>
        <span>高亮</span>
        <span class="shortcut">E</span>
      </div>

      <div class="context-menu-item" @click="quickFilter" v-if="contextMenu.hasSelection">
        <el-icon><Filter /></el-icon>
        <span>过滤</span>
        <span class="shortcut">F</span>
      </div>

      <div class="context-menu-item" @click="reverseFilter" v-if="contextMenu.hasSelection">
        <el-icon><Filter /></el-icon>
        <span>反过滤</span>
      </div>

      <div class="context-menu-divider" v-if="contextMenu.hasSelection"></div>

      <div class="context-menu-item" @click="addToTimeline" v-if="contextMenu.focusedLine">
        <el-icon><Clock /></el-icon>
        <span>添加到时间线</span>
        <span class="shortcut">T</span>
      </div>

      <div class="context-menu-item" @click="jumpToFocusedLine" v-if="contextMenu.focusedLine">
        <el-icon><Position /></el-icon>
        <span>跳转该行</span>
      </div>

      <div class="context-menu-divider"></div>

      <div class="context-menu-item" @click="selectAll">
        <el-icon><Select /></el-icon>
        <span>全选</span>
        <span class="shortcut">Ctrl+A</span>
      </div>
    </div>

    <!-- 点击遮罩层关闭菜单 -->
    <div
      v-if="contextMenu.visible"
      class="context-menu-overlay"
      @click="hideContextMenu"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useAppStore } from '@/stores/app'
import { ElMessage } from 'element-plus'
import {
  Close,
  DocumentCopy,
  Edit,
  Filter,
  Clock,
  Position,
  Select
} from '@element-plus/icons-vue'
import { ReadLogFile, GetFileInfo, ReadLogFileChunk } from 'wailsjs/go/main/App'

const appStore = useAppStore()

// 响应式数据
const logDisplayRef = ref<HTMLElement>()
const scrollBarRef = ref<HTMLElement>()
const currentLine = ref(1)
const currentColumn = ref(1)
const startLine = ref(0)
const visibleLineCount = ref(50)

// 滚动条拖拽相关变量已移至共用滚动条函数内部
const highlightedLines = ref<number[]>([])

// 搜索相关 - 这些变量在模板中通过 appStore 使用
// const searchResults = ref<Array<{ lineNumber: number, content: string }>>([])
// const currentSearchIndex = ref(0)

// 系统资源信息
const systemInfo = ref<any>(null)
let systemInfoTimer: number | null = null

// 日志数据 - 使用稀疏数组支持分块加载
const logLines = ref<(string | undefined)[]>([])
const isLoading = ref(false)

// 过滤相关
const filterWindows = ref<Array<{
  id: string
  name: string
  filter: string
  mode: string
  options: {
    isRegex: boolean
    caseSensitive: boolean
    wholeWord: boolean
  }
  height: number
  filteredLines: string[]
  originalLineNumbers: number[] // 保存原始行号映射
}>>([])
const activeFilterWindow = ref('')
const filterWindowsHeight = ref(300) // 过滤窗口区域的总高度

// 过滤窗口的状态变量（完全复用主窗口逻辑）
const filterStartLine = ref(0) // 过滤窗口起始行
const filterVisibleLineCount = ref(100) // 过滤窗口可见行数（增加到100行）
const highlightedFilterLines = ref<number[]>([]) // 过滤窗口高亮的行
const filterLogDisplayRef = ref<HTMLElement>() // 过滤窗口显示区域引用
// 过滤窗口不再有独立滚动条，使用共用滚动条

// 窗口聚焦和选择状态
const focusedWindow = ref<'main' | 'filter'>('main') // 当前聚焦的窗口
const selectedWindow = ref<'main' | 'filter'>('main') // 当前选择的窗口（来自Toolbar）

// 右键菜单状态
const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  hasSelection: false,
  selectedText: '',
  focusedLine: null as any
})

// 大文件优化相关
const isLargeFile = ref(false) // 是否为大文件
const totalLines = ref(0) // 文件总行数
const loadedChunks = ref(new Set<number>()) // 已加载的块
const chunkSize = 1000 // 每块行数
const loadBuffer = 2 // 预加载缓冲区（前后各2块）
const isLoadingChunk = ref(false) // 是否正在加载块
const preloadTimeout = ref<number | null>(null) // 预加载防抖定时器
const performanceStats = ref({
  chunksLoaded: 0,
  totalLoadTime: 0,
  averageLoadTime: 0
}) // 性能统计





// 计算属性 - 支持稀疏数组
const visibleLines = computed(() => {
  const totalLength = isLargeFile.value ? totalLines.value : logLines.value.length
  const end = Math.min(startLine.value + visibleLineCount.value, totalLength)

  if (isLargeFile.value) {
    // 大文件模式：从稀疏数组中获取数据，未加载的显示占位符
    const result = []
    for (let i = startLine.value; i < end; i++) {
      const line = logLines.value[i]
      result.push(line !== undefined ? line : `正在加载第 ${i + 1} 行...`)
    }
    return result
  } else {
    // 普通模式：直接切片
    return logLines.value.slice(startLine.value, end)
  }
})

// 当前活动的过滤窗口数据
const activeFilterWindowData = computed(() => {
  return filterWindows.value.find(w => w.id === activeFilterWindow.value)
})

// 过滤窗口可见行数据（完全复用主窗口逻辑）
const visibleFilterLines = computed(() => {
  if (!activeFilterWindowData.value) return []
  const end = Math.min(filterStartLine.value + filterVisibleLineCount.value, activeFilterWindowData.value.filteredLines.length)
  return activeFilterWindowData.value.filteredLines.slice(filterStartLine.value, end)
})

// 共用滚动条显示条件
const shouldShowScrollBar = computed(() => {
  if (focusedWindow.value === 'filter' && activeFilterWindowData.value) {
    // 过滤窗口聚焦时，检查过滤窗口是否需要滚动条
    return activeFilterWindowData.value.filteredLines.length > filterVisibleLineCount.value
  } else {
    // 主窗口聚焦时，检查主窗口是否需要滚动条
    const totalLength = isLargeFile.value ? totalLines.value : logLines.value.length
    return totalLength > visibleLineCount.value
  }
})

// 共用滚动条样式
const unifiedScrollThumbStyle = computed(() => {
  if (focusedWindow.value === 'filter' && activeFilterWindowData.value) {
    // 过滤窗口模式
    const totalLines = activeFilterWindowData.value.filteredLines.length
    const thumbHeight = Math.max((filterVisibleLineCount.value / totalLines) * 100, 5)
    const thumbTop = (filterStartLine.value / totalLines) * 100

    return {
      height: `${thumbHeight}%`,
      top: `${thumbTop}%`
    }
  } else {
    // 主窗口模式
    const totalLength = isLargeFile.value ? totalLines.value : logLines.value.length
    const thumbHeight = Math.max((visibleLineCount.value / totalLength) * 100, 5)
    const thumbTop = (startLine.value / totalLength) * 100

    return {
      height: `${thumbHeight}%`,
      top: `${thumbTop}%`
    }
  }
})









// 识别日志条目开始的函数
const isLogEntryStart = (line: string): boolean => {
  // 识别常见的日志格式：
  // 1. 包含时间戳的行 [日期时间] 或 日期时间
  // 2. 包含日志级别的行 (INFO, ERROR, WARN, DEBUG等)
  // 3. 包含IP地址的行
  // 4. 包含HTTP状态码的行
  const timePattern = /\[\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2}|\d{4}-\d{2}-\d{2}[\sT]\d{2}:\d{2}:\d{2}/
  const logLevelPattern = /\b(INFO|ERROR|WARN|DEBUG|TRACE|FATAL)\b/i
  const ipPattern = /\b(?:\d{1,3}\.){3}\d{1,3}\b/
  const httpStatusPattern = /\s[1-5]\d{2}\s/

  return timePattern.test(line) || logLevelPattern.test(line) ||
         ipPattern.test(line) || httpStatusPattern.test(line)
}

// 创建行号到日志条目编号的映射
const lineToLogEntryMap = computed(() => {
  // 大文件模式下跳过映射计算，避免稀疏数组问题
  if (isLargeFile.value) {
    return []
  }

  const map: number[] = []
  let currentLogEntry = 0

  for (let i = 0; i < logLines.value.length; i++) {
    const line = logLines.value[i]
    // 跳过未加载的行
    if (line === undefined) {
      map[i] = currentLogEntry || 1
      continue
    }

    if (isLogEntryStart(line)) {
      currentLogEntry++
    }
    map[i] = currentLogEntry || 1
  }

  return map
})

// 计算日志条数
const logEntryCount = computed(() => {
  return lineToLogEntryMap.value.length > 0 ? Math.max(...lineToLogEntryMap.value) : 0
})

// 获取真正的行号（基于日志记录而不是显示行）
const getLogLineNumber = (lineIndex: number): number => {
  // 大文件模式下，直接返回物理行号，避免稀疏数组问题
  if (isLargeFile.value) {
    return lineIndex + 1
  }

  // 小文件模式：如果启用了日志记录模式，返回日志记录编号
  if (lineToLogEntryMap.value.length > 0 && lineIndex < lineToLogEntryMap.value.length) {
    return lineToLogEntryMap.value[lineIndex] || lineIndex + 1
  }

  // 否则返回物理行号
  return lineIndex + 1
}

// 获取可见行的行号数组
const visibleLineNumbers = computed(() => {
  return visibleLines.value.map((_, index) => {
    const actualLineIndex = startLine.value + index
    return getLogLineNumber(actualLineIndex)
  })
})

// 强制刷新行号显示
const refreshLineNumbers = async () => {
  console.log('🔢 强制刷新行号显示:', {
    startLine: startLine.value,
    visibleCount: visibleLineCount.value,
    isLargeFile: isLargeFile.value,
    totalLines: totalLines.value,
    currentVisibleNumbers: visibleLineNumbers.value.slice(0, 3)
  })

  // 等待DOM更新
  await nextTick()

  // 强制重新渲染行号元素
  const lineNumberElements = document.querySelectorAll('.line-number-text')
  lineNumberElements.forEach((element, index) => {
    const actualLineIndex = startLine.value + index
    const correctLineNumber = getLogLineNumber(actualLineIndex)
    if (element.textContent !== correctLineNumber.toString()) {
      element.textContent = correctLineNumber.toString()
      console.log(`🔢 修正行号: 索引${index} -> ${correctLineNumber}`)
    }
  })
}

// 获取过滤窗口可见行的行号数组
const visibleFilterLineNumbers = computed(() => {
  if (!activeFilterWindowData.value) return []

  return visibleFilterLines.value.map((_, index) => {
    const filterLineIndex = filterStartLine.value + index
    const originalLineIndex = activeFilterWindowData.value!.originalLineNumbers[filterLineIndex]
    return getLogLineNumber(originalLineIndex)
  })
})



// 原有的scrollThumbStyle已被unifiedScrollThumbStyle替代

// 方法
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

// 当前聚焦的行信息
const currentFocusedLine = ref<{
  lineIndex: number
  lineNumber: number
  content: string
} | null>(null)

const selectLine = (lineIndex: number) => {
  currentLine.value = lineIndex + 1
  currentColumn.value = 1
  highlightedLines.value = [lineIndex]

  // 更新当前聚焦行信息
  currentFocusedLine.value = {
    lineIndex: lineIndex,
    lineNumber: lineIndex + 1,
    content: logLines.value[lineIndex] || ''
  }

  // 触发聚焦行变化事件
  window.dispatchEvent(new CustomEvent('focusedLineChanged', {
    detail: currentFocusedLine.value
  }))
}

// 过滤窗口选中行
const selectFilterLine = (filterLineIndex: number) => {
  if (!activeFilterWindowData.value) return

  // 获取原始行号
  const originalLineIndex = activeFilterWindowData.value.originalLineNumbers[filterLineIndex]
  const lineContent = activeFilterWindowData.value.filteredLines[filterLineIndex]

  // 更新状态栏显示
  currentLine.value = originalLineIndex + 1
  currentColumn.value = 1

  // 高亮过滤窗口中的选中行
  highlightedFilterLines.value = [filterLineIndex]

  // 更新当前聚焦行信息（使用原始行号）
  currentFocusedLine.value = {
    lineIndex: originalLineIndex,
    lineNumber: originalLineIndex + 1,
    content: lineContent || ''
  }

  // 触发聚焦行变化事件
  window.dispatchEvent(new CustomEvent('focusedLineChanged', {
    detail: currentFocusedLine.value
  }))
}

// 处理键盘事件
const handleKeyDown = (event: KeyboardEvent) => {
  // Ctrl+A 全选
  if (event.ctrlKey && event.key === 'a') {
    event.preventDefault()
    selectAllText()
    return
  }

  // 快捷键：t - 快速添加时间线
  if (event.key === 't' || event.key === 'T') {
    event.preventDefault()
    if (currentFocusedLine.value) {
      // 触发快速添加时间线事件
      window.dispatchEvent(new CustomEvent('quickAddTimeline', {
        detail: currentFocusedLine.value
      }))
    } else {
      ElMessage.warning('请先点击选择一行日志')
    }
    return
  }

  // 快捷键：f - 对选中词过滤
  if (event.key === 'f' || event.key === 'F') {
    event.preventDefault()
    const selection = window.getSelection()
    if (selection && selection.toString().trim()) {
      const selectedText = selection.toString().trim()
      // 触发快速过滤事件
      window.dispatchEvent(new CustomEvent('quickFilter', {
        detail: { text: selectedText }
      }))
    } else {
      ElMessage.warning('请先选中要过滤的文本')
    }
    return
  }
}

// 全选所有可见文本
const selectAllText = () => {
  const logContentLines = document.querySelector('.log-content-lines')
  if (logContentLines) {
    const range = document.createRange()
    range.selectNodeContents(logContentLines)
    const selection = window.getSelection()
    if (selection) {
      selection.removeAllRanges()
      selection.addRange(range)
    }
  }
}

// 处理右键菜单
const handleContextMenu = (event: MouseEvent) => {
  event.preventDefault()

  // 获取选中的文本
  const selection = window.getSelection()
  const selectedText = selection ? selection.toString().trim() : ''

  // 设置菜单位置和状态
  contextMenu.value = {
    visible: true,
    x: event.clientX,
    y: event.clientY,
    hasSelection: selectedText.length > 0,
    selectedText: selectedText,
    focusedLine: currentFocusedLine.value
  }

  // 确保菜单不会超出屏幕边界
  nextTick(() => {
    const menuElement = document.querySelector('.context-menu') as HTMLElement
    if (menuElement) {
      const rect = menuElement.getBoundingClientRect()
      const windowWidth = window.innerWidth
      const windowHeight = window.innerHeight

      if (rect.right > windowWidth) {
        contextMenu.value.x = windowWidth - rect.width - 10
      }
      if (rect.bottom > windowHeight) {
        contextMenu.value.y = windowHeight - rect.height - 10
      }
    }
  })
}

// 隐藏右键菜单
const hideContextMenu = () => {
  contextMenu.value.visible = false
}

// 右键菜单功能函数
const copySelectedText = async () => {
  if (contextMenu.value.selectedText) {
    try {
      await navigator.clipboard.writeText(contextMenu.value.selectedText)
      ElMessage.success('已复制到剪贴板')
    } catch (err) {
      // 降级到传统方法
      const textArea = document.createElement('textarea')
      textArea.value = contextMenu.value.selectedText
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      ElMessage.success('已复制到剪贴板')
    }
  }
  hideContextMenu()
}

const addHighlight = () => {
  if (contextMenu.value.selectedText) {
    // 触发添加高亮事件，复用现有的快捷键逻辑
    window.dispatchEvent(new CustomEvent('addHighlight', {
      detail: { text: contextMenu.value.selectedText }
    }))
  }
  hideContextMenu()
}

const quickFilter = () => {
  if (contextMenu.value.selectedText) {
    // 触发快速过滤事件
    window.dispatchEvent(new CustomEvent('quickFilter', {
      detail: { text: contextMenu.value.selectedText }
    }))
  }
  hideContextMenu()
}

const reverseFilter = () => {
  if (contextMenu.value.selectedText) {
    // 触发反向过滤事件
    window.dispatchEvent(new CustomEvent('quickFilter', {
      detail: {
        text: contextMenu.value.selectedText,
        reverse: true
      }
    }))
  }
  hideContextMenu()
}

const addToTimeline = () => {
  if (contextMenu.value.focusedLine) {
    // 触发快速添加时间线事件
    window.dispatchEvent(new CustomEvent('quickAddTimeline', {
      detail: contextMenu.value.focusedLine.lineNumber ? contextMenu.value.focusedLine : { lineNumber: contextMenu.value.focusedLine }
    }))
  }
  hideContextMenu()
}

const jumpToFocusedLine = () => {
  if (contextMenu.value.focusedLine) {
    // 跳转到指定行
    window.dispatchEvent(new CustomEvent('jumpToLine', {
      detail: {
        lineNumber: contextMenu.value.focusedLine.lineNumber || contextMenu.value.focusedLine,
        isSearchResult: false,
        targetWindow: focusedWindow.value
      }
    }))
  }
  hideContextMenu()
}

const selectAll = () => {
  selectAllText()
  hideContextMenu()
}

const highlightSyntax = (line: string, lineIndex: number): string => {
  // 简化的语法高亮，高亮时间、IP、请求方式、状态码
  let highlighted = line

  // 1. 搜索结果高亮 (优先级最高)
  const searchQuery = appStore.searchOptions.query
  if (searchQuery && searchQuery.trim()) {
    try {
      const regex = appStore.searchOptions.isRegex
        ? new RegExp(searchQuery, appStore.searchOptions.caseSensitive ? 'g' : 'gi')
        : new RegExp(escapeRegExp(searchQuery), appStore.searchOptions.caseSensitive ? 'g' : 'gi')

      // 检查当前行是否是当前搜索结果
      const isCurrentSearchLine = appStore.searchResults.length > 0 &&
        appStore.currentSearchIndex < appStore.searchResults.length &&
        appStore.searchResults[appStore.currentSearchIndex]?.lineNumber === lineIndex + 1

      const highlightClass = isCurrentSearchLine ? 'search-highlight-current' : 'search-highlight'
      highlighted = highlighted.replace(regex, `<span class="${highlightClass}">$&</span>`)
    } catch (error) {
      // 正则表达式错误时忽略
    }
  }

  // 2. 高亮词高亮 (在搜索高亮之后，语法高亮之前)
  if (appStore.highlightWords && appStore.highlightWords.length > 0) {
    appStore.highlightWords.forEach((highlightWord, index) => {
      try {
        // 转义特殊字符，确保作为字面量匹配
        const escapedText = escapeRegExp(highlightWord.text)
        const regex = new RegExp(`\\b${escapedText}\\b`, 'gi')

        // 使用唯一的类名避免冲突
        const className = `highlight-word-${index}`
        highlighted = highlighted.replace(regex, `<span class="${className}" style="background-color: ${highlightWord.color}; color: ${getContrastColor(highlightWord.color)}; padding: 1px 2px; border-radius: 2px; font-weight: 500;">$&</span>`)
      } catch (error) {
        // 忽略正则表达式错误
      }
    })
  }

  // 3. 时间戳高亮 (方括号内的日期时间，更精确的匹配)
  highlighted = highlighted.replace(
    /\[\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2}\s[+-]\d{4}\]/g,
    '<span class="syntax-timestamp">$&</span>'
  )

  // 4. IP地址高亮 (精确匹配IP地址格式，避免在User-Agent中误判)
  highlighted = highlighted.replace(
    /\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(?::\d+)?\b/g,
    '<span class="syntax-ip">$&</span>'
  )

  // 5. HTTP方法高亮 (在引号开头的HTTP方法)
  highlighted = highlighted.replace(
    /"(GET|POST|PUT|DELETE|HEAD|OPTIONS|PATCH)\s/g,
    '"<span class="syntax-method">$1</span> '
  )

  // 6. HTTP状态码高亮 (在HTTP版本后面的3位数字，避免User-Agent中的版本号)
  // 匹配模式：HTTP/1.1" 200 或 HTTP/2" 404 等
  highlighted = highlighted.replace(
    /(HTTP\/[\d.]+"\s)(\d{3})(\s)/g,
    '$1<span class="syntax-status">$2</span>$3'
  )

  return highlighted
}

// 转义正则表达式特殊字符
const escapeRegExp = (string: string): string => {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

// 获取对比色（用于文字颜色）
const getContrastColor = (hexColor: string): string => {
  // 移除#号
  const hex = hexColor.replace('#', '')

  // 转换为RGB
  const r = parseInt(hex.substr(0, 2), 16)
  const g = parseInt(hex.substr(2, 2), 16)
  const b = parseInt(hex.substr(4, 2), 16)

  // 计算亮度
  const brightness = (r * 299 + g * 587 + b * 114) / 1000

  // 返回黑色或白色
  return brightness > 128 ? '#000000' : '#ffffff'
}

// 格式化内存使用率
const formatMemoryUsage = (usage: number): string => {
  return usage.toFixed(1)
}

// 格式化字节大小
const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

// 获取系统资源信息
const fetchSystemInfo = async () => {
  try {
    const { GetSystemInfo } = await import('wailsjs/go/main/App')
    const info = await GetSystemInfo()
    systemInfo.value = info
  } catch (error) {
    // 静默处理错误，不在控制台打印
    // 在开发模式下，如果API不可用，使用模拟数据
    if (import.meta.env.DEV) {
      const appMemory = Math.floor(Math.random() * 200 * 1024 * 1024) + 50 * 1024 * 1024 // 50-250MB
      const totalMemory = 32 * 1024 * 1024 * 1024 // 32GB
      systemInfo.value = {
        cpuUsage: Math.random() * 5 + 0.5, // 应用程序CPU通常较低：0.5-5.5%
        memoryUsage: (appMemory / totalMemory) * 100, // 应用内存占总内存的百分比
        memoryUsed: Math.floor(Math.random() * 16 * 1024 * 1024 * 1024), // 系统已用内存
        memoryTotal: totalMemory,
        memoryAvailable: Math.floor(Math.random() * 16 * 1024 * 1024 * 1024),
        appMemoryUsage: appMemory
      }
    }
  }
}

// 启动系统信息监控
const startSystemInfoMonitoring = () => {
  // 立即获取一次
  fetchSystemInfo()

  // 每3秒更新一次
  systemInfoTimer = window.setInterval(fetchSystemInfo, 3000)
}

// 停止系统信息监控
const stopSystemInfoMonitoring = () => {
  if (systemInfoTimer) {
    clearInterval(systemInfoTimer)
    systemInfoTimer = null
  }
}

const handleScroll = (event: WheelEvent) => {
  event.preventDefault()
  const delta = event.deltaY > 0 ? 3 : -3
  const totalLength = isLargeFile.value ? totalLines.value : logLines.value.length
  const newStartLine = Math.max(0, Math.min(startLine.value + delta, totalLength - visibleLineCount.value))

  if (newStartLine !== startLine.value) {
    startLine.value = newStartLine

    // 大文件模式：触发预加载和内存清理
    if (isLargeFile.value && appStore.currentFile) {
      preloadVisibleChunks(appStore.currentFile.path)

      // 每隔一段时间清理远程块
      if (Math.random() < 0.1) { // 10%的概率触发清理
        cleanupDistantChunks()
      }
    }
  }

  // 滚动后同步行号高度
  if (appStore.wordWrap) {
    nextTick(() => {
      syncMainWindowLineHeights()
    })
  }
}

// 旧的滚动条函数已被共用滚动条函数替代

// 过滤窗口滚动函数
const handleFilterScroll = (event: WheelEvent) => {
  if (!activeFilterWindowData.value) return
  event.preventDefault()
  const delta = event.deltaY > 0 ? 3 : -3
  const newStartLine = Math.max(0, Math.min(filterStartLine.value + delta, activeFilterWindowData.value.filteredLines.length - filterVisibleLineCount.value))
  filterStartLine.value = newStartLine

  // 滚动后同步行号高度
  if (appStore.wordWrap) {
    nextTick(() => {
      syncFilterWindowLineHeights()
    })
  }
}

// 共用滚动条点击处理
const startUnifiedScrollDrag = (event: MouseEvent) => {
  if (!scrollBarRef.value) return

  const rect = scrollBarRef.value.getBoundingClientRect()
  const clickY = event.clientY - rect.top
  const scrollBarHeight = rect.height

  if (focusedWindow.value === 'filter' && activeFilterWindowData.value) {
    // 过滤窗口模式
    const totalLines = activeFilterWindowData.value.filteredLines.length
    const maxScrollLines = totalLines - filterVisibleLineCount.value
    const targetLine = (clickY / scrollBarHeight) * totalLines
    const newStartLine = Math.max(0, Math.min(targetLine, maxScrollLines))
    filterStartLine.value = Math.round(newStartLine)

    console.log('🎚️ 共用滚动条控制过滤窗口:', { targetLine: newStartLine, totalLines })
  } else {
    // 主窗口模式
    const totalLength = isLargeFile.value ? totalLines.value : logLines.value.length
    const maxScrollLines = totalLength - visibleLineCount.value
    const targetLine = (clickY / scrollBarHeight) * totalLength
    const newStartLine = Math.max(0, Math.min(targetLine, maxScrollLines))
    startLine.value = Math.round(newStartLine)

    console.log('🎚️ 共用滚动条控制主窗口:', { targetLine: newStartLine, totalLength })

    // 大文件模式：触发预加载
    if (isLargeFile.value && appStore.currentFile) {
      preloadVisibleChunks(appStore.currentFile.path)
    }
  }
}

// 共用滚动条拖拽处理
const startUnifiedThumbDrag = (event: MouseEvent) => {
  const startY = event.clientY

  if (focusedWindow.value === 'filter' && activeFilterWindowData.value) {
    // 过滤窗口模式
    const startScrollLine = filterStartLine.value
    const totalLines = activeFilterWindowData.value.filteredLines.length

    const handleThumbDrag = (e: MouseEvent) => {
      const deltaY = e.clientY - startY
      const scrollBarHeight = scrollBarRef.value?.getBoundingClientRect().height || 100
      const deltaLines = (deltaY / scrollBarHeight) * totalLines
      const newStartLine = Math.max(0, Math.min(startScrollLine + deltaLines, totalLines - filterVisibleLineCount.value))
      filterStartLine.value = Math.round(newStartLine)
    }

    const stopThumbDrag = () => {
      document.removeEventListener('mousemove', handleThumbDrag)
      document.removeEventListener('mouseup', stopThumbDrag)
    }

    document.addEventListener('mousemove', handleThumbDrag)
    document.addEventListener('mouseup', stopThumbDrag)
  } else {
    // 主窗口模式
    const startScrollLine = startLine.value
    const totalLength = isLargeFile.value ? totalLines.value : logLines.value.length

    const handleThumbDrag = (e: MouseEvent) => {
      const deltaY = e.clientY - startY
      const scrollBarHeight = scrollBarRef.value?.getBoundingClientRect().height || 100
      const deltaLines = (deltaY / scrollBarHeight) * totalLength
      const newStartLine = Math.max(0, Math.min(startScrollLine + deltaLines, totalLength - visibleLineCount.value))
      startLine.value = Math.round(newStartLine)

      // 大文件模式：触发预加载
      if (isLargeFile.value && appStore.currentFile) {
        preloadVisibleChunks(appStore.currentFile.path)
      }
    }

    const stopThumbDrag = () => {
      document.removeEventListener('mousemove', handleThumbDrag)
      document.removeEventListener('mouseup', stopThumbDrag)
    }

    document.addEventListener('mousemove', handleThumbDrag)
    document.addEventListener('mouseup', stopThumbDrag)
  }

  event.preventDefault()
}

const handleResize = () => {
  if (logDisplayRef.value) {
    const height = logDisplayRef.value.clientHeight
    const lineHeight = 20 // 假设每行高度为20px
    visibleLineCount.value = Math.floor(height / lineHeight)
  }
}

// 同步行号和内容行高度
const syncLineHeights = () => {
  nextTick(() => {
    if (!appStore.showLineNumbers) return

    // 主窗口行号同步
    syncMainWindowLineHeights()

    // 过滤窗口行号同步
    syncFilterWindowLineHeights()
  })
}

// 同步主窗口行号高度
const syncMainWindowLineHeights = () => {
  const mainLineNumbers = document.querySelectorAll('.log-display .line-number-container') as NodeListOf<HTMLElement>
  const mainLogLines = document.querySelectorAll('.log-display .log-line') as NodeListOf<HTMLElement>

  mainLineNumbers.forEach((lineNumber, index) => {
    const logLine = mainLogLines[index]
    if (logLine) {
      if (appStore.wordWrap) {
        // 自动换行模式：行号高度跟随内容高度
        const logLineHeight = logLine.offsetHeight
        lineNumber.style.height = `${logLineHeight}px`
        lineNumber.style.minHeight = '20px'
      } else {
        // 固定高度模式：重置为默认高度
        lineNumber.style.height = '20px'
        lineNumber.style.minHeight = '20px'
      }
    }
  })
}

// 同步过滤窗口行号高度
const syncFilterWindowLineHeights = () => {
  const filterLineNumbers = document.querySelectorAll('.filter-log-display .line-number-container') as NodeListOf<HTMLElement>
  const filterLogLines = document.querySelectorAll('.filter-log-display .log-line') as NodeListOf<HTMLElement>

  filterLineNumbers.forEach((lineNumber, index) => {
    const logLine = filterLogLines[index]
    if (logLine) {
      if (appStore.wordWrap) {
        // 自动换行模式：行号高度跟随内容高度
        const logLineHeight = logLine.offsetHeight
        lineNumber.style.height = `${logLineHeight}px`
        lineNumber.style.minHeight = '20px'
      } else {
        // 固定高度模式：重置为默认高度
        lineNumber.style.height = '20px'
        lineNumber.style.minHeight = '20px'
      }
    }
  })
}

// 清理之前文件的状态
const cleanupPreviousFileState = async () => {
  console.log('🧹 清理之前文件的状态:', {
    isLargeFile: isLargeFile.value,
    isLoadingChunk: isLoadingChunk.value,
    loadedChunks: loadedChunks.value.size,
    logLinesLength: logLines.value.length
  })

  // 中断正在进行的块加载
  if (isLoadingChunk.value) {
    console.log('⏹️ 中断正在进行的块加载')
    isLoadingChunk.value = false
  }

  // 清理大文件状态
  if (isLargeFile.value) {
    console.log('🗑️ 清理大文件状态')
    isLargeFile.value = false
    totalLines.value = 0
    loadedChunks.value.clear()

    // 清理性能统计
    performanceStats.value = {
      chunksLoaded: 0,
      totalLoadTime: 0,
      averageLoadTime: 0
    }
  }

  // 清空日志内容
  logLines.value = []

  // 重置滚动位置
  startLine.value = 0

  // 清理过滤窗口状态
  filterWindows.value = []
  activeFilterWindow.value = ''
  filterStartLine.value = 0

  // 清理高亮状态
  highlightedLines.value = []
  highlightedFilterLines.value = []

  // 等待DOM更新
  await nextTick()

  console.log('✅ 文件状态清理完成')
}

// 检查文件大小并决定加载策略
const loadLogFile = async (filePath: string) => {
  if (!filePath) return

  isLoading.value = true
  try {
    // 检查是否为临时文件（剪切板导入）
    if (filePath.startsWith('temp://')) {
      console.log('📋 检测到临时文件，直接加载内容')
      appStore.updateLoadingProgress(70, '正在加载剪切板内容...')
      await loadTempFile()
      appStore.updateLoadingProgress(85, '正在初始化显示...')
      return
    }

    // 先获取文件信息
    appStore.updateLoadingProgress(60, '正在分析文件结构...')
    const fileInfo = await GetFileInfo(filePath)
    const fileSizeMB = fileInfo.size / (1024 * 1024)

    console.log('📁 文件信息:', {
      path: filePath.split('/').pop(),
      size: `${fileSizeMB.toFixed(2)} MB`,
      lines: fileInfo.lines || '未知'
    })

    // 判断是否为大文件（超过200MB或10万行）
    if (fileSizeMB > 200 || (fileInfo.lines && fileInfo.lines > 100000)) {
      console.log('📊 检测到大文件，启用分块加载模式')
      appStore.updateLoadingProgress(70, `正在加载大文件 (${fileSizeMB.toFixed(1)} MB)...`)
      await loadLargeFile(filePath, fileInfo)
    } else {
      console.log('📄 普通文件，使用标准加载模式')
      appStore.updateLoadingProgress(70, '正在读取文件内容...')
      await loadSmallFile(filePath)
    }

    appStore.updateLoadingProgress(85, '正在初始化显示...')
  } catch (error) {
    console.error('加载日志文件失败:', error)
    ElMessage.error('加载日志文件失败')
    logLines.value = []
  } finally {
    isLoading.value = false
  }
}

// 加载临时文件（剪切板导入）
const loadTempFile = async () => {
  const currentFile = appStore.currentFile
  if (!currentFile || !currentFile.content) {
    throw new Error('临时文件内容不存在')
  }

  const lines = currentFile.content.split('\n')
  isLargeFile.value = false
  totalLines.value = lines.length
  logLines.value = lines
  appStore.setLogContent(lines)
  ElMessage.success(`已加载剪切板内容: ${lines.length} 行`)

  syncLineHeights()
}

// 加载小文件（原有逻辑）
const loadSmallFile = async (filePath: string) => {
  const content = await ReadLogFile(filePath)
  if (content && content.lines) {
    isLargeFile.value = false
    totalLines.value = content.lines.length
    logLines.value = content.lines
    appStore.setLogContent(content.lines)
    ElMessage.success(`已加载 ${content.total} 行日志`)

    syncLineHeights()
  }
}

// 加载大文件（分块加载）
const loadLargeFile = async (filePath: string, fileInfo: any) => {
  isLargeFile.value = true
  totalLines.value = fileInfo.lines || 0

  // 初始化稀疏数组
  logLines.value = new Array(totalLines.value)
  loadedChunks.value.clear()

  // 先加载第一块数据
  await loadChunk(filePath, 0)

  ElMessage.success(`大文件模式：共 ${totalLines.value} 行，已加载首屏数据`)

  syncLineHeights()
}

// 加载指定块的数据
const loadChunk = async (filePath: string, chunkIndex: number) => {
  if (loadedChunks.value.has(chunkIndex) || isLoadingChunk.value) {
    return
  }

  // 检查是否应该中断加载（文件已切换）
  if (!appStore.currentFile || appStore.currentFile.path !== filePath) {
    console.log(`⏹️ 文件已切换，中断块 ${chunkIndex} 的加载`)
    return
  }

  isLoadingChunk.value = true
  const startTime = performance.now()

  try {
    const startLine = chunkIndex * chunkSize
    const endLine = Math.min(startLine + chunkSize, totalLines.value)

    console.log(`📦 加载块 ${chunkIndex}: 行 ${startLine}-${endLine}`)

    // 再次检查是否应该中断
    if (!appStore.currentFile || appStore.currentFile.path !== filePath) {
      console.log(`⏹️ 加载过程中文件已切换，中断块 ${chunkIndex} 的加载`)
      return
    }

    // 这里需要后端支持分块读取API
    const chunkData = await ReadLogFileChunk(filePath, startLine, endLine)

    // 加载完成后再次检查文件是否还是当前文件
    if (!appStore.currentFile || appStore.currentFile.path !== filePath) {
      console.log(`⏹️ 加载完成后文件已切换，丢弃块 ${chunkIndex} 的数据`)
      return
    }

    if (chunkData && chunkData.lines) {
      // 将数据填入稀疏数组
      for (let i = 0; i < chunkData.lines.length; i++) {
        logLines.value[startLine + i] = chunkData.lines[i]
      }

      loadedChunks.value.add(chunkIndex)

      // 更新性能统计
      const loadTime = performance.now() - startTime
      performanceStats.value.chunksLoaded++
      performanceStats.value.totalLoadTime += loadTime
      performanceStats.value.averageLoadTime = performanceStats.value.totalLoadTime / performanceStats.value.chunksLoaded

      console.log(`✅ 块 ${chunkIndex} 加载完成 (${loadTime.toFixed(2)}ms)`)
    }
  } catch (error) {
    console.error(`❌ 加载块 ${chunkIndex} 失败:`, error)
  } finally {
    isLoadingChunk.value = false
  }
}

// 预加载可见区域周围的块（带防抖）
const preloadVisibleChunks = (filePath: string) => {
  if (!isLargeFile.value || !filePath) return

  // 清除之前的定时器
  if (preloadTimeout.value) {
    clearTimeout(preloadTimeout.value)
  }

  // 防抖：200ms后执行预加载
  preloadTimeout.value = setTimeout(async () => {
    // 检查文件是否还是当前文件
    if (!appStore.currentFile || appStore.currentFile.path !== filePath) {
      console.log('⏹️ 文件已切换，中断预加载')
      return
    }

    const visibleStartLine = startLine.value
    const visibleEndLine = startLine.value + visibleLineCount.value

    const startChunk = Math.floor(visibleStartLine / chunkSize)
    const endChunk = Math.ceil(visibleEndLine / chunkSize)

    console.log(`🔄 预加载块范围: ${startChunk}-${endChunk}`)

    // 预加载当前可见区域前后的块
    for (let chunk = Math.max(0, startChunk - loadBuffer);
         chunk <= Math.min(Math.ceil(totalLines.value / chunkSize) - 1, endChunk + loadBuffer);
         chunk++) {
      // 再次检查文件是否还是当前文件
      if (!appStore.currentFile || appStore.currentFile.path !== filePath) {
        console.log('⏹️ 预加载过程中文件已切换，中断')
        return
      }

      if (!loadedChunks.value.has(chunk)) {
        await loadChunk(filePath, chunk)
      }
    }
  }, 200)
}

// 清理远离可见区域的块以节省内存
const cleanupDistantChunks = () => {
  if (!isLargeFile.value) return

  const visibleStartLine = startLine.value
  const visibleEndLine = startLine.value + visibleLineCount.value

  const startChunk = Math.floor(visibleStartLine / chunkSize)
  const endChunk = Math.ceil(visibleEndLine / chunkSize)

  const keepBuffer = loadBuffer * 3 // 保留更大的缓冲区

  loadedChunks.value.forEach(chunkIndex => {
    // 如果块距离可见区域太远，清理它
    if (chunkIndex < startChunk - keepBuffer || chunkIndex > endChunk + keepBuffer) {
      console.log(`🧹 清理远程块: ${chunkIndex}`)

      // 清理稀疏数组中的数据
      const startLine = chunkIndex * chunkSize
      const endLine = Math.min(startLine + chunkSize, totalLines.value)
      for (let i = startLine; i < endLine; i++) {
        logLines.value[i] = undefined
      }

      loadedChunks.value.delete(chunkIndex)
    }
  })
}

// 监听当前文件变化
watch(() => appStore.currentFile, async (newFile, oldFile) => {
  console.log('📁 LogViewer: 当前文件变化', {
    oldFile: oldFile?.name,
    newFile: newFile?.name,
    isLargeFile: isLargeFile.value,
    loadingChunk: isLoadingChunk.value
  })

  // 立即清理之前文件的状态，避免状态混乱
  await cleanupPreviousFileState()

  if (newFile) {
    await loadLogFile(newFile.path)
  } else {
    // 清空日志内容
    logLines.value = []
    isLargeFile.value = false
    totalLines.value = 0
    loadedChunks.value.clear()
  }
}, { immediate: true })

// 监听换行模式变化
watch(() => appStore.wordWrap, () => {
  syncLineHeights()
})

// 监听行号显示变化
watch(() => appStore.showLineNumbers, () => {
  syncLineHeights()
})

// 跳转到指定行 - 支持大文件
const jumpToLine = async (lineNumber: number, isSearchResult: boolean = false) => {
  const targetLine = Math.max(0, lineNumber - 1) // 转换为0基索引
  const totalLength = isLargeFile.value ? totalLines.value : logLines.value.length
  const maxStartLine = Math.max(0, totalLength - visibleLineCount.value)

  // 计算最佳的起始行，让目标行显示在视口中间
  const idealStartLine = Math.max(0, targetLine - Math.floor(visibleLineCount.value / 2))
  startLine.value = Math.min(idealStartLine, maxStartLine)

  // 大文件模式：确保目标行已加载
  if (isLargeFile.value && appStore.currentFile) {
    const targetChunk = Math.floor(targetLine / chunkSize)
    if (!loadedChunks.value.has(targetChunk)) {
      console.log(`🎯 跳转到行 ${lineNumber}，需要加载块 ${targetChunk}`)
      await loadChunk(appStore.currentFile.path, targetChunk)
    }

    // 预加载周围的块
    preloadVisibleChunks(appStore.currentFile.path)
  }

  // 如果是搜索结果跳转，不需要额外的行高亮，因为搜索高亮已经足够明显
  if (!isSearchResult) {
    // 高亮目标行
    highlightedLines.value = [targetLine]

    // 3秒后清除高亮
    setTimeout(() => {
      highlightedLines.value = []
    }, 3000)
  }

  // 强制更新行号显示和行高度同步
  await nextTick()

  // 在大文件模式下强制刷新行号显示
  if (isLargeFile.value) {
    refreshLineNumbers()
  }

  syncLineHeights()

  console.log('✅ 跳转完成:', {
    targetLine: lineNumber,
    startLine: startLine.value,
    visibleRange: `${startLine.value + 1}-${startLine.value + visibleLineCount.value}`,
    visibleLineNumbers: visibleLineNumbers.value.slice(0, 5) // 显示前5个行号作为调试
  })
}

// 初始化过滤窗口位置
const initializeFilterWindowPosition = () => {
  // 默认高度为300px
  filterWindowsHeight.value = 300
}

// 监听活动过滤窗口变化
watch(activeFilterWindow, () => {
  // 切换过滤窗口时重置状态
  filterStartLine.value = 0
  highlightedFilterLines.value = []

  // 通知状态变化
  notifyFilterWindowsUpdated()
})

// 监听过滤窗口高度变化
watch(filterWindowsHeight, () => {
  // 通知状态变化
  notifyFilterWindowsUpdated()
})

// 监听过滤窗口滚动位置变化
watch(filterStartLine, () => {
  // 延迟通知，避免频繁更新
  clearTimeout(scrollUpdateTimeout.value)
  scrollUpdateTimeout.value = window.setTimeout(() => {
    notifyFilterWindowsUpdated()
  }, 200)
})

// 滚动更新防抖
const scrollUpdateTimeout = ref<number | null>(null)

// 生命周期
onMounted(() => {
  nextTick(() => {
    handleResize()
    if (logDisplayRef.value) {
      logDisplayRef.value.addEventListener('wheel', handleScroll, { passive: false })
    }
    if (filterLogDisplayRef.value) {
      filterLogDisplayRef.value.addEventListener('wheel', handleFilterScroll, { passive: false })
    }
    window.addEventListener('resize', handleResize)

    // 监听全局点击事件来关闭右键菜单
    document.addEventListener('click', hideContextMenu)

    // 监听跳转到行的事件
    window.addEventListener('jumpToLine', (event: any) => {
      const { lineNumber, isSearchResult, targetWindow, requireLoading } = event.detail
      console.log('🎯 收到跳转请求:', {
        lineNumber,
        isSearchResult,
        targetWindow,
        requireLoading,
        currentSelectedWindow: selectedWindow.value
      })

      // 根据目标窗口执行跳转
      if (targetWindow === 'filter' && filterWindows.value.length > 0) {
        jumpToFilterLine(lineNumber, isSearchResult)
      } else {
        // 如果明确要求加载支持或者是大文件模式，使用带加载的跳转
        if (requireLoading || isLargeFile.value) {
          jumpToLineWithLoading(lineNumber, isSearchResult)
        } else {
          jumpToLine(lineNumber, isSearchResult)
        }
      }
    })

    // 监听跳转到尾行的事件
    window.addEventListener('jumpToLastLine', (event: any) => {
      const { targetWindow } = event.detail
      console.log('🎯 收到跳转到尾行请求:', { targetWindow })

      if (targetWindow === 'filter' && filterWindows.value.length > 0) {
        jumpToFilterLastLine()
      } else {
        jumpToLastLineWithLoading()
      }
    })

    // 监听获取当前行的事件
    window.addEventListener('getCurrentLine', () => {
      // 计算当前视口中间的行号
      const middleLine = startLine.value + Math.floor(visibleLineCount.value / 2) + 1
      window.dispatchEvent(new CustomEvent('currentLineResponse', { detail: { lineNumber: middleLine } }))
    })

    // 监听应用过滤的事件
    window.addEventListener('applyFilter', handleApplyFilter)

    // 监听获取过滤窗口列表的事件
    window.addEventListener('getFilterWindows', handleGetFilterWindows)

    // 监听恢复过滤窗口状态事件
    window.addEventListener('restoreFilterWindows', handleRestoreFilterWindows)

    // 监听高亮词更新事件
    window.addEventListener('highlightWordsUpdated', handleHighlightWordsUpdated)

    // 监听窗口选择变化事件
    window.addEventListener('windowChanged', handleWindowChanged)

    // 监听执行搜索的事件
    const handlePerformSearch = (event: any) => {
      const { results, targetWindow } = event.detail
      console.log('🔍 LogViewer: 收到搜索结果:', {
        resultsCount: results.length,
        targetWindow
      })
      
      // 这里可以添加对搜索结果的额外处理逻辑
      // 例如：高亮显示、统计信息等
    }
    window.addEventListener('performSearch', handlePerformSearch)

    // 启动系统信息监控
    startSystemInfoMonitoring()

    // 初始化过滤窗口位置
    initializeFilterWindowPosition()
  })
})

onUnmounted(() => {
  // 清理定时器
  if (scrollUpdateTimeout.value) {
    clearTimeout(scrollUpdateTimeout.value)
  }

  if (logDisplayRef.value) {
    logDisplayRef.value.removeEventListener('wheel', handleScroll)
  }
  if (filterLogDisplayRef.value) {
    filterLogDisplayRef.value.removeEventListener('wheel', handleFilterScroll)
  }
  window.removeEventListener('resize', handleResize)
  document.removeEventListener('click', hideContextMenu)

  // 拖拽事件监听器现在在共用滚动条函数内部管理

  // 清理跳转事件监听器
  window.removeEventListener('jumpToLine', jumpToLine as any)

  // 清理过滤事件监听器
  window.removeEventListener('applyFilter', handleApplyFilter)
  window.removeEventListener('getFilterWindows', handleGetFilterWindows)
  window.removeEventListener('restoreFilterWindows', handleRestoreFilterWindows)
  window.removeEventListener('highlightWordsUpdated', handleHighlightWordsUpdated)
  window.removeEventListener('windowChanged', handleWindowChanged)
  
  // 清理搜索事件监听器
  // window.removeEventListener('performSearch', handlePerformSearch)

  // 停止系统信息监控
  stopSystemInfoMonitoring()
})

// 处理来自 Toolbar 的过滤事件
const handleApplyFilter = async (event: any) => {
  const { filter, mode, options, sourceWindow } = event.detail

  console.log('🔍 LogViewer: 收到过滤请求:', {
    filter,
    mode,
    sourceWindow,
    currentSelectedWindow: selectedWindow.value,
    isLargeFile: isLargeFile.value,
    totalLines: totalLines.value,
    loadedChunks: loadedChunks.value.size
  })

  try {
    // 根据源窗口决定过滤的数据源
    let sourceLines = logLines.value
    let sourceDescription = '主窗口'

    if (sourceWindow === 'filter' && activeFilterWindowData.value) {
      sourceLines = activeFilterWindowData.value.filteredLines
      sourceDescription = `过滤窗口(${activeFilterWindowData.value.name})`
      console.log('🔍 基于过滤窗口进行二次过滤:', sourceDescription)
    } else if (isLargeFile.value && sourceWindow === 'main') {
      // 大文件模式下，统一使用完整过滤
      const loadedLines = logLines.value.filter(line => line !== undefined).length
      const totalLinesCount = totalLines.value

      if (loadedLines < totalLinesCount) {
        console.log(`📊 大文件检测: 已加载 ${loadedLines}/${totalLinesCount} 行，开始完整过滤`)

        // 显示过滤提示
        ElMessage.info({
          message: `检测到大文件，将加载并过滤全部 ${totalLinesCount.toLocaleString()} 行数据`,
          duration: 3000
        })

        appStore.setGlobalLoading(true, '正在加载全部数据进行过滤...', 10)

        try {
          // 加载全部数据
          sourceLines = await loadAllDataForFilter(appStore.currentFile.path)
          sourceDescription = '主窗口(全部数据)'

          // 检查是否被取消
          if (!appStore.isGlobalLoading) {
            console.log('⏹️ 操作已被用户取消')
            return
          }
        } catch (error) {
          console.error('❌ 加载全部数据失败:', error)
          appStore.setGlobalLoading(false)
          ElMessage.error('加载全部数据失败，无法进行完整过滤')

          // 抛出错误，中断过滤操作
          throw new Error('加载全部数据失败，过滤操作已中断')
        }
      }
    }

    // 分批处理过滤，避免界面卡死
    if (!appStore.isGlobalLoading) {
      appStore.setGlobalLoading(true, '正在执行过滤...', 70)
    } else {
      appStore.updateLoadingProgress(70, '正在执行过滤...')
    }

    const filterResult = await filterLogLinesWithProgress(sourceLines, filter, mode, options)

    // 创建新的过滤窗口标签页
    const filterId = `filter_${Date.now()}`
    const windowName = sourceWindow === 'filter' ?
      `${filter} (基于${sourceDescription})` :
      `${filter}`

    const filterWindow = {
      id: filterId,
      name: windowName,
      filter: filter,
      mode: mode,
      options: { ...options },
      sourceWindow: sourceWindow || 'main', // 记录源窗口
      height: 300, // 保留这个属性以兼容现有代码
      filteredLines: filterResult.filteredLines,
      originalLineNumbers: filterResult.originalLineNumbers
    }

    filterWindows.value.push(filterWindow)
    activeFilterWindow.value = filterId

    // 显示过滤结果统计
    // 修复源行数计算问题
    let sourceLineCount = 0;
    if (Array.isArray(sourceLines)) {
      if (isLargeFile.value) {
        // 对于大文件，使用totalLines.value
        sourceLineCount = totalLines.value;
      } else {
        // 对于普通文件，计算非undefined行数
        sourceLineCount = (sourceLines as any[]).filter((line: any) => line !== undefined).length;
      }
    } else {
      sourceLineCount = (sourceLines as any[]).length;
    }

    console.log('📊 过滤完成统计:', {
      sourceLines: sourceLineCount,
      filteredLines: filterResult.filteredLines.length,
      percentage: ((filterResult.filteredLines.length / sourceLineCount) * 100).toFixed(2) + '%'
    })

    ElMessage.success(
      `过滤完成：从 ${sourceLineCount.toLocaleString()} 行中找到 ${filterResult.filteredLines.length.toLocaleString()} 条匹配记录`
    )

    // 通知 Toolbar 更新过滤窗口列表和保存状态
    notifyFilterWindowsUpdated()

    // 关闭加载状态（如果有）
    if (appStore.isGlobalLoading) {
      appStore.setGlobalLoading(false)
    }

    // 设置过滤窗口焦点和事件监听器
    nextTick(() => {
      if (filterLogDisplayRef.value) {
        filterLogDisplayRef.value.focus()
        // 添加滚动事件监听器
        filterLogDisplayRef.value.addEventListener('wheel', handleFilterScroll, { passive: false })
      }
      // 同步过滤窗口行号高度
      syncFilterWindowLineHeights()
    })

    ElMessage.success(`过滤完成，找到 ${filterResult.filteredLines.length} 条匹配记录`)
  } catch (error) {
    console.error('过滤操作失败:', error)

    // 关闭加载状态
    if (appStore.isGlobalLoading) {
      appStore.setGlobalLoading(false)
    }

    const errorMessage = (error as Error).message

    // 如果是用户取消操作，显示友好提示
    if (errorMessage.includes('操作已取消') || errorMessage.includes('取消')) {
      ElMessage.info('过滤操作已取消')
      return
    }

    // 其他错误显示具体错误信息
    ElMessage.error('过滤失败: ' + errorMessage)
  }
}

// 处理获取过滤窗口列表的事件
const handleGetFilterWindows = () => {
  notifyFilterWindowsUpdated()
}



// 处理恢复过滤窗口状态事件
const handleRestoreFilterWindows = (event: any) => {
  const { windows, activeWindow, filterWindowsHeight: height, filterStartLine: startLine } = event.detail

  console.log('📂 LogViewer: 恢复过滤窗口状态', {
    windowsCount: windows.length,
    activeWindow,
    height,
    startLine
  })

  // 清空现有状态，避免重复
  filterWindows.value = []
  activeFilterWindow.value = ''

  // 等待DOM更新
  nextTick(async () => {
    // 重新执行过滤逻辑以生成filteredLines和originalLineNumbers
    const restoredWindows = []
    for (const window of windows) {
      console.log('🔍 重新执行过滤:', window.name, 'for', logLines.value.length, 'lines')

      if (logLines.value.length === 0 && !isLargeFile.value) {
        console.warn('⚠️ 日志内容为空，跳过过滤')
        restoredWindows.push({
          ...window,
          filteredLines: [],
          originalLineNumbers: []
        })
        continue
      }

      try {
        let filterResult;
        if (isLargeFile.value && appStore.currentFile) {
          // 对于大文件，加载全部数据进行过滤
          const allLines = await loadAllDataForFilter(appStore.currentFile.path)
          filterResult = filterLogLines(allLines, window.filter, window.mode, window.options)
        } else {
          // 对于普通文件，直接过滤
          filterResult = filterLogLines(logLines.value, window.filter, window.mode, window.options)
        }
        
        console.log('🔍 过滤结果:', window.name, '找到', filterResult.filteredLines.length, '条记录')
        restoredWindows.push({
          ...window,
          filteredLines: filterResult.filteredLines,
          originalLineNumbers: filterResult.originalLineNumbers
        })
      } catch (error) {
        console.error('❌ 恢复过滤窗口时过滤失败:', error)
        // 如果过滤失败，使用空结果
        restoredWindows.push({
          ...window,
          filteredLines: [],
          originalLineNumbers: []
        })
      }
    }

    filterWindows.value = restoredWindows
    activeFilterWindow.value = activeWindow

    if (height !== undefined) {
      filterWindowsHeight.value = height
    }

    if (startLine !== undefined) {
      filterStartLine.value = startLine
    } else {
      filterStartLine.value = 0
    }

    // 重置高亮状态
    highlightedFilterLines.value = []

    console.log('✅ 过滤窗口状态恢复完成:', {
      windowsCount: filterWindows.value.length,
      activeWindow: activeFilterWindow.value,
      firstWindowLines: filterWindows.value[0]?.filteredLines.length || 0
    })

    // 如果有过滤窗口，设置焦点并同步行号
    if (restoredWindows.length > 0) {
      nextTick(() => {
        if (filterLogDisplayRef.value) {
          filterLogDisplayRef.value.focus()
          // 重新添加滚动事件监听器
          filterLogDisplayRef.value.addEventListener('wheel', handleFilterScroll, { passive: false })
        }
        // 同步过滤窗口行号高度
        syncFilterWindowLineHeights()
      })
    }
  })
}

// 通知过滤窗口列表更新
const notifyFilterWindowsUpdated = () => {

  window.dispatchEvent(new CustomEvent('filterWindowsUpdated', {
    detail: {
      windows: filterWindows.value, // 发送完整的过滤窗口数据
      activeWindow: activeFilterWindow.value,
      filterWindowsHeight: filterWindowsHeight.value,
      filterStartLine: filterStartLine.value
    }
  }))
}

const filterLogLines = (lines: (string | undefined)[], filter: string, mode: string, options: any): {filteredLines: string[], originalLineNumbers: number[]} => {
  if (!filter.trim()) return {filteredLines: [], originalLineNumbers: []}

  try {
    const filteredLines: string[] = []
    const originalLineNumbers: number[] = []

    // 使用for循环而不是forEach以更好地处理稀疏数组
    for (let index = 0; index < lines.length; index++) {
      const line = lines[index];
      // 跳过未加载的行（大文件模式）
      if (line === undefined) continue

      const matches = evaluateFilterExpression(line, filter, options)
      const shouldInclude = mode === 'include' ? matches : !matches
      if (shouldInclude) {
        filteredLines.push(line)
        originalLineNumbers.push(index)
      }
    }

    return { filteredLines, originalLineNumbers }
  } catch (error) {
    console.error('过滤表达式错误:', error)
    return {filteredLines: [], originalLineNumbers: []}
  }
}

// 过滤操作状态管理（暂时未使用，但保留以备将来使用）
// let currentFilterOperation = null as { cancelled: boolean } | null

// 取消过滤操作（暂时未使用，但保留以备将来使用）
// const cancelFilterOperation = () => {
//   if (currentFilterOperation) {
//     currentFilterOperation.cancelled = true
//     console.log('🛑 用户取消过滤操作')
//   }
// }

// 带进度的过滤函数（用于大文件）
const filterLogLinesWithProgress = async (lines: (string | undefined)[], filter: string, mode: string, options: any): Promise<{filteredLines: string[], originalLineNumbers: number[]}> => {
  if (!filter.trim()) return {filteredLines: [], originalLineNumbers: []}

  try {
    // 创建新的过滤操作状态（暂时未使用）
    // currentFilterOperation = { cancelled: false }

    const filteredLines: string[] = []
    const originalLineNumbers: number[] = []
    const batchSize = 10000 // 每批处理1万行
    // 对于大文件，使用totalLines.value而不是lines.length
    const totalLineCount = isLargeFile.value ? totalLines.value : lines.length

    console.log(`🔍 开始分批过滤，总行数: ${totalLineCount.toLocaleString()}，批大小: ${batchSize.toLocaleString()}`)

    for (let i = 0; i < totalLineCount; i += batchSize) {
      // 暂时移除取消检查，让过滤正常执行
      // TODO: 后续可以添加更好的取消机制

      const endIndex = Math.min(i + batchSize, totalLineCount)

      // 处理当前批次，直接遍历而不是使用slice（避免稀疏数组问题）
      for (let batchIndex = 0; batchIndex < batchSize && i + batchIndex < totalLineCount; batchIndex++) {
        try {
          const actualIndex = i + batchIndex
          const line = lines[actualIndex]
          
          // 跳过未加载的行（大文件模式）
          if (line === undefined) continue

          const matches = evaluateFilterExpression(line, filter, options)
          const shouldInclude = mode === 'include' ? matches : !matches
          if (shouldInclude) {
            filteredLines.push(line)
            originalLineNumbers.push(actualIndex)
          }
        } catch (lineError) {
          // 单行处理错误不应该中断整个过滤过程
          console.warn(`处理第 ${i + batchIndex + 1} 行时出错:`, lineError)
        }
      }

      // 更新进度
      const progress = 70 + (endIndex / totalLineCount) * 20 // 70-90%
      const processedLines = endIndex.toLocaleString()
      const totalLinesStr = totalLineCount.toLocaleString()
      const foundLines = filteredLines.length.toLocaleString()

      appStore.updateLoadingProgress(
        progress,
        `正在过滤... ${processedLines}/${totalLinesStr} 行 (找到 ${foundLines} 条匹配)`
      )

      // 给UI一些时间更新，避免界面卡死
      if (i % (batchSize * 5) === 0) { // 每5批次暂停一次
        await new Promise(resolve => setTimeout(resolve, 5))
      }
    }

    console.log(`✅ 过滤完成: 找到 ${filteredLines.length.toLocaleString()} 条匹配记录`)
    appStore.updateLoadingProgress(90, `过滤完成，找到 ${filteredLines.length.toLocaleString()} 条匹配记录`)

    return { filteredLines, originalLineNumbers }
  } catch (error) {
    console.error('过滤表达式错误:', error)
    return {filteredLines: [], originalLineNumbers: []}
  } finally {
    // 清理过滤操作状态（暂时未使用）
    // currentFilterOperation = null
  }
}

// 解析和评估过滤表达式
const evaluateFilterExpression = (line: string, expression: string, options: any): boolean => {
  try {
    if (options.isRegex) {
      // 正则表达式模式
      const flags = options.caseSensitive ? 'g' : 'gi'
      const regex = new RegExp(expression, flags)
      return regex.test(line)
    }

    // 解析逻辑表达式
    return parseLogicalExpression(line, expression, options)
  } catch (error) {
    console.warn('过滤表达式评估错误:', error, '表达式:', expression)
    // 表达式错误时返回false，不匹配该行
    return false
  }
}

// 解析逻辑表达式（支持 AND, OR, &&, ||, 括号）
const parseLogicalExpression = (line: string, expression: string, options: any): boolean => {
  try {
    // 预处理：标准化操作符
    let normalizedExpr = expression
      .replace(/\s+AND\s+/gi, ' AND ')
      .replace(/\s+OR\s+/gi, ' OR ')
      .replace(/\s+and\s+/gi, ' AND ')
      .replace(/\s+or\s+/gi, ' OR ')
      .replace(/\s*&&\s*/g, ' AND ')  // 支持 && 操作符
      .replace(/\s*\|\|\s*/g, ' OR ')  // 支持 || 操作符
      .replace(/\s*\(\s*/g, '(')
      .replace(/\s*\)\s*/g, ')')

    // 检查括号是否匹配
    const openCount = (normalizedExpr.match(/\(/g) || []).length
    const closeCount = (normalizedExpr.match(/\)/g) || []).length
    if (openCount !== closeCount) {
      throw new Error('括号不匹配')
    }

    // 递归解析括号，限制递归深度防止无限循环
    let recursionDepth = 0
    const maxRecursionDepth = 10

    while (normalizedExpr.includes('(') && recursionDepth < maxRecursionDepth) {
      const innerMatch = normalizedExpr.match(/\([^()]+\)/)
      if (!innerMatch) break

      const innerExpr = innerMatch[0].slice(1, -1) // 去掉括号
      const innerResult = evaluateSimpleExpression(line, innerExpr, options)
      normalizedExpr = normalizedExpr.replace(innerMatch[0], innerResult.toString())
      recursionDepth++
    }

    if (recursionDepth >= maxRecursionDepth) {
      throw new Error('表达式过于复杂，递归深度超限')
    }

    return evaluateSimpleExpression(line, normalizedExpr, options)
  } catch (error) {
    console.warn('逻辑表达式解析错误:', error, '表达式:', expression)
    // 解析错误时，尝试简单的字符串匹配
    return matchKeyword(line, expression, options)
  }
}

// 评估简单表达式（不含括号）
const evaluateSimpleExpression = (line: string, expression: string, options: any): boolean => {
  try {
    // 处理空表达式
    if (!expression || !expression.trim()) {
      return false
    }

    // 处理 OR 操作符（优先级较低）
    if (expression.includes(' OR ')) {
      return expression.split(' OR ').some(part => {
        const trimmedPart = part.trim()
        return trimmedPart ? evaluateSimpleExpression(line, trimmedPart, options) : false
      })
    }

    // 处理 AND 操作符（优先级较高）
    if (expression.includes(' AND ')) {
      return expression.split(' AND ').every(part => {
        const trimmedPart = part.trim()
        return trimmedPart ? evaluateSimpleExpression(line, trimmedPart, options) : false
      })
    }

    // 处理布尔值（来自括号内的计算结果）
    if (expression === 'true') return true
    if (expression === 'false') return false

    // 单个关键词匹配
    return matchKeyword(line, expression, options)
  } catch (error) {
    console.warn('简单表达式评估错误:', error, '表达式:', expression)
    return false
  }
}

// 匹配单个关键词
const matchKeyword = (line: string, keyword: string, options: any): boolean => {
  try {
    if (!keyword || !keyword.trim()) return false
    if (!line) return false

    let searchText = line
    let searchKeyword = keyword.trim()

    if (!options.caseSensitive) {
      searchText = line.toLowerCase()
      searchKeyword = searchKeyword.toLowerCase()
    }

    if (options.wholeWord) {
      const escapedKeyword = escapeRegExp(searchKeyword)
      const regex = new RegExp(`\\b${escapedKeyword}\\b`, options.caseSensitive ? 'g' : 'gi')
      return regex.test(line)
    }

    return searchText.includes(searchKeyword)
  } catch (error) {
    console.warn('关键词匹配错误:', error, '关键词:', keyword)
    return false
  }
}



// 关闭过滤窗口
const closeFilterWindow = (filterId: string) => {
  const index = filterWindows.value.findIndex(w => w.id === filterId)
  if (index > -1) {
    filterWindows.value.splice(index, 1)
    if (activeFilterWindow.value === filterId) {
      activeFilterWindow.value = filterWindows.value.length > 0 ? filterWindows.value[0].id : ''
      // 重置状态
      filterStartLine.value = 0
      highlightedFilterLines.value = []
    }
    // 通知 Toolbar 更新过滤窗口列表
    notifyFilterWindowsUpdated()
  }
}







// 从表达式中提取关键词
const extractKeywordsFromExpression = (expression: string): string[] => {
  // 移除逻辑操作符和括号，提取关键词
  return expression
    .replace(/\s+AND\s+/gi, ' ')
    .replace(/\s+OR\s+/gi, ' ')
    .replace(/\s*&&\s*/g, ' ')  // 移除 && 操作符
    .replace(/\s*\|\|\s*/g, ' ')  // 移除 || 操作符
    .replace(/[()]/g, ' ')
    .split(/\s+/)
    .filter(word => word.trim() && !['AND', 'OR', 'and', 'or'].includes(word.trim()))
}

// 过滤窗口键盘事件处理（完全复用主窗口逻辑）
const handleFilterKeyDown = (event: KeyboardEvent) => {
  if (event.ctrlKey && event.key === 'a') {
    // Ctrl+A 全选 - 使用浏览器原生选择
    event.preventDefault()
    const selection = window.getSelection()
    if (selection && filterLogDisplayRef.value) {
      const range = document.createRange()
      range.selectNodeContents(filterLogDisplayRef.value)
      selection.removeAllRanges()
      selection.addRange(range)
    }
    return
  }

  // 快捷键：t - 快速添加时间线
  if (event.key === 't' || event.key === 'T') {
    event.preventDefault()
    if (currentFocusedLine.value) {
      // 触发快速添加时间线事件
      window.dispatchEvent(new CustomEvent('quickAddTimeline', {
        detail: currentFocusedLine.value
      }))
    } else {
      ElMessage.warning('请先点击选择一行日志')
    }
    return
  }

  // 快捷键：f - 对选中词过滤
  if (event.key === 'f' || event.key === 'F') {
    event.preventDefault()
    const selection = window.getSelection()
    if (selection && selection.toString().trim()) {
      const selectedText = selection.toString().trim()
      // 触发快速过滤事件
      window.dispatchEvent(new CustomEvent('quickFilter', {
        detail: { text: selectedText }
      }))
    } else {
      ElMessage.warning('请先选中要过滤的文本')
    }
    return
  }
}







// 过滤窗口语法高亮（复用主窗口逻辑）
const highlightFilterSyntax = (line: string, lineIndex: number): string => {
  // 复用主窗口的语法高亮逻辑
  let highlighted = highlightSyntax(line, lineIndex)

  // 添加过滤条件高亮
  if (activeFilterWindowData.value) {
    const filterWindow = activeFilterWindowData.value
    try {
      if (filterWindow.options.isRegex) {
        const flags = filterWindow.options.caseSensitive ? 'g' : 'gi'
        const regex = new RegExp(filterWindow.filter, flags)
        highlighted = highlighted.replace(regex, '<span class="filter-highlight">$&</span>')
      } else {
        const keywords = extractKeywordsFromExpression(filterWindow.filter)
        keywords.forEach(keyword => {
          if (keyword.trim()) {
            let regex: RegExp
            if (filterWindow.options.wholeWord) {
              regex = new RegExp(`\\b${escapeRegExp(keyword)}\\b`, filterWindow.options.caseSensitive ? 'g' : 'gi')
            } else {
              regex = new RegExp(escapeRegExp(keyword), filterWindow.options.caseSensitive ? 'g' : 'gi')
            }
            highlighted = highlighted.replace(regex, '<span class="filter-highlight">$&</span>')
          }
        })
      }
    } catch (error) {
      console.error('过滤高亮处理错误:', error)
    }
  }

  return highlighted
}











// 拖拽分隔条调整过滤窗口的大小和位置
const startFilterSplitterDrag = (event: MouseEvent) => {
  if (filterWindows.value.length === 0) return

  const startY = event.clientY
  const startHeight = filterWindowsHeight.value

  // 获取容器高度
  const viewerElement = document.querySelector('.log-viewer') as HTMLElement
  const containerHeight = viewerElement ? viewerElement.clientHeight : 800

  const handleMouseMove = (e: MouseEvent) => {
    const deltaY = e.clientY - startY

    // 简化逻辑：直接根据拖拽方向调整高度
    // 向上拖拽（deltaY < 0）：增加高度
    // 向下拖拽（deltaY > 0）：减少高度
    const newHeight = Math.max(50, Math.min(containerHeight - 50, startHeight - deltaY))

    filterWindowsHeight.value = newHeight
  }

  const handleMouseUp = () => {
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
  }

  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)
  event.preventDefault()
}

// 处理主窗口聚焦
const handleMainWindowFocus = () => {
  console.log('🎯 主窗口获得聚焦')
  focusedWindow.value = 'main'

  // 通知Toolbar更新窗口选择
  window.dispatchEvent(new CustomEvent('windowFocused', {
    detail: { windowType: 'main' }
  }))
}

// 处理过滤窗口聚焦
const handleFilterWindowFocus = () => {
  console.log('🎯 过滤窗口获得聚焦')
  focusedWindow.value = 'filter'

  // 通知Toolbar更新窗口选择
  window.dispatchEvent(new CustomEvent('windowFocused', {
    detail: { windowType: 'filter' }
  }))
}

// 处理窗口选择变化（来自Toolbar）
const handleWindowChanged = (event: any) => {
  const { windowId } = event.detail
  console.log('🪟 收到窗口选择变化:', windowId)
  selectedWindow.value = windowId === 'main' ? 'main' : 'filter'
}

// 处理高亮词更新事件
const handleHighlightWordsUpdated = (event: any) => {
  const { highlightWords } = event.detail
  console.log('🎨 LogViewer: 收到高亮词更新事件:', {
    highlightWordsCount: highlightWords.length,
    words: highlightWords.map((w: any) => w.text),
    appStoreHighlightWords: appStore.highlightWords.length
  })

  // 强制重新渲染以应用新的高亮词
  // 使用多种方法确保重新渲染
  nextTick(() => {
    // 方法1: 触发startLine变化
    const currentStartLine = startLine.value
    startLine.value = currentStartLine + 1

    nextTick(() => {
      startLine.value = currentStartLine

      // 方法2: 强制更新可见行数
      const currentVisibleCount = visibleLineCount.value
      visibleLineCount.value = currentVisibleCount + 1

      nextTick(() => {
        visibleLineCount.value = currentVisibleCount

        // 方法3: 触发窗口resize事件
        window.dispatchEvent(new Event('resize'))

        console.log('🎨 LogViewer: 高亮词更新完成，已强制重新渲染')
      })
    })
  })
}



// 加载全部数据用于过滤（分块加载，避免界面卡死）
const loadAllDataForFilter = async (filePath: string): Promise<string[]> => {
  console.log('📖 开始分块加载全部数据用于过滤...')

  try {
    const allLines: string[] = []
    const chunkSize = 50000 // 每次加载5万行
    let currentLine = 0
    let hasMoreData = true

    appStore.updateLoadingProgress(20, '正在分块加载文件数据...')

    while (hasMoreData) {
      // 检查是否被取消
      if (!appStore.isGlobalLoading) {
        console.log('⏹️ 数据加载被用户取消')
        throw new Error('操作已取消')
      }

      console.log(`📖 加载数据块: ${currentLine} - ${currentLine + chunkSize}`)

      try {
        // 使用分块读取API
        const endLine = Math.min(currentLine + chunkSize, totalLines.value)
        const chunk = await ReadLogFileChunk(filePath, currentLine, endLine)

        if (chunk && chunk.lines && chunk.lines.length > 0) {
          // 过滤掉undefined的行
          const validLines = chunk.lines.filter(line => line !== undefined) as string[]
          allLines.push(...validLines)

          currentLine += chunkSize

          // 更新进度
          const progress = Math.min(20 + (allLines.length / totalLines.value) * 40, 60)
          appStore.updateLoadingProgress(progress, `已加载 ${allLines.length.toLocaleString()} 行数据...`)

          // 给UI一些时间更新，避免界面卡死
          await new Promise(resolve => setTimeout(resolve, 10))

          // 如果返回的行数少于请求的数量，说明已经到文件末尾
          if (chunk.lines.length < chunkSize) {
            hasMoreData = false
          }
        } else {
          hasMoreData = false
        }
      } catch (chunkError) {
        console.warn(`⚠️ 加载数据块失败 (${currentLine}-${currentLine + chunkSize}):`, chunkError)
        // 如果某个块加载失败，停止加载并抛出错误
        throw new Error(`加载数据块失败 (${currentLine}-${currentLine + chunkSize}): ${chunkError.message}`)
      }
    }

    console.log(`✅ 全部数据加载完成: ${allLines.length.toLocaleString()} 行`)
    appStore.updateLoadingProgress(65, `数据加载完成，共 ${allLines.length.toLocaleString()} 行`)

    return allLines
  } catch (error) {
    console.error('❌ 加载全部数据失败:', error)
    throw error
  }
}

// 支持大文件模式的跳转到指定行
const jumpToLineWithLoading = async (lineNumber: number, isSearchResult: boolean = false) => {
  console.log('🎯 跳转到行（支持大文件）:', { lineNumber, isSearchResult, isLargeFile: isLargeFile.value })

  if (lineNumber < 1) {
    console.warn('❌ 无效的行号:', lineNumber)
    ElMessage.warning('行号必须大于0')
    return
  }

  // 检查行号是否超出范围
  const maxLines = isLargeFile.value ? totalLines.value : logLines.value.length
  if (lineNumber > maxLines) {
    console.warn('❌ 行号超出范围:', { lineNumber, maxLines })
    ElMessage.warning(`行号不能超过 ${maxLines}`)
    return
  }

  // 大文件模式下需要确保目标行已加载
  if (isLargeFile.value && appStore.currentFile) {
    const targetIndex = lineNumber - 1 // 转换为0基索引

    // 检查目标行是否已加载
    if (logLines.value[targetIndex] === undefined) {
      console.log('📖 目标行未加载，开始加载...')
      appStore.setGlobalLoading(true, `正在加载第 ${lineNumber} 行...`, 20)

      try {
        // 计算需要加载的块
        const targetChunk = Math.floor(targetIndex / chunkSize)
        console.log(`📦 需要加载块: ${targetChunk}`)

        // 加载目标块
        await loadChunk(appStore.currentFile.path, targetChunk)
        appStore.updateLoadingProgress(80, '正在定位目标行...')

        // 等待DOM更新
        await nextTick()

        appStore.setGlobalLoading(false)
        console.log('✅ 目标行加载完成')
      } catch (error) {
        console.error('❌ 加载目标行失败:', error)
        appStore.setGlobalLoading(false)
        ElMessage.error(`跳转到第 ${lineNumber} 行失败`)
        return
      }
    }
  }

  // 执行跳转
  await jumpToLine(lineNumber, isSearchResult)

  // 确保行号显示正确更新
  await nextTick()

  // 额外的行号刷新确保
  if (isLargeFile.value) {
    refreshLineNumbers()
    await nextTick()
  }

  console.log('🔢 行号显示更新完成')
}

// 支持大文件模式的跳转到尾行
const jumpToLastLineWithLoading = async () => {
  console.log('🎯 跳转到尾行（支持大文件）:', { isLargeFile: isLargeFile.value, totalLines: totalLines.value })

  if (isLargeFile.value && appStore.currentFile) {
    const lastLineNumber = totalLines.value

    if (lastLineNumber <= 0) {
      console.warn('❌ 无效的总行数:', lastLineNumber)
      return
    }

    console.log(`📖 跳转到尾行: ${lastLineNumber}`)
    appStore.setGlobalLoading(true, `正在加载尾行（第 ${lastLineNumber} 行）...`, 20)

    try {
      // 计算最后一块
      const lastChunk = Math.floor((lastLineNumber - 1) / chunkSize)
      console.log(`📦 需要加载最后一块: ${lastChunk}`)

      // 加载最后一块
      await loadChunk(appStore.currentFile.path, lastChunk)
      appStore.updateLoadingProgress(80, '正在定位尾行...')

      // 等待DOM更新
      await nextTick()

      appStore.setGlobalLoading(false)

      // 跳转到尾行
      jumpToLine(lastLineNumber, false)
    } catch (error) {
      console.error('❌ 加载尾行失败:', error)
      appStore.setGlobalLoading(false)
      ElMessage.error('跳转到尾行失败')
    }
  } else {
    // 小文件模式，直接跳转
    const lastLineNumber = logLines.value.length
    if (lastLineNumber > 0) {
      jumpToLine(lastLineNumber, false)
    }
  }
}

// 过滤窗口跳转到尾行
const jumpToFilterLastLine = () => {
  if (activeFilterWindowData.value && activeFilterWindowData.value.filteredLines.length > 0) {
    const lastLine = activeFilterWindowData.value.filteredLines.length
    jumpToFilterLine(lastLine, false)
  }
}

// 跳转到过滤窗口中的指定行
const jumpToFilterLine = (lineNumber: number, isSearchResult: boolean = false) => {
  if (!activeFilterWindowData.value || !activeFilterWindowData.value.filteredLines.length) {
    console.warn('⚠️ 过滤窗口无数据，无法跳转')
    return
  }

  // 在过滤结果中查找对应的原始行号
  const originalLineNumbers = activeFilterWindowData.value.originalLineNumbers
  const targetIndex = originalLineNumbers.findIndex(num => num === lineNumber)

  if (targetIndex === -1) {
    console.warn('⚠️ 在过滤结果中未找到行号:', lineNumber)
    return
  }

  console.log('🎯 跳转到过滤窗口行:', {
    originalLineNumber: lineNumber,
    filterIndex: targetIndex,
    totalFilterLines: activeFilterWindowData.value.filteredLines.length
  })

  // 计算过滤窗口的起始行
  const maxStartLine = Math.max(0, activeFilterWindowData.value.filteredLines.length - filterVisibleLineCount.value)
  const idealStartLine = Math.max(0, targetIndex - Math.floor(filterVisibleLineCount.value / 2))
  filterStartLine.value = Math.min(idealStartLine, maxStartLine)

  // 高亮目标行
  if (isSearchResult) {
    highlightedFilterLines.value = [targetIndex]
  }

  // 聚焦到过滤窗口
  focusedWindow.value = 'filter'
  selectedWindow.value = 'filter'

  nextTick(() => {
    if (filterLogDisplayRef.value) {
      filterLogDisplayRef.value.focus()
    }
  })
}

// 调试函数：检查行号显示
const debugLineNumbers = () => {
  console.log('🔍 行号显示调试信息:')
  console.log('当前显示范围:', startLine.value, '-', startLine.value + visibleLineCount.value)
  console.log('可见行数:', visibleLines.value.length)
  console.log('可见行号:', visibleLineNumbers.value)
  console.log('日志记录映射示例:', lineToLogEntryMap.value.slice(startLine.value, startLine.value + 5))

  if (activeFilterWindowData.value) {
    console.log('过滤窗口信息:')
    console.log('过滤显示范围:', filterStartLine.value, '-', filterStartLine.value + filterVisibleLineCount.value)
    console.log('过滤可见行数:', visibleFilterLines.value.length)
    console.log('过滤可见行号:', visibleFilterLineNumbers.value)
  }
}

// 调试函数：检查过滤状态
const debugFilterStatus = () => {
  console.log('🔍 过滤状态调试信息:')
  console.log('是否大文件模式:', isLargeFile.value)
  console.log('总行数:', totalLines.value)
  console.log('已加载块数:', loadedChunks.value.size)

  const loadedLines = logLines.value.filter(line => line !== undefined).length
  console.log('已加载行数:', loadedLines)
  console.log('加载百分比:', totalLines.value > 0 ? ((loadedLines / totalLines.value) * 100).toFixed(2) + '%' : 'N/A')

  console.log('过滤窗口数量:', filterWindows.value.length)
  filterWindows.value.forEach((window, index) => {
    console.log(`  过滤窗口 ${index + 1}:`, {
      name: window.name,
      filteredLines: window.filteredLines.length,
      filter: window.filter,
      mode: window.mode
    })
  })

  if (activeFilterWindowData.value) {
    const shouldShowScrollbar = activeFilterWindowData.value.filteredLines.length > filterVisibleLineCount.value
    console.log('当前活动过滤窗口:', {
      name: activeFilterWindowData.value.name,
      totalLines: activeFilterWindowData.value.filteredLines.length,
      visibleLines: filterVisibleLineCount.value,
      startLine: filterStartLine.value,
      shouldShowScrollbar: shouldShowScrollbar,
      scrollbarElement: scrollBarRef.value ? '存在' : '不存在' // 使用共用滚动条
    })

    if (shouldShowScrollbar) {
      console.log('📏 共用滚动条信息:', {
        thumbStyle: unifiedScrollThumbStyle.value,
        scrollBarVisible: !!scrollBarRef.value,
        focusedWindow: focusedWindow.value
      })
    }
  } else {
    console.log('❌ 没有活动的过滤窗口')
  }
}

// 在开发模式下暴露调试函数
if (import.meta.env.DEV) {
  const win = window as any
  win.debugLineNumbers = debugLineNumbers
  win.debugFilterStatus = debugFilterStatus
}

// 暴露方法给父组件
defineExpose({
  debugLineNumbers
})
</script>

<style scoped>
.log-viewer {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #ffffff;
  overflow: hidden; /* 防止内容溢出 */
  --log-font-size: 14px;
  --log-font-family: 'Courier New';
  --log-line-height: 1.5;
}

/* 主日志容器 - 可被过滤窗口挤压 */
.main-log-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 200px; /* 确保主窗口有最小高度 */
}

.file-tabs {
  display: flex;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e5e7eb;
  overflow-x: auto;
}

.file-tab {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-right: 1px solid #e5e7eb;
  cursor: pointer;
  white-space: nowrap;
  transition: background-color 0.2s;
}

.file-tab:hover {
  background-color: #e5e7eb;
}

.file-tab.active {
  background-color: #ffffff;
  border-bottom: 2px solid #3b82f6;
}

.tab-name {
  font-size: 13px;
  color: #374151;
}

.tab-close {
  opacity: 0;
  transition: opacity 0.2s;
  padding: 2px;
  width: 20px;
  height: 20px;
}

.file-tab:hover .tab-close {
  opacity: 1;
}

.log-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 200px;
}

.status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 16px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e5e7eb;
  font-size: 12px;
  color: #6b7280;
}

.status-left {
  display: flex;
  gap: 16px;
}



.file-info {
  font-weight: 500;
  color: #374151;
}

.file-size {
  font-weight: normal;
  color: #6b7280;
}

.log-display {
  flex: 1;
  position: relative;
  overflow: hidden;
  background-color: #ffffff;
  outline: none; /* 移除焦点时的边框 */
}

.log-display:focus {
  /* 可以添加焦点样式，比如淡淡的边框 */
  box-shadow: inset 0 0 0 1px rgba(107, 155, 210, 0.3);
}

.log-lines {
  display: flex;
  height: 100%;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 20px;
  align-items: flex-start;
}

.line-numbers {
  background-color: #6b9bd2;
  border-right: 1px solid #5a8bc4;
  padding: 8px 12px;
  text-align: right;
  color: #ffffff;
  user-select: none;
  min-width: 60px;
  font-weight: 500;
  transition: all 0.2s ease;
}

/* 聚焦状态的行号 - 蓝色 */
.line-numbers.focused {
  background-color: #2196f3;
  border-right-color: #1976d2;
  color: #ffffff;
}

/* 选择状态但未聚焦的行号 - 灰色 */
.line-numbers.selected:not(.focused) {
  background-color: #9e9e9e;
  border-right-color: #757575;
  color: #ffffff;
}

.line-number-container {
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  min-height: 20px;
  height: 20px; /* 默认高度，会被JS动态调整 */
  position: relative;
  transition: height 0.1s ease; /* 平滑的高度变化 */
}

.line-number-content {
  height: 100%;
  min-height: 20px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  position: sticky;
  top: 0;
  background-color: inherit;
  padding: 0 4px;
}

.line-number-text {
  padding: 0 4px;
  line-height: 20px;
}

.line-number-container.highlighted {
  background-color: #5a8bc4;
}

.line-number-container.highlighted .line-number-text {
  color: #ffffff;
  font-weight: bold;
}

.log-content-lines {
  flex: 1;
  padding: 8px 12px;
  overflow: hidden;
  user-select: text;
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
}

.log-line {
  height: calc(var(--log-font-size) * var(--log-line-height));
  display: flex;
  align-items: center;
  cursor: text;
  white-space: nowrap;
  overflow: hidden;
  padding: 0;
  line-height: calc(var(--log-font-size) * var(--log-line-height));
  font-size: var(--log-font-size);
  font-family: var(--log-font-family);
  user-select: text;
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
}

.log-line.word-wrap {
  white-space: pre-wrap;
  word-break: break-all;
  height: auto;
  min-height: 20px;
  align-items: flex-start;
  padding: 0;
}

.log-line.word-wrap span {
  display: block;
  width: 100%;
  text-align: left;
}

.log-line.highlighted {
  background-color: #dbeafe;
}

.log-line:hover {
  background-color: #f3f4f6;
}

.log-line.highlighted {
  background-color: #dbeafe !important;
  border-left: 3px solid #3b82f6;
  animation: highlight-flash 0.5s ease-in-out;
}

@keyframes highlight-flash {
  0% { background-color: #60a5fa; }
  100% { background-color: #dbeafe; }
}

/* 系统信息样式 */
.system-info {
  font-size: 12px;
  color: #059669;
  font-weight: 500;
  margin-right: 12px;
  padding: 2px 6px;
  background-color: #ecfdf5;
  border-radius: 4px;
  border: 1px solid #d1fae5;
}

.system-info-debug {
  font-size: 12px;
  color: #f59e0b;
  font-weight: 500;
  margin-right: 12px;
  padding: 2px 6px;
  background-color: #fef3c7;
  border-radius: 4px;
  border: 1px solid #fde68a;
}

.scroll-info {
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 12px;
  background-color: #f1f5f9;
  border-left: 1px solid #e2e8f0; /* 添加边框增强可见性 */
  cursor: pointer;
  z-index: 700; /* 确保滚动条在时间线面板之上 */
}

.scroll-thumb {
  position: absolute;
  right: 2px;
  width: 8px;
  background-color: #94a3b8; /* 更深的颜色增强可见性 */
  border-radius: 4px;
  border: 1px solid #64748b; /* 添加边框 */
  transition: background-color 0.2s;
  cursor: grab;
  min-height: 20px;
}

.scroll-thumb:hover {
  background-color: #94a3b8;
}

.scroll-thumb:active {
  background-color: #64748b;
  cursor: grabbing;
}

/* 语法高亮样式 - 简化版本 */
:deep(.syntax-timestamp) {
  color: #7c3aed;
  font-weight: 500;
}

:deep(.syntax-ip) {
  color: #059669;
  font-weight: 500;
}

:deep(.syntax-method) {
  color: #dc2626;
  font-weight: 600;
}

:deep(.syntax-status) {
  color: #ea580c;
  font-weight: 600;
}

/* 搜索结果高亮样式 */
:deep(.search-highlight) {
  background-color: #fef08a;
  color: #92400e;
  padding: 1px 2px;
  border-radius: 2px;
  font-weight: 600;
}

/* 当前搜索结果高亮样式 */
:deep(.search-highlight-current) {
  background-color: #f97316;
  color: #ffffff;
  padding: 1px 2px;
  border-radius: 2px;
  font-weight: 700;
  box-shadow: 0 0 0 1px #ea580c;
}

/* 过滤窗口样式 - 挤压式布局 */
.filter-windows {
  border-top: 1px solid #e5e7eb;
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
  z-index: 100;
  flex-shrink: 0; /* 防止被压缩 */
  min-height: 100px; /* 最小高度 */
  max-height: 60%; /* 最大高度不超过容器的60% */
}

.filter-splitter {
  height: 8px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e5e7eb;
  cursor: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M8 6L12 2L16 6"/><path d="M2 12L22 12"/><path d="M2 12L22 12" transform="translate(0,2)"/><path d="M8 18L12 22L16 18"/></svg>') 12 12, ns-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
  flex-shrink: 0;
  position: relative;
}

.filter-splitter:hover {
  background-color: #e5e7eb;
}

.filter-splitter:hover .filter-splitter-line {
  background-color: #6b7280;
}

.filter-splitter-line {
  width: 60px;
  height: 3px;
  background-color: #9ca3af;
  border-radius: 2px;
  transition: background-color 0.2s;
  position: relative;
}

.filter-splitter-line::before {
  content: '';
  position: absolute;
  top: -8px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 4px solid transparent;
  border-right: 4px solid transparent;
  border-bottom: 6px solid #9ca3af;
}

.filter-splitter-line::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 4px solid transparent;
  border-right: 4px solid transparent;
  border-top: 6px solid #9ca3af;
}

.filter-splitter:hover .filter-splitter-line::before,
.filter-splitter:hover .filter-splitter-line::after {
  border-bottom-color: #6b7280;
  border-top-color: #6b7280;
}

.filter-container {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.filter-window-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e5e7eb;
  padding: 0;
  height: 32px;
  flex-shrink: 0;
}

.filter-window-tabs {
  display: flex;
  overflow-x: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.filter-window-tabs::-webkit-scrollbar {
  display: none;
}

.filter-window-tab {
  display: flex;
  align-items: center;
  padding: 6px 12px;
  background-color: #e5e7eb;
  border-right: 1px solid #d1d5db;
  cursor: pointer;
  font-size: 12px;
  gap: 6px;
  white-space: nowrap;
  flex-shrink: 0;
  transition: background-color 0.2s;
}

.filter-window-tab:hover {
  background-color: #d1d5db;
}

.filter-window-tab.active {
  background-color: #ffffff;
  border-bottom: 1px solid #ffffff;
}

.filter-tab-name {
  color: #374151;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.filter-tab-count {
  color: #6b7280;
  font-size: 11px;
}

.filter-tab-close {
  opacity: 0;
  transition: opacity 0.2s;
  padding: 2px;
  width: 16px;
  height: 16px;
}

.filter-window-tab:hover .filter-tab-close {
  opacity: 1;
}

.filter-window-controls {
  padding: 0 12px;
  font-size: 12px;
  color: #6b7280;
}

.filter-window-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  position: relative; /* 为绝对定位的滚动条提供定位上下文 */
}

/* 过滤窗口复用主窗口样式 */
.filter-log-display {
  height: 100%;
  display: flex;
  background-color: #ffffff;
  position: relative;
  outline: none;
  overflow: hidden;
}

/* 过滤窗口特定样式 */
.filter-log-display .line-number-container {
  height: 20px; /* 默认高度，会被JS动态调整 */
  min-height: 20px;
  transition: height 0.1s ease;
}

.filter-log-display .line-number-content {
  position: sticky;
  top: 0;
  height: 100%;
  min-height: 20px;
  background-color: inherit;
  padding: 0 4px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.filter-log-display .log-line.word-wrap {
  align-items: flex-start;
}

.filter-log-display .line-number-container .line-number-text {
  line-height: 20px;
}





/* 过滤结果高亮样式 */
:deep(.filter-highlight) {
  background-color: #fef08a;
  color: #92400e;
  padding: 1px 2px;
  border-radius: 2px;
  font-weight: 600;
}

/* 性能信息样式 */
.performance-info {
  color: #888;
  font-size: 11px;
  margin-left: 8px;
}

/* 右键菜单样式 */
.context-menu-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  background: transparent;
}

.context-menu {
  position: fixed;
  background: #ffffff;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  padding: 4px 0;
  min-width: 160px;
  z-index: 1001;
  font-size: 14px;
}

.context-menu-item {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  cursor: pointer;
  transition: background-color 0.2s;
  gap: 8px;
}

.context-menu-item:hover {
  background-color: #f5f7fa;
}

.context-menu-item .el-icon {
  font-size: 16px;
  color: #606266;
}

.context-menu-item span {
  flex: 1;
  color: #303133;
}

.context-menu-item .shortcut {
  flex: none;
  font-size: 12px;
  color: #909399;
  background: #f0f2f5;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: monospace;
}

.context-menu-divider {
  height: 1px;
  background-color: #e4e7ed;
  margin: 4px 0;
}
</style>
