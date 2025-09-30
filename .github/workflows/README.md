# GitHub Actions 工作流说明

这个项目使用GitHub Actions来自动化构建和发布应用程序到Windows、Linux和macOS平台。

## 工作流触发条件

1. 当向`main`或`master`分支推送代码时
2. 当创建以`v`开头的标签时（例如`v1.0.0`）
3. 当创建Pull Request到`main`或`master`分支时
4. 手动触发（通过GitHub界面）

## 构建过程

工作流会并行在三个平台上构建应用程序：
- Ubuntu (Linux)
- macOS
- Windows

每个平台的构建产物会被上传为artifact。

## 发布过程

当创建以`v`开头的标签时（例如`v1.0.0`），工作流会自动创建一个GitHub Release，并将所有平台的构建产物附加到Release中。

## 缓存优化

为了加速构建过程，工作流使用了以下缓存机制：
- Go模块缓存
- Go构建缓存
- Node.js模块缓存

## 本地测试

你可以通过以下命令在本地测试构建过程：

```bash
# 在Linux/macOS上
./build.sh all

# 在Windows上
build.bat all
```