<template>
  <!-- 时间线侧边栏（关闭时显示） -->
  <div
    class="timeline-sidebar"
    :class="{ 'timeline-sidebar-hidden': isVisible }"
    @click="openPanel"
    v-if="timelineEntries.length > 0"
  >
    <div class="sidebar-content">
      <div class="sidebar-title">时间线</div>
      <div class="sidebar-count">{{ timelineEntries.length }}</div>
    </div>
  </div>

  <!-- 时间线面板（展开时显示） -->
  <div class="timeline-panel" :class="{ 'timeline-panel-visible': isVisible }">
    <div class="timeline-header">
      <div class="header-left">
        <el-icon class="timeline-icon"><Clock /></el-icon>
        <span class="timeline-title">时间线</span>
        <span class="timeline-count">({{ timelineEntries.length }})</span>
      </div>
      <div class="header-actions">
        <el-button
          class="sort-btn"
          size="small"
          text
          @click="() => sortTimelineEntries(true)"
          title="按时间排序"
        >
          <el-icon><Sort /></el-icon>
        </el-button>
        <el-dropdown @command="handleExportCommand" trigger="click">
          <el-button
            class="export-btn"
            size="small"
            text
            title="导出流程图"
          >
            <el-icon><Download /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="preview">🎨 预览和渲染流程图</el-dropdown-item>
              <el-dropdown-item divided command="mermaid">📄 导出 Mermaid 文件</el-dropdown-item>
              <el-dropdown-item command="copy-mermaid">📋 复制 Mermaid 代码</el-dropdown-item>
              <el-dropdown-item divided command="flowchart">📄 导出 Flowchart.js 文件</el-dropdown-item>
              <el-dropdown-item command="copy-flowchart">📋 复制 Flowchart 代码</el-dropdown-item>
              <el-dropdown-item divided command="export-detailed">📊 导出详细信息表格</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button
          class="close-btn"
          size="small"
          text
          @click="closePanel"
          title="关闭时间线"
        >
          <el-icon><Close /></el-icon>
        </el-button>
      </div>
    </div>

    <div class="timeline-content">
      <div v-if="timelineEntries.length === 0" class="empty-timeline">
        <el-icon class="empty-icon"><Clock /></el-icon>
        <p>暂无时间线条目</p>
        <p class="empty-tip">点击工具栏的"+"按钮添加时间线条目</p>
      </div>

      <div v-else class="timeline-list">
        <div
          v-for="(entry, index) in timelineEntries"
          :key="entry.id"
          class="timeline-item"
        >
          <!-- 时间线节点 -->
          <div class="timeline-node">
            <div
              class="timeline-dot"
              :style="{ backgroundColor: entry.color }"
            ></div>
            <div class="timeline-line" v-if="index < timelineEntries.length - 1"></div>
          </div>

          <!-- 时间线内容 -->
          <div class="timeline-content-wrapper">
            <div class="timeline-card">
              <div class="timeline-header-info">
                <div class="timeline-title-section">
                  <span class="timeline-line-number" v-if="entry.lineNumber">行 {{ entry.lineNumber }}</span>
                  <span class="timeline-formatted-time" v-if="entry.logTimestamp">
                    {{ formatLogTime(entry.logTimestamp) }}
                  </span>
                </div>
                <div class="timeline-actions">
                  <el-button
                    size="small"
                    text
                    @click="editEntry(index)"
                    title="编辑"
                  >
                    <el-icon><Edit /></el-icon>
                  </el-button>
                  <el-button
                    size="small"
                    text
                    @click="deleteEntry(index)"
                    title="删除"
                  >
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
              </div>

              <!-- 日志内容 -->
              <div class="timeline-log-content" v-if="entry.logContent">
                <div class="log-content">{{ entry.logContent }}</div>
              </div>

              <!-- 备注内容 - 支持内联编辑 -->
              <div class="timeline-note">
                <!-- 编辑状态 -->
                <div v-if="editingIndex === index" class="note-editing">
                  <el-input
                    v-model="editingNote"
                    type="textarea"
                    :rows="2"
                    placeholder="输入备注..."
                    @keydown="handleKeyDown"
                    @blur="confirmEdit"
                    ref="noteInputRef"
                    class="note-input"
                  />
                  <div class="note-actions">
                    <el-button size="small" type="primary" @click="confirmEdit">保存</el-button>
                    <el-button size="small" @click="cancelEdit">取消</el-button>
                  </div>
                </div>
                <!-- 显示状态 -->
                <div v-else-if="entry.note" class="note-content" @click="editEntry(index)">
                  {{ entry.note }}
                </div>
                <!-- 空状态 -->
                <div v-else class="note-placeholder" @click="editEntry(index)">
                  点击添加备注...
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- 流程图预览组件 -->
  <FlowchartPreview
    v-model="showPreview"
    :timeline-entries="timelineEntries"
    :mermaid-code="previewMermaidCode"
    :flowchart-code="previewFlowchartCode"
    :markdown-table="previewMarkdownTable"
    @copy="handlePreviewCopy"
    @download="handlePreviewDownload"
  />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Close, Clock, Edit, Delete, Download, Sort } from '@element-plus/icons-vue'
import FlowchartPreview from './FlowchartPreview.vue'

const isVisible = ref(false)
const timelineEntries = ref([])

// 时间线排序函数
const sortTimelineEntries = (showMessage = false) => {
  const originalOrder = timelineEntries.value.map(entry => entry.id)

  timelineEntries.value.sort((a, b) => {
    // 首先尝试按日志时间戳排序
    if (a.logTimestamp && b.logTimestamp) {
      const timeA = parseLogTimestamp(a.logTimestamp)
      const timeB = parseLogTimestamp(b.logTimestamp)

      if (timeA && timeB) {
        return timeA.getTime() - timeB.getTime()
      }
    }

    // 如果没有日志时间戳，按行号排序
    if (a.lineNumber && b.lineNumber) {
      return a.lineNumber - b.lineNumber
    }

    // 如果都没有，按创建时间排序
    if (a.timestamp && b.timestamp) {
      return new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
    }

    // 最后按ID排序
    return (a.id || 0) - (b.id || 0)
  })

  // 检查是否有变化
  const newOrder = timelineEntries.value.map(entry => entry.id)
  const hasChanged = !originalOrder.every((id, index) => id === newOrder[index])

  if (showMessage) {
    if (hasChanged) {
      ElMessage.success('时间线已按时间顺序重新排序')
    } else {
      ElMessage.info('时间线已经是正确的时间顺序')
    }
  }
}

// 解析日志时间戳
const parseLogTimestamp = (logTimestamp: string): Date | null => {
  if (!logTimestamp) return null

  try {
    // 处理 Apache 日志格式：18/Jul/2025:00:53:14 +0800
    const apacheMatch = logTimestamp.match(/(\d{2})\/(\w{3})\/(\d{4}):(\d{2}):(\d{2}):(\d{2})\s*([+-]\d{4})?/)
    if (apacheMatch) {
      const [, day, monthStr, year, hour, minute, second] = apacheMatch
      const monthMap: { [key: string]: number } = {
        'Jan': 0, 'Feb': 1, 'Mar': 2, 'Apr': 3, 'May': 4, 'Jun': 5,
        'Jul': 6, 'Aug': 7, 'Sep': 8, 'Oct': 9, 'Nov': 10, 'Dec': 11
      }
      const month = monthMap[monthStr]
      if (month !== undefined) {
        return new Date(parseInt(year), month, parseInt(day), parseInt(hour), parseInt(minute), parseInt(second))
      }
    }

    // 处理 ISO 格式：2025-07-18T00:53:14 或 2025-07-18 00:53:14
    const isoMatch = logTimestamp.match(/(\d{4}-\d{2}-\d{2})[T\s](\d{2}:\d{2}:\d{2})/)
    if (isoMatch) {
      return new Date(`${isoMatch[1]}T${isoMatch[2]}`)
    }

    // 处理纯时间格式：00:53:14（假设是今天）
    const timeMatch = logTimestamp.match(/^(\d{2}):(\d{2}):(\d{2})$/)
    if (timeMatch) {
      const today = new Date()
      const [, hour, minute, second] = timeMatch
      return new Date(today.getFullYear(), today.getMonth(), today.getDate(), parseInt(hour), parseInt(minute), parseInt(second))
    }

    // 尝试直接解析
    const parsed = new Date(logTimestamp)
    return isNaN(parsed.getTime()) ? null : parsed
  } catch (error) {
    return null
  }
}

// 监听时间线更新事件
const handleTimelineUpdate = (event: CustomEvent) => {
  console.log('📅 TimelinePanel: 收到时间线更新事件:', {
    entriesCount: event.detail.entries.length,
    visible: event.detail.visible
  })

  timelineEntries.value = event.detail.entries
  // 自动排序时间线
  sortTimelineEntries()

  // 恢复面板显示状态
  if (event.detail.visible !== undefined) {
    isVisible.value = event.detail.visible
    console.log('📅 TimelinePanel: 设置面板可见性:', isVisible.value)
  } else if (timelineEntries.value.length > 0) {
    isVisible.value = true
    console.log('📅 TimelinePanel: 自动显示面板 (有条目)')
  }

  // 强制触发响应式更新
  nextTick(() => {
    if (timelineEntries.value.length > 0 && !isVisible.value) {
      isVisible.value = true
      console.log('📅 TimelinePanel: 强制显示面板')
    }
  })

  console.log('📅 TimelinePanel: 更新完成:', {
    entriesCount: timelineEntries.value.length,
    isVisible: isVisible.value
  })
}

// 监听时间线清除事件
const handleTimelineClear = () => {
  timelineEntries.value = []
  isVisible.value = false
}

// 打开面板
const openPanel = () => {
  isVisible.value = true
}

// 关闭面板
const closePanel = () => {
  isVisible.value = false
  // 通知保存状态
  window.dispatchEvent(new CustomEvent('timelinePanelVisibilityChanged', {
    detail: { visible: false }
  }))
}



// 格式化日志时间
const formatLogTime = (logTimestamp: string) => {
  if (!logTimestamp) return ''

  try {
    // 处理多种日志时间格式
    let date: Date

    // 处理 Apache 日志格式: 18/Jul/2025:00:59:07 +0800
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
      // 尝试直接解析其他格式
      date = new Date(logTimestamp)
    }

    if (isNaN(date.getTime())) {
      return logTimestamp // 如果解析失败，返回原始字符串
    }

    // 格式化为中文日期时间
    const year = date.getFullYear()
    const month = date.getMonth() + 1
    const day = date.getDate()
    const hour = date.getHours().toString().padStart(2, '0')
    const minute = date.getMinutes().toString().padStart(2, '0')
    const second = date.getSeconds().toString().padStart(2, '0')

    return `${year}年${month}月${day}日 ${hour}:${minute}:${second}`
  } catch (error) {
    return logTimestamp // 解析失败时返回原始字符串
  }
}

// 内联编辑相关状态
const editingIndex = ref(-1)
const editingNote = ref('')
const noteInputRef = ref()

// 预览相关状态
const showPreview = ref(false)
const previewMermaidCode = ref('')
const previewFlowchartCode = ref('')
const previewMarkdownTable = ref('')

// 开始编辑条目
const editEntry = (index: number) => {
  const entry = timelineEntries.value[index]
  if (!entry) return

  editingIndex.value = index
  editingNote.value = entry.note || ''

  // 下一帧聚焦输入框
  nextTick(() => {
    if (noteInputRef.value) {
      noteInputRef.value.focus()
    }
  })
}

// 确认编辑
const confirmEdit = () => {
  if (editingIndex.value >= 0 && editingIndex.value < timelineEntries.value.length) {
    const entry = timelineEntries.value[editingIndex.value]
    entry.note = editingNote.value.trim()

    // 自动排序时间线
    sortTimelineEntries()

    // 触发时间线更新事件（这会通知Toolbar保存状态）
    window.dispatchEvent(new CustomEvent('timelineUpdated', {
      detail: { entries: timelineEntries.value }
    }))

    // 通知Toolbar保存文件状态
    window.dispatchEvent(new CustomEvent('saveFileState'))

    ElMessage.success('备注已更新')
  }

  cancelEdit()
}

// 取消编辑
const cancelEdit = () => {
  editingIndex.value = -1
  editingNote.value = ''
}

// 处理键盘事件
const handleKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    confirmEdit()
  } else if (event.key === 'Escape') {
    event.preventDefault()
    cancelEdit()
  }
}

// 安全的文本转义函数
const escapeMermaidText = (text: string): string => {
  if (!text) return ''

  return text
    .replace(/"/g, "'")         // 双引号转为单引号
    .replace(/\n/g, ' ')        // 换行符转为空格
    .replace(/\r/g, '')         // 移除回车符
    .replace(/\t/g, ' ')        // 制表符转为空格
    .replace(/\s+/g, ' ')       // 多个空格合并为一个
    .trim()
}

// 生成 Mermaid 流程图代码
const generateMermaidFlowchart = () => {
  if (timelineEntries.value.length === 0) {
    return 'flowchart TD\n    A[没有时间线数据]'
  }

  let mermaidCode = 'flowchart TD\n'

  // 添加开始节点
  mermaidCode += '    Start([开始分析])\n'

  timelineEntries.value.forEach((entry, index) => {
    const nodeId = `Step${index + 1}`

    // 构建简化的节点标签，避免复杂的HTML
    let nodeLabel = ''

    // 1. 时间信息
    if (entry.logTimestamp) {
      nodeLabel = formatLogTime(entry.logTimestamp)
    } else {
      nodeLabel = `步骤 ${index + 1}`
    }

    // 2. 行号信息
    if (entry.lineNumber) {
      nodeLabel += ` (行 ${entry.lineNumber})`
    }

    // 3. 备注信息
    if (entry.note) {
      nodeLabel += ` - ${escapeMermaidText(entry.note)}`
    }

    // 4. 完整的日志内容（不截断）
    if (entry.logContent) {
      let logContent = entry.logContent.trim()

      // 转义处理，但保留完整内容
      logContent = escapeMermaidText(logContent)
      nodeLabel += ` | ${logContent}`
    }

    // 使用简单的矩形节点，避免特殊形状的解析问题
    const nodeShape = `${nodeId}["${nodeLabel}"]`

    mermaidCode += `    ${nodeShape}\n`

    // 添加连接线
    if (index === 0) {
      mermaidCode += `    Start --> ${nodeId}\n`
    } else {
      mermaidCode += `    Step${index} --> ${nodeId}\n`
    }
  })

  // 添加结束节点
  mermaidCode += '    End([分析完成])\n'
  mermaidCode += `    Step${timelineEntries.value.length} --> End\n`

  // 简化的样式定义
  mermaidCode += '\n    %% 样式定义\n'
  mermaidCode += '    classDef default fill:#f9f9f9,stroke:#333,stroke-width:2px\n'
  mermaidCode += '    classDef startEnd fill:#e3f2fd,stroke:#1976d2,stroke-width:3px\n'
  mermaidCode += '    classDef error fill:#ffebee,stroke:#d32f2f,stroke-width:2px\n'
  mermaidCode += '    classDef success fill:#e8f5e8,stroke:#388e3c,stroke-width:2px\n'
  mermaidCode += '    classDef warning fill:#fff8e1,stroke:#ffa000,stroke-width:2px\n'

  // 应用样式
  mermaidCode += '\n    %% 应用样式\n'
  mermaidCode += '    class Start,End startEnd\n'

  // 为不同类型的节点应用样式
  timelineEntries.value.forEach((entry, index) => {
    const nodeId = `Step${index + 1}`
    if (entry.logContent) {
      const logContent = entry.logContent.toUpperCase()
      if (logContent.includes('ERROR')) {
        mermaidCode += `    class ${nodeId} error\n`
      } else if (logContent.includes('SUCCESS') || logContent.includes('OK')) {
        mermaidCode += `    class ${nodeId} success\n`
      } else if (logContent.includes('WARN')) {
        mermaidCode += `    class ${nodeId} warning\n`
      }
    }
  })

  return mermaidCode
}

// 生成 Flowchart.js 代码
const generateFlowchartJs = () => {
  if (timelineEntries.value.length === 0) {
    return 'start=>start: 开始\nend=>end: 没有时间线数据\nstart->end'
  }

  let flowchartCode = ''
  let connections = []

  // 添加开始节点
  flowchartCode += 'start=>start: 时间线分析开始\n'

  timelineEntries.value.forEach((entry, index) => {
    const nodeId = `step${index + 1}`

    // 构建完整的节点标签
    let nodeLabel = ''

    // 1. 时间作为主标题
    if (entry.logTimestamp) {
      nodeLabel = formatLogTime(entry.logTimestamp)
    } else {
      nodeLabel = `步骤 ${index + 1}`
    }

    // 2. 添加行号
    if (entry.lineNumber) {
      nodeLabel += `\\n第 ${entry.lineNumber} 行`
    }

    // 3. 添加备注
    if (entry.note) {
      nodeLabel += `\\n备注: ${entry.note}`
    }

    // 4. 添加日志内容
    if (entry.logContent) {
      let logContent = entry.logContent.trim()

      // 处理特殊字符，确保在Flowchart.js中正确显示
      logContent = logContent
        .replace(/\\/g, '\\\\')   // 转义反斜杠
        .replace(/\n/g, ' ')      // 换行符转为空格
        .replace(/\r/g, '')       // 移除回车符
        .replace(/\t/g, ' ')      // 制表符转为空格
        .replace(/\s+/g, ' ')     // 多个空格合并为一个
        .trim()                   // 去除首尾空格

      // 显示完整内容，不截断
      nodeLabel += `\\n日志: ${logContent}`
    }

    // 根据内容类型选择节点类型
    let nodeType = 'operation'
    if (entry.logContent) {
      const logContent = entry.logContent.toUpperCase()
      if (logContent.includes('ERROR')) {
        nodeType = 'condition'  // 错误用判断框
      } else if (logContent.includes('SUCCESS') || logContent.includes('OK')) {
        nodeType = 'subroutine'  // 成功用子程序框
      } else if (entry.note) {
        nodeType = 'inputoutput'  // 有备注用输入输出框
      }
    } else if (entry.note) {
      nodeType = 'inputoutput'
    }

    // 添加节点定义
    flowchartCode += `${nodeId}=>${nodeType}: ${nodeLabel}\n`

    // 记录连接关系
    if (index === 0) {
      connections.push(`start->${nodeId}`)
    } else {
      connections.push(`step${index}->${nodeId}`)
    }
  })

  // 添加结束节点
  flowchartCode += 'end=>end: 分析完成\n'
  connections.push(`step${timelineEntries.value.length}->end`)

  // 添加注释说明
  flowchartCode += '\n// 节点类型说明:\n'
  flowchartCode += '// operation: 普通操作步骤\n'
  flowchartCode += '// condition: 错误或判断节点\n'
  flowchartCode += '// inputoutput: 备注或输入输出\n'
  flowchartCode += '// subroutine: 成功操作\n'

  // 添加连接定义
  flowchartCode += '\n// 流程连接:\n'
  flowchartCode += connections.join('\n')

  return flowchartCode
}

// 生成详细信息表格（Markdown格式）
const generateDetailedTable = () => {
  if (timelineEntries.value.length === 0) {
    return '# 时间线详细信息\n\n暂无时间线数据'
  }

  let markdownContent = '# 时间线详细信息\n\n'
  markdownContent += `> 生成时间：${new Date().toLocaleString()}\n\n`

  // 添加统计信息
  const totalEntries = timelineEntries.value.length
  const entriesWithNotes = timelineEntries.value.filter(entry => entry.note).length
  const errorEntries = timelineEntries.value.filter(entry =>
    entry.logContent && entry.logContent.toUpperCase().includes('ERROR')
  ).length

  markdownContent += '## 📊 统计概览\n\n'
  markdownContent += `- **总步骤数**: ${totalEntries}\n`
  markdownContent += `- **包含备注**: ${entriesWithNotes}\n`
  markdownContent += `- **错误记录**: ${errorEntries}\n\n`

  // 添加详细表格
  markdownContent += '## 📋 详细步骤\n\n'
  markdownContent += '| 步骤 | 时间 | 行号 | 备注 | 日志内容 | 类型 |\n'
  markdownContent += '|------|------|------|------|----------|------|\n'

  timelineEntries.value.forEach((entry, index) => {
    const stepNum = index + 1
    const time = entry.logTimestamp ? formatLogTime(entry.logTimestamp) : '-'
    const lineNum = entry.lineNumber ? `第 ${entry.lineNumber} 行` : '-'
    const note = entry.note ? entry.note.replace(/\|/g, '\\|') : '-'

    // 处理日志内容
    let logContent = '-'
    let logType = '普通'

    if (entry.logContent) {
      // 处理特殊字符，确保在Markdown表格中正确显示
      logContent = entry.logContent
        .replace(/\|/g, '\\|')     // 转义管道符
        .replace(/\n/g, ' ')       // 换行符转为空格
        .replace(/\r/g, '')        // 移除回车符

      // 显示完整内容，不截断

      const upperContent = entry.logContent.toUpperCase()
      if (upperContent.includes('ERROR')) {
        logType = '❌ 错误'
      } else if (upperContent.includes('WARN')) {
        logType = '⚠️ 警告'
      } else if (upperContent.includes('SUCCESS') || upperContent.includes('OK')) {
        logType = '✅ 成功'
      } else if (upperContent.includes('INFO')) {
        logType = 'ℹ️ 信息'
      }
    }

    markdownContent += `| ${stepNum} | ${time} | ${lineNum} | ${note} | ${logContent} | ${logType} |\n`
  })

  // 添加流程图代码
  markdownContent += '\n## 🔄 Mermaid 流程图\n\n'
  markdownContent += '```mermaid\n'
  markdownContent += generateMermaidFlowchart()
  markdownContent += '\n```\n\n'

  // 添加使用说明
  markdownContent += '## 📖 使用说明\n\n'
  markdownContent += '1. **查看流程图**: 将上面的 Mermaid 代码复制到 [mermaid.live](https://mermaid.live) 查看渲染效果\n'
  markdownContent += '2. **导入文档**: 可以将此 Markdown 内容导入到支持 Mermaid 的文档工具中\n'
  markdownContent += '3. **数据分析**: 使用表格数据进行进一步的分析和处理\n\n'

  markdownContent += '---\n'
  markdownContent += '*此文档由 LogTrawl 时间线功能自动生成*'

  return markdownContent
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('代码已复制到剪贴板')
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    ElMessage.success('代码已复制到剪贴板')
  }
}

// 下载文件
const downloadFile = (content: string, filename: string) => {
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
  ElMessage.success(`文件 ${filename} 已下载`)
}



// 预览流程图
const previewFlowchart = () => {
  previewMermaidCode.value = generateMermaidFlowchart()
  previewFlowchartCode.value = generateFlowchartJs()
  previewMarkdownTable.value = generateDetailedTable()
  showPreview.value = true
}

// 处理预览组件的复制事件
const handlePreviewCopy = (content: string) => {
  copyToClipboard(content)
}

// 处理预览组件的下载事件
const handlePreviewDownload = (content: string, filename: string) => {
  downloadFile(content, filename)
}

// 处理导出命令
const handleExportCommand = (command: string) => {
  switch (command) {
    case 'preview':
      previewFlowchart()
      break
    case 'mermaid':
      const mermaidCode = generateMermaidFlowchart()
      downloadFile(mermaidCode, 'timeline-flowchart.mmd')
      break
    case 'copy-mermaid':
      const mermaidCodeCopy = generateMermaidFlowchart()
      copyToClipboard(mermaidCodeCopy)
      break
    case 'flowchart':
      const flowchartCode = generateFlowchartJs()
      downloadFile(flowchartCode, 'timeline-flowchart.txt')
      break
    case 'copy-flowchart':
      const flowchartCodeCopy = generateFlowchartJs()
      copyToClipboard(flowchartCodeCopy)
      break
    case 'export-detailed':
      const detailedContent = generateDetailedTable()
      downloadFile(detailedContent, 'timeline-detailed-report.md')
      break
  }
}

// 删除条目
const deleteEntry = (index: number) => {
  timelineEntries.value.splice(index, 1)

  // 自动排序时间线（虽然删除后不需要排序，但保持一致性）
  sortTimelineEntries()

  // 触发时间线更新事件
  window.dispatchEvent(new CustomEvent('timelineUpdated', {
    detail: { entries: timelineEntries.value }
  }))

  // 通知Toolbar保存文件状态
  window.dispatchEvent(new CustomEvent('saveFileState'))

  ElMessage.success('时间线条目已删除')
}

onMounted(() => {
  window.addEventListener('timelineUpdated', handleTimelineUpdate)
  window.addEventListener('timelineCleared', handleTimelineClear)
})

onUnmounted(() => {
  window.removeEventListener('timelineUpdated', handleTimelineUpdate)
  window.removeEventListener('timelineCleared', handleTimelineClear)
})
</script>

<style scoped>
/* 时间线侧边栏 */
.timeline-sidebar {
  position: fixed;
  top: 72px;
  right: 0;
  width: 40px;
  height: calc(100vh - 72px);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  cursor: pointer;
  transition: all 0.3s ease;
  z-index: 500; /* 降低z-index，避免遮挡滚动条 */
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.1);
}

.timeline-sidebar:hover {
  width: 50px;
  background: linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%);
}

.timeline-sidebar-hidden {
  right: -40px;
}

.sidebar-content {
  writing-mode: vertical-rl;
  text-orientation: mixed;
  color: white;
  font-weight: 600;
  text-align: center;
}

.sidebar-title {
  font-size: 14px;
  margin-bottom: 8px;
}

.sidebar-count {
  font-size: 12px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 10px;
  padding: 2px 6px;
  margin-top: 8px;
}

/* 时间线面板 */
.timeline-panel {
  position: fixed;
  top: 72px;
  right: -400px;
  width: 360px;
  height: calc(100vh - 72px);
  background: #ffffff;
  border-left: 1px solid #e1e5e9;
  box-shadow: -4px 0 20px rgba(0, 0, 0, 0.1);
  transition: right 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 500; /* 与侧边栏相同层级 */
  display: flex;
  flex-direction: column;
}

.timeline-panel-visible {
  right: 0;
}

/* 时间线头部 */
.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  min-height: 56px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.timeline-icon {
  font-size: 20px;
}

.timeline-title {
  font-size: 18px;
  font-weight: 600;
}

.timeline-count {
  font-size: 14px;
  opacity: 0.8;
}

.sort-btn,
.export-btn,
.close-btn {
  width: 32px;
  height: 32px;
  padding: 0;
  color: white;
  border-radius: 50%;
  transition: all 0.3s;
}

.sort-btn:hover,
.export-btn:hover,
.close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: scale(1.1);
}

/* 时间线内容 */
.timeline-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
  background: #f8fafc;
}

.empty-timeline {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  color: #94a3b8;
  text-align: center;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  color: #cbd5e1;
}

.empty-timeline p {
  margin: 8px 0;
  font-size: 16px;
}

.empty-tip {
  font-size: 14px;
  color: #94a3b8;
}

/* 时间线列表 */
.timeline-list {
  position: relative;
}

.timeline-item {
  display: flex;
  margin-bottom: 16px;
  position: relative;
}

/* 时间线节点 */
.timeline-node {
  position: relative;
  margin-right: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.timeline-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #3b82f6;
  border: 3px solid #ffffff;
  box-shadow: 0 0 0 2px #e2e8f0;
  z-index: 2;
}

.timeline-line {
  width: 2px;
  height: 100%;
  background: #e2e8f0;
  margin-top: 8px;
  min-height: 40px;
}

/* 时间线内容卡片 */
.timeline-content-wrapper {
  flex: 1;
  margin-top: -6px;
}

.timeline-card {
  background: white;
  border-radius: 8px;
  padding: 12px 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border: 1px solid #e2e8f0;
  transition: all 0.3s ease;
}

.timeline-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.timeline-header-info {
  display: flex;
  align-items: center; /* 改为居中对齐，确保时间和按钮在同一水平线 */
  justify-content: space-between;
  margin-bottom: 12px;
}

.timeline-title-section {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap; /* 允许在空间不足时换行 */
}

/* 备用样式 - 垂直布局时间信息 */
.timeline-title-row {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.timeline-line-number {
  display: inline-block;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 12px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 12px;
  flex-shrink: 0; /* 防止被压缩 */
}

.timeline-formatted-time {
  font-size: 14px;
  color: #1f2937;
  font-weight: 600;
  white-space: nowrap; /* 防止时间换行 */
  flex-shrink: 0; /* 防止被压缩 */
}

.timeline-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.3s;
}

.timeline-card:hover .timeline-actions {
  opacity: 1;
}

.timeline-actions .el-button {
  width: 28px;
  height: 28px;
  padding: 0;
  border-radius: 6px;
}

/* 日志内容 */
.timeline-log-content {
  margin-bottom: 10px;
}

.log-content {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  color: #374151;
  background: #f8fafc;
  padding: 10px 12px;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
  white-space: pre-wrap;
  word-break: break-all;
  line-height: 1.4;
  max-height: 100px;
  overflow-y: auto;
}

/* 备注内容 */
.timeline-note {
  margin-top: 8px;
}

.note-content {
  font-size: 14px;
  color: #065f46;
  line-height: 1.5;
  background: #ecfdf5;
  padding: 8px 12px;
  border-radius: 6px;
  border-left: 4px solid #10b981;
  border: 1px solid #a7f3d0;
  cursor: pointer;
  transition: all 0.2s ease;
}

.note-content:hover {
  background: #d1fae5;
  border-color: #6ee7b7;
}

.note-placeholder {
  font-size: 14px;
  color: #9ca3af;
  line-height: 1.5;
  background: #f9fafb;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px dashed #d1d5db;
  cursor: pointer;
  transition: all 0.2s ease;
}

.note-placeholder:hover {
  background: #f3f4f6;
  border-color: #9ca3af;
  color: #6b7280;
}

.note-editing {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.note-input {
  font-size: 14px;
}

.note-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}



/* 日志时间 */
.timeline-log-time {
  margin-bottom: 8px;
}

.log-time-label {
  font-size: 12px;
  color: #64748b;
  font-weight: 600;
  margin-bottom: 4px;
}

.log-time-content {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  color: #059669;
  background: #ecfdf5;
  padding: 6px 10px;
  border-radius: 4px;
  border: 1px solid #a7f3d0;
  display: inline-block;
}

/* 日志预览 */
.timeline-log-preview {
  background: #fefce8;
  border-radius: 6px;
  padding: 8px;
  border-left: 4px solid #eab308;
  border: 1px solid #fde047;
}

.log-preview-header {
  font-size: 12px;
  color: #a16207;
  font-weight: 600;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.log-preview-header::before {
  content: "📋";
  font-size: 14px;
}

.log-preview-content {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 12px;
  color: #374151;
  background: white;
  padding: 8px;
  border-radius: 4px;
  border: 1px solid #e5e7eb;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 150px;
  overflow-y: auto;
  line-height: 1.3;
}

/* 滚动条样式 */
.timeline-content::-webkit-scrollbar {
  width: 6px;
}

.timeline-content::-webkit-scrollbar-track {
  background: #f1f5f9;
}

.timeline-content::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}

.timeline-content::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}


</style>
