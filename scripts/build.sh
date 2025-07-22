#!/bin/bash

echo "🚀 LogTrawl 构建脚本"
echo "===================="

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

# 清理之前的构建
echo "🧹 清理之前的构建..."
rm -rf build/bin/*

# 构建应用
echo "🔨 开始构建应用..."
wails build --clean

if [ $? -eq 0 ]; then
    echo "✅ 构建成功！"
    echo "📍 可执行文件位置: build/bin/"
    ls -la build/bin/
else
    echo "❌ 构建失败"
    exit 1
fi

echo "🎉 构建完成！"
