<template>
  <div class="log-viewer">
    <!-- 文件标签页 -->
    <div class="file-tabs" v-if="appStore.openFiles.length > 1">
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
          <span class="line-count">共 {{ logLines.length }} 行，{{ logEntryCount }} 条日志</span>
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
      
      <!-- 日志显示区域 -->
      <div class="log-display" ref="logDisplayRef">
        <div class="log-lines">
          <!-- 行号列 -->
          <div class="line-numbers" v-if="appStore.showLineNumbers">
            <div
              v-for="(_, index) in visibleLines"
              :key="index + startLine"
              class="line-number-container"
              :class="{
                highlighted: highlightedLines.includes(index + startLine)
              }"
            >
              <div class="line-number-content">
                <span class="line-number-text">{{ index + startLine + 1 }}</span>
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
                v-if="appStore.syntaxHighlighting"
                v-html="highlightSyntax(line, startLine + index)"
              ></span>
              <span v-else>{{ line }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 滚动条信息 -->
      <div
        class="scroll-info"
        v-if="logLines.length > visibleLineCount"
        @mousedown="startScrollDrag"
        ref="scrollBarRef"
      >
        <div
          class="scroll-thumb"
          :style="scrollThumbStyle"
          @mousedown.stop="startThumbDrag"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useAppStore } from '@/stores/app'
import { ElMessage } from 'element-plus'
import { Close } from '@element-plus/icons-vue'
import { ReadLogFile } from 'wailsjs/go/main/App'

const appStore = useAppStore()

// 响应式数据
const logDisplayRef = ref<HTMLElement>()
const scrollBarRef = ref<HTMLElement>()
const currentLine = ref(1)
const currentColumn = ref(1)
const startLine = ref(0)
const visibleLineCount = ref(50)

// 滚动条拖拽相关
const isDragging = ref(false)
const dragStartY = ref(0)
const dragStartLine = ref(0)
const highlightedLines = ref<number[]>([])

// 搜索相关 - 这些变量在模板中通过 appStore 使用
// const searchResults = ref<Array<{ lineNumber: number, content: string }>>([])
// const currentSearchIndex = ref(0)

// 系统资源信息
const systemInfo = ref<any>(null)
let systemInfoTimer: number | null = null

// 日志数据
const logLines = ref<string[]>([])
const isLoading = ref(false)

// 计算属性
const visibleLines = computed(() => {
  const end = Math.min(startLine.value + visibleLineCount.value, logLines.value.length)
  return logLines.value.slice(startLine.value, end)
})

// 计算日志条数（识别典型的日志格式）
const logEntryCount = computed(() => {
  return logLines.value.filter(line => {
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
  }).length
})

const scrollThumbStyle = computed(() => {
  const totalLines = logLines.value.length
  const thumbHeight = Math.max((visibleLineCount.value / totalLines) * 100, 5)
  const thumbTop = (startLine.value / totalLines) * 100
  
  return {
    height: `${thumbHeight}%`,
    top: `${thumbTop}%`
  }
})

// 方法
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const selectLine = (lineIndex: number) => {
  currentLine.value = lineIndex + 1
  currentColumn.value = 1
  highlightedLines.value = [lineIndex]
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

  // 2. 时间戳高亮 (方括号内的日期时间，更精确的匹配)
  highlighted = highlighted.replace(
    /\[\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2}\s[+-]\d{4}\]/g,
    '<span class="syntax-timestamp">$&</span>'
  )

  // 3. IP地址高亮 (精确匹配IP地址格式，避免在User-Agent中误判)
  highlighted = highlighted.replace(
    /\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(?::\d+)?\b/g,
    '<span class="syntax-ip">$&</span>'
  )

  // 4. HTTP方法高亮 (在引号开头的HTTP方法)
  highlighted = highlighted.replace(
    /"(GET|POST|PUT|DELETE|HEAD|OPTIONS|PATCH)\s/g,
    '"<span class="syntax-method">$1</span> '
  )

  // 5. HTTP状态码高亮 (在HTTP版本后面的3位数字，避免User-Agent中的版本号)
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
  const newStartLine = Math.max(0, Math.min(startLine.value + delta, logLines.value.length - visibleLineCount.value))
  startLine.value = newStartLine
}

// 滚动条拖拽处理
const startThumbDrag = (event: MouseEvent) => {
  isDragging.value = true
  dragStartY.value = event.clientY
  dragStartLine.value = startLine.value

  document.addEventListener('mousemove', handleThumbDrag)
  document.addEventListener('mouseup', stopThumbDrag)
  event.preventDefault()
}

const handleThumbDrag = (event: MouseEvent) => {
  if (!isDragging.value || !scrollBarRef.value) return

  const deltaY = event.clientY - dragStartY.value
  const scrollBarHeight = scrollBarRef.value.clientHeight
  const totalLines = logLines.value.length
  const maxScrollLines = totalLines - visibleLineCount.value

  const deltaLines = (deltaY / scrollBarHeight) * totalLines
  const newStartLine = Math.max(0, Math.min(dragStartLine.value + deltaLines, maxScrollLines))

  startLine.value = Math.round(newStartLine)
}

const stopThumbDrag = () => {
  isDragging.value = false
  document.removeEventListener('mousemove', handleThumbDrag)
  document.removeEventListener('mouseup', stopThumbDrag)
}

const startScrollDrag = (event: MouseEvent) => {
  if (!scrollBarRef.value) return

  const rect = scrollBarRef.value.getBoundingClientRect()
  const clickY = event.clientY - rect.top
  const scrollBarHeight = rect.height
  const totalLines = logLines.value.length
  const maxScrollLines = totalLines - visibleLineCount.value

  const targetLine = (clickY / scrollBarHeight) * totalLines
  const newStartLine = Math.max(0, Math.min(targetLine, maxScrollLines))

  startLine.value = Math.round(newStartLine)
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
    if (!appStore.showLineNumbers || !appStore.wordWrap) return

    const lineNumbers = document.querySelectorAll('.line-number-container') as NodeListOf<HTMLElement>
    const logLines = document.querySelectorAll('.log-line') as NodeListOf<HTMLElement>

    lineNumbers.forEach((lineNumber, index) => {
      const logLine = logLines[index]
      if (logLine) {
        const logLineHeight = logLine.offsetHeight
        lineNumber.style.height = `${logLineHeight}px`
      }
    })
  })
}

// 加载日志文件内容
const loadLogFile = async (filePath: string) => {
  if (!filePath) return

  isLoading.value = true
  try {
    const content = await ReadLogFile(filePath)
    if (content && content.lines) {
      logLines.value = content.lines
      appStore.setLogContent(content.lines)
      ElMessage.success(`已加载 ${content.total} 行日志`)
      // 同步行号高度
      syncLineHeights()
    }
  } catch (error) {
    console.error('加载日志文件失败:', error)
    ElMessage.error('加载日志文件失败')
    logLines.value = []
  } finally {
    isLoading.value = false
  }
}

// 监听当前文件变化
watch(() => appStore.currentFile, (newFile) => {
  if (newFile) {
    loadLogFile(newFile.path)
  } else {
    logLines.value = []
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

// 跳转到指定行
const jumpToLine = (lineNumber: number, isSearchResult: boolean = false) => {
  const targetLine = Math.max(0, lineNumber - 1) // 转换为0基索引
  const maxStartLine = Math.max(0, logLines.value.length - visibleLineCount.value)

  // 计算最佳的起始行，让目标行显示在视口中间
  const idealStartLine = Math.max(0, targetLine - Math.floor(visibleLineCount.value / 2))
  startLine.value = Math.min(idealStartLine, maxStartLine)

  // 如果是搜索结果跳转，不需要额外的行高亮，因为搜索高亮已经足够明显
  if (!isSearchResult) {
    // 高亮目标行
    highlightedLines.value = [targetLine]

    // 3秒后清除高亮
    setTimeout(() => {
      highlightedLines.value = []
    }, 3000)
  }
}

// 生命周期
onMounted(() => {
  nextTick(() => {
    handleResize()
    if (logDisplayRef.value) {
      logDisplayRef.value.addEventListener('wheel', handleScroll, { passive: false })
    }
    window.addEventListener('resize', handleResize)

    // 监听跳转到行的事件
    window.addEventListener('jumpToLine', (event: any) => {
      jumpToLine(event.detail.lineNumber, event.detail.isSearchResult)
    })

    // 启动系统信息监控
    startSystemInfoMonitoring()
  })
})

onUnmounted(() => {
  if (logDisplayRef.value) {
    logDisplayRef.value.removeEventListener('wheel', handleScroll)
  }
  window.removeEventListener('resize', handleResize)

  // 清理拖拽事件监听器
  document.removeEventListener('mousemove', handleThumbDrag)
  document.removeEventListener('mouseup', stopThumbDrag)

  // 清理跳转事件监听器
  window.removeEventListener('jumpToLine', jumpToLine as any)

  // 停止系统信息监控
  stopSystemInfoMonitoring()
})
</script>

<style scoped>
.log-viewer {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #ffffff;
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
}

.line-number-container {
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  min-height: 20px;
  position: relative;
}

.line-number-content {
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  position: sticky;
  top: 0;
  background-color: inherit;
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
}

.log-line {
  height: 20px;
  display: flex;
  align-items: center;
  cursor: text;
  white-space: nowrap;
  overflow: hidden;
  padding: 0;
  line-height: 20px;
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
  cursor: pointer;
}

.scroll-thumb {
  position: absolute;
  right: 2px;
  width: 8px;
  background-color: #cbd5e1;
  border-radius: 4px;
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
</style>
