import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { GetFileInfo } from 'wailsjs/go/main/App'

export interface LogFile {
  id: string
  name: string
  path: string
  size: number
  lastModified: Date
  isOpen: boolean
  content?: string // 可选属性，用于临时文件（如剪切板内容）
}

export interface RecentFile {
  name: string
  path: string
  lastOpened: Date
}

export interface SearchOptions {
  query: string
  isRegex: boolean
  caseSensitive: boolean
  wholeWord: boolean
}

export interface FilterRule {
  id: string
  name: string
  pattern: string
  isRegex: boolean
  isEnabled: boolean
  color: string
}

export interface HighlightWord {
  text: string
  color: string
}

export interface AnalysisData {
  overview: {
    totalRequests: number
    uniqueIPs: number
    internalIPs: number
    externalIPs: number
  }
  ipAnalysis: {
    internalTop10: any[]
    externalTop10: any[]
  }
  urlAnalysis: {
    getTop10: any[]
    postTop10: any[]
  }
}

export interface AnalysisCache {
  filePath: string
  fileName: string
  data: AnalysisData
  timestamp: number
}

export const useAppStore = defineStore('app', () => {
  // 当前视图状态
  const currentView = ref<'welcome' | 'log-viewer' | 'analysis'>('welcome')
  
  // 文件相关状态
  const logFiles = ref<LogFile[]>([])
  const currentFileId = ref<string | null>(null)
  const recentFiles = ref<RecentFile[]>([])
  
  // 日志内容
  const logContent = ref<string[]>([])
  const filteredContent = ref<string[]>([])
  
  // 搜索和过滤状态
  const searchOptions = ref<SearchOptions>({
    query: '',
    isRegex: false,
    caseSensitive: false,
    wholeWord: false
  })
  
  const filterRules = ref<FilterRule[]>([])
  const activeFilters = ref<string[]>([])

  // 高亮词
  const highlightWords = ref<HighlightWord[]>([])

  // 搜索结果
  const searchResults = ref<Array<{ lineNumber: number, content: string }>>([])
  const currentSearchIndex = ref(0)
  
  // 界面状态
  const sidebarCollapsed = ref(false)
  const showLineNumbers = ref(true)
  const syntaxHighlight = ref(true)
  const wordWrap = ref(true)

  // 全局加载状态
  const isGlobalLoading = ref(false)
  const loadingMessage = ref('')
  const loadingProgress = ref(0)

  // 分析数据缓存
  const analysisCache = ref<Map<string, AnalysisCache>>(new Map())
  
  // 计算属性
  const currentFile = computed(() => {
    return logFiles.value.find(file => file.id === currentFileId.value)
  })
  
  const openFiles = computed(() => {
    return logFiles.value.filter(file => file.isOpen)
  })
  
  const hasOpenFiles = computed(() => {
    return openFiles.value.length > 0
  })
  
  // 最近文件管理方法
  const loadRecentFiles = () => {
    try {
      const stored = localStorage.getItem('recentFiles')
      if (stored) {
        const files = JSON.parse(stored)
        recentFiles.value = files.map((file: any) => ({
          ...file,
          lastOpened: new Date(file.lastOpened)
        })).sort((a: any, b: any) => b.lastOpened - a.lastOpened)
      }
    } catch (error) {
      console.error('加载最近文件失败:', error)
      recentFiles.value = []
    }
  }

  const saveRecentFiles = () => {
    try {
      localStorage.setItem('recentFiles', JSON.stringify(recentFiles.value))
    } catch (error) {
      console.error('保存最近文件失败:', error)
    }
  }

  const addToRecentFiles = (filePath: string, fileName: string) => {
    // 移除已存在的相同文件
    recentFiles.value = recentFiles.value.filter(f => f.path !== filePath)

    // 添加到列表开头
    recentFiles.value.unshift({
      name: fileName,
      path: filePath,
      lastOpened: new Date()
    })

    // 限制最多保存 10 个最近文件
    if (recentFiles.value.length > 10) {
      recentFiles.value = recentFiles.value.slice(0, 10)
    }

    saveRecentFiles()
  }

  const removeFromRecentFiles = (filePath: string) => {
    recentFiles.value = recentFiles.value.filter(f => f.path !== filePath)
    saveRecentFiles()
  }

  // 方法
  const addLogFile = (file: LogFile) => {
    const existingIndex = logFiles.value.findIndex(f => f.path === file.path)
    if (existingIndex >= 0) {
      logFiles.value[existingIndex] = { ...file, isOpen: true }
    } else {
      logFiles.value.push({ ...file, isOpen: true })
    }

    // 添加到最近文件列表（注意参数顺序：filePath, fileName）
    addToRecentFiles(file.path, file.name)

    setCurrentFile(file.id)
  }
  
  const removeLogFile = (fileId: string) => {
    const index = logFiles.value.findIndex(f => f.id === fileId)
    if (index >= 0) {
      logFiles.value.splice(index, 1)
      if (currentFileId.value === fileId) {
        const remainingFiles = openFiles.value
        if (remainingFiles.length > 0) {
          setCurrentFile(remainingFiles[0].id)
        } else {
          setCurrentFile(null)
        }
      }
    }
  }
  
  const setCurrentFile = (fileId: string | null) => {
    currentFileId.value = fileId
    currentView.value = fileId ? 'log-viewer' : 'welcome'

    // 当打开文件时自动隐藏侧边栏
    if (fileId) {
      sidebarCollapsed.value = true
    }
  }
  
  const closeFile = (fileId: string) => {
    const file = logFiles.value.find(f => f.id === fileId)
    if (file) {
      file.isOpen = false
      if (currentFileId.value === fileId) {
        const remainingFiles = openFiles.value
        if (remainingFiles.length > 0) {
          setCurrentFile(remainingFiles[0].id)
        } else {
          setCurrentFile(null)
        }
      }
    }
  }
  
  const updateSearchOptions = (options: Partial<SearchOptions>) => {
    searchOptions.value = { ...searchOptions.value, ...options }
  }

  const setSearchResults = (results: Array<{ lineNumber: number, content: string }>) => {
    searchResults.value = results
    currentSearchIndex.value = 0
  }

  const setCurrentSearchIndex = (index: number) => {
    if (index >= 0 && index < searchResults.value.length) {
      currentSearchIndex.value = index
    }
  }
  
  const addFilterRule = (rule: FilterRule) => {
    filterRules.value.push(rule)
  }
  
  const removeFilterRule = (ruleId: string) => {
    const index = filterRules.value.findIndex(r => r.id === ruleId)
    if (index >= 0) {
      filterRules.value.splice(index, 1)
    }
  }
  
  const toggleFilter = (ruleId: string) => {
    const index = activeFilters.value.indexOf(ruleId)
    if (index >= 0) {
      activeFilters.value.splice(index, 1)
    } else {
      activeFilters.value.push(ruleId)
    }
  }
  
  const setLogContent = (content: string[]) => {
    logContent.value = content
    filteredContent.value = content
  }
  
  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  const setHighlightWords = (words: HighlightWord[]) => {
    highlightWords.value = words
  }

  // 全局加载状态管理
  const setGlobalLoading = (loading: boolean, message: string = '', progress: number = 0) => {
    isGlobalLoading.value = loading
    loadingMessage.value = message
    loadingProgress.value = progress
  }

  const updateLoadingProgress = (progress: number, message?: string) => {
    loadingProgress.value = progress
    if (message) {
      loadingMessage.value = message
    }
  }

  // 分析数据缓存管理
  const setAnalysisData = (filePath: string, fileName: string, data: AnalysisData) => {
    const cache: AnalysisCache = {
      filePath,
      fileName,
      data,
      timestamp: Date.now()
    }
    analysisCache.value.set(filePath, cache)

    // 限制缓存大小，只保留最近5个分析结果
    if (analysisCache.value.size > 5) {
      const entries = Array.from(analysisCache.value.entries())
      entries.sort((a, b) => a[1].timestamp - b[1].timestamp)
      // 删除最旧的
      analysisCache.value.delete(entries[0][0])
    }
  }

  const getAnalysisData = (filePath: string): AnalysisCache | null => {
    return analysisCache.value.get(filePath) || null
  }

  const hasAnalysisData = (filePath: string): boolean => {
    return analysisCache.value.has(filePath)
  }

  const clearAnalysisData = (filePath?: string) => {
    if (filePath) {
      analysisCache.value.delete(filePath)
    } else {
      analysisCache.value.clear()
    }
  }

  const clearAllFiles = () => {
    logFiles.value = []
    currentFileId.value = null
    logContent.value = []
    filteredContent.value = []
    searchResults.value = []
    currentSearchIndex.value = 0
    currentView.value = 'welcome'
  }

  const openFile = async (filePath: string, encoding: string = 'utf-8') => {
    console.log('🔄 AppStore.openFile 被调用:', { filePath, encoding })

    try {
      // 如果有正在进行的全局加载，先清理
      if (isGlobalLoading.value) {
        console.log('🧹 清理之前的全局加载状态')
        setGlobalLoading(false)
      }

      console.log('📞 调用 GetFileInfo:', filePath)
      const fileInfo = await GetFileInfo(filePath)
      console.log('📋 GetFileInfo 返回:', fileInfo)

      if (fileInfo) {
        const logFile: LogFile = {
          id: fileInfo.id,
          name: fileInfo.name,
          path: fileInfo.path,
          size: fileInfo.size,
          lastModified: new Date(fileInfo.lastModified),
          isOpen: true
        }

        console.log('📝 创建 LogFile 对象:', logFile)

        addLogFile(logFile)
        // addLogFile 已经会调用 addToRecentFiles，这里不需要重复调用

        // 设置为当前文件
        setCurrentFile(logFile.id)

        // 切换到日志查看器视图
        currentView.value = 'log-viewer'

        console.log('✅ 文件打开完成:', logFile.name)
        return logFile
      } else {
        throw new Error('无法获取文件信息')
      }
    } catch (error) {
      console.error('❌ AppStore.openFile 失败:', error)
      // 确保加载状态被清理
      if (isGlobalLoading.value) {
        setGlobalLoading(false)
      }
      throw error
    }
  }

  return {
    // 状态
    currentView,
    logFiles,
    currentFileId,
    recentFiles,
    logContent,
    filteredContent,
    searchOptions,
    searchResults,
    currentSearchIndex,
    filterRules,
    activeFilters,
    highlightWords,
    sidebarCollapsed,
    showLineNumbers,
    syntaxHighlight,
    wordWrap,
    isGlobalLoading,
    loadingMessage,
    loadingProgress,
    analysisCache,

    // 计算属性
    currentFile,
    openFiles,
    hasOpenFiles,

    // 方法
    addLogFile,
    removeLogFile,
    setCurrentFile,
    closeFile,
    clearAllFiles,
    openFile,
    updateSearchOptions,
    setSearchResults,
    setCurrentSearchIndex,
    addFilterRule,
    removeFilterRule,
    toggleFilter,
    setLogContent,
    toggleSidebar,
    setHighlightWords,
    setGlobalLoading,
    updateLoadingProgress,

    // 最近文件管理方法
    loadRecentFiles,
    saveRecentFiles,
    addToRecentFiles,
    removeFromRecentFiles,

    // 分析数据管理方法
    setAnalysisData,
    getAnalysisData,
    hasAnalysisData,
    clearAnalysisData
  }
})
