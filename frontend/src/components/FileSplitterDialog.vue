<template>
  <el-dialog
    v-model="visible"
    width="900px"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    @close="handleClose"
    class="splitter-dialog-wrapper"
    :top="'8vh'"
  >
    <template #header>
      <div class="dialog-header">
        <div class="header-icon">
          <el-icon size="24" color="#fa8c16"><Scissor /></el-icon>
        </div>
        <div class="header-content">
          <h2 class="dialog-title">文件分片处理</h2>
          <p class="dialog-subtitle">将大文件分割为多个小文件，便于处理和分析</p>
        </div>
      </div>
    </template>

    <div class="splitter-dialog">
      <!-- 文件信息卡片 -->
      <div class="file-info-card">
        <div class="card-header">
          <el-icon class="card-icon"><Document /></el-icon>
          <h3 class="card-title">文件信息</h3>
        </div>
        <div class="file-details">
          <div class="detail-grid">
            <div class="detail-item">
              <span class="detail-label">文件名</span>
              <span class="detail-value">{{ fileName }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">文件大小</span>
              <span class="detail-value size-highlight">{{ formatFileSize(fileSize) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">文件行数</span>
              <span class="detail-value lines-highlight">{{ fileLines ? formatNumber(fileLines) : '计算中...' }}</span>
            </div>
          </div>

          <!-- 日期范围信息 -->
          <div v-if="fileDateRange && fileDateRange.hasDateInfo" class="date-range-section">
            <div class="detail-grid">
              <div class="detail-item">
                <span class="detail-label">开始日期</span>
                <span class="detail-value date-highlight">{{ formatDate(fileDateRange.startDate) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">结束日期</span>
                <span class="detail-value date-highlight">{{ formatDate(fileDateRange.endDate) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">时间跨度</span>
                <span class="detail-value duration-highlight">{{ formatDuration(fileDateRange.totalDays) }}</span>
              </div>
            </div>
          </div>

          <div class="detail-row">
            <div class="detail-item full-width">
              <span class="detail-label">文件路径</span>
              <span class="detail-value path-value" :title="filePath">{{ filePath }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 分片策略选择 -->
      <div class="strategy-card">
        <div class="card-header">
          <el-icon class="card-icon"><Setting /></el-icon>
          <h3 class="card-title">分片策略</h3>
        </div>
        <div class="strategy-options">
          <div class="strategy-grid">
            <div
              v-for="strategy in strategyOptions"
              :key="strategy.value"
              class="strategy-option"
              :class="{ active: splitOptions.strategy === strategy.value }"
              @click="selectStrategy(strategy.value)"
            >
              <div class="strategy-header">
                <div class="strategy-icon">
                  <el-icon :size="24" :color="splitOptions.strategy === strategy.value ? '#fa8c16' : '#909399'">
                    <component :is="strategy.icon" />
                  </el-icon>
                </div>
                <div class="strategy-radio">
                  <el-radio
                    v-model="splitOptions.strategy"
                    :label="strategy.value"
                    @change="handleStrategyChange"
                  />
                </div>
              </div>
              <div class="strategy-content">
                <h4 class="strategy-name">{{ strategy.name }}</h4>
                <p class="strategy-desc">{{ strategy.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 配置区域 -->
      <div class="config-card">
        <div class="card-header">
          <el-icon class="card-icon"><Tools /></el-icon>
          <h3 class="card-title">分片配置</h3>
        </div>

        <!-- 按日期分片配置 -->
        <div v-if="splitOptions.strategy === 'date'" class="config-content">
          <div v-if="!fileDateRange || !fileDateRange.hasDateInfo" class="no-date-info">
            <el-alert
              title="无法检测到日期信息"
              description="该文件中未找到可识别的日期格式，无法进行日期分片"
              type="warning"
              :closable="false"
            />
          </div>
          <div v-else class="config-grid">
            <div class="config-group">
              <label class="config-label">
                <el-icon><Calendar /></el-icon>
                检测到的日期格式
              </label>
              <el-input
                :value="fileDateRange.datePattern"
                readonly
                class="config-input"
                placeholder="自动检测"
              />
            </div>
            <div class="config-group">
              <label class="config-label">
                <el-icon><Clock /></el-icon>
                每个文件包含天数
              </label>
              <el-input-number
                v-model="splitOptions.daysPerFile"
                :min="0.1"
                :max="365"
                :step="0.1"
                :precision="1"
                class="config-input"
                controls-position="right"
                placeholder="支持小数，如0.5天"
              />
            </div>
          </div>
        </div>

        <!-- 按大小分片配置 -->
        <div v-if="splitOptions.strategy === 'size'" class="config-content">
          <div class="config-grid">
            <div class="config-group">
              <label class="config-label">
                <el-icon><DataBoard /></el-icon>
                每个文件大小
              </label>
              <div class="size-input-group">
                <el-input-number
                  v-model="sizeValue"
                  :min="1"
                  :max="1000"
                  @change="updateSizePerFile"
                  class="size-number"
                  controls-position="right"
                />
                <el-select
                  v-model="sizeUnit"
                  @change="updateSizePerFile"
                  class="size-unit"
                >
                  <el-option label="MB" value="MB" />
                  <el-option label="GB" value="GB" />
                </el-select>
              </div>
            </div>
          </div>
        </div>

        <!-- 按行数分片配置 -->
        <div v-if="splitOptions.strategy === 'lines'" class="config-content">
          <div class="config-grid">
            <div class="config-group">
              <label class="config-label">
                <el-icon><List /></el-icon>
                每个文件行数
              </label>
              <el-input-number
                v-model="splitOptions.linesPerFile"
                :min="1000"
                :max="10000000"
                :step="1000"
                class="config-input"
                controls-position="right"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 输出设置 -->
      <div class="output-card">
        <div class="card-header">
          <el-icon class="card-icon"><FolderOpened /></el-icon>
          <h3 class="card-title">输出设置</h3>
        </div>
        <div class="output-content">
          <div class="config-group">
            <label class="config-label">
              <el-icon><Folder /></el-icon>
              输出目录
            </label>
            <div class="output-path-group">
              <el-input
                v-model="splitOptions.outputDir"
                placeholder="选择输出目录"
                readonly
                class="output-path-input"
              />
              <el-button
                @click="selectOutputDir"
                :icon="FolderOpened"
                type="primary"
                class="select-dir-btn"
              >
                选择目录
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 预估结果 -->
      <div v-if="estimatedResult" class="estimate-card">
        <div class="card-header">
          <el-icon class="card-icon"><DataAnalysis /></el-icon>
          <h3 class="card-title">预估结果</h3>
        </div>
        <div class="estimate-content">
          <div class="estimate-grid">
            <div class="estimate-item">
              <div class="estimate-icon">
                <el-icon size="20" color="#52c41a"><Files /></el-icon>
              </div>
              <div class="estimate-info">
                <span class="estimate-label">预计文件数</span>
                <span class="estimate-value">{{ estimatedResult.fileCount }} 个</span>
              </div>
            </div>
            <div class="estimate-item">
              <div class="estimate-icon">
                <el-icon size="20" color="#1890ff"><DataBoard /></el-icon>
              </div>
              <div class="estimate-info">
                <span class="estimate-label">平均文件大小</span>
                <span class="estimate-value">{{ formatFileSize(estimatedResult.avgFileSize) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <div class="footer-info">
          <el-icon><InfoFilled /></el-icon>
          <span>分片后的文件将保存到指定目录</span>
        </div>
        <div class="footer-actions">
          <el-button @click="handleClose" size="large">
            取消
          </el-button>
          <el-button
            type="primary"
            @click="startSplit"
            :loading="isSplitting"
            :disabled="!canSplit"
            size="large"
            class="split-btn"
          >
            <el-icon v-if="!isSplitting"><Scissor /></el-icon>
            {{ isSplitting ? '分片中...' : '开始分片' }}
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  FolderOpened,
  Scissor,
  Document,
  Setting,
  Tools,
  Calendar,
  Clock,
  DataBoard,
  List,
  Folder,
  DataAnalysis,
  Files,
  InfoFilled
} from '@element-plus/icons-vue'
import { SplitFile, GetCommonDatePatterns, OpenDirectoryDialogForSplit, GetFileLineCount, GetFileDateRange } from '../../wailsjs/go/main/App'

// Props
interface Props {
  modelValue: boolean
  filePath: string
  fileName: string
  fileSize: number
}

const props = defineProps<Props>()
const emit = defineEmits(['update:modelValue', 'split-complete'])

// 响应式数据
const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const isSplitting = ref(false)
const datePatterns = ref([])
const sizeValue = ref(100)
const sizeUnit = ref('MB')

// 策略选项
const strategyOptions = [
  {
    value: 'date',
    name: '按日期分片',
    description: '根据日志中的日期模式自动分片',
    icon: 'Calendar'
  },
  {
    value: 'size',
    name: '按大小分片',
    description: '按指定文件大小进行分片',
    icon: 'DataBoard'
  },
  {
    value: 'lines',
    name: '按行数分片',
    description: '按指定行数进行分片',
    icon: 'List'
  }
]

// 分片选项
const splitOptions = ref({
  strategy: 'size',
  filePath: props.filePath,
  outputDir: '',
  datePattern: '',
  daysPerFile: 1,
  sizePerFile: 100 * 1024 * 1024, // 100MB
  linesPerFile: 100000
})

// 预估结果
const estimatedResult = ref(null)
const fileLines = ref(0)
const fileDateRange = ref(null)

// 计算属性
const canSplit = computed(() => {
  if (!splitOptions.value.outputDir) return false

  switch (splitOptions.value.strategy) {
    case 'date':
      return fileDateRange.value && fileDateRange.value.hasDateInfo && splitOptions.value.daysPerFile > 0
    case 'size':
      return splitOptions.value.sizePerFile > 0
    case 'lines':
      return splitOptions.value.linesPerFile > 0
    default:
      return false
  }
})

// 方法
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatNumber = (num: number): string => {
  return num.toLocaleString()
}

const formatDate = (dateStr: string): string => {
  if (!dateStr) return '未知'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatDuration = (days: number): string => {
  if (days < 1) {
    const hours = Math.round(days * 24 * 10) / 10
    return `${hours} 小时`
  } else if (days < 30) {
    return `${Math.round(days * 10) / 10} 天`
  } else if (days < 365) {
    const months = Math.round(days / 30 * 10) / 10
    return `${months} 个月`
  } else {
    const years = Math.round(days / 365 * 10) / 10
    return `${years} 年`
  }
}

const selectStrategy = (strategy: string) => {
  splitOptions.value.strategy = strategy
  handleStrategyChange()
}

const handleStrategyChange = () => {
  estimatedResult.value = null
  calculateEstimate()
}

const updateSizePerFile = () => {
  const multiplier = sizeUnit.value === 'GB' ? 1024 * 1024 * 1024 : 1024 * 1024
  splitOptions.value.sizePerFile = sizeValue.value * multiplier
  calculateEstimate()
}

const calculateEstimate = () => {
  if (!props.fileSize) return
  
  let fileCount = 1
  let avgFileSize = props.fileSize
  
  switch (splitOptions.value.strategy) {
    case 'size':
      if (splitOptions.value.sizePerFile > 0) {
        fileCount = Math.ceil(props.fileSize / splitOptions.value.sizePerFile)
        avgFileSize = splitOptions.value.sizePerFile
      }
      break
    case 'lines':
      // 假设平均每行100字节
      if (splitOptions.value.linesPerFile > 0) {
        const avgLineSize = 100
        const totalLines = Math.ceil(props.fileSize / avgLineSize)
        fileCount = Math.ceil(totalLines / splitOptions.value.linesPerFile)
        avgFileSize = splitOptions.value.linesPerFile * avgLineSize
      }
      break
    case 'date':
      // 日期分片难以预估，显示提示
      fileCount = 0
      avgFileSize = 0
      break
  }
  
  estimatedResult.value = {
    fileCount,
    avgFileSize
  }
}

const selectOutputDir = async () => {
  try {
    const dir = await OpenDirectoryDialogForSplit()
    if (dir) {
      splitOptions.value.outputDir = dir
    }
  } catch (error) {
    console.error('选择目录失败:', error)
    ElMessage.error('选择目录失败')
  }
}

const startSplit = async () => {
  if (!canSplit.value) {
    ElMessage.warning('请完善分片配置')
    return
  }

  isSplitting.value = true
  try {
    const result = await SplitFile(splitOptions.value)
    
    if (result.success) {
      ElMessage.success(result.message)
      emit('split-complete', result)
      handleClose()
    } else {
      ElMessage.error(result.message)
    }
  } catch (error) {
    console.error('分片失败:', error)
    ElMessage.error('分片失败: ' + error)
  } finally {
    isSplitting.value = false
  }
}

const handleClose = () => {
  visible.value = false
}

// 获取文件行数
const getFileLineCount = async () => {
  if (!props.filePath) return

  try {
    fileLines.value = 0 // 重置为0，显示"计算中..."
    const lines = await GetFileLineCount(props.filePath)
    fileLines.value = lines
  } catch (error) {
    console.error('获取文件行数失败:', error)
    fileLines.value = 0
  }
}

// 获取文件日期范围
const getFileDateRange = async () => {
  if (!props.filePath) return

  try {
    fileDateRange.value = null
    const dateRange = await GetFileDateRange(props.filePath)
    fileDateRange.value = dateRange

    // 如果检测到日期信息，设置默认分片天数
    if (dateRange && dateRange.hasDateInfo) {
      // 根据总天数智能设置默认分片天数
      if (dateRange.totalDays <= 1) {
        splitOptions.value.daysPerFile = 0.5 // 半天
      } else if (dateRange.totalDays <= 7) {
        splitOptions.value.daysPerFile = 1 // 1天
      } else if (dateRange.totalDays <= 30) {
        splitOptions.value.daysPerFile = 7 // 1周
      } else {
        splitOptions.value.daysPerFile = 30 // 1个月
      }

      // 设置检测到的日期模式
      splitOptions.value.datePattern = dateRange.datePattern
    }
  } catch (error) {
    console.error('获取文件日期范围失败:', error)
    fileDateRange.value = { hasDateInfo: false }
  }
}

// 设置默认输出目录为文件所在目录
const setDefaultOutputDir = (filePath: string) => {
  if (!filePath) return

  // 获取文件所在目录
  const lastSlashIndex = Math.max(filePath.lastIndexOf('/'), filePath.lastIndexOf('\\'))
  if (lastSlashIndex > 0) {
    const fileDir = filePath.substring(0, lastSlashIndex)
    splitOptions.value.outputDir = fileDir
    console.log('🗂️ 设置默认输出目录:', fileDir)
  }
}

// 监听器
watch(() => props.filePath, (newPath) => {
  splitOptions.value.filePath = newPath
  if (newPath) {
    getFileLineCount()
    getFileDateRange()
    setDefaultOutputDir(newPath)
  }
})

watch(() => splitOptions.value, () => {
  calculateEstimate()
}, { deep: true })

// 组件挂载时加载日期模式
onMounted(async () => {
  try {
    datePatterns.value = await GetCommonDatePatterns()
    if (datePatterns.value.length > 0) {
      splitOptions.value.datePattern = datePatterns.value[0].pattern
    }
  } catch (error) {
    console.error('加载日期模式失败:', error)
  }
  
  // 设置默认输出目录为文件所在目录
  if (props.filePath) {
    setDefaultOutputDir(props.filePath)
    getFileLineCount()
    getFileDateRange()
  }
  
  calculateEstimate()
})
</script>

<style scoped>
/* 对话框包装器 */
:deep(.splitter-dialog-wrapper) {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.splitter-dialog-wrapper .el-dialog__header) {
  padding: 0;
  margin: 0;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.splitter-dialog-wrapper .el-dialog__body) {
  padding: 16px 20px;
}

:deep(.splitter-dialog-wrapper .el-dialog__footer) {
  padding: 12px 20px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
}

/* 对话框头部 */
.dialog-header {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  background: linear-gradient(135deg, #fa8c16 0%, #ff9c3e 100%);
  color: white;
}

.header-icon {
  margin-right: 12px;
}

.dialog-title {
  margin: 0 0 2px 0;
  font-size: 18px;
  font-weight: 600;
}

.dialog-subtitle {
  margin: 0;
  font-size: 13px;
  opacity: 0.9;
}

/* 主要内容区域 */
.splitter-dialog {
  max-height: 70vh;
  overflow-y: auto;
}

/* 卡片通用样式 */
.file-info-card,
.strategy-card,
.config-card,
.output-card,
.estimate-card {
  margin-bottom: 16px;
  background: white;
  border-radius: 8px;
  border: 1px solid #e8e8e8;
  overflow: hidden;
  transition: all 0.3s ease;
}

.file-info-card:hover,
.strategy-card:hover,
.config-card:hover,
.output-card:hover,
.estimate-card:hover {
  border-color: #fa8c16;
  box-shadow: 0 4px 12px rgba(250, 140, 22, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
}

.card-icon {
  margin-right: 8px;
  color: #fa8c16;
}

.card-title {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

/* 文件信息卡片 */
.file-details {
  padding: 16px;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 16px;
  margin-bottom: 12px;
}

.detail-row {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-item.full-width {
  flex: none;
  width: 100%;
}

.detail-label {
  font-size: 12px;
  color: #999;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.detail-value {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.size-highlight {
  color: #fa8c16;
  font-weight: 600;
}

.lines-highlight {
  color: #52c41a;
  font-weight: 600;
}

.date-highlight {
  color: #1890ff;
  font-weight: 600;
}

.duration-highlight {
  color: #722ed1;
  font-weight: 600;
}

.date-range-section {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.no-date-info {
  padding: 16px 0;
}

.path-value {
  word-break: break-all;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  background: #f5f5f5;
  padding: 4px 8px;
  border-radius: 4px;
}

/* 策略选择 */
.strategy-options {
  padding: 16px;
}

.strategy-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 12px;
}

.strategy-option {
  display: flex;
  flex-direction: column;
  padding: 16px;
  border: 2px solid #f0f0f0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: center;
  min-height: 100px;
}

.strategy-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.strategy-option:hover {
  border-color: #fa8c16;
  background: #fff7e6;
}

.strategy-option.active {
  border-color: #fa8c16;
  background: #fff7e6;
  box-shadow: 0 2px 8px rgba(250, 140, 22, 0.2);
}

.strategy-icon {
  flex-shrink: 0;
}

.strategy-content {
  flex: 1;
  text-align: center;
}

.strategy-name {
  margin: 0 0 4px 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.strategy-desc {
  margin: 0;
  font-size: 12px;
  color: #666;
  line-height: 1.3;
}

.strategy-radio {
  flex-shrink: 0;
}

/* 配置区域 */
.config-content {
  padding: 16px;
}

.config-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.config-group {
  margin-bottom: 16px;
}

.config-group:last-child {
  margin-bottom: 0;
}

.config-label {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
  font-size: 13px;
  font-weight: 500;
  color: #333;
}

.config-input {
  width: 100%;
}

.size-input-group {
  display: flex;
  gap: 12px;
  align-items: center;
}

.size-number {
  flex: 1;
}

.size-unit {
  width: 80px;
}

.pattern-option {
  padding: 4px 0;
}

.pattern-name {
  font-weight: 500;
  color: #333;
}

.pattern-desc {
  font-size: 12px;
  color: #999;
  margin: 2px 0;
}

.pattern-example {
  font-size: 11px;
  color: #666;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

/* 输出设置 */
.output-content {
  padding: 16px;
}

.output-path-group {
  display: flex;
  gap: 8px;
  align-items: center;
}

.output-path-input {
  flex: 1;
}

.select-dir-btn {
  flex-shrink: 0;
}

/* 预估结果 */
.estimate-content {
  padding: 16px;
}

.estimate-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.estimate-item {
  display: flex;
  align-items: center;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

.estimate-icon {
  margin-right: 8px;
  flex-shrink: 0;
}

.estimate-info {
  flex: 1;
}

.estimate-label {
  display: block;
  font-size: 11px;
  color: #666;
  margin-bottom: 2px;
}

.estimate-value {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

/* 对话框底部 */
.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-info {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #666;
  font-size: 13px;
}

.footer-actions {
  display: flex;
  gap: 8px;
}

.split-btn {
  background: linear-gradient(135deg, #fa8c16 0%, #ff9c3e 100%);
  border: none;
  font-weight: 600;
}

.split-btn:hover {
  background: linear-gradient(135deg, #e67e22 0%, #f39c12 100%);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .strategy-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .config-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .detail-grid {
    grid-template-columns: 1fr 1fr;
    gap: 16px;
  }
}

@media (max-width: 768px) {
  .estimate-grid {
    grid-template-columns: 1fr;
  }

  .detail-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .detail-row {
    flex-direction: column;
    gap: 12px;
  }

  .size-input-group,
  .output-path-group {
    flex-direction: column;
    align-items: stretch;
  }

  .dialog-footer {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }

  .footer-actions {
    justify-content: center;
  }

  .strategy-option {
    min-height: auto;
    padding: 16px;
  }
}
</style>
