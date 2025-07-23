<template>
  <div v-if="appStore.isGlobalLoading" class="global-loading-overlay">
    <div class="loading-container">
      <div class="loading-content">
        <!-- 加载动画 -->
        <div class="loading-spinner">
          <div class="spinner-ring"></div>
          <div class="spinner-ring"></div>
          <div class="spinner-ring"></div>
        </div>
        
        <!-- 加载信息 -->
        <div class="loading-info">
          <h3 class="loading-title">{{ appStore.loadingMessage || '正在处理...' }}</h3>
          
          <!-- 进度条 -->
          <div v-if="appStore.loadingProgress > 0" class="progress-container">
            <div class="progress-bar">
              <div 
                class="progress-fill" 
                :style="{ width: `${appStore.loadingProgress}%` }"
              ></div>
            </div>
            <span class="progress-text">{{ appStore.loadingProgress.toFixed(1) }}%</span>
          </div>
          
          <!-- 提示信息 -->
          <p class="loading-hint">
            <template v-if="isLargeFile">
              正在处理大文件，请稍候...
            </template>
            <template v-else>
              请稍候，正在加载文件...
            </template>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 判断是否为大文件
const isLargeFile = computed(() => {
  return appStore.loadingMessage.includes('大文件') || 
         appStore.loadingMessage.includes('分块') ||
         appStore.loadingMessage.includes('MB')
})
</script>

<style scoped>
.global-loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  animation: fadeIn 0.3s ease-out;
}

.loading-container {
  background: white;
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  max-width: 400px;
  width: 90%;
  text-align: center;
  animation: slideIn 0.4s ease-out;
}

.loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

/* 加载动画 */
.loading-spinner {
  position: relative;
  width: 60px;
  height: 60px;
}

.spinner-ring {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 3px solid transparent;
  border-top: 3px solid #409eff;
  border-radius: 50%;
  animation: spin 1.2s linear infinite;
}

.spinner-ring:nth-child(2) {
  width: 80%;
  height: 80%;
  top: 10%;
  left: 10%;
  border-top-color: #67c23a;
  animation-delay: -0.4s;
}

.spinner-ring:nth-child(3) {
  width: 60%;
  height: 60%;
  top: 20%;
  left: 20%;
  border-top-color: #e6a23c;
  animation-delay: -0.8s;
}

/* 加载信息 */
.loading-info {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.loading-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

/* 进度条 */
.progress-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #409eff, #67c23a);
  border-radius: 4px;
  transition: width 0.3s ease;
  animation: progressPulse 2s ease-in-out infinite;
}

.progress-text {
  font-size: 14px;
  font-weight: 500;
  color: #606266;
  min-width: 50px;
}

.loading-hint {
  margin: 0;
  font-size: 14px;
  color: #909399;
  line-height: 1.5;
}

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes progressPulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}

/* 响应式设计 */
@media (max-width: 480px) {
  .loading-container {
    padding: 30px 20px;
    margin: 20px;
  }
  
  .loading-title {
    font-size: 16px;
  }
  
  .loading-spinner {
    width: 50px;
    height: 50px;
  }
}
</style>
