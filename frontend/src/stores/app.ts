import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface LogFile {
  id: string
  name: string
  path: string
  size: number
  lastModified: Date
  isOpen: boolean
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

export const useAppStore = defineStore('app', () => {
  // 当前视图状态
  const currentView = ref<'welcome' | 'log-viewer'>('welcome')
  
  // 文件相关状态
  const logFiles = ref<LogFile[]>([])
  const currentFileId = ref<string | null>(null)
  const recentFiles = ref<LogFile[]>([])
  
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

  // 搜索结果
  const searchResults = ref<Array<{ lineNumber: number, content: string }>>([])
  const currentSearchIndex = ref(0)
  
  // 界面状态
  const sidebarCollapsed = ref(false)
  const showLineNumbers = ref(true)
  const syntaxHighlighting = ref(true)
  const wordWrap = ref(true)
  
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
  
  // 方法
  const addLogFile = (file: LogFile) => {
    const existingIndex = logFiles.value.findIndex(f => f.path === file.path)
    if (existingIndex >= 0) {
      logFiles.value[existingIndex] = { ...file, isOpen: true }
    } else {
      logFiles.value.push({ ...file, isOpen: true })
    }
    
    // 添加到最近文件列表
    const recentIndex = recentFiles.value.findIndex(f => f.path === file.path)
    if (recentIndex >= 0) {
      recentFiles.value.splice(recentIndex, 1)
    }
    recentFiles.value.unshift(file)
    if (recentFiles.value.length > 10) {
      recentFiles.value = recentFiles.value.slice(0, 10)
    }
    
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
    sidebarCollapsed,
    showLineNumbers,
    syntaxHighlighting,
    wordWrap,
    
    // 计算属性
    currentFile,
    openFiles,
    hasOpenFiles,
    
    // 方法
    addLogFile,
    removeLogFile,
    setCurrentFile,
    closeFile,
    updateSearchOptions,
    setSearchResults,
    setCurrentSearchIndex,
    addFilterRule,
    removeFilterRule,
    toggleFilter,
    setLogContent,
    toggleSidebar
  }
})
