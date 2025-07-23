<template>
  <el-dialog
    v-model="visible"
    title="流程图预览"
    width="90%"
    :modal="true"
    class="flowchart-preview-dialog"
    @close="handleClose"
  >
    <div class="preview-container">
      <!-- 工具栏 -->
      <div class="preview-toolbar">
        <div class="toolbar-left">
          <el-button-group>
            <el-button
              :type="activeTab === 'mermaid' ? 'primary' : ''"
              @click="activeTab = 'mermaid'"
              size="small"
            >
              Mermaid 流程图
            </el-button>
            <el-button
              :type="activeTab === 'flowchart' ? 'primary' : ''"
              @click="activeTab = 'flowchart'"
              size="small"
            >
              Flowchart.js
            </el-button>
            <el-button
              :type="activeTab === 'table' ? 'primary' : ''"
              @click="activeTab = 'table'"
              size="small"
            >
              详细表格
            </el-button>
          </el-button-group>
        </div>
        <div class="toolbar-right">
          <el-button size="small" @click="copyCurrentContent">
            <el-icon><DocumentCopy /></el-icon>
            复制代码
          </el-button>
          <el-button size="small" type="primary" @click="downloadCurrentContent">
            <el-icon><Download /></el-icon>
            下载文件
          </el-button>
          <el-button
            v-if="activeTab === 'mermaid'"
            size="small"
            @click="openMermaidEditor"
          >
            <el-icon><View /></el-icon>
            在线预览
          </el-button>
        </div>
      </div>

      <!-- 内容区域 -->
      <div class="preview-content">
        <!-- Mermaid 代码预览 -->
        <div v-if="activeTab === 'mermaid'" class="code-panel">
          <div class="code-header">
            <h4>Mermaid 流程图</h4>
            <span class="code-info">{{ mermaidCode.split('\n').length }} 行代码</span>
            <div class="view-toggle">
              <el-button-group>
                <el-button
                  :type="mermaidViewMode === 'rendered' ? 'primary' : ''"
                  size="small"
                  @click="mermaidViewMode = 'rendered'"
                >
                  渲染视图
                </el-button>
                <el-button
                  :type="mermaidViewMode === 'code' ? 'primary' : ''"
                  size="small"
                  @click="mermaidViewMode = 'code'"
                >
                  代码视图
                </el-button>
              </el-button-group>
            </div>
          </div>

          <!-- 渲染视图 -->
          <div v-if="mermaidViewMode === 'rendered'" class="mermaid-render-container">
            <div class="render-controls">
              <el-button size="small" @click="zoomIn" title="放大">
                <el-icon><ZoomIn /></el-icon>
              </el-button>
              <el-button size="small" @click="zoomOut" title="缩小">
                <el-icon><ZoomOut /></el-icon>
              </el-button>
              <el-button size="small" @click="resetZoom" title="重置缩放">
                <el-icon><Refresh /></el-icon>
              </el-button>
              <span class="zoom-level">{{ Math.round(zoomLevel * 100) }}%</span>
            </div>
            <div class="mermaid-render-wrapper" :style="{ transform: `scale(${zoomLevel})` }">
              <div
                ref="mermaidContainer"
                class="mermaid-diagram"
                v-html="renderedMermaid"
              ></div>
            </div>
            <div v-if="renderError" class="render-error">
              <el-alert
                title="渲染失败"
                :description="renderError"
                type="error"
                show-icon
                :closable="false"
              />
            </div>
          </div>

          <!-- 代码视图 -->
          <div v-if="mermaidViewMode === 'code'" class="code-content">
            <pre><code>{{ mermaidCode }}</code></pre>
          </div>

          <div class="usage-tips">
            <h5>使用说明：</h5>
            <ul>
              <li>切换到"渲染视图"查看流程图效果</li>
              <li>复制代码到 <a href="https://mermaid-live.nodejs.cn" target="_blank">mermaid-live.nodejs.cn</a> 在线编辑</li>
              <li>在支持 Mermaid 的 Markdown 编辑器中使用</li>
              <li>集成到 GitBook、Notion、Typora 等文档工具</li>
            </ul>
          </div>
        </div>

        <!-- Flowchart.js 代码预览 -->
        <div v-if="activeTab === 'flowchart'" class="code-panel">
          <div class="code-header">
            <h4>Flowchart.js 代码</h4>
            <span class="code-info">{{ flowchartCode.split('\n').length }} 行代码</span>
          </div>
          <div class="code-content">
            <pre><code>{{ flowchartCode }}</code></pre>
          </div>
          <div class="usage-tips">
            <h5>使用说明：</h5>
            <ul>
              <li>使用 Flowchart.js 库在网页中渲染</li>
              <li>支持多种节点类型：operation、condition、inputoutput、subroutine</li>
              <li>可以自定义样式和主题</li>
            </ul>
          </div>
        </div>

        <!-- 详细表格预览 -->
        <div v-if="activeTab === 'table'" class="table-panel">
          <div class="table-header">
            <h4>时间线详细信息</h4>
            <span class="table-info">{{ timelineEntries.length }} 个步骤</span>
          </div>

          <!-- 统计信息 -->
          <div class="statistics">
            <div class="stat-item">
              <span class="stat-label">总步骤数</span>
              <span class="stat-value">{{ statistics.total }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">包含备注</span>
              <span class="stat-value">{{ statistics.withNotes }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">错误记录</span>
              <span class="stat-value">{{ statistics.errors }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">成功记录</span>
              <span class="stat-value">{{ statistics.success }}</span>
            </div>
          </div>

          <!-- 详细表格 -->
          <div class="table-content">
            <el-table :data="timelineEntries" stripe style="width: 100%">
              <el-table-column prop="index" label="步骤" width="80" align="center">
                <template #default="{ $index }">
                  {{ $index + 1 }}
                </template>
              </el-table-column>
              <el-table-column prop="logTimestamp" label="时间" width="180">
                <template #default="{ row }">
                  {{ row.logTimestamp ? formatLogTime(row.logTimestamp) : '-' }}
                </template>
              </el-table-column>
              <el-table-column prop="lineNumber" label="行号" width="100" align="center">
                <template #default="{ row }">
                  {{ row.lineNumber ? `第 ${row.lineNumber} 行` : '-' }}
                </template>
              </el-table-column>
              <el-table-column prop="note" label="备注" min-width="150">
                <template #default="{ row }">
                  <span v-if="row.note" class="note-content">{{ row.note }}</span>
                  <span v-else class="no-content">-</span>
                </template>
              </el-table-column>
              <el-table-column prop="logContent" label="日志内容" min-width="200">
                <template #default="{ row }">
                  <div v-if="row.logContent" class="log-content">
                    <span :class="getLogTypeClass(row.logContent)">
                      {{ getLogTypeIcon(row.logContent) }}
                    </span>
                    <span class="log-text">{{ row.logContent }}</span>
                  </div>
                  <span v-else class="no-content">-</span>
                </template>
              </el-table-column>
              <el-table-column prop="type" label="类型" width="100" align="center">
                <template #default="{ row }">
                  <el-tag :type="getLogTagType(row.logContent)" size="small">
                    {{ getLogTypeName(row.logContent) }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button @click="copyCurrentContent">复制代码</el-button>
        <el-button type="primary" @click="downloadCurrentContent">下载文件</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { DocumentCopy, Download, View, ZoomIn, ZoomOut, Refresh } from '@element-plus/icons-vue'
import mermaid from 'mermaid'

// Props
interface TimelineEntry {
  logTimestamp?: string
  lineNumber?: number
  note?: string
  logContent?: string
}

interface Props {
  modelValue: boolean
  timelineEntries: TimelineEntry[]
  mermaidCode: string
  flowchartCode: string
  markdownTable: string
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'copy': [content: string]
  'download': [content: string, filename: string]
}>()

// 响应式数据
const visible = ref(false)
const activeTab = ref('mermaid')

// Mermaid 相关状态
const mermaidViewMode = ref('rendered') // 'rendered' 或 'code'
const renderedMermaid = ref('')
const renderError = ref('')
const zoomLevel = ref(1)
const mermaidContainer = ref()

// 监听 modelValue 变化
watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
  if (newVal) {
    // 对话框打开时渲染图表
    nextTick(() => {
      renderMermaid()
    })
  }
})

// 监听 mermaidCode 变化
watch(() => props.mermaidCode, () => {
  if (visible.value && mermaidViewMode.value === 'rendered') {
    renderMermaid()
  }
})

// 监听视图模式变化
watch(mermaidViewMode, (newMode) => {
  if (newMode === 'rendered' && visible.value) {
    nextTick(() => {
      renderMermaid()
    })
  }
})

// 组件挂载时初始化
onMounted(() => {
  initMermaid()
})

// 统计信息
const statistics = computed(() => {
  const total = props.timelineEntries.length
  const withNotes = props.timelineEntries.filter(entry => entry.note).length
  const errors = props.timelineEntries.filter(entry =>
    entry.logContent && entry.logContent.toUpperCase().includes('ERROR')
  ).length
  const success = props.timelineEntries.filter(entry =>
    entry.logContent && (
      entry.logContent.toUpperCase().includes('SUCCESS') ||
      entry.logContent.toUpperCase().includes('OK')
    )
  ).length

  return { total, withNotes, errors, success }
})

// 当前内容
const currentContent = computed(() => {
  switch (activeTab.value) {
    case 'mermaid':
      return props.mermaidCode
    case 'flowchart':
      return props.flowchartCode
    case 'table':
      return props.markdownTable
    default:
      return ''
  }
})

// 当前文件名
const currentFilename = computed(() => {
  switch (activeTab.value) {
    case 'mermaid':
      return 'timeline-flowchart.mmd'
    case 'flowchart':
      return 'timeline-flowchart.txt'
    case 'table':
      return 'timeline-detailed-report.md'
    default:
      return 'timeline-export.txt'
  }
})

// 方法
const handleClose = () => {
  visible.value = false
}

const copyCurrentContent = () => {
  emit('copy', currentContent.value)
}

const downloadCurrentContent = () => {
  emit('download', currentContent.value, currentFilename.value)
}

const openMermaidEditor = () => {
  try {
    const encodedCode = btoa(unescape(encodeURIComponent(props.mermaidCode)))
    const url = `https://mermaid-live.nodejs.cn/edit#pako:${encodedCode}`
    window.open(url, '_blank')
  } catch (error) {
    window.open('https://mermaid-live.nodejs.cn/', '_blank')
    ElMessage.info('请手动粘贴代码到编辑器中')
  }
}

// 初始化 Mermaid
const initMermaid = () => {
  mermaid.initialize({
    startOnLoad: false,
    theme: 'default',
    securityLevel: 'loose',
    flowchart: {
      useMaxWidth: true,
      htmlLabels: true,
      curve: 'basis'
    },
    themeVariables: {
      primaryColor: '#409eff',
      primaryTextColor: '#ffffff',
      primaryBorderColor: '#409eff',
      lineColor: '#606266',
      sectionBkgColor: '#f5f7fa',
      altSectionBkgColor: '#ffffff',
      gridColor: '#e4e7ed',
      secondaryColor: '#67c23a',
      tertiaryColor: '#e6f7ff'
    }
  })
}

// 渲染 Mermaid 图表
const renderMermaid = async () => {
  if (!props.mermaidCode || !props.mermaidCode.trim()) {
    renderedMermaid.value = '<div class="empty-diagram">暂无流程图数据</div>'
    renderError.value = ''
    return
  }

  try {
    renderError.value = ''

    // 生成唯一ID
    const id = `mermaid-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`

    // 验证和渲染图表
    const { svg } = await mermaid.render(id, props.mermaidCode)
    renderedMermaid.value = svg

    // 等待DOM更新后调整样式
    await nextTick()
    adjustDiagramStyles()

  } catch (error) {
    console.error('Mermaid 渲染失败:', error)
    renderError.value = `渲染失败: ${error.message || '未知错误'}`
    renderedMermaid.value = ''
  }
}

// 调整图表样式
const adjustDiagramStyles = () => {
  if (mermaidContainer.value) {
    const svg = mermaidContainer.value.querySelector('svg')
    if (svg) {
      svg.style.maxWidth = '100%'
      svg.style.height = 'auto'
      svg.style.display = 'block'
      svg.style.margin = '0 auto'
    }
  }
}

// 缩放控制
const zoomIn = () => {
  if (zoomLevel.value < 3) {
    zoomLevel.value = Math.min(3, zoomLevel.value + 0.2)
  }
}

const zoomOut = () => {
  if (zoomLevel.value > 0.3) {
    zoomLevel.value = Math.max(0.3, zoomLevel.value - 0.2)
  }
}

const resetZoom = () => {
  zoomLevel.value = 1
}

// 格式化时间
const formatLogTime = (logTimestamp: string) => {
  if (!logTimestamp) return ''

  try {
    let date: Date

    // 处理 Apache 日志格式
    const apacheMatch = logTimestamp.match(/(\d{2})\/(\w{3})\/(\d{4}):(\d{2}):(\d{2}):(\d{2})\s*([+-]\d{4})?/)
    if (apacheMatch) {
      const [, day, monthStr, year, hour, minute, second] = apacheMatch
      const monthMap: { [key: string]: number } = {
        'Jan': 0, 'Feb': 1, 'Mar': 2, 'Apr': 3, 'May': 4, 'Jun': 5,
        'Jul': 6, 'Aug': 7, 'Sep': 8, 'Oct': 9, 'Nov': 10, 'Dec': 11
      }
      const month = monthMap[monthStr]
      if (month !== undefined) {
        date = new Date(parseInt(year), month, parseInt(day), parseInt(hour), parseInt(minute), parseInt(second))
      } else {
        date = new Date(logTimestamp)
      }
    } else {
      date = new Date(logTimestamp)
    }

    if (isNaN(date.getTime())) {
      return logTimestamp
    }

    const year = date.getFullYear()
    const month = date.getMonth() + 1
    const day = date.getDate()
    const hour = date.getHours().toString().padStart(2, '0')
    const minute = date.getMinutes().toString().padStart(2, '0')
    const second = date.getSeconds().toString().padStart(2, '0')

    return `${year}年${month}月${day}日 ${hour}:${minute}:${second}`
  } catch (error) {
    return logTimestamp
  }
}

// 获取日志类型图标
const getLogTypeIcon = (logContent: string) => {
  if (!logContent) return '📝'

  const upperContent = logContent.toUpperCase()
  if (upperContent.includes('ERROR')) return '❌'
  if (upperContent.includes('WARN')) return '⚠️'
  if (upperContent.includes('SUCCESS') || upperContent.includes('OK')) return '✅'
  if (upperContent.includes('INFO')) return 'ℹ️'
  return '📝'
}

// 获取日志类型样式类
const getLogTypeClass = (logContent: string) => {
  if (!logContent) return 'log-normal'

  const upperContent = logContent.toUpperCase()
  if (upperContent.includes('ERROR')) return 'log-error'
  if (upperContent.includes('WARN')) return 'log-warning'
  if (upperContent.includes('SUCCESS') || upperContent.includes('OK')) return 'log-success'
  if (upperContent.includes('INFO')) return 'log-info'
  return 'log-normal'
}

// 获取日志标签类型
const getLogTagType = (logContent: string) => {
  if (!logContent) return ''

  const upperContent = logContent.toUpperCase()
  if (upperContent.includes('ERROR')) return 'danger'
  if (upperContent.includes('WARN')) return 'warning'
  if (upperContent.includes('SUCCESS') || upperContent.includes('OK')) return 'success'
  if (upperContent.includes('INFO')) return 'info'
  return ''
}

// 获取日志类型名称
const getLogTypeName = (logContent: string) => {
  if (!logContent) return '普通'

  const upperContent = logContent.toUpperCase()
  if (upperContent.includes('ERROR')) return '错误'
  if (upperContent.includes('WARN')) return '警告'
  if (upperContent.includes('SUCCESS') || upperContent.includes('OK')) return '成功'
  if (upperContent.includes('INFO')) return '信息'
  return '普通'
}
</script>

<style scoped>
.flowchart-preview-dialog :deep(.el-dialog) {
  max-height: 90vh;
}

.flowchart-preview-dialog :deep(.el-dialog__body) {
  padding: 0;
  height: 70vh;
  overflow: hidden;
}

.preview-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.preview-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
  background: #f8f9fa;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.preview-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.code-panel,
.table-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.code-header,
.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
  background: #fff;
}

.view-toggle {
  display: flex;
  align-items: center;
  gap: 12px;
}

.code-header h4,
.table-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.code-info,
.table-info {
  font-size: 12px;
  color: #666;
  background: #f0f0f0;
  padding: 2px 8px;
  border-radius: 12px;
}

.code-content {
  flex: 1;
  overflow: auto;
  background: #f8f9fa;
}

.code-content pre {
  margin: 0;
  padding: 20px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  color: #333;
  background: transparent;
  white-space: pre-wrap;
  word-break: break-word;
}

.code-content.small pre {
  font-size: 12px;
  padding: 12px;
}

.usage-tips {
  padding: 16px;
  background: #fff;
  border-top: 1px solid #e5e7eb;
}

.usage-tips h5 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.usage-tips ul {
  margin: 0;
  padding-left: 20px;
}

.usage-tips li {
  margin-bottom: 8px;
  font-size: 13px;
  color: #666;
  line-height: 1.5;
}

.usage-tips a {
  color: #409eff;
  text-decoration: none;
}

.usage-tips a:hover {
  text-decoration: underline;
}

.statistics {
  display: flex;
  gap: 20px;
  padding: 16px;
  background: #f8f9fa;
  border-bottom: 1px solid #e5e7eb;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-label {
  font-size: 12px;
  color: #666;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.table-content {
  flex: 1;
  overflow: auto;
  padding: 16px;
}

.note-content {
  color: #065f46;
  background: #ecfdf5;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.no-content {
  color: #9ca3af;
  font-style: italic;
}

.log-content {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.log-text {
  flex: 1;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 12px;
  line-height: 1.4;
  word-break: break-all;
}

.log-error {
  color: #dc2626;
}

.log-warning {
  color: #d97706;
}

.log-success {
  color: #059669;
}

.log-info {
  color: #2563eb;
}

.log-normal {
  color: #374151;
}

.markdown-preview {
  margin-top: 16px;
  border-top: 1px solid #e5e7eb;
}

.markdown-preview .code-header h5 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* Mermaid 渲染相关样式 */
.mermaid-render-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.render-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #f8f9fa;
  border-bottom: 1px solid #e5e7eb;
}

.zoom-level {
  font-size: 12px;
  color: #666;
  margin-left: 8px;
  min-width: 40px;
}

.mermaid-render-wrapper {
  flex: 1;
  overflow: auto;
  padding: 20px;
  background: #fff;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  transform-origin: center top;
  transition: transform 0.3s ease;
}

.mermaid-diagram {
  max-width: 100%;
  overflow: visible;
}

.mermaid-diagram :deep(svg) {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 0 auto;
}

.empty-diagram {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
  color: #999;
  font-size: 14px;
  background: #f8f9fa;
  border: 2px dashed #ddd;
  border-radius: 8px;
}

.render-error {
  padding: 16px;
  background: #fff;
}

/* Mermaid 图表样式优化 */
.mermaid-diagram :deep(.node rect),
.mermaid-diagram :deep(.node circle),
.mermaid-diagram :deep(.node ellipse),
.mermaid-diagram :deep(.node polygon) {
  stroke-width: 2px;
}

.mermaid-diagram :deep(.edgePath path) {
  stroke-width: 2px;
}

.mermaid-diagram :deep(.edgeLabel) {
  background-color: #fff;
  border-radius: 4px;
  padding: 2px 4px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .mermaid-render-wrapper {
    padding: 10px;
  }

  .render-controls {
    padding: 8px 12px;
  }

  .render-controls .el-button {
    padding: 4px 8px;
  }
}
</style>