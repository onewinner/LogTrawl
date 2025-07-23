<template>
  <div class="toolbar">
    <!-- 左侧按钮组 -->
    <div class="toolbar-left">
      <!-- 主要操作区域 -->
      <div class="main-operations">
        <div class="main-buttons-row">
          <!-- 菜单按钮 -->
          <div class="main-button-group">
            <el-button
              :icon="Menu"
              size="large"
              text
              @click="appStore.toggleSidebar()"
              title="切换侧边栏"
              class="main-btn"
            />
            <span class="main-btn-label">菜单</span>
          </div>

          <!-- 文件按钮 -->
          <div class="main-button-group">
            <el-button
              :icon="FolderOpened"
              size="large"
              text
              @click="openFile"
              title="打开文件"
              class="main-btn"
            />
            <span class="main-btn-label">文件</span>
          </div>

          <!-- 文件分片按钮 -->
          <div class="main-button-group">
            <el-button
              :icon="ScissorsIcon"
              size="large"
              text
              @click="openFileSplitter"
              title="文件分片"
              class="main-btn splitter-btn"
            />
            <span class="main-btn-label">分片</span>
          </div>

          <!-- 最近文件按钮 -->
          <div class="main-button-group">
            <el-dropdown @command="openRecentFile" trigger="click" :disabled="appStore.recentFiles.length === 0">
              <el-button
                :icon="Document"
                size="large"
                text
                title="最近文件"
                class="main-btn"
                :disabled="appStore.recentFiles.length === 0"
              />
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item
                    v-for="file in appStore.recentFiles"
                    :key="file.path"
                    :command="file.path"
                  >
                    <div class="recent-file-option">
                      <span class="file-name">{{ file.name }}</span>
                      <span class="file-path">{{ file.path }}</span>
                    </div>
                  </el-dropdown-item>
                  <el-dropdown-item v-if="appStore.recentFiles.length === 0" disabled>
                    暂无最近文件
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <span class="main-btn-label">最近文件</span>
          </div>
        </div>
      </div>

      <!-- 分隔符 -->
      <el-divider direction="vertical" />

      <!-- 窗口选择和操作区域 -->
      <div class="window-operations">
        <!-- 第一行：窗口选择器 -->
        <div class="window-selector-row">
          <el-dropdown @command="handleWindowSelect" trigger="click">
            <el-button class="window-selector-btn">
              {{ currentWindowName }}
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="main">
                  <el-icon><Document /></el-icon>
                  主窗口
                </el-dropdown-item>
                <el-dropdown-item
                  v-if="filterWindows.length > 0"
                  command="filter"
                >
                  <el-icon><Filter /></el-icon>
                  过滤窗口
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>

        <!-- 第二行：刷新和编码选择 -->
        <div class="refresh-encoding-row">
          <el-button
            :icon="Refresh"
            size="small"
            @click="refreshFile"
            :disabled="!appStore.currentFile"
            title="刷新文件"
            class="refresh-btn"
          />

          <el-select
            v-model="encoding"
            size="small"
            class="encoding-select"
            title="文件编码"
          >
            <el-option label="UTF-8" value="utf-8" />
            <el-option label="GBK" value="gbk" />
            <el-option label="ASCII" value="ascii" />
          </el-select>
        </div>
      </div>

      <!-- 分隔符 -->
      <el-divider direction="vertical" />

      <!-- 搜索容器 -->
      <div class="search-container">
        <!-- 搜索输入框容器 -->
        <div class="search-input-container">
          <el-input
            v-model="searchQuery"
            placeholder="输入要搜索的词"
            class="search-input"
            @keyup.enter="performSearch"
            @clear="clearSearch"
            clearable
          />
          <!-- 搜索结果计数 - 显示在输入框内右侧 -->
          <div class="search-count-overlay" v-if="searchResults.total > 0">
            {{ searchResults.current }}/{{ searchResults.total }}
          </div>
        </div>

        <!-- 搜索控制按钮 -->
        <div class="search-controls">
          <el-tooltip content="区分大小写" placement="bottom">
            <el-button
              :type="searchOptions.caseSensitive ? 'primary' : ''"
              class="search-control-btn"
              size="small"
              @click="toggleSearchOption('caseSensitive')"
            >
              Aa
            </el-button>
          </el-tooltip>

          <el-tooltip content="正则表达式" placement="bottom">
            <el-button
              :type="searchOptions.isRegex ? 'primary' : ''"
              class="search-control-btn"
              size="small"
              @click="toggleSearchOption('isRegex')"
            >
              .*
            </el-button>
          </el-tooltip>

          <el-tooltip content="全词匹配" placement="bottom">
            <el-button
              :type="searchOptions.wholeWord ? 'primary' : ''"
              class="search-control-btn"
              size="small"
              @click="toggleSearchOption('wholeWord')"
            >
              ab
            </el-button>
          </el-tooltip>

          <el-tooltip content="下一个" placement="bottom">
            <el-button
              :icon="ArrowDown"
              class="search-control-btn"
              size="small"
              @click="nextMatch"
              :disabled="searchResults.total <= 0 || searchResults.current >= searchResults.total"
            />
          </el-tooltip>

          <el-tooltip content="上一个" placement="bottom">
            <el-button
              :icon="ArrowUp"
              class="search-control-btn"
              size="small"
              @click="previousMatch"
              :disabled="searchResults.total <= 0 || searchResults.current <= 1"
            />
          </el-tooltip>

          <el-tooltip content="清除搜索高亮" placement="bottom">
            <el-button
              :icon="Close"
              class="search-control-btn"
              size="small"
              @click="clearSearch"
            />
          </el-tooltip>
        </div>

       
      </div>

      <!-- 行号跳转区域 -->
      <div class="line-jump-section">
        <!-- 行号跳转按钮（收起时显示） -->
        <div class="line-jump-button-group" v-if="!showLineJumpInput">
          <el-button
            :icon="Position"
            size="large"
            text
            @click="toggleLineJumpInput"
            title="行号定位"
            class="line-jump-toggle-btn"
          />
          <span class="line-jump-btn-label">行号</span>
        </div>

        <!-- 行号跳转容器（展开时显示） -->
        <div class="line-jump-container" v-if="showLineJumpInput">
          <!-- 关闭按钮 -->
          <el-button
            class="close-line-jump-btn"
            size="small"
            text
            @click="toggleLineJumpInput"
            title="关闭"
          >
            <el-icon><Close /></el-icon>
          </el-button>

          <!-- 第一行：行号输入和加减按钮 -->
          <div class="line-number-input-row">
            <el-button
              class="line-control-btn"
              @click="jumpToPreviousLine"
              size="small"
              title="上一行"
            >
              −
            </el-button>
            <el-input
              v-model="targetLineNumber"
              class="line-number-input"
              @keyup.enter="jumpToLine"
              placeholder="行号"
              size="small"
            />
            <el-button
              class="line-control-btn"
              @click="jumpToNextLine"
              size="small"
              title="下一行"
            >
              +
            </el-button>
          </div>

          <!-- 第二行：导航按钮 -->
          <div class="line-navigation-controls">
            <el-tooltip content="跳转到首行" placement="bottom">
              <el-button
                class="line-nav-btn"
                @click="jumpToFirstLine"
                size="small"
              >
                ⇈
              </el-button>
            </el-tooltip>
            <el-tooltip content="跳转到尾行" placement="bottom">
              <el-button
                class="line-nav-btn"
                @click="jumpToLastLine"
                size="small"
              >
                ⇊
              </el-button>
            </el-tooltip>
            <el-tooltip content="上一行" placement="bottom">
              <el-button
                class="line-nav-btn"
                @click="jumpToPreviousLine"
                size="small"
              >
                ←
              </el-button>
            </el-tooltip>
            <el-tooltip content="下一行" placement="bottom">
              <el-button
                class="line-nav-btn"
                @click="jumpToNextLine"
                size="small"
              >
                →
              </el-button>
            </el-tooltip>
            <el-tooltip content="追踪到当前行" placement="bottom">
              <el-button
                class="line-nav-btn"
                @click="trackToCurrentLine"
                size="small"
              >
                🎯
              </el-button>
            </el-tooltip>
          </div>
        </div>
      </div>

      <!-- 分隔符 -->
      <el-divider direction="vertical" />

      <!-- 过滤器区域 -->
      <div class="filter-section">
        <!-- 过滤器按钮（未展开时显示） -->
        <div class="filter-button-group" v-if="!showFilterInput">
          <el-button
            :icon="Filter"
            size="large"
            text
            @click="toggleFilterInput"
            title="过滤器"
            class="filter-toggle-btn"
          />
          <span class="filter-btn-label">过滤</span>
        </div>

        <!-- 过滤器展开内容（展开时显示） -->
        <div class="filter-container" v-if="showFilterInput">
          <!-- 过滤输入行 -->
          <div class="filter-input-row">
            <el-select
              v-model="filterMode"
              class="filter-mode-select"
              size="small"
            >
              <el-option label="包含" value="include" />
              <el-option label="排除" value="exclude" />
            </el-select>
            <el-input
              v-model="filterInput"
              placeholder="输入要过滤的内容，支持 AND/&&、OR/|| 操作符"
              class="filter-input"
              @keyup.enter="applyFilter"
              @clear="clearFilter"
              clearable
              size="small"
            >
              <template #suffix>
                <el-tooltip
                  content="查看搜索示例"
                  placement="bottom"
                >
                  <el-button
                    :icon="InfoFilled"
                    size="small"
                    text
                    @click="showFilterExamples = true"
                    class="filter-help-btn"
                  />
                </el-tooltip>
              </template>
            </el-input>
          </div>

          <!-- 过滤控制按钮 -->
          <div class="filter-controls">
            <el-tooltip content="区分大小写" placement="bottom">
              <el-button
                :type="filterOptions.caseSensitive ? 'primary' : ''"
                class="filter-control-btn"
                size="small"
                @click="toggleFilterOption('caseSensitive')"
              >
                Aa
              </el-button>
            </el-tooltip>

            <el-tooltip content="正则表达式" placement="bottom">
              <el-button
                :type="filterOptions.isRegex ? 'primary' : ''"
                class="filter-control-btn"
                size="small"
                @click="toggleFilterOption('isRegex')"
              >
                .*
              </el-button>
            </el-tooltip>

            <el-tooltip content="全词匹配" placement="bottom">
              <el-button
                :type="filterOptions.wholeWord ? 'primary' : ''"
                class="filter-control-btn"
                size="small"
                @click="toggleFilterOption('wholeWord')"
              >
                Ab
              </el-button>
            </el-tooltip>

            <el-button
              type="primary"
              class="filter-action-btn"
              size="small"
              @click="applyFilter"
            >
              确认
            </el-button>

            <el-button
              class="filter-action-btn"
              size="small"
              @click="cancelFilter"
            >
              取消
            </el-button>
          </div>
        </div>
      </div>

      <!-- 高亮词管理区域 -->
      <div class="highlight-words-section">
        <!-- 添加高亮词按钮 -->
        <div class="add-highlight-btn" v-if="highlightWords.length === 0" @click="showAddHighlightDialog = true">
          <el-icon><Plus /></el-icon>
          <span>添加高亮单词</span>
        </div>

        <!-- 高亮词标签列表 -->
        <div class="highlight-words-list" v-if="highlightWords.length > 0">
          <div class="highlight-words-container">
            <!-- 高亮词标签 -->
            <el-dropdown
              v-for="(word, index) in highlightWords"
              :key="index"
              trigger="click"
              @command="(command) => handleHighlightWordCommand(command, index)"
            >
              <div
                class="highlight-word-tag"
                :style="{ backgroundColor: word.color, color: getContrastColor(word.color) }"
              >
                {{ word.text }}
                <el-icon class="remove-icon" @click.stop="removeHighlightWord(index)">
                  <Close />
                </el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="filter">
                    <el-icon><Filter /></el-icon>
                    添加到过滤
                  </el-dropdown-item>
                  <el-dropdown-item command="changeColor">
                    <el-icon><Edit /></el-icon>
                    更改颜色
                  </el-dropdown-item>
                  <el-dropdown-item command="remove" divided>
                    <el-icon><Delete /></el-icon>
                    删除高亮
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>

            <!-- 添加按钮（当有高亮词时显示） -->
            <el-button
              class="add-more-btn"
              size="small"
              circle
              @click="showAddHighlightDialog = true"
            >
              <el-icon><Plus /></el-icon>
            </el-button>
          </div>
        </div>
      </div>

      <!-- 时间线管理区域 -->
      <div class="timeline-section">
        <!-- 时间线按钮组 -->
        <div class="timeline-buttons">
          <!-- 第一行按钮 -->
          <div class="timeline-row">
            <el-tooltip content="增加时间线" placement="bottom">
              <el-button
                class="timeline-btn"
                size="small"
                @click="addTimeline"
                title="增加时间线"
              >
                <el-icon><Plus /></el-icon>
              </el-button>
            </el-tooltip>

            <el-tooltip content="清除时间线" placement="bottom">
              <el-button
                class="timeline-btn"
                size="small"
                @click="clearTimeline"
                title="清除时间线"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
            </el-tooltip>
          </div>

          <!-- 第二行按钮 -->
          <div class="timeline-row">
            <el-dropdown @command="handleTimelineCopyCommand" trigger="click" :disabled="timelineEntries.length === 0">
              <el-button
                class="timeline-btn"
                size="small"
                title="复制时间线"
                :disabled="timelineEntries.length === 0"
              >
                <el-icon><Document /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="detailed">📋 详细报告</el-dropdown-item>
                  <el-dropdown-item command="simple">📝 简单列表</el-dropdown-item>
                  <el-dropdown-item command="csv">📊 CSV格式</el-dropdown-item>
                  <el-dropdown-item command="json">🔧 JSON格式</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>

            <el-tooltip content="颜色自定义" placement="bottom">
              <el-button
                class="timeline-btn"
                size="small"
                @click="showTimelineColorDialog = true"
                title="颜色自定义"
              >
                <el-icon><Edit /></el-icon>
              </el-button>
            </el-tooltip>
          </div>
        </div>
      </div>
    </div>

    <!-- 中间区域：空白 -->
    <div class="toolbar-center">
    </div>
    
    <!-- 右侧按钮组 -->
    <div class="toolbar-right" ref="toolbarRightRef">
      <!-- 分隔符 -->
      <el-divider direction="vertical" />

      <!-- 右侧主按钮行 -->
      <div class="right-buttons-row">
        <!-- 始终显示的核心按钮 -->
        <div class="core-buttons">
          <!-- 分析按钮 -->
          <div class="main-button-group" :class="{ 'compact': isCompactMode }">
            <el-button
              :icon="DataAnalysis"
              :size="isCompactMode ? 'default' : 'large'"
              text
              @click="openAnalysisPage"
              :disabled="!appStore.currentFile"
              title="数据分析"
              class="main-btn"
            />
            <span v-if="!isCompactMode" class="main-btn-label">分析</span>
          </div>

          <!-- 设置按钮 -->
          <div class="main-button-group" :class="{ 'compact': isCompactMode }">
            <el-button
              :icon="Setting"
              :size="isCompactMode ? 'default' : 'large'"
              text
              @click="showSettings = true"
              title="设置"
              class="main-btn"
            />
            <span v-if="!isCompactMode" class="main-btn-label">设置</span>
          </div>
        </div>

        <!-- 响应式显示的按钮 -->
        <div class="responsive-buttons" v-if="showResponsiveButtons">
          <!-- 保存项目按钮 -->
          <div class="main-button-group" :class="{ 'compact': isCompactMode }">
            <el-button
              :icon="FolderOpened"
              :size="isCompactMode ? 'default' : 'large'"
              text
              @click="saveProject"
              :disabled="!hasProjectData"
              title="保存项目"
              class="main-btn"
            />
            <span v-if="!isCompactMode" class="main-btn-label">保存项目</span>
          </div>

          <!-- 视图选项按钮 -->
          <div class="main-button-group" :class="{ 'compact': isCompactMode }">
            <el-dropdown @command="handleViewCommand" trigger="click">
              <el-button
                :icon="View"
                :size="isCompactMode ? 'default' : 'large'"
                text
                title="视图选项"
                class="main-btn"
              />
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item
                    :icon="appStore.showLineNumbers ? Check : ''"
                    command="toggleLineNumbers"
                  >
                    显示行号
                  </el-dropdown-item>
                  <el-dropdown-item
                    :icon="appStore.syntaxHighlight ? Check : ''"
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
            <span v-if="!isCompactMode" class="main-btn-label">视图</span>
          </div>

          <!-- 导出按钮 -->
          <div class="main-button-group" :class="{ 'compact': isCompactMode }">
            <el-button
              :icon="Download"
              :size="isCompactMode ? 'default' : 'large'"
              text
              @click="exportData"
              :disabled="!appStore.currentFile"
              title="导出"
              class="main-btn"
            />
            <span v-if="!isCompactMode" class="main-btn-label">导出</span>
          </div>
        </div>

        <!-- 更多菜单按钮 -->
        <div class="main-button-group more-menu" v-if="showMoreMenu" :class="{ 'compact': isCompactMode }">
          <el-dropdown @command="handleMoreCommand" trigger="click" placement="bottom-end">
            <el-button
              :icon="MoreFilled"
              :size="isCompactMode ? 'default' : 'large'"
              text
              title="更多功能"
              class="main-btn"
            />
            <template #dropdown>
              <el-dropdown-menu>
                <!-- 当响应式按钮被隐藏时，显示在更多菜单中 -->
                <template v-if="!showResponsiveButtons">
                  <el-dropdown-item
                    :icon="FolderOpened"
                    command="saveProject"
                    :disabled="!hasProjectData"
                  >
                    保存项目
                  </el-dropdown-item>
                  <el-dropdown-item
                    :icon="View"
                    command="viewOptions"
                  >
                    视图选项
                  </el-dropdown-item>
                  <el-dropdown-item
                    :icon="Download"
                    command="exportData"
                    :disabled="!appStore.currentFile"
                  >
                    导出
                  </el-dropdown-item>
                  <el-dropdown-item divided />
                </template>

                <!-- 其他功能 -->
                <el-dropdown-item
                  :icon="InfoFilled"
                  command="about"
                >
                  关于
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <span v-if="!isCompactMode" class="main-btn-label">更多</span>
        </div>
      </div>
    </div>
  </div>

  <!-- 过滤搜索示例对话框 -->
  <el-dialog
    v-model="showFilterExamples"
    title="过滤搜索示例"
    width="600px"
    :show-close="true"
  >
    <div class="filter-examples">
      <div class="example-section">
        <h4>基础搜索</h4>
        <div class="example-list">
          <div class="example-item" @click="useExample('error')">
            <code>error</code>
            <span class="example-desc">包含 "error" 的行</span>
          </div>
          <div class="example-item" @click="useExample('404')">
            <code>404</code>
            <span class="example-desc">包含 "404" 的行</span>
          </div>
          <div class="example-item" @click="useExample('192.168.1.1')">
            <code>192.168.1.1</code>
            <span class="example-desc">包含特定IP地址的行</span>
          </div>
        </div>
      </div>

      <div class="example-section">
        <h4>AND 操作符（同时包含）</h4>
        <div class="example-list">
          <div class="example-item" @click="useExample('error AND 404')">
            <code>error AND 404</code>
            <span class="example-desc">同时包含 "error" 和 "404" 的行</span>
          </div>
          <div class="example-item" @click="useExample('error && POST')">
            <code>error && POST</code>
            <span class="example-desc">同时包含 "error" 和 "POST" 的行</span>
          </div>
          <div class="example-item" @click="useExample('login AND success AND admin')">
            <code>login AND success AND admin</code>
            <span class="example-desc">同时包含三个关键词的行</span>
          </div>
        </div>
      </div>

      <div class="example-section">
        <h4>OR 操作符（包含任一）</h4>
        <div class="example-list">
          <div class="example-item" @click="useExample('error OR warning')">
            <code>error OR warning</code>
            <span class="example-desc">包含 "error" 或 "warning" 的行</span>
          </div>
          <div class="example-item" @click="useExample('404 || 500')">
            <code>404 || 500</code>
            <span class="example-desc">包含 "404" 或 "500" 的行</span>
          </div>
          <div class="example-item" @click="useExample('GET OR POST OR PUT')">
            <code>GET OR POST OR PUT</code>
            <span class="example-desc">包含任一HTTP方法的行</span>
          </div>
        </div>
      </div>

      <div class="example-section">
        <h4>复合条件（括号分组）</h4>
        <div class="example-list">
          <div class="example-item" @click="useExample('(error OR warning) AND 404')">
            <code>(error OR warning) AND 404</code>
            <span class="example-desc">包含错误信息且状态码为404的行</span>
          </div>
          <div class="example-item" @click="useExample('login && (success || failed)')">
            <code>login && (success || failed)</code>
            <span class="example-desc">登录相关的成功或失败记录</span>
          </div>
          <div class="example-item" @click="useExample('(GET || POST) AND /api/ AND (200 OR 201)')">
            <code>(GET || POST) AND /api/ AND (200 OR 201)</code>
            <span class="example-desc">API请求的成功响应</span>
          </div>
        </div>
      </div>

      <div class="example-section">
        <h4>实际应用场景</h4>
        <div class="example-list">
          <div class="example-item" @click="useExample('SQL AND (injection OR attack)')">
            <code>SQL AND (injection OR attack)</code>
            <span class="example-desc">SQL注入攻击检测</span>
          </div>
          <div class="example-item" @click="useExample('(failed || error) AND login')">
            <code>(failed || error) AND login</code>
            <span class="example-desc">登录失败记录</span>
          </div>
          <div class="example-item" @click="useExample('(500 OR 502 OR 503) AND /api/')">
            <code>(500 OR 502 OR 503) AND /api/</code>
            <span class="example-desc">API服务器错误</span>
          </div>
          <div class="example-item" @click="useExample('timeout OR slow OR (response && time)')">
            <code>timeout OR slow OR (response && time)</code>
            <span class="example-desc">性能问题相关日志</span>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="showFilterExamples = false">关闭</el-button>
      </div>
    </template>
  </el-dialog>

  <!-- 添加高亮词对话框 -->
  <el-dialog
    v-model="showAddHighlightDialog"
    title="添加高亮单词"
    width="400px"
  >
    <div class="add-highlight-content">
      <el-form label-width="80px">
        <el-form-item label="单词">
          <el-input
            v-model="newHighlightWord.text"
            placeholder="输入要高亮的单词"
            @keyup.enter="addHighlightWord"
          />
        </el-form-item>
        <el-form-item label="颜色">
          <div class="color-selection">
            <!-- 预设颜色 -->
            <div class="preset-colors">
              <div
                v-for="color in presetColors"
                :key="color"
                class="color-option"
                :class="{ active: newHighlightWord.color === color }"
                :style="{ backgroundColor: color }"
                @click="newHighlightWord.color = color"
              ></div>
            </div>
            <!-- 自定义颜色 -->
            <el-color-picker
              v-model="newHighlightWord.color"
              show-alpha
              :predefine="presetColors"
            />
          </div>
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <el-button @click="showAddHighlightDialog = false">取消</el-button>
      <el-button type="primary" @click="addHighlightWord" :disabled="!newHighlightWord.text.trim()">添加</el-button>
    </template>
  </el-dialog>

  <!-- 颜色选择对话框 -->
  <el-dialog
    v-model="showColorDialog"
    title="更改颜色"
    width="300px"
  >
    <div class="color-change-content">
      <div class="preset-colors">
        <div
          v-for="color in presetColors"
          :key="color"
          class="color-option large"
          :class="{ active: selectedColor === color }"
          :style="{ backgroundColor: color }"
          @click="selectedColor = color"
        ></div>
      </div>
      <el-color-picker
        v-model="selectedColor"
        show-alpha
        :predefine="presetColors"
        style="margin-top: 16px;"
      />
    </div>
    <template #footer>
      <el-button @click="showColorDialog = false">取消</el-button>
      <el-button type="primary" @click="changeHighlightColor">确定</el-button>
    </template>
  </el-dialog>

  <!-- 时间线颜色自定义对话框 -->
  <el-dialog
    v-model="showTimelineColorDialog"
    title="时间线颜色自定义"
    width="400px"
    :modal="true"
    :close-on-click-modal="false"
  >
    <div class="timeline-color-content">
      <div class="color-section">
        <label>选择时间线颜色：</label>
        <div class="color-picker-container">
          <el-color-picker
            v-model="timelineColor"
            :predefine="presetColors"
            show-alpha
          />
          <span class="color-preview" :style="{ backgroundColor: timelineColor }">
            {{ timelineColor }}
          </span>
        </div>
      </div>
    </div>
    <template #footer>
      <el-button @click="showTimelineColorDialog = false">取消</el-button>
      <el-button type="primary" @click="showTimelineColorDialog = false">确定</el-button>
    </template>
  </el-dialog>

  <!-- 添加时间线对话框 -->
  <el-dialog
    v-model="showAddTimelineDialog"
    title="添加时间线条目"
    width="600px"
    :modal="true"
    :close-on-click-modal="false"
  >
    <div class="add-timeline-content">
      <!-- 当前聚焦行信息 -->
      <div class="timeline-field" v-if="currentFocusedLogLine">
        <label class="field-label">当前聚焦行</label>
        <div class="focused-line-info">
          <span class="line-number-badge">行 {{ currentFocusedLogLine.lineNumber }}</span>
          <div class="line-content-preview">{{ currentFocusedLogLine.content }}</div>
        </div>
      </div>

      <!-- 备注输入 -->
      <div class="timeline-field">
        <label class="field-label">备注信息 <span class="required">*</span></label>
        <el-input
          v-model="newTimelineEntry.note"
          type="textarea"
          :rows="3"
          placeholder="请输入备注信息..."
          maxlength="200"
          show-word-limit
        />
      </div>

      <!-- 行号信息 -->
      <div class="timeline-field" v-if="newTimelineEntry.lineNumber">
        <label class="field-label">行号</label>
        <el-input
          :model-value="`第 ${newTimelineEntry.lineNumber} 行`"
          readonly
        />
      </div>

      <!-- 日志时间戳 -->
      <div class="timeline-field" v-if="newTimelineEntry.logTimestamp">
        <label class="field-label">日志时间</label>
        <el-input
          v-model="newTimelineEntry.logTimestamp"
          placeholder="日志时间戳"
          readonly
        />
      </div>

      <!-- 日志内容 -->
      <div class="timeline-field" v-if="newTimelineEntry.logContent">
        <label class="field-label">相关日志</label>
        <el-input
          v-model="newTimelineEntry.logContent"
          type="textarea"
          :rows="6"
          placeholder="日志内容"
          readonly
        />
      </div>

      <!-- 颜色选择 -->
      <div class="timeline-field">
        <label class="field-label">标记颜色</label>
        <div class="color-picker-container">
          <el-color-picker
            v-model="timelineColor"
            :predefine="presetColors"
          />
          <span class="color-preview" :style="{ backgroundColor: timelineColor }">
            {{ timelineColor }}
          </span>
        </div>
      </div>
    </div>
    <template #footer>
      <el-button @click="showAddTimelineDialog = false">取消</el-button>
      <el-button type="primary" @click="confirmAddTimeline">添加</el-button>
    </template>
  </el-dialog>

  <!-- 文件分片对话框 -->
  <FileSplitterDialog
    v-model="showFileSplitter"
    :file-path="splitterFilePath"
    :file-name="splitterFileName"
    :file-size="splitterFileSize"
    @split-complete="handleSplitComplete"
  />

  <!-- 设置对话框 -->
  <SettingsDialog v-model="showSettings" />
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted, onUnmounted, watch } from 'vue'
import { useAppStore, type HighlightWord } from '@/stores/app'
import { ElMessage } from 'element-plus'
import {
  Menu,
  FolderOpened,
  Refresh,
  ArrowUp,
  ArrowDown,
  Filter,
  Setting,
  View,
  Download,
  Check,
  Close,
  Document,
  Plus,
  Edit,
  Delete,
  Position,
  DataAnalysis,
  Scissor as ScissorsIcon,
  InfoFilled,
  MoreFilled
} from '@element-plus/icons-vue'
import { OpenFileDialog, GetFileInfo, SearchInFile, SaveFileDialog, ExportLogLines } from 'wailsjs/go/main/App'
import FileSplitterDialog from './FileSplitterDialog.vue'
import SettingsDialog from './SettingsDialog.vue'

const appStore = useAppStore()

// 搜索相关
const searchQuery = ref('')
const encoding = ref('utf-8')

// 文件分片相关
const showFileSplitter = ref(false)
const splitterFilePath = ref('')
const splitterFileName = ref('')
const splitterFileSize = ref(0)

// 设置对话框
const showSettings = ref(false)

// 响应式工具栏
const toolbarRightRef = ref<HTMLElement>()
const windowWidth = ref(window.innerWidth)
const isCompactMode = computed(() => windowWidth.value < 1200)
const showResponsiveButtons = computed(() => windowWidth.value >= 900)
const showMoreMenu = computed(() => windowWidth.value < 900)

// 项目保存相关
const hasProjectData = computed(() => {
  return appStore.openFiles.length > 0 ||
         Object.keys(fileStates.value).length > 0 ||
         highlightWords.value.length > 0
})

// 高亮词相关
const highlightWords = ref<HighlightWord[]>([])
const showAddHighlightDialog = ref(false)
const showColorDialog = ref(false)
const selectedColor = ref('#337ecc')
const currentEditingIndex = ref(-1)
const newHighlightWord = ref({
  text: '',
  color: '#337ecc'
})

// 预设颜色
const presetColors = [
  '#337ecc', '#79bbff', '#529b2e', '#95d475',
  '#b88230', '#eebe77', '#c45656', '#f89898',
  '#73767a'
]

// 时间线相关
const timelineEntries = ref([])
const showTimelineColorDialog = ref(false)
const timelineColor = ref('#409eff')
const showAddTimelineDialog = ref(false)
const newTimelineEntry = ref({
  note: '',
  logContent: '',
  logTimestamp: '',
  lineNumber: null
})

// 当前聚焦的日志行信息
const currentFocusedLogLine = ref<{
  lineIndex: number
  lineNumber: number
  content: string
} | null>(null)

// 文件状态管理 - 每个文件独立的状态
const fileStates = ref<Map<string, {
  timelineEntries: any[]
  filterWindows: any[]
  activeFilterWindow: string
  filterInput: string
  filterMode: string
  filterOptions: any
  showFilterInput: boolean
  timelinePanelVisible: boolean
  filterWindowsHeight: number
  filterStartLine: number
  highlightWords: any[]  // 新增：每个文件独立的高亮词
  lastAccessed: number
}>>(new Map())

// 当前文件路径
const currentFilePath = ref<string>('')

// 最大缓存文件数量（LRU缓存）
const MAX_CACHED_FILES = 10

// LRU缓存清理
const cleanupOldFileStates = () => {
  if (fileStates.value.size <= MAX_CACHED_FILES) return

  // 按最后访问时间排序，删除最旧的文件状态
  const sortedEntries = Array.from(fileStates.value.entries())
    .sort((a, b) => a[1].lastAccessed - b[1].lastAccessed)

  const toDelete = sortedEntries.slice(0, sortedEntries.length - MAX_CACHED_FILES)
  toDelete.forEach(([filePath]) => {
    fileStates.value.delete(filePath)
  })
}

// 获取当前文件的状态，如果不存在则创建默认状态
const getCurrentFileState = () => {
  const filePath = appStore.currentFile?.path || ''
  if (!fileStates.value.has(filePath)) {
    // 清理旧的缓存
    cleanupOldFileStates()

    fileStates.value.set(filePath, {
      timelineEntries: [],
      filterWindows: [],
      activeFilterWindow: '',
      filterInput: '',
      filterMode: 'include',
      filterOptions: {
        caseSensitive: false,
        wholeWord: false,
        isRegex: false
      },
      showFilterInput: false,
      timelinePanelVisible: false,
      filterWindowsHeight: 300,
      filterStartLine: 0,
      highlightWords: [],  // 新增：初始化为空数组
      lastAccessed: Date.now()
    })
  } else {
    // 更新访问时间
    const state = fileStates.value.get(filePath)!
    state.lastAccessed = Date.now()
  }
  return fileStates.value.get(filePath)!
}

// 防止重复保存的标志
const isSaving = ref(false)

// 保存当前状态到文件状态中
const saveCurrentState = () => {
  const filePath = currentFilePath.value
  if (!filePath || isSaving.value) return

  isSaving.value = true

  try {
    const state = fileStates.value.get(filePath)
    if (state) {
      state.timelineEntries = [...timelineEntries.value]
      state.filterWindows = [...filterWindows.value]
      state.activeFilterWindow = activeFilterWindow.value
      state.filterInput = filterInput.value
      state.filterMode = filterMode.value
      state.filterOptions = { ...filterOptions.value }
      state.showFilterInput = showFilterInput.value
      state.highlightWords = [...highlightWords.value]  // 新增：保存当前文件的高亮词
      state.lastAccessed = Date.now()

      console.log('💾 保存文件状态:', {
        filePath: filePath.split('/').pop(),
        timelineCount: state.timelineEntries.length,
        filterWindowsCount: state.filterWindows.length,
        activeFilterWindow: state.activeFilterWindow
      })
    }
  } finally {
    // 使用 nextTick 确保在下一个事件循环中重置标志
    nextTick(() => {
      isSaving.value = false
    })
  }
}

// 防止重复恢复的标志
const isRestoring = ref(false)

// 恢复文件状态
const restoreFileState = () => {
  if (isRestoring.value) return

  isRestoring.value = true
  const state = getCurrentFileState()

  console.log('📂 开始恢复文件状态:', {
    filePath: currentFilePath.value?.split('/').pop(),
    timelineCount: state.timelineEntries.length,
    filterWindowsCount: state.filterWindows.length,
    timelinePanelVisible: state.timelinePanelVisible,
    showFilterInput: state.showFilterInput
  })

  // 恢复状态到当前变量
  timelineEntries.value = [...state.timelineEntries]
  filterWindows.value = [...state.filterWindows]
  activeFilterWindow.value = state.activeFilterWindow
  filterInput.value = state.filterInput
  filterMode.value = state.filterMode
  filterOptions.value = { ...state.filterOptions }
  showFilterInput.value = state.showFilterInput
  highlightWords.value = [...(state.highlightWords || [])]  // 新增：恢复高亮词

  // 同步高亮词到appStore
  appStore.setHighlightWords([...highlightWords.value])

  console.log('📋 状态变量已更新:', {
    timelineEntriesCount: timelineEntries.value.length,
    filterWindowsCount: filterWindows.value.length,
    showFilterInput: showFilterInput.value,
    highlightWordsCount: highlightWords.value.length
  })

  // 立即触发相关更新事件
  try {
    // 恢复时间线状态
    if (state.timelineEntries.length > 0) {
      console.log('📅 触发时间线更新事件')
      window.dispatchEvent(new CustomEvent('timelineUpdated', {
        detail: {
          entries: state.timelineEntries,
          visible: state.timelinePanelVisible
        }
      }))
    } else {
      console.log('🧹 触发时间线清空事件')
      window.dispatchEvent(new CustomEvent('timelineCleared'))
    }

    // 通知LogViewer恢复过滤窗口状态
    if (state.filterWindows.length > 0) {
      console.log('🔍 触发过滤窗口恢复事件')
      window.dispatchEvent(new CustomEvent('restoreFilterWindows', {
        detail: {
          windows: state.filterWindows,
          activeWindow: state.activeFilterWindow,
          filterWindowsHeight: state.filterWindowsHeight,
          filterStartLine: state.filterStartLine
        }
      }))
    }

    // 更新窗口选择器（不触发保存）
    console.log('🪟 触发窗口选择器更新事件')
    window.dispatchEvent(new CustomEvent('filterWindowsUpdated', {
      detail: {
        windows: state.filterWindows,
        activeWindow: state.activeFilterWindow,
        skipSave: true // 标记跳过保存
      }
    }))

    console.log('✅ 文件状态恢复事件已触发')
  } finally {
    // 延迟重置标志，确保所有事件处理完成
    setTimeout(() => {
      isRestoring.value = false
      console.log('🔓 状态恢复标志已重置')
    }, 200)
  }
}

// 过滤相关
const showFilterInput = ref(false)
const showFilterExamples = ref(false)
const filterInput = ref('')
const filterMode = ref('include') // include 或 exclude
const filterOptions = ref({
  isRegex: false,
  caseSensitive: false,
  wholeWord: false
})

// 调试信息
console.log('🔍 Toolbar组件已加载，过滤相关变量初始化完成:', {
  showFilterInput: showFilterInput.value,
  filterInput: filterInput.value,
  filterMode: filterMode.value
})

// 窗口选择相关
const currentWindow = ref('main') // 当前选中的窗口
const filterWindows = ref<Array<{
  id: string
  name: string
  filter: string
  mode: string
  options: any
  height?: number
  filteredLines: string[]
  originalLineNumbers: number[]
}>>([])

// 活动过滤窗口
const activeFilterWindow = ref('')

// 计算当前窗口名称
const currentWindowName = computed(() => {
  if (currentWindow.value === 'main') {
    return '主窗口'
  } else if (currentWindow.value === 'filter') {
    return '过滤窗口'
  }
  return '主窗口'
})



// 搜索结果
const searchResults = ref({
  total: 0,
  current: 0
})

// 行号跳转相关
const showLineJumpInput = ref(false)
const targetLineNumber = ref('')

// 计算属性
const searchOptions = computed(() => appStore.searchOptions)

// 方法
const openFile = async () => {
  try {
    // 如果有正在进行的加载，先中断
    if (appStore.isGlobalLoading) {
      console.log('🛑 中断之前的文件加载操作')
      appStore.setGlobalLoading(false)
      // 给一点时间让之前的操作清理
      await new Promise(resolve => setTimeout(resolve, 100))
    }

    appStore.setGlobalLoading(true, '正在选择文件...', 5)

    const filePath = await OpenFileDialog()
    if (filePath) {
      appStore.updateLoadingProgress(15, '正在检查文件信息...')

      const fileInfo = await GetFileInfo(filePath)
      const fileSizeMB = fileInfo.size / (1024 * 1024)

      if (fileInfo) {
        if (fileSizeMB > 100) {
          appStore.updateLoadingProgress(25, `正在处理大文件 (${fileSizeMB.toFixed(1)} MB)...`)
        } else {
          appStore.updateLoadingProgress(25, '正在打开文件...')
        }

        // 使用 AppStore 的 openFile 方法
        await appStore.openFile(filePath)

        // 自动隐藏侧边栏
        if (!appStore.sidebarCollapsed) {
          appStore.toggleSidebar()
        }

        appStore.setGlobalLoading(false)
        ElMessage.success(`已打开文件: ${fileInfo.name}`)
      }
    } else {
      appStore.setGlobalLoading(false)
    }
  } catch (error) {
    console.error('打开文件失败:', error)
    appStore.setGlobalLoading(false)
    ElMessage.error('打开文件失败')
  }
}

// 打开最近文件
const openRecentFile = async (filePath: string) => {
  if (!filePath) return

  try {
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
    } else {
      ElMessage.error('文件不存在或无法访问')
      // 从最近文件列表中移除无效文件
      appStore.removeFromRecentFiles(filePath)
    }
  } catch (error) {
    console.error('打开最近文件失败:', error)
    ElMessage.error('打开文件失败')
  }
}

// 打开文件分片器
const openFileSplitter = async () => {
  try {
    const filePath = await OpenFileDialog()
    if (filePath) {
      const fileInfo = await GetFileInfo(filePath)
      if (fileInfo) {
        showFileSplitter.value = true
        splitterFilePath.value = filePath
        splitterFileName.value = fileInfo.name
        splitterFileSize.value = fileInfo.size
      }
    }
  } catch (error) {
    console.error('选择文件失败:', error)
    ElMessage.error('选择文件失败')
  }
}

// 处理分片完成
const handleSplitComplete = (result) => {
  console.log('分片完成:', result)
  ElMessage.success(`文件已分片为 ${result.totalFiles} 个文件`)
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
    console.log('🔍 执行搜索:', {
      query: searchQuery.value,
      currentWindow: currentWindow.value,
      windowName: currentWindowName.value
    })

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

      // 通知LogViewer执行搜索，指定目标窗口
      window.dispatchEvent(new CustomEvent('performSearch', {
        detail: {
          query: searchQuery.value,
          results: results,
          targetWindow: currentWindow.value,
          caseSensitive: searchOptions.value.caseSensitive
        }
      }))

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
    // 触发 LogViewer 跳转到对应行，指定目标窗口
    const lineNumber = appStore.searchResults[index].lineNumber
    console.log('🎯 跳转到搜索结果:', {
      lineNumber,
      targetWindow: currentWindow.value,
      index: index + 1,
      total: appStore.searchResults.length
    })

    // 搜索结果跳转也需要支持大文件模式
    window.dispatchEvent(new CustomEvent('jumpToLine', {
      detail: {
        lineNumber,
        isSearchResult: true,
        targetWindow: currentWindow.value,
        requireLoading: true // 标记需要加载支持
      }
    }))
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
      appStore.syntaxHighlight = !appStore.syntaxHighlight
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

// 保存项目
const saveProject = async () => {
  try {
    // 确保当前状态已保存
    if (currentFilePath.value && !isSaving.value) {
      saveCurrentState()
    }

    // 生成项目数据
    const projectData = generateProjectData()

    // 生成文件名
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, 19)
    const fileName = `LogTrawl-Project-${timestamp}.ltproj`

    // 创建并下载文件
    const blob = new Blob([JSON.stringify(projectData, null, 2)], {
      type: 'application/json'
    })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = fileName
    link.click()
    URL.revokeObjectURL(url)

    ElMessage.success(`项目已保存为 ${fileName}`)
  } catch (error) {
    console.error('保存项目失败:', error)
    ElMessage.error('保存项目失败: ' + (error as Error).message)
  }
}

// 生成项目数据
const generateProjectData = () => {
  const projectData = {
    version: '1.0.0',
    createdAt: new Date().toISOString(),
    metadata: {
      appVersion: '1.0.0',
      totalFiles: appStore.openFiles.length,
      totalHighlightWords: highlightWords.value.length,
      description: '日志分析项目文件'
    },
    openFiles: appStore.openFiles.map(file => ({
      path: file.path,
      name: file.name,
      encoding: 'utf-8',
      isActive: file.path === appStore.currentFile?.path
    })),
    fileStates: Object.fromEntries(
      Array.from(fileStates.value.entries()).map(([path, state]) => [
        path,
        {
          timelineEntries: state.timelineEntries,
          filterWindows: state.filterWindows.map(fw => ({
            id: fw.id,
            name: fw.name,
            filter: fw.filter,
            mode: fw.mode,
            options: fw.options
          })),
          activeFilterWindow: state.activeFilterWindow,
          filterInput: state.filterInput,
          filterMode: state.filterMode,
          filterOptions: state.filterOptions,
          showFilterInput: state.showFilterInput,
          timelinePanelVisible: state.timelinePanelVisible,
          filterWindowsHeight: state.filterWindowsHeight,
          filterStartLine: state.filterStartLine,
          highlightWords: state.highlightWords || []  // 新增：保存每个文件的高亮词
        }
      ])
    ),
    globalSettings: {
      currentWindow: currentWindow.value,
      showLineNumbers: appStore.showLineNumbers,
      wordWrap: appStore.wordWrap,
      syntaxHighlight: appStore.syntaxHighlight
      // 注意：高亮词现在保存在每个文件的状态中，不再是全局设置
    }
  }

  console.log('📦 生成项目数据:', {
    version: projectData.version,
    totalFiles: projectData.metadata.totalFiles,
    totalHighlightWords: projectData.metadata.totalHighlightWords,
    fileStatesCount: Object.keys(projectData.fileStates).length
  })

  return projectData
}

// 显示过滤输入框
const toggleFilterInput = () => {
  console.log('🔍 切换过滤输入框，当前状态:', showFilterInput.value)
  showFilterInput.value = !showFilterInput.value
  console.log('🔍 切换后状态:', showFilterInput.value)

  if (showFilterInput.value) {
    nextTick(() => {
      const input = document.querySelector('.filter-input input') as HTMLInputElement
      if (input) {
        input.focus()
        console.log('🔍 过滤输入框已聚焦')
      } else {
        console.warn('🔍 未找到过滤输入框元素')
      }
    })
  }
}

// 显示行号跳转输入框
const toggleLineJumpInput = () => {
  showLineJumpInput.value = !showLineJumpInput.value
  if (showLineJumpInput.value) {
    nextTick(() => {
      const input = document.querySelector('.line-number-input input') as HTMLInputElement
      if (input) {
        input.focus()
        input.select()
      }
    })
  }
}

// 切换过滤选项
const toggleFilterOption = (option: keyof typeof filterOptions.value) => {
  filterOptions.value[option] = !filterOptions.value[option]
}

// 应用过滤
const applyFilter = () => {
  try {
    if (!filterInput.value.trim()) {
      ElMessage.warning('请输入过滤条件')
      return
    }

    // 验证当前文件是否存在
    if (!appStore.currentFile) {
      ElMessage.warning('请先打开一个文件')
      return
    }

    console.log('🔍 应用过滤:', {
      filter: filterInput.value,
      mode: filterMode.value,
      currentWindow: currentWindow.value,
      windowName: currentWindowName.value,
      shouldKeepWindow: true // 标记应该保持窗口选择
    })

    // 通过事件通知 LogViewer 应用过滤，指定目标窗口
    const filterData = {
      filter: filterInput.value.trim(),
      mode: filterMode.value,
      options: { ...filterOptions.value },
      sourceWindow: currentWindow.value  // 新增：指定源窗口
    }

    window.dispatchEvent(new CustomEvent('applyFilter', { detail: filterData }))

    // 保存到文件状态
    saveCurrentState()

    // 保持过滤界面显示，不自动隐藏
    ElMessage.success(`过滤已应用到${currentWindowName.value}`)

    // 确保窗口选择保持不变
    console.log('🔍 过滤完成，保持窗口选择:', currentWindow.value)
  } catch (error) {
    console.error('应用过滤失败:', error)
    ElMessage.error('应用过滤失败: ' + (error as Error).message)
  }
}

// 清除过滤
const clearFilter = () => {
  console.log('🧹 清除过滤，当前窗口:', currentWindow.value)

  filterInput.value = ''
  filterWindows.value = []
  activeFilterWindow.value = ''

  // 如果当前选择的是过滤窗口，切换到主窗口（因为没有过滤窗口了）
  if (currentWindow.value === 'filter') {
    console.log('🪟 清除过滤后切换到主窗口')
    currentWindow.value = 'main'
  }

  // 保存到文件状态
  saveCurrentState()

  window.dispatchEvent(new CustomEvent('clearFilter'))
}

// 使用示例
const useExample = (example: string) => {
  filterInput.value = example
  showFilterExamples.value = false
  // 自动显示过滤输入框
  showFilterInput.value = true
  ElMessage.success(`已应用示例: ${example}`)
}

// 取消过滤
const cancelFilter = () => {
  console.log('❌ 取消过滤，当前窗口:', currentWindow.value)

  showFilterInput.value = false
  filterInput.value = ''
  filterWindows.value = []
  activeFilterWindow.value = ''

  // 如果当前选择的是过滤窗口，切换到主窗口（因为没有过滤窗口了）
  if (currentWindow.value === 'filter') {
    console.log('🪟 取消过滤后切换到主窗口')
    currentWindow.value = 'main'
  }

  // 保存到文件状态
  saveCurrentState()

  // 清除过滤
  window.dispatchEvent(new CustomEvent('clearFilter'))
}

// 行号跳转相关方法
const jumpToLine = () => {
  const lineNumber = parseInt(targetLineNumber.value)
  if (isNaN(lineNumber) || lineNumber < 1) {
    ElMessage.warning('请输入有效的行号')
    return
  }

  console.log('🎯 请求跳转到行:', lineNumber)
  // 传递 isSearchResult: false 来触发行高亮，并指定目标窗口
  window.dispatchEvent(new CustomEvent('jumpToLine', {
    detail: {
      lineNumber,
      isSearchResult: false,
      targetWindow: currentWindow.value
    }
  }))
}

const jumpToFirstLine = () => {
  targetLineNumber.value = '1'
  console.log('🎯 请求跳转到首行')
  window.dispatchEvent(new CustomEvent('jumpToLine', {
    detail: {
      lineNumber: 1,
      isSearchResult: false,
      targetWindow: currentWindow.value
    }
  }))
}

const jumpToLastLine = () => {
  // 通知LogViewer跳转到尾行，让它处理大文件模式
  console.log('🎯 请求跳转到尾行')
  window.dispatchEvent(new CustomEvent('jumpToLastLine', {
    detail: {
      targetWindow: currentWindow.value
    }
  }))
}

const jumpToPreviousLine = () => {
  // 获取当前行号并减1
  const current = parseInt(targetLineNumber.value) || 1
  if (current > 1) {
    const newLine = current - 1
    targetLineNumber.value = newLine.toString()
    console.log('🎯 请求跳转到上一行:', newLine)
    window.dispatchEvent(new CustomEvent('jumpToLine', {
      detail: {
        lineNumber: newLine,
        isSearchResult: false,
        targetWindow: currentWindow.value
      }
    }))
  }
}

const jumpToNextLine = () => {
  // 获取当前行号并加1
  const current = parseInt(targetLineNumber.value) || 0
  const newLine = current + 1
  targetLineNumber.value = newLine.toString()
  console.log('🎯 请求跳转到下一行:', newLine)
  window.dispatchEvent(new CustomEvent('jumpToLine', {
    detail: {
      lineNumber: newLine,
      isSearchResult: false,
      targetWindow: currentWindow.value
    }
  }))
}

const trackToCurrentLine = () => {
  // 触发事件获取当前显示的行号
  window.dispatchEvent(new CustomEvent('getCurrentLine'))
}

// 监听当前行响应事件
window.addEventListener('currentLineResponse', (event: any) => {
  targetLineNumber.value = event.detail.lineNumber.toString()
  // 跳转到该行并高亮
  console.log('🎯 跟踪到当前行:', event.detail.lineNumber)
  window.dispatchEvent(new CustomEvent('jumpToLine', {
    detail: {
      lineNumber: event.detail.lineNumber,
      isSearchResult: false,
      targetWindow: currentWindow.value
    }
  }))
})

// 高亮词管理方法
const addHighlightWord = () => {
  if (!newHighlightWord.value.text.trim()) {
    ElMessage.warning('请输入要高亮的单词')
    return
  }

  // 检查是否已存在
  const exists = highlightWords.value.some(word => word.text === newHighlightWord.value.text.trim())
  if (exists) {
    ElMessage.warning('该单词已存在')
    return
  }

  const newWord = {
    text: newHighlightWord.value.text.trim(),
    color: newHighlightWord.value.color
  }

  highlightWords.value.push(newWord)

  // 同步到appStore
  appStore.setHighlightWords([...highlightWords.value])

  // 重置表单
  newHighlightWord.value = {
    text: '',
    color: '#337ecc'
  }

  showAddHighlightDialog.value = false
  ElMessage.success('高亮词添加成功')
}

const removeHighlightWord = (index: number) => {
  highlightWords.value.splice(index, 1)
  // 同步到appStore
  appStore.setHighlightWords([...highlightWords.value])
  ElMessage.success('高亮词已删除')
}

const handleHighlightWordCommand = (command: string, index: number) => {
  const word = highlightWords.value[index]

  switch (command) {
    case 'filter':
      // 添加到过滤
      filterInput.value = word.text
      filterMode.value = 'include'
      showFilterInput.value = true
      applyFilter()
      break
    case 'changeColor':
      // 更改颜色
      currentEditingIndex.value = index
      selectedColor.value = word.color
      showColorDialog.value = true
      break
    case 'remove':
      // 删除高亮
      removeHighlightWord(index)
      break
  }
}

const changeHighlightColor = () => {
  if (currentEditingIndex.value >= 0) {
    highlightWords.value[currentEditingIndex.value].color = selectedColor.value
    // 同步到appStore
    appStore.setHighlightWords([...highlightWords.value])
    showColorDialog.value = false
    currentEditingIndex.value = -1
    ElMessage.success('颜色已更改')
  }
}

// 获取对比色（用于文字颜色）
const getContrastColor = (hexColor: string) => {
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

// 时间线相关方法
const addTimeline = () => {
  // 优先使用当前聚焦的日志行信息
  if (currentFocusedLogLine.value) {
    const focusedLine = currentFocusedLogLine.value
    const logContent = focusedLine.content

    // 从日志内容中提取时间戳
    let logTimestamp = ''
    const timestampPatterns = [
      /(\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2}\s*[+-]\d{4})/,  // Apache格式
      /(\d{4}-\d{2}-\d{2}[T\s]\d{2}:\d{2}:\d{2})/,           // ISO格式
      /(\d{2}:\d{2}:\d{2})/                                   // 时间格式
    ]

    for (const pattern of timestampPatterns) {
      const match = logContent.match(pattern)
      if (match) {
        logTimestamp = match[1]
        break
      }
    }

    // 设置对话框的初始值
    newTimelineEntry.value = {
      note: '',
      logContent: logContent,
      logTimestamp: logTimestamp,
      lineNumber: focusedLine.lineNumber
    }
  } else {
    // 如果没有聚焦行，尝试使用选中文本
    const selection = window.getSelection()
    let selectedText = ''
    let logTimestamp = ''

    if (selection && selection.toString().trim()) {
      selectedText = selection.toString().trim()

      // 尝试从选中文本中提取时间戳
      const timestampPatterns = [
        /(\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2}\s*[+-]\d{4})/,  // Apache格式
        /(\d{4}-\d{2}-\d{2}[T\s]\d{2}:\d{2}:\d{2})/,           // ISO格式
        /(\d{2}:\d{2}:\d{2})/                                   // 时间格式
      ]

      for (const pattern of timestampPatterns) {
        const match = selectedText.match(pattern)
        if (match) {
          logTimestamp = match[1]
          break
        }
      }
    }

    // 设置对话框的初始值
    newTimelineEntry.value = {
      note: '',
      logContent: selectedText,
      logTimestamp: logTimestamp,
      lineNumber: null
    }
  }

  showAddTimelineDialog.value = true
}

// 快速添加时间线（快捷键t触发）
const quickAddTimeline = (focusedLineInfo: any) => {
  const logContent = focusedLineInfo.content

  // 从日志内容中提取时间戳
  let logTimestamp = ''
  const timestampPatterns = [
    /(\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2}\s*[+-]\d{4})/,  // Apache格式
    /(\d{4}-\d{2}-\d{2}[T\s]\d{2}:\d{2}:\d{2})/,           // ISO格式
    /(\d{2}:\d{2}:\d{2})/                                   // 时间格式
  ]

  for (const pattern of timestampPatterns) {
    const match = logContent.match(pattern)
    if (match) {
      logTimestamp = match[1]
      break
    }
  }

  // 直接添加时间线条目，使用默认备注
  const currentTime = new Date().toISOString()
  const newEntry = {
    id: Date.now(),
    timestamp: currentTime,
    note: `行 ${focusedLineInfo.lineNumber} 的日志`,
    logContent: logContent,
    logTimestamp: logTimestamp,
    lineNumber: focusedLineInfo.lineNumber,
    color: timelineColor.value
  }

  timelineEntries.value.push(newEntry)

  // 保存到文件状态
  saveCurrentState()

  // 触发时间线更新事件
  window.dispatchEvent(new CustomEvent('timelineUpdated', {
    detail: { entries: timelineEntries.value }
  }))

  ElMessage.success(`已快速添加第 ${focusedLineInfo.lineNumber} 行到时间线`)
}

// 快速过滤（快捷键f触发）
const quickFilter = (filterInfo: any) => {
  try {
    const filterText = filterInfo.text
    const isReverse = filterInfo.reverse || false

    // 验证过滤文本
    if (!filterText || !filterText.trim()) {
      ElMessage.warning('过滤文本不能为空')
      return
    }

    // 设置过滤器值
    filterInput.value = filterText.trim()
    filterMode.value = isReverse ? 'exclude' : 'include'

    // 显示过滤界面
    showFilterInput.value = true

    // 应用过滤
    applyFilter()

    const action = isReverse ? '反向过滤' : '过滤'
    ElMessage.success(`已对 "${filterText}" 应用${action}`)
  } catch (error) {
    console.error('快速过滤失败:', error)
    ElMessage.error('快速过滤失败: ' + (error as Error).message)
  }
}

// 确认添加时间线条目
const confirmAddTimeline = () => {
  if (!newTimelineEntry.value.note.trim()) {
    ElMessage.warning('请输入备注信息')
    return
  }

  const currentTime = new Date().toISOString()
  const newEntry = {
    id: Date.now(),
    timestamp: currentTime,
    note: newTimelineEntry.value.note.trim(),
    logContent: newTimelineEntry.value.logContent,
    logTimestamp: newTimelineEntry.value.logTimestamp,
    lineNumber: newTimelineEntry.value.lineNumber,
    color: timelineColor.value
  }



  timelineEntries.value.push(newEntry)

  // 保存到文件状态
  saveCurrentState()

  // 触发时间线更新事件
  window.dispatchEvent(new CustomEvent('timelineUpdated', {
    detail: { entries: timelineEntries.value }
  }))

  showAddTimelineDialog.value = false
  ElMessage.success('时间线条目已添加')
}

const clearTimeline = () => {
  timelineEntries.value = []

  // 触发时间线清除事件
  window.dispatchEvent(new CustomEvent('timelineCleared'))

  ElMessage.success('时间线已清除')
}

// 生成时间线文本的函数
const generateTimelineText = () => {
  const sortedEntries = [...timelineEntries.value].sort((a, b) => {
    // 优先按日志时间戳排序，如果没有则按创建时间排序
    const timeA = a.logTimestamp || a.timestamp
    const timeB = b.logTimestamp || b.timestamp
    return new Date(timeA).getTime() - new Date(timeB).getTime()
  })

  const lines = []
  lines.push('='.repeat(60))
  lines.push('时间线报告')
  lines.push('='.repeat(60))
  lines.push(`生成时间: ${new Date().toLocaleString('zh-CN')}`)
  lines.push(`条目数量: ${sortedEntries.length}`)
  lines.push('')

  sortedEntries.forEach((entry, index) => {
    lines.push(`${index + 1}. ${entry.note}`)
    lines.push(`   时间: ${entry.logTimestamp || '未知'}`)
    if (entry.lineNumber) {
      lines.push(`   行号: ${entry.lineNumber}`)
    }
    lines.push(`   创建: ${new Date(entry.timestamp).toLocaleString('zh-CN')}`)
    if (entry.logContent) {
      lines.push(`   内容: ${entry.logContent}`)
    }
    lines.push('')
  })

  return lines.join('\n')
}

// 通用复制到剪贴板函数
const copyToClipboard = async (text, type = '内容') => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(`${type}已复制到剪贴板`)
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = text
    textArea.style.position = 'fixed'
    textArea.style.left = '-999999px'
    textArea.style.top = '-999999px'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()

    try {
      document.execCommand('copy')
      ElMessage.success(`${type}已复制到剪贴板`)
    } catch (fallbackErr) {
      ElMessage.error('复制失败，请手动复制')
      console.error('复制失败:', fallbackErr)
    } finally {
      document.body.removeChild(textArea)
    }
  }
}

// 处理时间线复制命令
const handleTimelineCopyCommand = (command) => {
  if (timelineEntries.value.length === 0) {
    ElMessage.warning('时间线为空，无法复制')
    return
  }

  let content = ''
  let type = ''

  switch (command) {
    case 'detailed':
      content = generateTimelineText()
      type = '详细报告'
      break
    case 'simple':
      content = generateSimpleTimelineText()
      type = '简单列表'
      break
    case 'csv':
      content = generateTimelineCSV()
      type = 'CSV格式'
      break
    case 'json':
      content = generateTimelineJSON()
      type = 'JSON格式'
      break
    default:
      content = generateTimelineText()
      type = '时间线'
  }

  copyToClipboard(content, type)
}

// 生成简单时间线文本
const generateSimpleTimelineText = () => {
  const sortedEntries = [...timelineEntries.value].sort((a, b) => {
    const timeA = a.logTimestamp || a.timestamp
    const timeB = b.logTimestamp || b.timestamp
    return new Date(timeA).getTime() - new Date(timeB).getTime()
  })

  return sortedEntries.map((entry, index) => {
    const time = entry.logTimestamp || new Date(entry.timestamp).toLocaleString('zh-CN')
    const line = entry.lineNumber ? ` (行${entry.lineNumber})` : ''
    return `${index + 1}. ${time}${line}: ${entry.note}`
  }).join('\n')
}

// 生成CSV格式
const generateTimelineCSV = () => {
  const sortedEntries = [...timelineEntries.value].sort((a, b) => {
    const timeA = a.logTimestamp || a.timestamp
    const timeB = b.logTimestamp || b.timestamp
    return new Date(timeA).getTime() - new Date(timeB).getTime()
  })

  const headers = ['序号', '备注', '日志时间', '行号', '创建时间', '日志内容']
  const csvLines = [headers.join(',')]

  sortedEntries.forEach((entry, index) => {
    const row = [
      index + 1,
      `"${entry.note.replace(/"/g, '""')}"`,
      `"${entry.logTimestamp || ''}"`,
      entry.lineNumber || '',
      `"${new Date(entry.timestamp).toLocaleString('zh-CN')}"`,
      `"${(entry.logContent || '').replace(/"/g, '""')}"`
    ]
    csvLines.push(row.join(','))
  })

  return csvLines.join('\n')
}

// 生成JSON格式
const generateTimelineJSON = () => {
  const sortedEntries = [...timelineEntries.value].sort((a, b) => {
    const timeA = a.logTimestamp || a.timestamp
    const timeB = b.logTimestamp || b.timestamp
    return new Date(timeA).getTime() - new Date(timeB).getTime()
  })

  const exportData = {
    exportTime: new Date().toISOString(),
    totalEntries: sortedEntries.length,
    entries: sortedEntries.map((entry, index) => ({
      index: index + 1,
      id: entry.id,
      note: entry.note,
      logTimestamp: entry.logTimestamp,
      lineNumber: entry.lineNumber,
      createdAt: entry.timestamp,
      logContent: entry.logContent,
      color: entry.color
    }))
  }

  return JSON.stringify(exportData, null, 2)
}

// 监听键盘事件（E键添加高亮）
const handleKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'e' || event.key === 'E') {
    const selection = window.getSelection()
    if (selection && selection.toString().trim()) {
      const selectedText = selection.toString().trim()

      // 检查是否已存在
      const exists = highlightWords.value.some(word => word.text === selectedText)
      if (!exists) {
        const newWord = {
          text: selectedText,
          color: presetColors[highlightWords.value.length % presetColors.length]
        }
        highlightWords.value.push(newWord)
        // 同步到appStore
        appStore.setHighlightWords([...highlightWords.value])
        ElMessage.success(`已添加高亮词: ${selectedText}`)
      } else {
        ElMessage.warning('该单词已存在')
      }

      // 清除选择
      selection.removeAllRanges()
    }
  }
}

// 文件切换防抖
const fileChangeTimeout = ref<number | null>(null)

// 监听文件切换
watch(() => appStore.currentFile, (newFile, oldFile) => {
  // 清除之前的定时器
  if (fileChangeTimeout.value) {
    clearTimeout(fileChangeTimeout.value)
  }

  // 防抖处理，避免快速切换时的问题
  fileChangeTimeout.value = window.setTimeout(() => {
    console.log('🔄 文件切换:', {
      oldFile: oldFile?.name,
      newFile: newFile?.name,
      currentFilePath: currentFilePath.value
    })

    // 保存旧文件的状态
    if (oldFile?.path && currentFilePath.value && !isRestoring.value) {
      console.log('💾 保存旧文件状态:', oldFile.name)
      saveCurrentState()
    }

    // 更新当前文件路径
    currentFilePath.value = newFile?.path || ''

    // 恢复新文件的状态
    if (newFile?.path) {
      console.log('📂 恢复新文件状态:', newFile.name)
      restoreFileState()
    } else {
      console.log('🧹 清空所有状态')
      // 如果没有文件，清空所有状态
      timelineEntries.value = []
      filterWindows.value = []
      activeFilterWindow.value = ''
      filterInput.value = ''
      filterMode.value = 'include'
      filterOptions.value = {
        caseSensitive: false,
        wholeWord: false,
        isRegex: false
      }
      showFilterInput.value = false
      highlightWords.value = []  // 新增：清空高亮词

      // 同步清空appStore的高亮词
      appStore.setHighlightWords([])

      // 触发清空事件
      nextTick(() => {
        window.dispatchEvent(new CustomEvent('timelineCleared'))
        window.dispatchEvent(new CustomEvent('restoreFilterWindows', {
          detail: { windows: [], activeWindow: '', filterWindowsHeight: 300, filterStartLine: 0 }
        }))
        window.dispatchEvent(new CustomEvent('filterWindowsUpdated', {
          detail: { windows: [], activeWindow: '', skipSave: true }
        }))
      })
    }
  }, 50) // 50ms 防抖
}, { immediate: true })

// 窗口大小监听
const handleResize = () => {
  windowWidth.value = window.innerWidth
}

// 更多菜单命令处理
const handleMoreCommand = (command: string) => {
  switch (command) {
    case 'saveProject':
      saveProject()
      break
    case 'viewOptions':
      // 可以打开视图选项对话框或者直接显示子菜单
      break
    case 'exportData':
      exportData()
      break
    case 'about':
      showSettings.value = true
      // 可以切换到关于标签页
      break
  }
}

// 添加键盘事件监听
onMounted(() => {
  document.addEventListener('keydown', handleKeyDown)
  window.addEventListener('resize', handleResize)
  // 从appStore加载高亮词
  highlightWords.value = [...appStore.highlightWords]

  // 监听聚焦行变化事件
  window.addEventListener('focusedLineChanged', (event: any) => {
    currentFocusedLogLine.value = event.detail
  })

  // 监听快速添加时间线事件
  window.addEventListener('quickAddTimeline', (event: any) => {
    quickAddTimeline(event.detail)
  })

  // 监听快速过滤事件
  window.addEventListener('quickFilter', (event: any) => {
    quickFilter(event.detail)
  })

  // 监听添加高亮事件（来自右键菜单）
  window.addEventListener('addHighlight', (event: any) => {
    const selectedText = event.detail.text
    if (selectedText) {
      // 检查是否已存在
      const exists = highlightWords.value.some(word => word.text === selectedText)
      if (!exists) {
        const newWord = {
          text: selectedText,
          color: presetColors[highlightWords.value.length % presetColors.length]
        }
        highlightWords.value.push(newWord)
        // 同步到appStore
        appStore.setHighlightWords([...highlightWords.value])
        ElMessage.success(`已添加高亮词: ${selectedText}`)
      } else {
        ElMessage.warning('该单词已存在')
      }
    }
  })

  // 监听保存文件状态事件
  window.addEventListener('saveFileState', () => {
    saveCurrentState()
  })

  // 监听时间线面板显示状态变化
  window.addEventListener('timelinePanelVisibilityChanged', (event: any) => {
    const state = getCurrentFileState()
    state.timelinePanelVisible = event.detail.visible
    saveCurrentState()
  })

  // 监听强制刷新UI事件
  window.addEventListener('forceRefreshUI', () => {
    console.log('🔄 收到强制刷新UI事件')

    // 强制触发响应式更新
    const currentState = getCurrentFileState()

    // 重新设置所有状态以触发响应式更新
    timelineEntries.value = [...currentState.timelineEntries]
    filterWindows.value = [...currentState.filterWindows]
    activeFilterWindow.value = currentState.activeFilterWindow
    filterInput.value = currentState.filterInput
    filterMode.value = currentState.filterMode
    filterOptions.value = { ...currentState.filterOptions }
    showFilterInput.value = currentState.showFilterInput

    // 确保高亮词同步到appStore（现在是文件独立的）
    appStore.setHighlightWords([...highlightWords.value])
    console.log('🎨 强制刷新时同步高亮词到appStore:', highlightWords.value.length)

    // 触发相关组件的更新
    nextTick(() => {
      // 触发时间线更新
      if (timelineEntries.value.length > 0) {
        window.dispatchEvent(new CustomEvent('timelineUpdated', {
          detail: {
            entries: timelineEntries.value,
            visible: true
          }
        }))
      }

      // 触发过滤窗口更新
      if (filterWindows.value.length > 0) {
        window.dispatchEvent(new CustomEvent('restoreFilterWindows', {
          detail: {
            windows: filterWindows.value,
            activeWindow: activeFilterWindow.value,
            filterWindowsHeight: 300,
            filterStartLine: 0
          }
        }))
      }

      // 触发高亮词更新
      if (highlightWords.value.length > 0) {
        window.dispatchEvent(new CustomEvent('highlightWordsUpdated', {
          detail: { highlightWords: highlightWords.value }
        }))
      }
    })

    console.log('🔄 强制刷新完成:', {
      timelineCount: timelineEntries.value.length,
      filterCount: filterWindows.value.length,
      showFilterInput: showFilterInput.value,
      highlightWordsCount: highlightWords.value.length
    })
  })

  // 监听项目设置恢复事件
  window.addEventListener('restoreProjectSettings', (event: any) => {
    const { currentWindow: projectCurrentWindow } = event.detail

    console.log('🎨 开始恢复项目设置:', {
      currentWindow: projectCurrentWindow
    })

    // 恢复当前窗口
    currentWindow.value = projectCurrentWindow

    console.log('🎨 项目设置恢复完成 (高亮词现在是文件独立的)')
  })

  // 监听项目文件状态恢复事件
  window.addEventListener('restoreProjectFileStates', (event: any) => {
    const { fileStates: projectFileStates, activeFilePath } = event.detail

    console.log('📥 收到项目文件状态恢复事件:', {
      fileStatesCount: Object.keys(projectFileStates).length,
      activeFilePath
    })

    // 清空当前文件状态
    fileStates.value.clear()

    // 恢复所有文件状态
    Object.entries(projectFileStates).forEach(([filePath, state]: [string, any]) => {
      const restoredState = {
        timelineEntries: state.timelineEntries || [],
        filterWindows: state.filterWindows || [],
        activeFilterWindow: state.activeFilterWindow || '',
        filterInput: state.filterInput || '',
        filterMode: state.filterMode || 'include',
        filterOptions: state.filterOptions || {
          caseSensitive: false,
          wholeWord: false,
          isRegex: false
        },
        showFilterInput: state.showFilterInput || false,
        timelinePanelVisible: state.timelinePanelVisible || false,
        filterWindowsHeight: state.filterWindowsHeight || 300,
        filterStartLine: state.filterStartLine || 0,
        highlightWords: state.highlightWords || [],  // 新增：恢复每个文件的高亮词
        lastAccessed: Date.now()
      }

      fileStates.value.set(filePath, restoredState)

      console.log(`📄 恢复文件状态 ${filePath.split('/').pop()}:`, {
        timelineCount: restoredState.timelineEntries.length,
        filterCount: restoredState.filterWindows.length,
        timelinePanelVisible: restoredState.timelinePanelVisible
      })
    })

    console.log('📁 项目文件状态恢复完成:', {
      fileStatesCount: fileStates.value.size,
      activeFilePath
    })

    // 如果有活动文件，延迟恢复其状态，确保日志内容已加载
    if (activeFilePath) {
      currentFilePath.value = activeFilePath
      console.log('🔄 准备恢复活动文件状态:', activeFilePath.split('/').pop())

      // 延迟恢复状态，确保日志内容已加载
      setTimeout(() => {
        console.log('🔄 开始恢复活动文件状态:', activeFilePath.split('/').pop())
        restoreFileState()
      }, 500) // 给日志内容加载留出时间
    }
  })

  // 监听来自LogViewer的过滤窗口状态更新
  window.addEventListener('filterWindowsUpdated', (event: any) => {
    const { windows, activeWindow, skipSave } = event.detail

    console.log('📥 收到过滤窗口更新事件:', {
      windowsCount: windows.length,
      activeWindow,
      skipSave,
      isRestoring: isRestoring.value
    })

    // 如果是恢复过程中的更新，跳过保存
    if (skipSave || isRestoring.value) {
      console.log('⏭️ 跳过过滤窗口状态保存')
      return
    }

    // 同步过滤窗口状态到当前变量
    filterWindows.value = [...windows]
    if (activeWindow !== undefined) {
      activeFilterWindow.value = activeWindow
    }

    // 保存状态（带防重复机制）
    if (!isSaving.value) {
      console.log('💾 保存过滤窗口状态变化')
      saveCurrentState()
    }
  })

  // 初始化当前文件状态
  if (appStore.currentFile?.path) {
    currentFilePath.value = appStore.currentFile.path
    restoreFileState()
  }
})

// 清理事件监听
onUnmounted(() => {
  // 清理定时器
  if (fileChangeTimeout.value) {
    clearTimeout(fileChangeTimeout.value)
  }

  // 保存当前状态
  if (currentFilePath.value && !isRestoring.value) {
    saveCurrentState()
  }

  document.removeEventListener('keydown', handleKeyDown)
  window.removeEventListener('resize', handleResize)
  window.removeEventListener('focusedLineChanged', () => {})
  window.removeEventListener('quickAddTimeline', () => {})
  window.removeEventListener('quickFilter', () => {})
  window.removeEventListener('addHighlight', () => {})
  window.removeEventListener('saveFileState', () => {})
  window.removeEventListener('timelinePanelVisibilityChanged', () => {})
  window.removeEventListener('filterWindowsUpdated', () => {})
})

// 处理窗口选择
const handleWindowSelect = (windowId: string) => {
  console.log('🪟 切换窗口:', windowId)
  currentWindow.value = windowId

  // 通知LogViewer窗口切换，用于聚焦管理
  window.dispatchEvent(new CustomEvent('windowChanged', {
    detail: {
      windowId,
      windowName: windowId === 'main' ? '主窗口' : '过滤窗口'
    }
  }))

  // 保存当前文件状态
  saveCurrentState()
}



// 监听过滤窗口列表更新
window.addEventListener('filterWindowsUpdated', (event: any) => {
  const previousWindows = filterWindows.value
  filterWindows.value = event.detail.windows

  console.log('📋 过滤窗口列表更新:', {
    currentWindow: currentWindow.value,
    previousCount: previousWindows.length,
    newCount: filterWindows.value.length,
    windowIds: filterWindows.value.map(w => w.id)
  })

  // 只有在以下情况下才切换到主窗口：
  // 1. 当前选择的是具体的过滤窗口ID，且该窗口被删除
  // 2. 当前选择的是'filter'，但没有任何过滤窗口
  if (currentWindow.value !== 'main') {
    if (currentWindow.value === 'filter') {
      // 如果选择的是通用过滤窗口，但没有过滤窗口了，才切换到主窗口
      if (filterWindows.value.length === 0) {
        console.log('🪟 没有过滤窗口了，切换到主窗口')
        currentWindow.value = 'main'
      } else {
        console.log('🪟 保持过滤窗口选择，当前有', filterWindows.value.length, '个过滤窗口')
      }
    } else {
      // 如果选择的是具体的过滤窗口ID，检查该窗口是否还存在
      if (!filterWindows.value.find(w => w.id === currentWindow.value)) {
        console.log('🪟 选中的过滤窗口被删除，切换到主窗口')
        currentWindow.value = 'main'
      }
    }
  }
})

// 打开分析页面
const openAnalysisPage = () => {
  if (!appStore.currentFile) {
    ElMessage.warning('请先打开一个日志文件')
    return
  }

  console.log('🔍 打开数据分析页面:', appStore.currentFile.name)

  // 通知主布局切换到分析页面
  window.dispatchEvent(new CustomEvent('openAnalysisPage', {
    detail: {
      filePath: appStore.currentFile.path,
      fileName: appStore.currentFile.name
    }
  }))
}

// 组件挂载时加载最近文件
onMounted(() => {
  appStore.loadRecentFiles()
})

// 调试：获取文件状态信息
const getFileStatesInfo = () => {
  const info = Array.from(fileStates.value.entries()).map(([filePath, state]) => ({
    filePath: filePath.split('/').pop() || filePath, // 只显示文件名
    timelineCount: state.timelineEntries.length,
    filterCount: state.filterWindows.length,
    lastAccessed: new Date(state.lastAccessed).toLocaleTimeString()
  }))
  console.log('文件状态缓存:', info)
  return info
}

// 调试：生成测试项目数据
const generateTestProject = () => {
  const testProject = {
    version: '1.0.0',
    createdAt: new Date().toISOString(),
    metadata: {
      appVersion: '1.0.0',
      totalFiles: 1,
      totalHighlightWords: 2,
      description: '测试项目文件'
    },
    openFiles: [
      {
        path: 'C:\\test\\sample.log',
        name: 'sample.log',
        encoding: 'utf-8',
        isActive: true
      }
    ],
    fileStates: {
      'C:\\test\\sample.log': {
        timelineEntries: [
          {
            id: 'test-1',
            note: '测试时间线条目',
            logContent: '这是一个测试日志行',
            logTimestamp: '2024-07-24 10:00:00',
            lineNumber: 1,
            timestamp: new Date().toISOString(),
            color: '#ff0000'
          }
        ],
        filterWindows: [],
        activeFilterWindow: '',
        filterInput: '',
        filterMode: 'include',
        filterOptions: {
          caseSensitive: false,
          wholeWord: false,
          isRegex: false
        },
        showFilterInput: false,
        timelinePanelVisible: true,
        filterWindowsHeight: 300,
        filterStartLine: 0
      }
    },
    globalSettings: {
      highlightWords: [
        { text: 'ERROR', color: '#ff0000' },
        { text: 'WARNING', color: '#ffaa00' }
      ],
      currentWindow: 'main',
      showLineNumbers: true,
      wordWrap: true,
      syntaxHighlight: true
    }
  }

  console.log('🧪 测试项目数据:', testProject)
  return testProject
}

// 调试：检查当前UI状态
const checkCurrentUIState = () => {
  console.log('🔍 当前UI状态检查:')
  console.log('📁 当前文件:', appStore.currentFile?.name)
  console.log('📂 打开的文件:', appStore.openFiles.map(f => f.name))
  console.log('📅 时间线条目:', timelineEntries.value.length)
  console.log('🔍 过滤窗口:', filterWindows.value.length)
  console.log('🎨 当前文件高亮词:', highlightWords.value.length)
  console.log('🪟 当前窗口:', currentWindow.value)
  console.log('👁️ 显示过滤输入:', showFilterInput.value)
  console.log('📊 文件状态缓存:', fileStates.value.size)

  // 检查时间线面板状态
  const timelinePanel = document.querySelector('.timeline-panel')
  console.log('📅 时间线面板DOM:', timelinePanel ? '存在' : '不存在')

  // 检查过滤窗口DOM
  const filterWindowElements = document.querySelectorAll('.filter-window')
  console.log('🔍 过滤窗口DOM:', filterWindowElements.length, '个')

  return {
    currentFile: appStore.currentFile?.name,
    openFiles: appStore.openFiles.length,
    timelineEntries: timelineEntries.value.length,
    filterWindows: filterWindows.value.length,
    highlightWords: highlightWords.value.length,
    currentWindow: currentWindow.value,
    showFilterInput: showFilterInput.value,
    fileStatesCount: fileStates.value.size
  }
}

// 在开发模式下暴露调试方法
if (import.meta.env.DEV) {
  const win = window as any
  win.getFileStatesInfo = getFileStatesInfo
  win.generateTestProject = generateTestProject
  win.checkCurrentUIState = checkCurrentUIState
  win.forceRefreshUI = () => {
    window.dispatchEvent(new CustomEvent('forceRefreshUI'))
  }
  win.manualTriggerUpdates = () => {
    console.log('🔧 手动触发界面更新...')

    // 触发时间线更新
    if (timelineEntries.value.length > 0) {
      window.dispatchEvent(new CustomEvent('timelineUpdated', {
        detail: {
          entries: timelineEntries.value,
          visible: true
        }
      }))
      console.log('📅 已触发时间线更新')
    }

    // 触发过滤窗口更新
    if (filterWindows.value.length > 0) {
      window.dispatchEvent(new CustomEvent('restoreFilterWindows', {
        detail: {
          windows: filterWindows.value,
          activeWindow: activeFilterWindow.value,
          filterWindowsHeight: 300,
          filterStartLine: 0
        }
      }))
      console.log('🔍 已触发过滤窗口更新')
    }

    // 触发高亮词更新（现在是文件独立的）
    // 确保同步到appStore
    appStore.setHighlightWords([...highlightWords.value])

    if (highlightWords.value.length > 0) {
      window.dispatchEvent(new CustomEvent('highlightWordsUpdated', {
        detail: { highlightWords: highlightWords.value }
      }))
      console.log('🎨 已触发高亮词更新并同步到appStore')
    } else {
      console.log('🎨 当前文件无高亮词')
    }

    console.log('✅ 手动触发完成')
  }
  win.debugFilterWindows = () => {
    console.log('🔍 过滤窗口调试信息:')
    console.log('Toolbar 过滤窗口数量:', filterWindows.value.length)
    filterWindows.value.forEach((window, index) => {
      console.log(`  ${index + 1}. ${window.name}:`, {
        filter: window.filter,
        mode: window.mode,
        hasFilteredLines: !!window.filteredLines,
        filteredLinesCount: window.filteredLines?.length || 0
      })
    })
    console.log('活动过滤窗口:', activeFilterWindow.value)

    // 检查LogViewer中的状态
    const logViewerWindows = document.querySelectorAll('.filter-window-tab')
    console.log('DOM中的过滤窗口标签数量:', logViewerWindows.length)

    return {
      toolbarWindows: filterWindows.value.length,
      activeWindow: activeFilterWindow.value,
      domTabs: logViewerWindows.length
    }
  }
  win.debugHighlightWords = () => {
    console.log('🎨 高亮词调试信息:')
    console.log('Toolbar 高亮词数量:', highlightWords.value.length)
    console.log('AppStore 高亮词数量:', appStore.highlightWords.length)

    console.log('Toolbar 高亮词:')
    highlightWords.value.forEach((word, index) => {
      console.log(`  ${index + 1}. "${word.text}" (${word.color})`)
    })

    console.log('AppStore 高亮词:')
    appStore.highlightWords.forEach((word, index) => {
      console.log(`  ${index + 1}. "${word.text}" (${word.color})`)
    })

    // 检查DOM中是否有高亮
    const highlightedElements = document.querySelectorAll('.highlight-word')
    console.log('DOM中的高亮元素数量:', highlightedElements.length)

    return {
      toolbarWords: highlightWords.value.length,
      appStoreWords: appStore.highlightWords.length,
      domHighlights: highlightedElements.length,
      synced: highlightWords.value.length === appStore.highlightWords.length
    }
  }
  win.forceHighlightWordsSync = () => {
    console.log('🔧 强制同步高亮词...')
    appStore.setHighlightWords([...highlightWords.value])
    window.dispatchEvent(new CustomEvent('highlightWordsUpdated', {
      detail: { highlightWords: highlightWords.value }
    }))
    console.log('✅ 高亮词强制同步完成')
  }
  win.debugFileHighlightWords = () => {
    console.log('🎨 文件独立高亮词调试信息:')
    console.log('当前文件:', appStore.currentFile?.name)
    console.log('当前文件高亮词数量:', highlightWords.value.length)
    console.log('AppStore 高亮词数量:', appStore.highlightWords.length)

    // 显示所有文件的高亮词
    console.log('所有文件的高亮词状态:')
    fileStates.value.forEach((state, filePath) => {
      const fileName = filePath.split('/').pop() || filePath
      console.log(`  📄 ${fileName}:`, {
        highlightWords: state.highlightWords?.length || 0,
        words: state.highlightWords?.map(w => w.text) || []
      })
    })

    return {
      currentFile: appStore.currentFile?.name,
      currentFileHighlightWords: highlightWords.value.length,
      appStoreHighlightWords: appStore.highlightWords.length,
      allFilesHighlightWords: Object.fromEntries(
        Array.from(fileStates.value.entries()).map(([path, state]) => [
          path.split('/').pop() || path,
          state.highlightWords?.length || 0
        ])
      )
    }
  }
  win.saveTestProject = () => {
    const testData = generateTestProject()
    const blob = new Blob([JSON.stringify(testData, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'test-project.ltproj'
    link.click()
    URL.revokeObjectURL(url)
    console.log('💾 测试项目已下载')
  }
}
</script>

<style scoped>
.toolbar {
  height: 72px;
  display: flex;
  align-items: center;
  padding: 8px 16px;
  background-color: #ffffff;
  border-bottom: 1px solid #e5e7eb;
  gap: 8px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

/* 主要操作区域样式 */
.main-operations {
  display: flex;
  flex-direction: column;
  gap: 2px;
  height: 56px;
}

.main-buttons-row {
  display: flex;
  gap: 2px;
}

.main-button-group {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.main-btn {
  width: 56px;
  height: 40px;
  padding: 0;
  border-radius: 6px;
  font-size: 28px;
}

.main-btn.splitter-btn {
  color: #fa8c16;
}

.main-btn.splitter-btn:hover {
  color: #d46b08;
  background-color: #fff7e6;
}

.main-btn-label {
  font-size: 11px;
  color: #606266;
  text-align: center;
  line-height: 1;
  white-space: nowrap;
}

.recent-file-option {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 200px;
}

.file-name {
  font-weight: 500;
  color: #303133;
  font-size: 13px;
}

.file-path {
  font-size: 11px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 180px;
}

.toolbar-center {
  flex: 1;
  display: flex;
  justify-content: flex-start;
  align-items: center;
  min-width: 0;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.right-buttons-row {
  display: flex;
  gap: 2px;
  align-items: center;
}

.core-buttons {
  display: flex;
  align-items: center;
  gap: 2px;
}

.responsive-buttons {
  display: flex;
  align-items: center;
  gap: 2px;
}

.more-menu {
  margin-left: 4px;
}

/* 紧凑模式样式 */
.main-button-group.compact {
  flex-direction: column;
  gap: 2px;
  min-width: auto;
  height: auto;
}

.main-button-group.compact .main-btn {
  padding: 6px 8px;
  min-height: 32px;
}

.main-button-group.compact .main-btn-label {
  display: none;
}

/* 高亮词管理区域样式 */
.highlight-words-section {
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: flex-start;
  height: 56px;
  justify-content: center;
  width: 200px;
}

.add-highlight-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  border: 1px dashed #d9d9d9;
  border-radius: 4px;
  cursor: pointer;
  color: #666;
  font-size: 12px;
  transition: all 0.3s;
  width: 100%;
  height: 100%;
  justify-content: center;
  background-color: #ffffff;
}

.add-highlight-btn:hover {
  border-color: #409eff;
  color: #409eff;
}

.highlight-words-list {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 2px;
  border: 1px dashed #d9d9d9;
  border-radius: 4px;
  padding: 6px 8px;
  background-color: #ffffff;
  transition: all 0.3s;
}

.highlight-words-list:hover {
  border-color: #409eff;
}

.highlight-words-container {
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
  align-items: center;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.highlight-word-tag {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  padding: 2px 6px;
  border-radius: 8px;
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  max-width: 80px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  height: 18px;
  flex-shrink: 0;
}

.highlight-word-tag:hover {
  transform: scale(1.05);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.highlight-word-tag .remove-icon {
  font-size: 12px;
  cursor: pointer;
  opacity: 0.7;
  transition: opacity 0.3s;
}

.highlight-word-tag .remove-icon:hover {
  opacity: 1;
}

.add-more-btn {
  width: 18px;
  height: 18px;
  font-size: 10px;
  flex-shrink: 0;
  min-width: 18px;
}

/* 时间线区域样式 */
.timeline-section {
  display: flex;
  align-items: center;
  height: 56px;
  width: 80px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 4px;
  background-color: #ffffff;
}

.timeline-buttons {
  display: flex;
  flex-direction: column;
  gap: 2px;
  width: 100%;
  height: 100%;
}

.timeline-row {
  display: flex;
  gap: 2px;
  justify-content: space-between;
  height: 50%;
}

.timeline-btn {
  width: 32px;
  height: 20px;
  padding: 0;
  font-size: 12px;
  border-radius: 3px;
  flex: 1;
}

.timeline-btn .el-icon {
  font-size: 12px;
}

/* 对话框样式 */
.add-highlight-content {
  padding: 16px 0;
}

.color-selection {
  display: flex;
  align-items: center;
  gap: 12px;
}

.preset-colors {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.color-option {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.3s;
}

.color-option:hover {
  transform: scale(1.1);
}

.color-option.active {
  border-color: #409eff;
  transform: scale(1.1);
}

.color-option.large {
  width: 32px;
  height: 32px;
}

.color-change-content {
  padding: 16px 0;
  text-align: center;
}

/* 窗口选择和操作区域样式 */
.window-operations {
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: flex-start;
  min-width: 120px;
  height: 56px;
  justify-content: center;
}

/* 窗口选择器行样式 */
.window-selector-row {
  display: flex;
  justify-content: flex-start;
  width: 100%;
}

.window-selector-btn {
  width: 120px;
  height: 32px;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  background-color: #ffffff;
  color: #606266;
  font-weight: 500;
  font-size: 13px;
  transition: all 0.2s ease;
}

.window-selector-btn:hover {
  border-color: #409eff;
  color: #409eff;
  background-color: #ecf5ff;
}

/* 刷新和编码选择行样式 */
.refresh-encoding-row {
  display: flex;
  gap: 8px;
  align-items: center;
  width: 100%;
  justify-content: flex-start;
  height: 22px;
}

.refresh-btn {
  width: 32px;
  height: 22px;
  padding: 0;
  border-radius: 4px;
}

.encoding-select {
  width: 80px;
  height: 22px;
}

:deep(.encoding-select .el-input__wrapper) {
  height: 22px !important;
  min-height: 22px !important;
}

:deep(.encoding-select .el-input__inner) {
  height: 22px !important;
  line-height: 22px !important;
  font-size: 11px !important;
}

.search-container {
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: flex-start;
  height: 56px;
  justify-content: center;
  width: auto;
}

.search-input {
  width: 202px;
  height: 32px;
}

.search-controls {
  display: flex;
  gap: 2px;
  align-items: center;
  justify-content: flex-start;
  width: 100%;
  height: 22px;
  margin: 0;
  padding: 0;
}

.search-control-btn {
  width: 32px;
  height: 22px;
  padding: 0 !important;
  margin: 0 !important;
  font-size: 11px;
  font-weight: bold;
  border-radius: 4px;
  flex-shrink: 0;
}

:deep(.search-control-btn) {
  margin: 0 !important;
  padding: 0 !important;
}

:deep(.el-button) {
  margin: 0 !important;
}

:deep(.el-button + .el-button) {
  margin-left: 0 !important;
}

/* 行号跳转区域样式 */
.line-jump-section {
  display: flex;
  align-items: center;
  width: auto;
  height: 56px;
  flex-shrink: 0;
}

.line-jump-button-group {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  width: 56px;
}

.line-jump-toggle-btn {
  width: 56px;
  height: 40px;
  padding: 0;
  border-radius: 6px;
  font-size: 28px;
}

.line-jump-btn-label {
  font-size: 11px;
  color: #606266;
  text-align: center;
  line-height: 1;
  white-space: nowrap;
}

.line-jump-container {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: flex-start;
  height: 56px;
  justify-content: center;
  width: 120px;
}

.close-line-jump-btn {
  position: absolute;
  top: -2px;
  right: -2px;
  width: 16px;
  height: 16px;
  padding: 0;
  font-size: 12px;
  z-index: 10;
}

.line-number-input-row {
  display: flex;
  gap: 2px;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  height: 32px;
}

.line-control-btn {
  width: 24px;
  height: 32px;
  padding: 0 !important;
  margin: 0 !important;
  font-size: 14px;
  font-weight: bold;
  border-radius: 4px;
  flex-shrink: 0;
}

.line-number-input {
  width: 70px;
  height: 32px;
}

:deep(.line-number-input .el-input__wrapper) {
  height: 32px !important;
  min-height: 32px !important;
}

:deep(.line-number-input .el-input__inner) {
  height: 32px !important;
  line-height: 32px !important;
  font-size: 12px !important;
  text-align: center;
}

.line-navigation-controls {
  display: flex;
  gap: 2px;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  height: 22px;
}

.line-nav-btn {
  width: 22px;
  height: 22px;
  padding: 0 !important;
  margin: 0 !important;
  font-size: 11px;
  border-radius: 4px;
  flex-shrink: 0;
}

:deep(.line-nav-btn) {
  margin: 0 !important;
  padding: 0 !important;
}

.search-control-btn.active {
  background-color: #409eff;
  border-color: #409eff;
  color: #ffffff;
}

/* 搜索输入框容器样式 */
.search-input-container {
  position: relative;
  width: 202px;
  height: 32px;
}

/* 搜索结果计数覆盖层样式 - 显示在输入框内右侧 */
.search-count-overlay {
  position: absolute;
  right: 30px; /* 留出清除按钮的空间 */
  top: 50%;
  transform: translateY(-50%);
  font-size: 11px;
  color: #6b7280;
  background-color: rgba(255, 255, 255, 0.9);
  padding: 2px 6px;
  border-radius: 3px;
  pointer-events: none; /* 防止阻挡输入框交互 */
  z-index: 10;
  text-align: center;
  white-space: nowrap;
}

/* 原来的搜索计数样式 - 保持兼容性 */
.search-count {
  font-size: 11px;
  color: #6b7280;
  text-align: center;
  margin-top: 2px;
}

/* 过滤器区域样式 */
.filter-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  height: 56px;
}

.filter-button-group {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.filter-toggle-btn {
  width: 56px;
  height: 40px;
  padding: 0;
  border-radius: 6px;
  font-size: 28px;
}

.filter-btn-label {
  font-size: 11px;
  color: #606266;
  text-align: center;
  line-height: 1;
  white-space: nowrap;
}

.filter-container {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 200px;
  align-items: flex-start;
  height: 56px;
  justify-content: center;
}

.filter-input-row {
  display: flex;
  gap: 0;
  align-items: center;
  width: 100%;
  height: 32px;
}

.filter-mode-select {
  width: 60px;
  height: 32px;
  flex-shrink: 0;
}

:deep(.filter-mode-select .el-input__wrapper) {
  height: 32px !important;
  border-top-right-radius: 0 !important;
  border-bottom-right-radius: 0 !important;
}

:deep(.filter-mode-select .el-input__inner) {
  height: 32px !important;
  line-height: 32px !important;
  font-size: 12px !important;
  text-align: center;
}

.filter-input {
  width: 140px;
  height: 32px;
  flex: 1;
}

:deep(.filter-input .el-input__wrapper) {
  height: 32px !important;
  border-top-left-radius: 0 !important;
  border-bottom-left-radius: 0 !important;
  border-left: 0 !important;
}

:deep(.filter-input .el-input__inner) {
  height: 32px !important;
  line-height: 32px !important;
  font-size: 12px !important;
}

.filter-controls {
  display: flex;
  gap: 0;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  height: 22px;
}

.filter-control-btn {
  width: 32px;
  height: 22px;
  padding: 0;
  font-size: 11px;
  font-weight: bold;
  border-radius: 4px;
  flex-shrink: 0;
}

.filter-action-btn {
  width: 40px;
  height: 22px;
  padding: 0;
  font-size: 11px;
  border-radius: 4px;
  flex-shrink: 0;
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

/* 确保下拉菜单在所有面板之上 */
:deep(.el-dropdown-menu) {
  z-index: 2000 !important;
}

/* 确保弹出层在所有面板之上 */
:deep(.el-popper) {
  z-index: 2000 !important;
}

/* 添加时间线对话框样式 */
.add-timeline-content {
  padding: 16px 0;
}

.timeline-field {
  margin-bottom: 20px;
}

.field-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #606266;
  font-size: 14px;
}

.required {
  color: #f56c6c;
}

/* 聚焦行信息样式 */
.focused-line-info {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 12px;
}

.line-number-badge {
  display: inline-block;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 12px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 12px;
  margin-bottom: 8px;
}

.line-content-preview {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  color: #374151;
  background: white;
  padding: 8px 10px;
  border-radius: 4px;
  border: 1px solid #d1d5db;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 100px;
  overflow-y: auto;
}

.timeline-field .color-picker-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.timeline-field .color-preview {
  margin-left: 12px;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: white;
  background-color: #409eff;
}

/* 时间线颜色对话框样式 */
.timeline-color-content {
  padding: 16px 0;
}

.timeline-color-content .color-section {
  margin-bottom: 16px;
}

.timeline-color-content label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #606266;
}

.timeline-color-content .color-picker-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* 分割线样式 */
:deep(.el-divider--vertical) {
  height: 48px !important;
  margin: 0 8px !important;
}

/* 过滤示例对话框样式 */
.filter-examples {
  max-height: 500px;
  overflow-y: auto;
}

.example-section {
  margin-bottom: 24px;
}

.example-section h4 {
  margin: 0 0 12px 0;
  color: #409eff;
  font-size: 14px;
  font-weight: 600;
  border-bottom: 1px solid #e4e7ed;
  padding-bottom: 8px;
}

.example-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.example-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.example-item:hover {
  background: #e3f2fd;
  border-color: #409eff;
  transform: translateY(-1px);
}

.example-item code {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  background: #fff;
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid #ddd;
  color: #e83e8c;
  font-size: 13px;
  font-weight: 500;
  flex-shrink: 0;
  margin-right: 12px;
}

.example-desc {
  color: #666;
  font-size: 13px;
  flex: 1;
  text-align: right;
}

.filter-help-btn {
  padding: 0 !important;
  margin-left: 4px !important;
  color: #909399 !important;
}

.filter-help-btn:hover {
  color: #409eff !important;
}

/* 响应式媒体查询 */
@media (max-width: 1400px) {
  .toolbar {
    padding: 0 12px;
  }

  .main-button-group {
    min-width: 60px;
  }

  .main-btn-label {
    font-size: 11px;
  }
}

@media (max-width: 1200px) {
  .toolbar {
    padding: 0 8px;
  }

  .main-button-group {
    min-width: 50px;
    gap: 1px;
  }

  .main-btn {
    padding: 6px 8px;
  }

  .main-btn-label {
    font-size: 10px;
  }
}

@media (max-width: 900px) {
  .toolbar {
    padding: 0 6px;
  }

  .main-operations .main-buttons-row {
    gap: 4px;
  }

  .window-operations {
    gap: 4px;
  }

  .right-buttons-row {
    gap: 1px;
  }
}

/* 搜索输入框容器样式 */
.search-input-container {
  position: relative;
  width: 202px;
  height: 32px;
}

/* 搜索结果计数覆盖层样式 - 显示在输入框内右侧 */
.search-count-overlay {
  position: absolute;
  right: 30px; /* 留出清除按钮的空间 */
  top: 50%;
  transform: translateY(-50%);
  font-size: 11px;
  color: #6b7280;
  background-color: rgba(255, 255, 255, 0.9);
  padding: 2px 6px;
  border-radius: 3px;
  pointer-events: none; /* 防止阻挡输入框交互 */
  z-index: 10;
  text-align: center;
  white-space: nowrap;
}

@media (max-width: 768px) {
  .toolbar {
    padding: 0 4px;
    min-height: 50px;
  }

  .main-button-group {
    min-width: 40px;
  }

  .main-btn {
    padding: 4px 6px;
    min-height: 28px;
  }

  .search-container {
    max-width: 150px;
  }

  .search-input {
    font-size: 12px;
  }
}
</style>
