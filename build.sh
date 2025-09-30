#!/bin/bash

# LogTrawl 构建脚本
# 支持Windows、Linux和macOS平台构建

echo "========================================"
echo "LogTrawl 应用程序构建脚本"
echo "========================================"

# 检查是否安装了Wails
if ! command -v wails &> /dev/null
then
    echo "错误: 未找到 Wails 命令。请先安装 Wails v2:"
    echo "请运行: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

# 获取当前目录
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$PROJECT_DIR/dist"

# 创建构建目录
mkdir -p "$BUILD_DIR"

# 显示帮助信息
show_help() {
    echo ""
    echo "用法:"
    echo "  ./build.sh [平台]"
    echo ""
    echo "平台选项:"
    echo "  win     - 构建 Windows 应用程序"
    echo "  linux   - 构建 Linux 应用程序"
    echo "  mac     - 构建 macOS 应用程序"
    echo "  all     - 构建所有平台应用程序"
    echo "  clean   - 清理构建目录"
    echo ""
    echo "示例:"
    echo "  ./build.sh win"
    echo "  ./build.sh all"
    echo ""
}

# 处理命令参数
if [ $# -eq 0 ]; then
    show_help
    exit 0
fi

# 清理构建目录
clean_build() {
    echo "清理构建目录..."
    if [ -d "$BUILD_DIR" ]; then
        rm -rf "$BUILD_DIR"
        echo "构建目录已清理"
    else
        echo "构建目录不存在"
    fi
}

# 构建Windows应用程序
build_windows() {
    echo ""
    echo "========================================"
    echo "构建 Windows 应用程序"
    echo "========================================"
    cd "$PROJECT_DIR"
    
    # 构建Windows应用程序
    wails build -platform windows/amd64 -o "$BUILD_DIR/LogTrawl-windows-amd64.exe"
    
    if [ $? -eq 0 ]; then
        echo "Windows 应用程序构建成功: $BUILD_DIR/LogTrawl-windows-amd64.exe"
    else
        echo "Windows 应用程序构建失败"
        exit 1
    fi
}

# 构建Linux应用程序
build_linux() {
    echo ""
    echo "========================================"
    echo "构建 Linux 应用程序"
    echo "========================================"
    cd "$PROJECT_DIR"
    
    # 构建Linux应用程序
    wails build -platform linux/amd64 -o "$BUILD_DIR/LogTrawl-linux-amd64"
    
    if [ $? -eq 0 ]; then
        echo "Linux 应用程序构建成功: $BUILD_DIR/LogTrawl-linux-amd64"
    else
        echo "Linux 应用程序构建失败"
        exit 1
    fi
}

# 构建macOS应用程序
build_macos() {
    echo ""
    echo "========================================"
    echo "构建 macOS 应用程序"
    echo "========================================"
    cd "$PROJECT_DIR"
    
    # 构建macOS应用程序
    wails build -platform darwin/amd64 -o "$BUILD_DIR/LogTrawl-darwin-amd64.app"
    
    if [ $? -eq 0 ]; then
        echo "macOS 应用程序构建成功: $BUILD_DIR/LogTrawl-darwin-amd64.app"
    else
        echo "macOS 应用程序构建失败"
        exit 1
    fi
}

# 处理命令参数
case "$1" in
    clean)
        clean_build
        ;;
    win)
        build_windows
        ;;
    linux)
        build_linux
        ;;
    mac)
        build_macos
        ;;
    all)
        build_windows
        build_linux
        build_macos
        echo ""
        echo "========================================"
        echo "所有平台构建完成!"
        echo "========================================"
        ;;
    *)
        echo "未知的平台选项: $1"
        show_help
        exit 1
        ;;
esac