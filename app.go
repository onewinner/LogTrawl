package main

import (
	"context"

	"LogTrawl/backend/models"
	"LogTrawl/backend/services"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx                 context.Context
	fileService         *services.FileService
	searchService       *services.SearchService
	recentFilesService  *services.RecentFilesService
	dialogService       *services.DialogService
	systemService       *services.SystemService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// 初始化服务
	a.fileService = services.NewFileService()
	a.searchService = services.NewSearchService(a.fileService)
	a.recentFilesService = services.NewRecentFilesService()
	a.dialogService = services.NewDialogService(ctx)
	a.systemService = services.NewSystemService()
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return "Hello " + name + ", It's show time!"
}

// OpenFileDialog opens a file dialog and returns the selected file path
func (a *App) OpenFileDialog() (string, error) {
	return a.dialogService.OpenFileDialog()
}

// OpenDirectoryDialog opens a directory dialog and returns the selected directory path
func (a *App) OpenDirectoryDialog() (string, error) {
	return a.dialogService.OpenDirectoryDialog()
}

// SaveFileDialog opens a save file dialog
func (a *App) SaveFileDialog() (string, error) {
	return a.dialogService.SaveFileDialog()
}

// GetFileInfo returns file information for the given path
func (a *App) GetFileInfo(filePath string) (*models.LogFile, error) {
	return a.fileService.GetFileInfo(filePath)
}

// ReadLogFile reads the content of a log file
func (a *App) ReadLogFile(filePath string) (*models.LogContent, error) {
	return a.fileService.ReadLogFile(filePath)
}

// GetFilesInDirectory returns all log files in a directory
func (a *App) GetFilesInDirectory(dirPath string) ([]models.LogFile, error) {
	return a.fileService.GetFilesInDirectory(dirPath)
}

// SearchInFile searches for a pattern in the log file
func (a *App) SearchInFile(filePath, pattern string, caseSensitive bool) ([]models.SearchResult, error) {
	return a.searchService.SearchInFile(filePath, pattern, caseSensitive, false)
}

// SearchInFileWithRegex searches for a regex pattern in the log file
func (a *App) SearchInFileWithRegex(filePath, pattern string, caseSensitive bool) ([]models.SearchResult, error) {
	return a.searchService.SearchInFile(filePath, pattern, caseSensitive, true)
}

// GetRecentFiles returns a list of recently opened files
func (a *App) GetRecentFiles() ([]models.LogFile, error) {
	return a.recentFilesService.GetRecentFiles()
}

// AddRecentFile adds a file to the recent files list
func (a *App) AddRecentFile(filePath string) error {
	return a.recentFilesService.AddRecentFile(filePath)
}

// RemoveRecentFile removes a file from the recent files list
func (a *App) RemoveRecentFile(filePath string) error {
	return a.recentFilesService.RemoveRecentFile(filePath)
}

// ClearRecentFiles clears all recent files
func (a *App) ClearRecentFiles() error {
	return a.recentFilesService.ClearRecentFiles()
}

// ExportLogLines exports filtered log lines to a file
func (a *App) ExportLogLines(lines []string, exportPath string) error {
	return a.fileService.ExportLogLines(lines, exportPath)
}

// GetSystemInfo 获取系统资源信息
func (a *App) GetSystemInfo() (*models.SystemInfo, error) {
	return a.systemService.GetSystemInfo()
}

// HandleFileDrop 处理文件拖拽
func (a *App) HandleFileDrop(filePath string) error {
	// 使用现有的 ReadLogFile 方法读取文件
	logFile, err := a.ReadLogFile(filePath)
	if err != nil {
		// 发送错误事件到前端
		runtime.EventsEmit(a.ctx, "file-drop-error", map[string]interface{}{
			"error": err.Error(),
			"path":  filePath,
		})
		return err
	}

	// 发送成功事件到前端
	runtime.EventsEmit(a.ctx, "file-drop-success", map[string]interface{}{
		"file": logFile,
		"path": filePath,
	})

	return nil
}
