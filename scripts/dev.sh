#!/bin/bash

echo "🚀 LogTrawl 开发模式"
echo "==================="

# 切换到项目根目录
cd "$(dirname "$0")/.."

echo "📁 当前目录: $(pwd)"

# 检查依赖
echo "🔍 检查依赖..."
if ! command -v wails &> /dev/null; then
    echo "❌ Wails CLI 未安装，请先安装 Wails"
    echo "   go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go 1.21+"
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo "❌ Node.js/npm 未安装，请先安装 Node.js 16+"
    exit 1
fi

echo "✅ 依赖检查通过"

# 安装前端依赖
echo "📦 安装前端依赖..."
cd frontend
npm install
cd ..

# 启动开发模式
echo "🔥 启动开发模式..."
echo "📝 提示: 应用将在开发模式下运行，支持热重载"
echo "🌐 前端开发服务器: http://localhost:5173"
echo "🖥️  应用程序将自动打开"
echo ""

wails dev

echo "👋 开发模式已退出"
