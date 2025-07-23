<template>
  <div class="analysis-page">
    <!-- 页面头部 -->
    <div class="analysis-header">
      <div class="header-left">
        <el-button
          :icon="ArrowLeft"
          @click="goBack"
          type="text"
          class="back-btn"
        >
          返回
        </el-button>
        <h1 class="page-title">数据分析</h1>
        <span class="file-info">{{ currentFileName }}</span>
      </div>
      <div class="header-right">
        <el-dropdown @command="handleExport" :disabled="!hasAnalysisData">
          <el-button
            :icon="Download"
            :disabled="!hasAnalysisData"
            type="success"
          >
            导出数据
            <el-icon class="el-icon--right"><arrow-down /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="json">导出为 JSON</el-dropdown-item>
              <el-dropdown-item command="csv">导出为 CSV</el-dropdown-item>
              <el-dropdown-item command="excel">导出为 Excel</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button
          :icon="FolderOpened"
          @click="selectFile"
          :disabled="isAnalyzing"
          type="default"
        >
          选择文件
        </el-button>
        <el-button
          :icon="ScissorsIcon"
          @click="showSplitterDialog"
          :disabled="!currentFilePath || isAnalyzing"
          type="warning"
        >
          文件分片
        </el-button>
        <el-button
          :icon="Refresh"
          @click="refreshAnalysis"
          :loading="isAnalyzing"
          :disabled="!currentFilePath"
          type="primary"
        >
          {{ isAnalyzing ? '分析中...' : '重新分析' }}
        </el-button>
      </div>
    </div>

    <!-- 分析内容 -->
    <div class="analysis-content" v-loading="isAnalyzing" :element-loading-text="analysisStatus">
      <!-- 分析进度 -->
      <div v-if="isAnalyzing" class="analysis-progress">
        <div class="progress-info">
          <h3>{{ analysisStatus || '正在分析文件...' }}</h3>
          <p v-if="analysisProgress > 0">
            进度: {{ analysisProgress }}%
            <span v-if="analysisProgress < 80">- 正在处理大文件，请耐心等待...</span>
            <span v-else>- 即将完成，正在生成结果...</span>
          </p>
          <p v-else>正在初始化分析...</p>
        </div>
        <el-progress
          :percentage="analysisProgress"
          :stroke-width="12"
          :status="analysisProgress >= 100 ? 'success' : ''"
          :show-text="true"
          class="progress-bar"
        />
        <div class="progress-actions" v-if="canCancel">
          <el-button
            @click="cancelAnalysis"
            type="danger"
            size="small"
            :icon="Close"
          >
            取消分析
          </el-button>
        </div>
      </div>
      <!-- 概览统计 -->
      <div class="overview-section">
        <h2 class="section-title">概览统计</h2>
        <div class="overview-cards">
          <div class="stat-card">
            <div class="stat-value">{{ analysisData.overview.totalRequests.toLocaleString() }}</div>
            <div class="stat-label">总请求数</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ analysisData.overview.uniqueIPs.toLocaleString() }}</div>
            <div class="stat-label">独立IP数</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ analysisData.overview.internalIPs.toLocaleString() }}</div>
            <div class="stat-label">内网IP数</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ analysisData.overview.externalIPs.toLocaleString() }}</div>
            <div class="stat-label">外网IP数</div>
          </div>
        </div>
      </div>

      <!-- IP分析 -->
      <div class="ip-analysis-section">
        <h2 class="section-title">IP访问分析</h2>
        <el-tabs v-model="activeIPTab" class="analysis-tabs" @tab-change="handleIPTabChange">
          <!-- 内网IP TOP10 -->
          <el-tab-pane label="内网IP TOP10" name="internal">
            <div class="chart-container">
              <div ref="internalIPChart" class="chart"></div>
            </div>
            <div class="table-container">
              <el-table :data="analysisData.ipAnalysis.internalTop10" stripe style="width: 100%">
                <el-table-column prop="rank" label="排名" width="80" align="center" />
                <el-table-column prop="ip" label="IP地址" align="center" />
                <el-table-column prop="count" label="访问次数" width="120" align="center" />
                <el-table-column prop="firstAccess" label="首次访问" align="center" />
                <el-table-column prop="lastAccess" label="最后访问" align="center" />
              </el-table>
            </div>
          </el-tab-pane>

          <!-- 外网IP TOP10 -->
          <el-tab-pane label="外网IP TOP10" name="external">
            <div class="chart-container">
              <div ref="externalIPChart" class="chart"></div>
            </div>
            <div class="table-container">
              <el-table :data="analysisData.ipAnalysis.externalTop10" stripe style="width: 100%">
                <el-table-column prop="rank" label="排名" width="80" align="center" />
                <el-table-column prop="ip" label="IP地址" align="center" />
                <el-table-column prop="count" label="访问次数" width="120" align="center" />
                <el-table-column prop="firstAccess" label="首次访问" align="center" />
                <el-table-column prop="lastAccess" label="最后访问" align="center" />
              </el-table>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- URL分析 -->
      <div class="url-analysis-section">
        <h2 class="section-title">URL访问分析</h2>
        <el-tabs v-model="activeURLTab" class="analysis-tabs" @tab-change="handleURLTabChange">
          <!-- GET请求 TOP10 -->
          <el-tab-pane label="GET请求 TOP10" name="get">
            <div class="chart-container">
              <div ref="getURLChart" class="chart"></div>
            </div>
            <div class="table-container">
              <el-table :data="analysisData.urlAnalysis.getTop10" stripe style="width: 100%">
                <el-table-column prop="rank" label="排名" width="80" align="center" />
                <el-table-column prop="url" label="URL" show-overflow-tooltip />
                <el-table-column prop="count" label="访问次数" width="120" align="center" />
                <el-table-column prop="percentage" label="占比" width="100" align="center">
                  <template #default="scope">
                    {{ scope.row.percentage.toFixed(2) }}%
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-tab-pane>

          <!-- POST请求 TOP10 -->
          <el-tab-pane label="POST请求 TOP10" name="post">
            <div class="chart-container">
              <div ref="postURLChart" class="chart"></div>
            </div>
            <div class="table-container">
              <el-table :data="analysisData.urlAnalysis.postTop10" stripe style="width: 100%">
                <el-table-column prop="rank" label="排名" width="80" align="center" />
                <el-table-column prop="url" label="URL" show-overflow-tooltip />
                <el-table-column prop="count" label="访问次数" width="120" align="center" />
                <el-table-column prop="percentage" label="占比" width="100" align="center">
                  <template #default="scope">
                    {{ scope.row.percentage.toFixed(2) }}%
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 单独IP分析 -->
      <div class="single-ip-section">
        <h2 class="section-title">单独IP分析</h2>
        <div class="ip-search">
          <el-input
            v-model="searchIP"
            placeholder="请输入要分析的IP地址"
            class="ip-input"
            @keyup.enter="analyzeSpecificIP"
          >
            <template #append>
              <el-button @click="analyzeSpecificIP" :loading="isAnalyzingIP">分析</el-button>
            </template>
          </el-input>
        </div>
        
        <div v-if="specificIPData" class="ip-result">
          <div class="ip-info-cards">
            <div class="stat-card">
              <div class="stat-value">{{ specificIPData.totalRequests.toLocaleString() }}</div>
              <div class="stat-label">总访问次数</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ specificIPData.firstAccess }}</div>
              <div class="stat-label">首次访问</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ specificIPData.lastAccess }}</div>
              <div class="stat-label">最后访问</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ specificIPData.ipType }}</div>
              <div class="stat-label">IP类型</div>
            </div>
          </div>

          <div class="ip-url-analysis">
            <h3>访问URL分析</h3>
            <el-tabs v-model="activeSpecificIPTab" class="analysis-tabs" @tab-change="handleSpecificIPTabChange">
              <!-- GET请求 TOP10 -->
              <el-tab-pane label="GET请求 TOP10" name="get">
                <div class="chart-container">
                  <div ref="specificIPGetChart" class="chart"></div>
                </div>
                <div class="table-container">
                  <el-table :data="specificIPData.urlAnalysis.getTop10" stripe style="width: 100%">
                    <el-table-column prop="rank" label="排名" width="80" align="center" />
                    <el-table-column prop="url" label="URL" show-overflow-tooltip />
                    <el-table-column prop="count" label="访问次数" width="120" align="center" />
                    <el-table-column prop="percentage" label="占比" width="100" align="center">
                      <template #default="scope">
                        {{ scope.row.percentage.toFixed(2) }}%
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </el-tab-pane>

              <!-- POST请求 TOP10 -->
              <el-tab-pane label="POST请求 TOP10" name="post">
                <div class="chart-container">
                  <div ref="specificIPPostChart" class="chart"></div>
                </div>
                <div class="table-container">
                  <el-table :data="specificIPData.urlAnalysis.postTop10" stripe style="width: 100%">
                    <el-table-column prop="rank" label="排名" width="80" align="center" />
                    <el-table-column prop="url" label="URL" show-overflow-tooltip />
                    <el-table-column prop="count" label="访问次数" width="120" align="center" />
                    <el-table-column prop="percentage" label="占比" width="100" align="center">
                      <template #default="scope">
                        {{ scope.row.percentage.toFixed(2) }}%
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </div>
      </div>
    </div>

    <!-- 文件分片对话框 -->
    <FileSplitterDialog
      v-model="showSplitter"
      :file-path="currentFilePath"
      :file-name="currentFileName"
      :file-size="currentFileSize"
      @split-complete="handleSplitComplete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Refresh, FolderOpened, Close, Download, ArrowDown, Scissor as ScissorsIcon } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { AnalyzeSpecificIP, OpenFileDialog, AnalyzeLogFileWithProgress, StartAnalysis, CancelAnalysis, GetFileInfo } from '../../wailsjs/go/main/App'
import { useAppStore } from '../stores/app'
import FileSplitterDialog from './FileSplitterDialog.vue'

// Props
interface Props {
  filePath: string
  fileName: string
}

const props = defineProps<Props>()

// Store
const appStore = useAppStore()

// 响应式数据
const isAnalyzing = ref(false)
const isAnalyzingIP = ref(false)
const activeIPTab = ref('internal')
const activeURLTab = ref('get')
const activeSpecificIPTab = ref('get')
const searchIP = ref('')

// 文件管理
const currentFilePath = ref(props.filePath || '')
const currentFileName = ref(props.fileName || '')
const currentFileSize = ref(0)
const analysisProgress = ref(0)
const analysisStatus = ref('')
const sessionID = ref('')
const canCancel = ref(false)
const showSplitter = ref(false)

// 分析数据
const analysisData = ref({
  overview: {
    totalRequests: 0,
    uniqueIPs: 0,
    internalIPs: 0,
    externalIPs: 0
  },
  ipAnalysis: {
    internalTop10: [],
    externalTop10: []
  },
  urlAnalysis: {
    getTop10: [],
    postTop10: []
  }
})

const specificIPData = ref(null)

// 图表引用
const internalIPChart = ref()
const externalIPChart = ref()
const getURLChart = ref()
const postURLChart = ref()
const specificIPGetChart = ref()
const specificIPPostChart = ref()

// 图表实例管理
const chartInstances = {
  externalIP: null,
  postURL: null,
  specificIPGet: null,
  specificIPPost: null
}

// 计算属性
const hasAnalysisData = computed(() => {
  return analysisData.value.overview.totalRequests > 0
})

// 返回主页面
const goBack = () => {
  window.dispatchEvent(new CustomEvent('closeAnalysisPage'))
}

// 导出数据
const handleExport = async (format: string) => {
  if (!hasAnalysisData.value) {
    ElMessage.warning('没有可导出的分析数据')
    return
  }

  try {
    const timestamp = new Date().toISOString().slice(0, 19).replace(/:/g, '-')
    const fileName = `分析报告_${currentFileName.value}_${timestamp}`

    switch (format) {
      case 'json':
        await exportAsJSON(fileName)
        break
      case 'csv':
        await exportAsCSV(fileName)
        break
      case 'excel':
        await exportAsExcel(fileName)
        break
      default:
        ElMessage.error('不支持的导出格式')
    }
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败: ' + error)
  }
}

// 导出为JSON
const exportAsJSON = async (fileName: string) => {
  const exportData = {
    metadata: {
      fileName: currentFileName.value,
      filePath: currentFilePath.value,
      exportTime: new Date().toISOString(),
      version: '1.0'
    },
    analysis: analysisData.value
  }

  const jsonStr = JSON.stringify(exportData, null, 2)
  downloadFile(jsonStr, `${fileName}.json`, 'application/json')
  ElMessage.success('JSON 文件导出成功')
}

// 导出为CSV
const exportAsCSV = async (fileName: string) => {
  let csvContent = ''

  // 概览统计
  csvContent += '概览统计\n'
  csvContent += '指标,数值\n'
  csvContent += `总请求数,${analysisData.value.overview.totalRequests}\n`
  csvContent += `独立IP数,${analysisData.value.overview.uniqueIPs}\n`
  csvContent += `内网IP数,${analysisData.value.overview.internalIPs}\n`
  csvContent += `外网IP数,${analysisData.value.overview.externalIPs}\n\n`

  // 内网IP TOP10
  csvContent += '内网IP TOP10\n'
  csvContent += '排名,IP地址,访问次数,首次访问,最后访问\n'
  analysisData.value.ipAnalysis.internalTop10.forEach(item => {
    csvContent += `${item.rank},${item.ip},${item.count},${item.firstAccess},${item.lastAccess}\n`
  })
  csvContent += '\n'

  // 外网IP TOP10
  csvContent += '外网IP TOP10\n'
  csvContent += '排名,IP地址,访问次数,首次访问,最后访问\n'
  analysisData.value.ipAnalysis.externalTop10.forEach(item => {
    csvContent += `${item.rank},${item.ip},${item.count},${item.firstAccess},${item.lastAccess}\n`
  })
  csvContent += '\n'

  // GET请求 TOP10
  csvContent += 'GET请求 TOP10\n'
  csvContent += '排名,URL,访问次数,占比(%)\n'
  analysisData.value.urlAnalysis.getTop10.forEach(item => {
    csvContent += `${item.rank},${item.url},${item.count},${item.percentage.toFixed(2)}\n`
  })
  csvContent += '\n'

  // POST请求 TOP10
  csvContent += 'POST请求 TOP10\n'
  csvContent += '排名,URL,访问次数,占比(%)\n'
  analysisData.value.urlAnalysis.postTop10.forEach(item => {
    csvContent += `${item.rank},${item.url},${item.count},${item.percentage.toFixed(2)}\n`
  })

  // 添加BOM以支持中文
  const csvWithBOM = '\uFEFF' + csvContent
  downloadFile(csvWithBOM, `${fileName}.csv`, 'text/csv')
  ElMessage.success('CSV 文件导出成功')
}

// 导出为Excel (使用HTML表格格式，Excel可以打开)
const exportAsExcel = async (fileName: string) => {
  let htmlContent = `
    <html>
    <head>
      <meta charset="utf-8">
      <style>
        table { border-collapse: collapse; width: 100%; margin-bottom: 20px; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; font-weight: bold; }
        h2 { color: #333; margin-top: 20px; }
      </style>
    </head>
    <body>
      <h1>日志分析报告</h1>
      <p>文件名: ${currentFileName.value}</p>
      <p>导出时间: ${new Date().toLocaleString()}</p>

      <h2>概览统计</h2>
      <table>
        <tr><th>指标</th><th>数值</th></tr>
        <tr><td>总请求数</td><td>${analysisData.value.overview.totalRequests}</td></tr>
        <tr><td>独立IP数</td><td>${analysisData.value.overview.uniqueIPs}</td></tr>
        <tr><td>内网IP数</td><td>${analysisData.value.overview.internalIPs}</td></tr>
        <tr><td>外网IP数</td><td>${analysisData.value.overview.externalIPs}</td></tr>
      </table>

      <h2>内网IP TOP10</h2>
      <table>
        <tr><th>排名</th><th>IP地址</th><th>访问次数</th><th>首次访问</th><th>最后访问</th></tr>
        ${analysisData.value.ipAnalysis.internalTop10.map(item =>
          `<tr><td>${item.rank}</td><td>${item.ip}</td><td>${item.count}</td><td>${item.firstAccess}</td><td>${item.lastAccess}</td></tr>`
        ).join('')}
      </table>

      <h2>外网IP TOP10</h2>
      <table>
        <tr><th>排名</th><th>IP地址</th><th>访问次数</th><th>首次访问</th><th>最后访问</th></tr>
        ${analysisData.value.ipAnalysis.externalTop10.map(item =>
          `<tr><td>${item.rank}</td><td>${item.ip}</td><td>${item.count}</td><td>${item.firstAccess}</td><td>${item.lastAccess}</td></tr>`
        ).join('')}
      </table>

      <h2>GET请求 TOP10</h2>
      <table>
        <tr><th>排名</th><th>URL</th><th>访问次数</th><th>占比(%)</th></tr>
        ${analysisData.value.urlAnalysis.getTop10.map(item =>
          `<tr><td>${item.rank}</td><td>${item.url}</td><td>${item.count}</td><td>${item.percentage.toFixed(2)}</td></tr>`
        ).join('')}
      </table>

      <h2>POST请求 TOP10</h2>
      <table>
        <tr><th>排名</th><th>URL</th><th>访问次数</th><th>占比(%)</th></tr>
        ${analysisData.value.urlAnalysis.postTop10.map(item =>
          `<tr><td>${item.rank}</td><td>${item.url}</td><td>${item.count}</td><td>${item.percentage.toFixed(2)}</td></tr>`
        ).join('')}
      </table>
    </body>
    </html>
  `

  downloadFile(htmlContent, `${fileName}.xls`, 'application/vnd.ms-excel')
  ElMessage.success('Excel 文件导出成功')
}

// 下载文件
const downloadFile = (content: string, fileName: string, mimeType: string) => {
  const blob = new Blob([content], { type: mimeType })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = fileName
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

// 选择文件
const selectFile = async () => {
  try {
    const filePath = await OpenFileDialog()
    if (filePath) {
      // 检查文件大小
      try {
        const fileInfo = await GetFileInfo(filePath)
        const fileSizeGB = fileInfo.size / (1024 * 1024 * 1024)

        if (fileSizeGB > 1) {
          // 超过1GB，提示使用分片功能
          ElMessage.warning('文件大小超过1GB，建议使用分片功能处理')
          currentFilePath.value = filePath
          currentFileName.value = filePath.split(/[/\\]/).pop() || ''
          currentFileSize.value = fileInfo.size
          showSplitter.value = true
          return
        }

        currentFileSize.value = fileInfo.size
      } catch (error) {
        console.warn('获取文件信息失败，继续处理:', error)
      }

      currentFilePath.value = filePath
      currentFileName.value = filePath.split(/[/\\]/).pop() || ''
      ElMessage.success(`已选择文件: ${currentFileName.value}`)
      // 自动开始分析
      await performAnalysis()
    }
  } catch (error) {
    console.error('选择文件失败:', error)
    ElMessage.error('选择文件失败')
  }
}

// 显示分片对话框
const showSplitterDialog = () => {
  if (!currentFilePath.value) {
    ElMessage.warning('请先选择文件')
    return
  }
  showSplitter.value = true
}

// 处理分片完成
const handleSplitComplete = (result) => {
  console.log('分片完成:', result)
  ElMessage.success(`文件已分片为 ${result.totalFiles} 个文件`)
  // 可以在这里提供选择分片文件进行分析的选项
}

// 刷新分析
const refreshAnalysis = async () => {
  await performAnalysis(true) // 强制重新分析
}

// 取消分析
const cancelAnalysis = async () => {
  if (sessionID.value) {
    try {
      await CancelAnalysis(sessionID.value)
      ElMessage.info('分析已取消')
    } catch (error) {
      console.error('取消分析失败:', error)
    }
  }
  isAnalyzing.value = false
  canCancel.value = false
  analysisProgress.value = 0
  analysisStatus.value = ''
}

// 生成会话ID
const generateSessionID = () => {
  return 'analysis_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
}

// 执行分析
const performAnalysis = async (forceRefresh = false) => {
  if (!currentFilePath.value) {
    ElMessage.warning('请先选择要分析的文件')
    return
  }

  // 检查缓存
  if (!forceRefresh && appStore.hasAnalysisData(currentFilePath.value)) {
    console.log('📋 使用缓存的分析数据:', currentFilePath.value)
    const cachedData = appStore.getAnalysisData(currentFilePath.value)
    if (cachedData) {
      analysisData.value = cachedData.data
      await nextTick()
      renderCharts()
      ElMessage.success('已加载缓存的分析结果')
      return
    }
  }

  isAnalyzing.value = true
  analysisProgress.value = 0
  analysisStatus.value = '正在初始化分析...'
  canCancel.value = true
  sessionID.value = generateSessionID()

  try {
    console.log('🔍 开始分析文件:', currentFilePath.value)

    // 启动分析会话
    await StartAnalysis(sessionID.value)

    // 模拟进度更新（实际应该通过WebSocket或轮询获取）
    const progressInterval = setInterval(() => {
      if (analysisProgress.value < 95 && isAnalyzing.value) {
        analysisProgress.value += Math.random() * 5
        if (analysisProgress.value < 30) {
          analysisStatus.value = '正在读取文件...'
        } else if (analysisProgress.value < 70) {
          analysisStatus.value = '正在解析日志...'
        } else if (analysisProgress.value < 90) {
          analysisStatus.value = '正在统计数据...'
        } else {
          analysisStatus.value = '正在生成结果...'
        }
      }
    }, 500)

    // 调用后端分析API
    const result = await AnalyzeLogFileWithProgress(currentFilePath.value, sessionID.value)

    clearInterval(progressInterval)

    if (result) {
      // 最终进度更新
      analysisProgress.value = 100
      analysisStatus.value = '分析完成'

      // 确保数据结构完整性
      analysisData.value = {
        overview: result.overview || {
          totalRequests: 0,
          uniqueIPs: 0,
          internalIPs: 0,
          externalIPs: 0
        },
        ipAnalysis: {
          internalTop10: result.ipAnalysis?.internalTop10 || [],
          externalTop10: result.ipAnalysis?.externalTop10 || []
        },
        urlAnalysis: {
          getTop10: result.urlAnalysis?.getTop10 || [],
          postTop10: result.urlAnalysis?.postTop10 || []
        }
      }

      console.log('📊 分析数据:', analysisData.value)

      // 保存到缓存
      appStore.setAnalysisData(currentFilePath.value, currentFileName.value, analysisData.value)

      // 等待DOM更新后渲染图表
      await nextTick()
      renderCharts()

      // 延迟显示成功消息，让用户看到100%进度
      setTimeout(() => {
        ElMessage.success('分析完成')
      }, 500)
    }
  } catch (error) {
    console.error('分析失败:', error)
    ElMessage.error('分析失败: ' + error)
    if (sessionID.value) {
      await CancelAnalysis(sessionID.value)
    }
  } finally {
    isAnalyzing.value = false
    canCancel.value = false
    // 保持进度显示一段时间
    setTimeout(() => {
      if (!isAnalyzing.value) {
        analysisProgress.value = 0
        analysisStatus.value = ''
      }
    }, 2000)
  }
}

// 分析特定IP
const analyzeSpecificIP = async () => {
  if (!searchIP.value.trim()) {
    ElMessage.warning('请输入IP地址')
    return
  }
  
  isAnalyzingIP.value = true
  try {
    console.log('🔍 分析特定IP:', searchIP.value)
    
    const result = await AnalyzeSpecificIP(currentFilePath.value, searchIP.value.trim())

    if (result) {
      // 确保数据结构完整性
      specificIPData.value = {
        ip: result.ip || searchIP.value,
        ipType: result.ipType || '未知',
        totalRequests: result.totalRequests || 0,
        firstAccess: result.firstAccess || '未知',
        lastAccess: result.lastAccess || '未知',
        urlAnalysis: {
          getTop10: result.urlAnalysis?.getTop10 || [],
          postTop10: result.urlAnalysis?.postTop10 || []
        }
      }

      console.log('🔍 特定IP分析数据:', specificIPData.value)

      await nextTick()
      renderSpecificIPCharts()

      ElMessage.success('IP分析完成')
    }
  } catch (error) {
    console.error('IP分析失败:', error)
    ElMessage.error('IP分析失败: ' + error)
  } finally {
    isAnalyzingIP.value = false
  }
}

// 渲染图表
const renderCharts = () => {
  try {
  // 渲染内网IP图表
  if (internalIPChart.value && analysisData.value.ipAnalysis.internalTop10.length > 0) {
    const chart = echarts.init(internalIPChart.value)
    chart.setOption({
      title: {
        text: '内网IP访问次数 TOP10',
        left: 'center',
        textStyle: { fontSize: 16 }
      },
      tooltip: {
        trigger: 'axis',
        formatter: '{b}<br/>访问次数: {c}'
      },
      grid: {
        left: '10%',
        right: '10%',
        bottom: '15%',
        top: '15%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: analysisData.value.ipAnalysis.internalTop10.map(item => item.ip),
        axisLabel: {
          rotate: 45,
          fontSize: 12
        }
      },
      yAxis: {
        type: 'value',
        name: '访问次数',
        nameTextStyle: { fontSize: 12 }
      },
      series: [{
        data: analysisData.value.ipAnalysis.internalTop10.map(item => item.count),
        type: 'bar',
        itemStyle: { color: '#409EFF' },
        label: {
          show: true,
          position: 'top',
          fontSize: 11
        },
        barWidth: '60%'
      }]
    })

    // 响应式调整
    window.addEventListener('resize', () => chart.resize())
  }

  // 外网IP图表在标签页切换时渲染，这里跳过

  // 渲染GET URL图表
  if (getURLChart.value && analysisData.value.urlAnalysis.getTop10.length > 0) {
    const chart = echarts.init(getURLChart.value)
    chart.setOption({
      title: {
        text: 'GET请求URL TOP10',
        left: 'center',
        textStyle: { fontSize: 16 }
      },
      tooltip: {
        trigger: 'axis',
        formatter: function(params) {
          const data = params[0]
          const originalUrl = analysisData.value.urlAnalysis.getTop10[data.dataIndex].url
          const percentage = analysisData.value.urlAnalysis.getTop10[data.dataIndex].percentage
          return `${originalUrl}<br/>访问次数: ${data.value}<br/>占比: ${percentage.toFixed(2)}%`
        }
      },
      grid: {
        left: '10%',
        right: '10%',
        bottom: '20%',
        top: '15%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: analysisData.value.urlAnalysis.getTop10.map(item =>
          item.url.length > 25 ? item.url.substring(0, 25) + '...' : item.url
        ),
        axisLabel: {
          rotate: 45,
          fontSize: 11
        }
      },
      yAxis: {
        type: 'value',
        name: '访问次数',
        nameTextStyle: { fontSize: 12 }
      },
      series: [{
        data: analysisData.value.urlAnalysis.getTop10.map(item => item.count),
        type: 'bar',
        itemStyle: { color: '#E6A23C' },
        label: {
          show: true,
          position: 'top',
          fontSize: 10
        },
        barWidth: '60%'
      }]
    })

    // 响应式调整
    window.addEventListener('resize', () => chart.resize())
  }

  // POST URL图表在标签页切换时渲染，这里跳过

  } catch (error) {
    console.error('渲染图表失败:', error)
    ElMessage.error('图表渲染失败')
  }
}

// 渲染特定IP的所有图表
const renderSpecificIPCharts = () => {
  // 默认渲染GET图表
  if (activeSpecificIPTab.value === 'get') {
    renderSpecificIPGetChart()
  } else if (activeSpecificIPTab.value === 'post') {
    renderSpecificIPPostChart()
  }
}

// 渲染特定IP的GET请求图表
const renderSpecificIPGetChart = () => {
  if (specificIPGetChart.value && specificIPData.value?.urlAnalysis.getTop10.length > 0) {
    // 如果图表已存在，先销毁
    if (chartInstances.specificIPGet) {
      chartInstances.specificIPGet.dispose()
    }

    chartInstances.specificIPGet = echarts.init(specificIPGetChart.value)
    chartInstances.specificIPGet.setOption({
      title: {
        text: `${searchIP.value} GET请求 TOP10`,
        left: 'center',
        textStyle: { fontSize: 16 }
      },
      tooltip: {
        trigger: 'axis',
        formatter: function(params) {
          const data = params[0]
          const originalUrl = specificIPData.value.urlAnalysis.getTop10[data.dataIndex].url
          const percentage = specificIPData.value.urlAnalysis.getTop10[data.dataIndex].percentage
          return `${originalUrl}<br/>访问次数: ${data.value}<br/>占比: ${percentage.toFixed(2)}%`
        }
      },
      grid: {
        left: '10%',
        right: '10%',
        bottom: '20%',
        top: '15%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: specificIPData.value.urlAnalysis.getTop10.map(item =>
          item.url.length > 25 ? item.url.substring(0, 25) + '...' : item.url
        ),
        axisLabel: {
          rotate: 45,
          fontSize: 11
        }
      },
      yAxis: {
        type: 'value',
        name: '访问次数',
        nameTextStyle: { fontSize: 12 }
      },
      series: [{
        data: specificIPData.value.urlAnalysis.getTop10.map(item => item.count),
        type: 'bar',
        itemStyle: { color: '#52c41a' },
        label: {
          show: true,
          position: 'top',
          fontSize: 10
        },
        barWidth: '60%'
      }]
    })

    // 响应式调整
    window.addEventListener('resize', () => chartInstances.specificIPGet.resize())
  }
}

// 渲染特定IP的POST请求图表
const renderSpecificIPPostChart = () => {
  if (specificIPPostChart.value && specificIPData.value?.urlAnalysis.postTop10.length > 0) {
    // 如果图表已存在，先销毁
    if (chartInstances.specificIPPost) {
      chartInstances.specificIPPost.dispose()
    }

    chartInstances.specificIPPost = echarts.init(specificIPPostChart.value)
    chartInstances.specificIPPost.setOption({
      title: {
        text: `${searchIP.value} POST请求 TOP10`,
        left: 'center',
        textStyle: { fontSize: 16 }
      },
      tooltip: {
        trigger: 'axis',
        formatter: function(params) {
          const data = params[0]
          const originalUrl = specificIPData.value.urlAnalysis.postTop10[data.dataIndex].url
          const percentage = specificIPData.value.urlAnalysis.postTop10[data.dataIndex].percentage
          return `${originalUrl}<br/>访问次数: ${data.value}<br/>占比: ${percentage.toFixed(2)}%`
        }
      },
      grid: {
        left: '10%',
        right: '10%',
        bottom: '20%',
        top: '15%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: specificIPData.value.urlAnalysis.postTop10.map(item =>
          item.url.length > 25 ? item.url.substring(0, 25) + '...' : item.url
        ),
        axisLabel: {
          rotate: 45,
          fontSize: 11
        }
      },
      yAxis: {
        type: 'value',
        name: '访问次数',
        nameTextStyle: { fontSize: 12 }
      },
      series: [{
        data: specificIPData.value.urlAnalysis.postTop10.map(item => item.count),
        type: 'bar',
        itemStyle: { color: '#fa8c16' },
        label: {
          show: true,
          position: 'top',
          fontSize: 10
        },
        barWidth: '60%'
      }]
    })

    // 响应式调整
    window.addEventListener('resize', () => chartInstances.specificIPPost.resize())
  }
}

// 处理IP标签页切换
const handleIPTabChange = async (tabName) => {
  await nextTick()
  if (tabName === 'external' && externalIPChart.value && analysisData.value.ipAnalysis.externalTop10.length > 0) {
    // 如果图表已存在，先销毁
    if (chartInstances.externalIP) {
      chartInstances.externalIP.dispose()
    }

    // 创建新的图表实例
    chartInstances.externalIP = echarts.init(externalIPChart.value)
    chartInstances.externalIP.setOption({
      title: {
        text: '外网IP访问次数 TOP10',
        left: 'center',
        textStyle: { fontSize: 16 }
      },
      tooltip: {
        trigger: 'axis',
        formatter: '{b}<br/>访问次数: {c}'
      },
      grid: {
        left: '10%',
        right: '10%',
        bottom: '15%',
        top: '15%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: analysisData.value.ipAnalysis.externalTop10.map(item => item.ip),
        axisLabel: {
          rotate: 45,
          fontSize: 11
        }
      },
      yAxis: {
        type: 'value',
        name: '访问次数',
        nameTextStyle: { fontSize: 12 }
      },
      series: [{
        data: analysisData.value.ipAnalysis.externalTop10.map(item => item.count),
        type: 'bar',
        itemStyle: { color: '#67C23A' },
        label: {
          show: true,
          position: 'top',
          fontSize: 10
        },
        barWidth: '60%'
      }]
    })

    // 响应式调整
    window.addEventListener('resize', () => chartInstances.externalIP.resize())
  }
}

// 处理特定IP标签页切换
const handleSpecificIPTabChange = async (tabName) => {
  await nextTick()
  if (tabName === 'get' && specificIPGetChart.value && specificIPData.value?.urlAnalysis.getTop10.length > 0) {
    renderSpecificIPGetChart()
  } else if (tabName === 'post' && specificIPPostChart.value && specificIPData.value?.urlAnalysis.postTop10.length > 0) {
    renderSpecificIPPostChart()
  }
}

// 处理URL标签页切换
const handleURLTabChange = async (tabName) => {
  await nextTick()
  if (tabName === 'post' && postURLChart.value && analysisData.value.urlAnalysis.postTop10.length > 0) {
    // 如果图表已存在，先销毁
    if (chartInstances.postURL) {
      chartInstances.postURL.dispose()
    }

    // 创建新的图表实例
    chartInstances.postURL = echarts.init(postURLChart.value)
    chartInstances.postURL.setOption({
      title: {
        text: 'POST请求URL TOP10',
        left: 'center',
        textStyle: { fontSize: 16 }
      },
      tooltip: {
        trigger: 'axis',
        formatter: function(params) {
          const data = params[0]
          const originalUrl = analysisData.value.urlAnalysis.postTop10[data.dataIndex].url
          const percentage = analysisData.value.urlAnalysis.postTop10[data.dataIndex].percentage
          return `${originalUrl}<br/>访问次数: ${data.value}<br/>占比: ${percentage.toFixed(2)}%`
        }
      },
      grid: {
        left: '10%',
        right: '10%',
        bottom: '20%',
        top: '15%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: analysisData.value.urlAnalysis.postTop10.map(item =>
          item.url.length > 25 ? item.url.substring(0, 25) + '...' : item.url
        ),
        axisLabel: {
          rotate: 45,
          fontSize: 11
        }
      },
      yAxis: {
        type: 'value',
        name: '访问次数',
        nameTextStyle: { fontSize: 12 }
      },
      series: [{
        data: analysisData.value.urlAnalysis.postTop10.map(item => item.count),
        type: 'bar',
        itemStyle: { color: '#F56C6C' },
        label: {
          show: true,
          position: 'top',
          fontSize: 10
        },
        barWidth: '60%'
      }]
    })

    // 响应式调整
    window.addEventListener('resize', () => chartInstances.postURL.resize())
  }
}

// 监听文件路径变化
watch(() => props.filePath, (newPath) => {
  if (newPath) {
    currentFilePath.value = newPath
    currentFileName.value = props.fileName || newPath.split(/[/\\]/).pop() || ''
    performAnalysis() // 会自动检查缓存
  }
}, { immediate: true })

// 组件挂载时开始分析
onMounted(() => {
  if (currentFilePath.value) {
    performAnalysis()
  }
})
</script>

<style scoped>
.analysis-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

.analysis-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: white;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-btn {
  font-size: 16px;
}

.page-title {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.file-info {
  color: #909399;
  font-size: 14px;
}

.analysis-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.analysis-progress {
  background: white;
  border-radius: 8px;
  padding: 32px;
  margin-bottom: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.progress-info {
  margin-bottom: 24px;
}

.progress-info h3 {
  color: #409eff;
  margin-bottom: 8px;
  font-size: 18px;
}

.progress-info p {
  color: #666;
  margin: 0;
}

.progress-bar {
  margin: 16px 0;
}

.progress-actions {
  margin-top: 16px;
  text-align: center;
}

.section-title {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

/* 概览统计 */
.overview-section {
  margin-bottom: 32px;
}

.overview-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.stat-card {
  background: white;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #409EFF;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #606266;
}

/* 分析区域 */
.ip-analysis-section,
.url-analysis-section,
.single-ip-section {
  margin-bottom: 32px;
  background: white;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.analysis-tabs {
  margin-top: 16px;
}

.chart-container {
  margin-bottom: 32px;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.chart {
  width: 100%;
  height: 450px;
  min-height: 400px;
}

.table-container {
  margin-top: 24px;
  overflow-x: auto;
  width: 100%;
}

.table-container .el-table {
  font-size: 14px;
  width: 100% !important;
}

.table-container .el-table th {
  background-color: #f5f7fa;
  font-weight: 600;
  text-align: center;
}

.table-container .el-table td {
  padding: 12px 0;
  text-align: center;
}

.table-container .el-table .cell {
  padding: 0 8px;
}

/* IP搜索 */
.ip-search {
  margin-bottom: 24px;
}

.ip-input {
  max-width: 400px;
}

.ip-result {
  margin-top: 24px;
}

.ip-info-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.ip-url-analysis {
  margin-top: 32px;
}

.ip-url-analysis h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .chart {
    height: 400px;
  }

  .overview-cards {
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }
}

@media (max-width: 768px) {
  .analysis-header {
    padding: 12px 16px;
    flex-direction: column;
    gap: 12px;
  }

  .header-left {
    gap: 12px;
  }

  .page-title {
    font-size: 20px;
  }

  .analysis-content {
    padding: 16px;
  }

  .overview-cards {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 12px;
  }

  .stat-card {
    padding: 16px;
  }

  .stat-value {
    font-size: 24px;
  }

  .chart {
    height: 350px;
  }

  .chart-container {
    padding: 12px;
    margin-bottom: 24px;
  }

  .table-container {
    font-size: 12px;
  }

  .ip-input {
    max-width: 100%;
  }
}

@media (max-width: 480px) {
  .analysis-header {
    padding: 8px 12px;
  }

  .page-title {
    font-size: 18px;
  }

  .analysis-content {
    padding: 12px;
  }

  .overview-cards {
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
  }

  .stat-card {
    padding: 12px;
  }

  .stat-value {
    font-size: 20px;
  }

  .chart {
    height: 300px;
  }

  .chart-container {
    padding: 8px;
    margin-bottom: 16px;
  }
}
</style>
