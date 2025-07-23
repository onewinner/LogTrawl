package services

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"LogTrawl/backend/models"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

// GetFileInfo returns file information for the given path
func (fs *FileService) GetFileInfo(filePath string) (*models.LogFile, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	// 对于大文件，计算行数
	var lineCount int
	if info.Size() > 20*1024*1024 { // 大于10MB的文件才计算行数
		lineCount, _ = fs.countLines(filePath)
	}

	return &models.LogFile{
		ID:           fmt.Sprintf("%d", time.Now().UnixNano()),
		Name:         info.Name(),
		Path:         filePath,
		Size:         info.Size(),
		Lines:        lineCount,
		LastModified: info.ModTime(),
		IsOpen:       false,
	}, nil
}

// GetFileLineCount returns the number of lines in a file
func (fs *FileService) GetFileLineCount(filePath string) (int, error) {
	return fs.countLines(filePath)
}

// GetFileDateRange returns the date range of a log file
func (fs *FileService) GetFileDateRange(filePath string) (*models.DateRange, error) {
	return fs.extractDateRange(filePath)
}

// countLines counts the number of lines in a file efficiently
func (fs *FileService) countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 使用更大的缓冲区提高性能
	const bufSize = 64 * 1024 // 64KB buffer
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, bufSize), bufSize)

	lineCount := 0
	for scanner.Scan() {
		lineCount++
		// 每1000行检查一次，避免长时间阻塞
		if lineCount%1000 == 0 {
			// 这里可以添加进度回调，但目前保持简单
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

// ReadLogFile reads the content of a log file
func (fs *FileService) ReadLogFile(filePath string) (*models.LogContent, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	
	// 限制最大行数以避免内存问题
	maxLines := 50000000
	lineCount := 0
	
	for scanner.Scan() && lineCount < maxLines {
		lines = append(lines, scanner.Text())
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &models.LogContent{
		Lines: lines,
		Total: len(lines),
	}, nil
}

// ReadLogFileChunk reads a specific chunk of lines from a log file
func (fs *FileService) ReadLogFileChunk(filePath string, startLine, endLine int) (*models.LogContent, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	currentLine := 0

	// 跳过到起始行
	for currentLine < startLine && scanner.Scan() {
		currentLine++
	}

	// 读取指定范围的行
	for currentLine < endLine && scanner.Scan() {
		lines = append(lines, scanner.Text())
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &models.LogContent{
		Lines: lines,
		Total: len(lines),
	}, nil
}

// GetFilesInDirectory returns all log files in a directory
func (fs *FileService) GetFilesInDirectory(dirPath string) ([]models.LogFile, error) {
	var logFiles []models.LogFile

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".log" || ext == ".txt" || ext == ".out" {
				logFiles = append(logFiles, models.LogFile{
					ID:           fmt.Sprintf("%d_%s", time.Now().UnixNano(), info.Name()),
					Name:         info.Name(),
					Path:         path,
					Size:         info.Size(),
					LastModified: info.ModTime(),
					IsOpen:       false,
				})
			}
		}

		return nil
	})

	return logFiles, err
}

// ExportLogLines exports filtered log lines to a file
func (fs *FileService) ExportLogLines(lines []string, exportPath string) error {
	file, err := os.Create(exportPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// extractDateRange extracts the date range from a log file
func (fs *FileService) extractDateRange(filePath string) (*models.DateRange, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 常用日期模式
	datePatterns := []struct {
		name    string
		pattern *regexp.Regexp
		layout  string
	}{
		{"YYYY-MM-DD HH:MM:SS", regexp.MustCompile(`(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})`), "2006-01-02 15:04:05"},
		{"YYYY-MM-DD", regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`), "2006-01-02"},
		{"DD/MMM/YYYY:HH:MM:SS", regexp.MustCompile(`(\d{2}/\w{3}/\d{4}:\d{2}:\d{2}:\d{2})`), "02/Jan/2006:15:04:05"},
		{"MM/DD/YYYY HH:MM:SS", regexp.MustCompile(`(\d{2}/\d{2}/\d{4}\s+\d{2}:\d{2}:\d{2})`), "01/02/2006 15:04:05"},
	}

	var startDate, endDate *time.Time
	var detectedPattern string
	scanner := bufio.NewScanner(file)

	// 读取前100行寻找第一个日期
	lineCount := 0
	for scanner.Scan() && lineCount < 100 {
		line := scanner.Text()
		lineCount++

		for _, dp := range datePatterns {
			if matches := dp.pattern.FindStringSubmatch(line); len(matches) > 1 {
				if parsedTime, err := time.Parse(dp.layout, matches[1]); err == nil {
					startDate = &parsedTime
					detectedPattern = dp.name
					break
				}
			}
		}

		if startDate != nil {
			break
		}
	}

	// 如果没找到日期，返回无日期信息的结果
	if startDate == nil {
		return &models.DateRange{
			HasDateInfo: false,
		}, nil
	}

	// 寻找最后一个日期 - 从文件末尾开始读取
	endDate = startDate // 默认为开始日期

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err == nil && fileInfo.Size() > 1024*1024 { // 大于1MB的文件才从末尾读取
		// 从文件末尾读取最后几KB
		const tailSize = 8192 // 8KB
		file.Seek(-tailSize, io.SeekEnd)

		tailScanner := bufio.NewScanner(file)
		var lastDate *time.Time

		for tailScanner.Scan() {
			line := tailScanner.Text()

			for _, dp := range datePatterns {
				if dp.name == detectedPattern { // 使用相同的模式
					if matches := dp.pattern.FindStringSubmatch(line); len(matches) > 1 {
						if parsedTime, err := time.Parse(dp.layout, matches[1]); err == nil {
							lastDate = &parsedTime
						}
					}
					break
				}
			}
		}

		if lastDate != nil {
			endDate = lastDate
		}
	}

	// 计算总天数
	totalDays := 0.0
	if endDate != nil && startDate != nil {
		duration := endDate.Sub(*startDate)
		totalDays = duration.Hours() / 24.0
		if totalDays < 0.1 {
			totalDays = 0.1 // 最小0.1天
		}
	}

	return &models.DateRange{
		StartDate:   startDate,
		EndDate:     endDate,
		DatePattern: detectedPattern,
		TotalDays:   totalDays,
		HasDateInfo: true,
	}, nil
}
