package main

import (
	"context"

	"LogTrawl/backend/models"
	"LogTrawl/backend/services"
)

// App struct
type App struct {
	ctx                 context.Context
	fileService         *services.FileService
	searchService       *services.SearchService
	recentFilesService  *services.RecentFilesService
	dialogService       *services.DialogService
	systemService       *services.SystemService
	analysisService     *services.AnalysisService
	fileSplitterService *services.FileSplitterService
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
	a.analysisService = services.NewAnalysisService()
	a.fileSplitterService = services.NewFileSplitterService()
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

// GetFileLineCount returns the number of lines in a file
func (a *App) GetFileLineCount(filePath string) (int, error) {
	return a.fileService.GetFileLineCount(filePath)
}

// GetFileDateRange returns the date range of a log file
func (a *App) GetFileDateRange(filePath string) (*models.DateRange, error) {
	return a.fileService.GetFileDateRange(filePath)
}

// ReadLogFile reads the content of a log file
func (a *App) ReadLogFile(filePath string) (*models.LogContent, error) {
	return a.fileService.ReadLogFile(filePath)
}

// ReadLogFileChunk reads a specific chunk of a log file
func (a *App) ReadLogFileChunk(filePath string, startLine, endLine int) (*models.LogContent, error) {
	return a.fileService.ReadLogFileChunk(filePath, startLine, endLine)
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

// AnalyzeLogFile 分析日志文件
func (a *App) AnalyzeLogFile(filePath string) (*models.AnalysisResult, error) {
	return a.analysisService.AnalyzeLogFile(filePath)
}

// AnalyzeSpecificIP 分析特定IP
func (a *App) AnalyzeSpecificIP(filePath, targetIP string) (*models.SpecificIPResult, error) {
	return a.analysisService.AnalyzeSpecificIP(filePath, targetIP)
}

// AnalyzeLogFileWithProgress 分析日志文件（支持进度）
func (a *App) AnalyzeLogFileWithProgress(filePath, sessionID string) (*models.AnalysisResult, error) {
	return a.analysisService.AnalyzeLogFileWithProgress(filePath, sessionID)
}

// StartAnalysis 开始分析会话
func (a *App) StartAnalysis(sessionID string) {
	a.analysisService.StartAnalysis(sessionID)
}

// CancelAnalysis 取消分析
func (a *App) CancelAnalysis(sessionID string) {
	a.analysisService.CancelAnalysis(sessionID)
}

// SplitFile 分片文件
func (a *App) SplitFile(options models.FileSplitOptions) (*models.FileSplitResult, error) {
	splitterOptions := services.SplitOptions{
		Strategy:     services.SplitStrategy(options.Strategy),
		FilePath:     options.FilePath,
		OutputDir:    options.OutputDir,
		DatePattern:  options.DatePattern,
		DaysPerFile:  options.DaysPerFile,
		SizePerFile:  options.SizePerFile,
		LinesPerFile: options.LinesPerFile,
	}

	result, err := a.fileSplitterService.SplitFile(splitterOptions)
	if err != nil {
		return nil, err
	}

	return &models.FileSplitResult{
		Success:     result.Success,
		Message:     result.Message,
		OutputFiles: result.OutputFiles,
		TotalFiles:  result.TotalFiles,
		TotalSize:   result.TotalSize,
	}, nil
}

// GetCommonDatePatterns 获取常用日期模式
func (a *App) GetCommonDatePatterns() []models.DatePattern {
	return a.fileSplitterService.GetCommonDatePatterns()
}

// OpenDirectoryDialogForSplit 为分片选择输出目录
func (a *App) OpenDirectoryDialogForSplit() (string, error) {
	return a.dialogService.OpenDirectoryDialog()
}
