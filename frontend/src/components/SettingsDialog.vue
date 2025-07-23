<template>
  <el-dialog
    v-model="visible"
    title="设置"
    width="800px"
    top="8vh"
    :show-close="true"
    class="settings-dialog"
  >
    <div class="settings-content">
      <el-tabs v-model="activeTab" class="settings-tabs">
        <!-- 显示设置 -->
        <el-tab-pane label="显示设置" name="display">
          <div class="settings-section">
            <div class="setting-group">
              <h4>日志字体</h4>
              <div class="setting-grid">
                <div class="setting-item">
                  <span class="setting-label">字体大小</span>
                  <div class="setting-control">
                    <el-input-number
                      v-model="settings.fontSize"
                      :min="10"
                      :max="24"
                      size="small"
                      controls-position="right"
                      style="width: 100px"
                      @change="updateFontSize"
                    />
                  </div>
                </div>

                <div class="setting-item">
                  <span class="setting-label">字体族</span>
                  <div class="setting-control">
                    <el-select v-model="settings.fontFamily" @change="updateFontFamily" style="width: 140px">
                      <el-option label="Courier New" value="Courier New" />
                      <el-option label="Consolas" value="Consolas" />
                      <el-option label="Monaco" value="Monaco" />
                      <el-option label="Menlo" value="Menlo" />
                      <el-option label="Ubuntu Mono" value="Ubuntu Mono" />
                      <el-option label="Source Code Pro" value="Source Code Pro" />
                    </el-select>
                  </div>
                </div>
              </div>
            </div>

            <div class="setting-group">
              <h4>编码与行距</h4>
              <div class="setting-grid">
                <div class="setting-item">
                  <span class="setting-label">默认编码</span>
                  <div class="setting-control">
                    <el-select v-model="settings.defaultEncoding" @change="updateEncoding" style="width: 120px">
                      <el-option label="UTF-8" value="utf-8" />
                      <el-option label="GBK" value="gbk" />
                      <el-option label="ASCII" value="ascii" />
                      <el-option label="ISO-8859-1" value="iso-8859-1" />
                    </el-select>
                  </div>
                </div>

                <div class="setting-item">
                  <span class="setting-label">行高倍数</span>
                  <div class="setting-control">
                    <el-input-number
                      v-model="settings.lineHeight"
                      :min="1.0"
                      :max="3.0"
                      :step="0.1"
                      :precision="1"
                      size="small"
                      controls-position="right"
                      style="width: 100px"
                      @change="updateLineHeight"
                    />
                  </div>
                </div>
              </div>
            </div>

            <div class="setting-group">
              <h4>显示选项</h4>
              <div class="setting-grid">
                <div class="setting-item">
                  <span class="setting-label">显示行号</span>
                  <div class="setting-control">
                    <el-switch v-model="settings.showLineNumbers" @change="updateShowLineNumbers" />
                  </div>
                </div>

                <div class="setting-item">
                  <span class="setting-label">自动换行</span>
                  <div class="setting-control">
                    <el-switch v-model="settings.wordWrap" @change="updateWordWrap" />
                  </div>
                </div>

                <div class="setting-item">
                  <span class="setting-label">语法高亮</span>
                  <div class="setting-control">
                    <el-switch v-model="settings.syntaxHighlight" @change="updateSyntaxHighlight" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>

        <!-- 快捷键 -->
        <el-tab-pane label="快捷键" name="shortcuts">
          <div class="shortcuts-section">
            <div class="shortcuts-grid">
              <div class="shortcuts-category">
                <h4>🎨 高亮操作</h4>
                <div class="shortcut-list">
                  <div class="shortcut-item">
                    <span class="shortcut-desc">添加高亮词</span>
                    <div class="shortcut-steps">
                      <span class="step">1. 选中文本</span>
                      <span class="step">2. 按 <kbd class="shortcut-key">E</kbd></span>
                    </div>
                  </div>
                </div>
              </div>

              <div class="shortcuts-category">
                <h4>⏰ 时间线操作</h4>
                <div class="shortcut-list">
                  <div class="shortcut-item">
                    <span class="shortcut-desc">快速添加时间线</span>
                    <div class="shortcut-steps">
                      <span class="step">1. 点击选择日志行</span>
                      <span class="step">2. 按 <kbd class="shortcut-key">T</kbd></span>
                    </div>
                  </div>
                </div>
              </div>

              <div class="shortcuts-category">
                <h4>🔍 过滤操作</h4>
                <div class="shortcut-list">
                  <div class="shortcut-item">
                    <span class="shortcut-desc">快速过滤</span>
                    <div class="shortcut-steps">
                      <span class="step">1. 选中要过滤的文本</span>
                      <span class="step">2. 按 <kbd class="shortcut-key">F</kbd></span>
                    </div>
                  </div>
                </div>
              </div>

              <div class="shortcuts-category">
                <h4>📋 文本操作</h4>
                <div class="shortcut-list">
                  <div class="shortcut-item">
                    <span class="shortcut-desc">全选当前窗口文本</span>
                    <div class="shortcut-steps">
                      <span class="step">按 <kbd class="shortcut-key">Ctrl + A</kbd></span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="shortcuts-tips">
              <h4>💡 使用提示</h4>
              <ul class="tips-list">
                <li><strong>高亮词：</strong>选中任意文本后按 E 键即可快速添加高亮，支持多种颜色自动分配</li>
                <li><strong>时间线：</strong>点击日志行使其获得焦点，然后按 T 键可快速添加到时间线</li>
                <li><strong>过滤：</strong>选中关键词后按 F 键可快速创建包含该词的过滤窗口</li>
                <li><strong>窗口切换：</strong>可在主窗口和过滤窗口之间切换，快捷键在各窗口中都有效</li>
              </ul>
            </div>
          </div>
        </el-tab-pane>

        <!-- 关于 -->
        <el-tab-pane label="关于" name="about">
          <div class="about-section">
            <div class="about-header">
              <div class="app-logo">
                <el-icon size="48" color="#409eff">
                  <Document />
                </el-icon>
              </div>
              <div class="app-basic-info">
                <h2 class="app-name">LogTrawl</h2>
                <p class="app-version">版本 v1.0.0</p>
                <p class="app-description">
                  专业的日志文件分析工具，支持大文件处理、智能过滤、语法高亮等功能
                </p>
              </div>
            </div>

            <div class="about-content">
              <div class="info-section">
                <div class="author-info">
                  <h4>👨‍💻 开发者信息</h4>
                  <div class="info-grid">
                    <div class="info-item">
                      <span class="info-label">作者</span>
                      <span class="info-value">onewin</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">GitHub</span>
                      <a href="https://github.com/onewinner" target="_blank" class="info-link">
                        github.com/onewinner
                      </a>
                    </div>
                  </div>
                </div>

                <div class="features-info">
                  <h4>✨ 主要功能</h4>
                  <div class="features-grid">
                    <div class="feature-item">
                      <span class="feature-icon">🚀</span>
                      <span class="feature-text">大文件快速处理</span>
                    </div>
                    <div class="feature-item">
                      <span class="feature-icon">🔍</span>
                      <span class="feature-text">智能搜索过滤</span>
                    </div>
                    <div class="feature-item">
                      <span class="feature-icon">🎨</span>
                      <span class="feature-text">语法高亮</span>
                    </div>
                    <div class="feature-item">
                      <span class="feature-icon">📊</span>
                      <span class="feature-text">数据分析统计</span>
                    </div>
                    <div class="feature-item">
                      <span class="feature-icon">✂️</span>
                      <span class="feature-text">文件分片</span>
                    </div>
                    <div class="feature-item">
                      <span class="feature-icon">📝</span>
                      <span class="feature-text">时间线标注</span>
                    </div>
                    <div class="feature-item">
                      <span class="feature-icon">💾</span>
                      <span class="feature-text">项目保存</span>
                    </div>
                    <div class="feature-item">
                      <span class="feature-icon">⚡</span>
                      <span class="feature-text">快捷键操作</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="resetToDefaults">恢复默认</el-button>
        <el-button type="primary" @click="visible = false">确定</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAppStore } from '@/stores/app'
import { ElMessage } from 'element-plus'
import { Document } from '@element-plus/icons-vue'

// Props
interface Props {
  modelValue: boolean
}

const props = defineProps<Props>()
const emit = defineEmits(['update:modelValue'])

const appStore = useAppStore()

// 响应式数据
const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const activeTab = ref('display')

// 设置数据
const settings = ref({
  fontSize: 14,
  fontFamily: 'Courier New',
  defaultEncoding: 'utf-8',
  lineHeight: 1.5,
  showLineNumbers: true,
  wordWrap: true,
  syntaxHighlight: true
})

// 初始化设置
onMounted(() => {
  loadSettings()
})

// 加载设置
const loadSettings = () => {
  const savedSettings = localStorage.getItem('logtrawl-settings')
  if (savedSettings) {
    try {
      const parsed = JSON.parse(savedSettings)
      settings.value = { ...settings.value, ...parsed }
    } catch (error) {
      console.error('加载设置失败:', error)
    }
  }
  
  // 同步到appStore
  appStore.showLineNumbers = settings.value.showLineNumbers
  appStore.wordWrap = settings.value.wordWrap
  appStore.syntaxHighlight = settings.value.syntaxHighlight
}

// 保存设置
const saveSettings = () => {
  localStorage.setItem('logtrawl-settings', JSON.stringify(settings.value))
  applySettings()
}

// 应用设置
const applySettings = () => {
  // 应用字体设置
  const logViewer = document.querySelector('.log-viewer') as HTMLElement
  if (logViewer) {
    logViewer.style.setProperty('--log-font-size', `${settings.value.fontSize}px`)
    logViewer.style.setProperty('--log-font-family', settings.value.fontFamily)
    logViewer.style.setProperty('--log-line-height', settings.value.lineHeight.toString())
  }
  
  // 同步到appStore
  appStore.showLineNumbers = settings.value.showLineNumbers
  appStore.wordWrap = settings.value.wordWrap
  appStore.syntaxHighlight = settings.value.syntaxHighlight
  
  // 触发重新渲染
  window.dispatchEvent(new CustomEvent('settingsChanged', {
    detail: settings.value
  }))
}

// 字体大小调整
const updateFontSize = () => {
  saveSettings()
}

// 字体族更新
const updateFontFamily = () => {
  saveSettings()
}

// 编码更新
const updateEncoding = () => {
  saveSettings()
  ElMessage.success('编码设置已更新')
}

// 行高调整
const updateLineHeight = () => {
  saveSettings()
}

// 显示选项更新
const updateShowLineNumbers = () => {
  saveSettings()
}

const updateWordWrap = () => {
  saveSettings()
}

const updateSyntaxHighlight = () => {
  saveSettings()
}

// 恢复默认设置
const resetToDefaults = () => {
  settings.value = {
    fontSize: 14,
    fontFamily: 'Courier New',
    defaultEncoding: 'utf-8',
    lineHeight: 1.5,
    showLineNumbers: true,
    wordWrap: true,
    syntaxHighlight: true
  }
  saveSettings()
  ElMessage.success('已恢复默认设置')
}
</script>

<style scoped>
.settings-dialog {
  --el-dialog-border-radius: 12px;
}

.settings-content {
  padding: 0;
}

.settings-tabs {
  --el-tabs-header-height: 48px;
}

.settings-section {
  padding: 16px 0;
}

.setting-group {
  margin-bottom: 24px;
}

.setting-group h4 {
  margin: 0 0 12px 0;
  font-size: 15px;
  font-weight: 600;
  color: #409eff;
  border-bottom: 1px solid #e4e7ed;
  padding-bottom: 6px;
}

.setting-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px 24px;
  align-items: center;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: #f8f9fa;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.setting-item:hover {
  background: #e9ecef;
}

.setting-label {
  font-size: 13px;
  color: #333;
  font-weight: 500;
  white-space: nowrap;
  margin-right: 12px;
}

.setting-control {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 快捷键样式 */
.shortcuts-section {
  padding: 16px 0;
}

.shortcuts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.shortcuts-category {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 16px;
}

.shortcuts-category h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #409eff;
  display: flex;
  align-items: center;
  gap: 6px;
}

.shortcut-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.shortcut-item {
  padding: 8px 0;
}

.shortcut-desc {
  font-size: 13px;
  color: #333;
  font-weight: 500;
  margin-bottom: 4px;
}

.shortcut-steps {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.step {
  font-size: 12px;
  color: #666;
}

.shortcut-key {
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 3px;
  padding: 2px 6px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 11px;
  color: #666;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  display: inline-block;
}

.shortcuts-tips {
  background: #fff7e6;
  border: 1px solid #ffd591;
  border-radius: 8px;
  padding: 16px;
}

.shortcuts-tips h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #fa8c16;
}

.tips-list {
  margin: 0;
  padding-left: 16px;
  list-style: none;
}

.tips-list li {
  margin-bottom: 8px;
  font-size: 13px;
  color: #666;
  line-height: 1.5;
  position: relative;
}

.tips-list li::before {
  content: '•';
  color: #fa8c16;
  position: absolute;
  left: -12px;
}

/* 关于页面样式 */
.about-section {
  padding: 16px 0;
}

.about-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border-radius: 12px;
}

.app-logo {
  flex-shrink: 0;
}

.app-basic-info {
  flex: 1;
}

.app-name {
  margin: 0 0 4px 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.app-version {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #409eff;
  font-weight: 500;
}

.app-description {
  margin: 0;
  font-size: 13px;
  color: #666;
  line-height: 1.5;
}

.about-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.author-info,
.features-info {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 16px;
}

.author-info h4,
.features-info h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #409eff;
}

.info-grid {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 0;
}

.info-label {
  font-size: 13px;
  color: #333;
  font-weight: 500;
}

.info-value {
  font-size: 13px;
  color: #666;
}

.info-link {
  font-size: 13px;
  color: #409eff;
  text-decoration: none;
}

.info-link:hover {
  text-decoration: underline;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  background: #fff;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.feature-item:hover {
  background: #e9ecef;
}

.feature-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.feature-text {
  font-size: 12px;
  color: #666;
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
