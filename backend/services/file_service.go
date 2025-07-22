package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

	return &models.LogFile{
		ID:           fmt.Sprintf("%d", time.Now().UnixNano()),
		Name:         info.Name(),
		Path:         filePath,
		Size:         info.Size(),
		LastModified: info.ModTime(),
		IsOpen:       false,
	}, nil
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
	maxLines := 50000
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
